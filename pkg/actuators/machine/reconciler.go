package machine

import (
	"fmt"
	"time"

	alibabacloudproviderv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1beta1"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"
	"github.com/openshift/machine-api-operator/pkg/metrics"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	errorutil "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog"
)

const (
	requeueAfterSeconds      = 20
	requeueAfterFatalSeconds = 180
	masterLabel              = "node-role.kubernetes.io/master"
)

// Reconciler runs the logic to reconciles a machine resource towards its desired state
type Reconciler struct {
	*machineScope
}

func newReconciler(scope *machineScope) *Reconciler {
	return &Reconciler{
		machineScope: scope,
	}
}

// create creates machine if it does not exists.
func (r *Reconciler) create() error {
	klog.Infof("%s: creating machine", r.machine.Name)

	if err := validateMachine(*r.machine); err != nil {
		return fmt.Errorf("%v: failed validating machine provider spec: %w", r.machine.GetName(), err)
	}

	// TODO: remove 45 - 59, this logic is not needed anymore
	// We explicitly do NOT want to remove stopped masters.
	isMaster, err := r.isMaster()
	if err != nil {
		// Unable to determine if a machine is a master machine.
		// Yet, it's only used to delete stopped machines that are not masters.
		// So we can safely continue to create a new machine since in the worst case
		// we just don't delete any stopped machine.
		klog.Errorf("%s: Error determining if machine is master: %v", r.machine.Name, err)
	} else {
		if !isMaster {
			// Prevent having a lot of stopped nodes sitting around.
			if err = removeStoppedMachine(r.machine, r.aliClient); err != nil {
				return fmt.Errorf("unable to remove stopped machines: %w", err)
			}
		}
	}

	userData, err := r.machineScope.getUserData()
	if err != nil {
		return fmt.Errorf("failed to get user data: %w", err)
	}

	instance, err := createInstance(r.machine, r.providerSpec, userData, r.aliClient)
	if err != nil {
		klog.Errorf("%s: error creating machine: %v", r.machine.Name, err)
		conditionFailed := conditionFailed()
		conditionFailed.Message = err.Error()
		r.machineScope.setProviderStatus(nil, conditionFailed)
		return fmt.Errorf("failed to launch instance: %w", err)
	}

	// if err = r.updateLoadBalancers(instance); err != nil {
	// 	metrics.RegisterFailedInstanceCreate(&metrics.MachineLabels{
	// 		Name:      r.machine.Name,
	// 		Namespace: r.machine.Namespace,
	// 		Reason:    err.Error(),
	// 	})
	// 	return fmt.Errorf("failed to updated update load balancers: %w", err)
	// }

	klog.Infof("Created Machine %v", r.machine.Name)
	if err = r.setProviderID(instance); err != nil {
		return fmt.Errorf("failed to update machine object with providerID: %w", err)
	}

	if err = r.setMachineCloudProviderSpecifics(instance); err != nil {
		return fmt.Errorf("failed to set machine cloud provider specifics: %w", err)
	}

	r.machineScope.setProviderStatus(instance, conditionSuccess())

	return r.requeueIfInstanceStarting(instance)
}

// delete deletes machine
func (r *Reconciler) delete() error {
	klog.Infof("%s: deleting machine", r.machine.Name)

	// Get all instances not terminated.
	existingInstances, err := r.getMachineInstances()
	if err != nil {
		klog.Errorf("%s: error getting existing instances: %v", r.machine.Name, err)
		return err
	}

	existingLen := len(existingInstances)
	klog.Infof("%s: found %d existing instances for machine", r.machine.Name, existingLen)
	if existingLen == 0 {
		klog.Warningf("%s: no instances found to delete for machine", r.machine.Name)
		return nil
	}

	err = deleteInstances(r.aliClient, existingInstances)
	if err != nil {
		metrics.RegisterFailedInstanceDelete(&metrics.MachineLabels{
			Name:      r.machine.Name,
			Namespace: r.machine.Namespace,
			Reason:    err.Error(),
		})
		return fmt.Errorf("failed to delete instaces: %w", err)
	}

	// if len(terminatingInstances) == 1 {
	// 	if terminatingInstances[0] != nil && terminatingInstances[0].CurrentState != nil && terminatingInstances[0].CurrentState.Name != nil {
	// 		r.machine.Annotations[machinecontroller.MachineInstanceStateAnnotationName] = StringValue(terminatingInstances[0].CurrentState.Name)
	// 	}
	// }

	klog.Infof("Deleted machine %v", r.machine.Name)

	return nil
}

// update finds a vm and reconciles the machine resource status against it.
func (r *Reconciler) update() error {
	klog.Infof("%s: updating machine", r.machine.Name)

	if err := validateMachine(*r.machine); err != nil {
		return fmt.Errorf("%v: failed validating machine provider spec: %v", r.machine.GetName(), err)
	}

	// Get all instances not terminated.
	existingInstances, err := r.getMachineInstances()
	if err != nil {
		klog.Errorf("%s: error getting existing instances: %v", r.machine.Name, err)
		return err
	}

	existingLen := len(existingInstances)
	if existingLen == 0 {
		if r.machine.Spec.ProviderID != nil && *r.machine.Spec.ProviderID != "" && (r.machine.Status.LastUpdated == nil || r.machine.Status.LastUpdated.Add(requeueAfterSeconds*time.Second).After(time.Now())) {
			klog.Infof("%s: Possible eventual-consistency discrepancy; returning an error to requeue", r.machine.Name)
			return &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
		}

		klog.Warningf("%s: attempted to update machine but no instances found", r.machine.Name)

		// Update status to clear out machine details.
		r.machineScope.setProviderStatus(nil, conditionSuccess())
		// This is an unrecoverable error condition.  We should delay to
		// minimize unnecessary API calls.
		return &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterFatalSeconds * time.Second}
	}

	sortInstances(existingInstances)
	runningInstances := getRunningFromInstances(existingInstances)
	runningLen := len(runningInstances)
	var newestInstance *ecs.Instance

	if runningLen > 0 {
		// It would be very unusual to have more than one here, but it is
		// possible if someone manually provisions a machine with same tag name.
		klog.Infof("%s: found %d running instances for machine", r.machine.Name, runningLen)
		newestInstance = runningInstances[0]

		// err = r.updateLoadBalancers(newestInstance)
		if err != nil {
			metrics.RegisterFailedInstanceUpdate(&metrics.MachineLabels{
				Name:      r.machine.Name,
				Namespace: r.machine.Namespace,
				Reason:    err.Error(),
			})
			return fmt.Errorf("failed to updated update load balancers: %w", err)
		}
	} else {
		// Didn't find any running instances, just newest existing one.
		// In most cases, there should only be one existing Instance.
		newestInstance = existingInstances[0]
	}

	if err = r.setProviderID(newestInstance); err != nil {
		return fmt.Errorf("failed to update machine object with providerID: %w", err)
	}

	if err = r.setMachineCloudProviderSpecifics(newestInstance); err != nil {
		return fmt.Errorf("failed to set machine cloud provider specifics: %w", err)
	}

	klog.Infof("Updated machine %s", r.machine.Name)

	r.machineScope.setProviderStatus(newestInstance, conditionSuccess())

	return r.requeueIfInstanceStarting(newestInstance)
}

// exists returns true if machine exists.
func (r *Reconciler) exists() (bool, error) {
	// Get all instances not terminated.
	existingInstances, err := r.getMachineInstances()
	if err != nil {
		klog.Errorf("%s: error getting existing instances: %v", r.machine.Name, err)
		return false, err
	}

	if len(existingInstances) == 0 {
		if r.machine.Spec.ProviderID != nil && *r.machine.Spec.ProviderID != "" && (r.machine.Status.LastUpdated == nil || r.machine.Status.LastUpdated.Add(requeueAfterSeconds*time.Second).After(time.Now())) {
			klog.Infof("%s: Possible eventual-consistency discrepancy; returning an error to requeue", r.machine.Name)
			return false, &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
		}

		klog.Infof("%s: Instance does not exist", r.machine.Name)
		return false, nil
	}

	return existingInstances[0] != nil, err
}

// isMaster returns true if the machine is part of a cluster's control plane
func (r *Reconciler) isMaster() (bool, error) {
	if r.machine.Status.NodeRef == nil {
		klog.Errorf("NodeRef not found in machine %s", r.machine.Name)
		return false, nil
	}
	node := &corev1.Node{}
	nodeKey := types.NamespacedName{
		Namespace: r.machine.Status.NodeRef.Namespace,
		Name:      r.machine.Status.NodeRef.Name,
	}

	err := r.client.Get(r.Context, nodeKey, node)
	if err != nil {
		return false, fmt.Errorf("failed to get node from machine %s", r.machine.Name)
	}

	if _, exists := node.Labels[masterLabel]; exists {
		return true, nil
	}
	return false, nil
}

// updateLoadBalancers adds a given machine instance to the load balancers specified in its provider config
func (r *Reconciler) updateLoadBalancers(instance *ecs.Instance) error {
	if len(r.providerSpec.LoadBalancers) == 0 {
		klog.V(4).Infof("%s: Instance %q has no load balancers configured. Skipping", r.machine.Name, instance.InstanceId)
		return nil
	}
	errs := []error{}
	classicLoadBalancerNames := []string{}
	networkLoadBalancerNames := []string{}
	for _, loadBalancerRef := range r.providerSpec.LoadBalancers {
		switch loadBalancerRef.Type {
		case alibabacloudproviderv1.NetworkLoadBalancerType:
			networkLoadBalancerNames = append(networkLoadBalancerNames, loadBalancerRef.Name)
		case alibabacloudproviderv1.ClassicLoadBalancerType:
			classicLoadBalancerNames = append(classicLoadBalancerNames, loadBalancerRef.Name)
		}
	}

	var err error
	if len(classicLoadBalancerNames) > 0 {
		//err := registerWithClassicLoadBalancers(r.awsClient, classicLoadBalancerNames, instance)
		if err != nil {
			klog.Errorf("%s: Failed to register classic load balancers: %v", r.machine.Name, err)
			errs = append(errs, err)
		}
	}
	if len(networkLoadBalancerNames) > 0 {
		// err = registerWithNetworkLoadBalancers(r.awsClient, networkLoadBalancerNames, instance)
		if err != nil {
			klog.Errorf("%s: Failed to register network load balancers: %v", r.machine.Name, err)
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errorutil.NewAggregate(errs)
	}
	return nil
}

// setProviderID adds providerID in the machine spec
func (r *Reconciler) setProviderID(instance *ecs.Instance) error {
	existingProviderID := r.machine.Spec.ProviderID
	if instance == nil {
		return nil
	}
	regionId := ""
	if instance.RegionId != "" {
		regionId = instance.RegionId
	}
	providerID := fmt.Sprintf("alicloud://%s.%s", regionId, instance.InstanceId)

	if existingProviderID != nil && *existingProviderID == providerID {
		klog.Infof("%s: ProviderID already set in the machine Spec with value:%s", r.machine.Name, *existingProviderID)
		return nil
	}
	r.machine.Spec.ProviderID = &providerID
	klog.Infof("%s: ProviderID set at machine spec: %s", r.machine.Name, providerID)
	return nil
}

func (r *Reconciler) setMachineCloudProviderSpecifics(instance *ecs.Instance) error {
	if instance == nil {
		return nil
	}

	if r.machine.Labels == nil {
		r.machine.Labels = make(map[string]string)
	}

	if r.machine.Spec.Labels == nil {
		r.machine.Spec.Labels = make(map[string]string)
	}

	if r.machine.Annotations == nil {
		r.machine.Annotations = make(map[string]string)
	}

	// Reaching to machine provider config since the region is not directly
	// providing by *ec2.Instance object
	machineProviderConfig, err := alibabacloudproviderv1.ProviderSpecFromRawExtension(r.machine.Spec.ProviderSpec.Value)
	if err != nil {
		return fmt.Errorf("error decoding MachineProviderConfig: %w", err)
	}

	r.machine.Labels[machinecontroller.MachineRegionLabelName] = machineProviderConfig.Placement.Region

	if instance.ZoneId != "" {
		r.machine.Labels[machinecontroller.MachineAZLabelName] = instance.ZoneId
	}

	if instance.InstanceType != "" {
		r.machine.Labels[machinecontroller.MachineInstanceTypeLabelName] = instance.InstanceType
	}

	if instance.Status != "" {
		r.machine.Annotations[machinecontroller.MachineInstanceStateAnnotationName] = instance.Status
	}

	if instance.IsSpot {
		// Label on the Spec so that it is propogated to the Node
		r.machine.Spec.Labels[machinecontroller.MachineInterruptibleInstanceLabelName] = ""
	}

	return nil
}

func (r *Reconciler) requeueIfInstanceStarting(instance *ecs.Instance) error {
	// If machine state is still pending, we will return an error to keep the controllers
	// attempting to update status until it hits a more permanent state. This will ensure
	// we get a public IP populated more quickly.
	if instance.Status != "" && instance.Status == InstanceStateNameStarting {
		klog.Infof("%s: Instance state still pending, returning an error to requeue", r.machine.Name)
		return &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
	}

	return nil
}

func (r *Reconciler) getMachineInstances() ([]*ecs.Instance, error) {
	// If there is a non-empty instance ID, search using that, otherwise
	// fallback to filtering based on tags.
	if r.providerStatus.InstanceID != nil && *r.providerStatus.InstanceID != "" {
		i, err := getExistingInstanceByID(*r.providerStatus.InstanceID, r.aliClient)
		if err != nil {
			klog.Warningf("%s: Failed to find existing instance by id %s: %v", r.machine.Name, *r.providerStatus.InstanceID, err)
		} else {
			klog.Infof("%s: Found instance by id: %s", r.machine.Name, *r.providerStatus.InstanceID)
			return []*ecs.Instance{i}, nil
		}
	}

	return getExistingInstances(r.machine, r.aliClient)
}
