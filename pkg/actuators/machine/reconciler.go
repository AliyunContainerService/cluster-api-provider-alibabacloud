/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package machine

import (
	"context"
	"fmt"
	"time"

	"k8s.io/klog"

	"github.com/openshift/machine-api-operator/pkg/metrics"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"
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

// NewReconciler creating new Reconciler instance
func NewReconciler(scope *machineScope) *Reconciler {
	return &Reconciler{
		machineScope: scope,
	}
}

// Create creates machine if and only if machine exists, handled by cluster-api
func (r *Reconciler) Create(ctx context.Context) error {
	klog.Infof("%s: creating machine ", r.machine.Name)

	instance, err := r.CreateMachine(ctx)
	if err != nil {
		return err
	}

	klog.Infof("Created Machine %v", r.machine.Name)
	if err = r.setProviderID(instance); err != nil {
		return fmt.Errorf("failed to update machine object with providerID: %w", err)
	}

	if err = r.setMachineCloudProviderSpecifics(instance); err != nil {
		return fmt.Errorf("failed to set machine cloud provider specifics: %w", err)
	}

	_ = r.machineScope.setProviderStatus(instance, conditionSuccess())

	return nil
}

func (r *Reconciler) CreateMachine(ctx context.Context) (*ecs.Instance, error) {
	if err := validateMachine(*r.machine); err != nil {
		return nil, fmt.Errorf("%v: failed validating machine provider spec: %w", r.machine.GetName(), err)
	}

	userData, err := r.machineScope.getUserData()
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	instance, err := runInstances(r.machine, r.providerSpec, userData, r.alibabacloudClient)
	if err != nil {
		klog.Errorf("%s: error creating machine: %v", r.machine.Name, err)
		conditionFailed := conditionFailed()
		conditionFailed.Message = err.Error()
		_ = r.machineScope.setProviderStatus(nil, conditionFailed)
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	return instance, nil
}

// Update updates machine if and only if machine exists, handled by cluster-api
func (r *Reconciler) Update(ctx context.Context) error {
	klog.Infof("%s: updating machine", r.machine.Name)

	instance, err := r.UpdateMachine(ctx)
	if err != nil {
		return err
	}

	if err = r.setProviderID(instance); err != nil {
		return fmt.Errorf("failed to update machine object with providerID: %w", err)
	}

	if err = r.setMachineCloudProviderSpecifics(instance); err != nil {
		return fmt.Errorf("failed to set machine cloud provider specifics: %w", err)
	}

	if err = correctExistingTags(r.machine, r.providerSpec.RegionID, instance, r.alibabacloudClient); err != nil {
		return fmt.Errorf("failed to correct existing instance tags: %w", err)
	}

	klog.Infof("Updated machine %s", r.machine.Name)

	r.machineScope.setProviderStatus(instance, conditionSuccess())

	return r.requeueIfInstancePending(instance)
}

func (r *Reconciler) UpdateMachine(ctx context.Context) (*ecs.Instance, error) {
	if err := validateMachine(*r.machine); err != nil {
		return nil, fmt.Errorf("%v: failed validating machine provider spec: %v", r.machine.GetName(), err)
	}

	// Get all instances not deleted.
	existingInstances, err := r.getMachineInstances()
	if err != nil {
		metrics.RegisterFailedInstanceUpdate(&metrics.MachineLabels{
			Name:      r.machine.Name,
			Namespace: r.machine.Namespace,
			Reason:    err.Error(),
		})
		klog.Errorf("%s: error getting existing instances: %v", r.machine.Name, err)
		return nil, err
	}

	existingLen := len(existingInstances)
	if existingLen == 0 {
		if r.machine.Spec.ProviderID != nil && *r.machine.Spec.ProviderID != "" && (r.machine.Status.LastUpdated == nil || r.machine.Status.LastUpdated.Add(requeueAfterSeconds*time.Second).After(time.Now())) {
			klog.Infof("%s: Possible eventual-consistency discrepancy; returning an error to requeue", r.machine.Name)
			return nil, &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
		}

		klog.Warningf("%s: attempted to update machine but no instances found", r.machine.Name)

		// Update status to clear out machine details.
		r.machineScope.setProviderStatus(nil, conditionSuccess())

		return nil, &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterFatalSeconds * time.Second}
	}

	sortInstances(existingInstances)
	runningInstances := getRunningFromInstances(existingInstances)
	runningLen := len(runningInstances)
	var newestInstance *ecs.Instance

	if runningLen > 0 {
		klog.Infof("%s: found %d running instances for machine", r.machine.Name, runningLen)
		newestInstance = runningInstances[0]
	} else {
		newestInstance = existingInstances[0]
	}

	return newestInstance, nil
}

func (r *Reconciler) requeueIfInstancePending(instance *ecs.Instance) error {
	// If machine state is still pending, we will return an error to keep the controllers
	// attempting to update status until it hits a more permanent state. This will ensure
	// we get a public IP populated more quickly.
	if instance.Status == ECSInstanceStatusPending {
		klog.Infof("%s: Instance state still pending, returning an error to requeue", r.machine.Name)
		return &machinecontroller.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
	}

	return nil
}

// Delete deletes machine
func (r *Reconciler) Delete(ctx context.Context) error {
	klog.Infof("%s: deleting machine", r.machine.Name)

	if err := r.DeleteMachine(ctx); err != nil {
		return err
	}

	klog.Infof("Deleted machine %v", r.machine.Name)
	return nil
}

func (r *Reconciler) DeleteMachine(ctx context.Context) error {
	// Get all instances not terminated.
	existingInstances, err := r.getMachineInstances()
	if err != nil {
		metrics.RegisterFailedInstanceDelete(&metrics.MachineLabels{
			Name:      r.machine.Name,
			Namespace: r.machine.Namespace,
			Reason:    err.Error(),
		})
		klog.Errorf("%s: error getting existing instances: %v", r.machine.Name, err)
		return err
	}

	existingLen := len(existingInstances)
	klog.Infof("%s: found %d existing instances for machine", r.machine.Name, existingLen)
	if existingLen < 1 {
		klog.Warningf("%s: no instances found to delete for machine", r.machine.Name)
		return nil
	}

	// stopInstances stop all running instances ,if instance stauts not running ,skip stop it
	stoppedInstances, err := stopInstances(r.alibabacloudClient, r.providerSpec.RegionID, existingInstances)
	if err != nil {
		metrics.RegisterFailedInstanceDelete(&metrics.MachineLabels{
			Name:      r.machine.Name,
			Namespace: r.machine.Namespace,
			Reason:    err.Error(),
		})
		klog.Errorf("failed to stop instances %v error %v", existingInstances, err)
		return fmt.Errorf("failed to stop instaces: %w", err)
	}

	if len(stoppedInstances) == 1 {
		if stoppedInstances[0].Code == "200" && stoppedInstances[0].CurrentStatus != "" {
			r.machine.Annotations[machinecontroller.MachineInstanceStateAnnotationName] = stoppedInstances[0].CurrentStatus
		}
	}

	existingInstancesIds := make([]string, 0)
	for _, instance := range existingInstances {
		existingInstancesIds = append(existingInstancesIds, instance.InstanceId)
	}

	// wait for all instances stopped
	// Query the status of the instance until Stopped
	_, err = waitForInstancesStatus(r.alibabacloudClient, r.providerSpec.RegionID, existingInstancesIds, ECSInstanceStatusStopped, InstanceDefaultTimeout)
	if err != nil {
		metrics.RegisterFailedInstanceDelete(&metrics.MachineLabels{
			Name:      r.machine.Name,
			Namespace: r.machine.Namespace,
			Reason:    err.Error(),
		})
		klog.Errorf("failed to wait for  instances %v stopped: %v", existingInstancesIds, err)
		return fmt.Errorf("failed to wait for  instances stopped: %v", err)
	}

	// delete stoppted instances
	for _, instanceID := range existingInstancesIds {
		klog.Infof("delete %v instance", instanceID)
	}

	deleteInstancesRequest := ecs.CreateDeleteInstancesRequest()
	deleteInstancesRequest.Scheme = "https"
	deleteInstancesRequest.RegionId = r.providerSpec.RegionID
	deleteInstancesRequest.InstanceId = &existingInstancesIds

	deleteInstancsResponse, err := r.alibabacloudClient.DeleteInstances(deleteInstancesRequest)
	if err != nil {
		metrics.RegisterFailedInstanceDelete(&metrics.MachineLabels{
			Name:      r.machine.Name,
			Namespace: r.machine.Namespace,
			Reason:    err.Error(),
		})
		klog.Errorf("failed to delete instances %v error %v", deleteInstancesRequest, err)
		return fmt.Errorf("failed to delete instaces: %v", err)
	}

	klog.V(3).Infof("Delete instance response %v", deleteInstancsResponse)
	return nil
}

// Exists checks if machine exists
func (r *Reconciler) Exists(ctx context.Context) (bool, error) {
	// Get all instances not terminated.
	existingInstances, err := r.getMachineInstances()
	if err != nil {
		// Reporting as update here, as successfull return value from the method
		// later indicases that an instance update flow will be executed.
		metrics.RegisterFailedInstanceUpdate(&metrics.MachineLabels{
			Name:      r.machine.Name,
			Namespace: r.machine.Namespace,
			Reason:    err.Error(),
		})
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

// setProviderID adds providerID in the machine spec
func (r *Reconciler) setProviderID(instance *ecs.Instance) error {
	existingProviderID := r.machine.Spec.ProviderID
	if instance == nil {
		return nil
	}

	providerID := fmt.Sprintf("alibabacloud:///%s/%s/%s", instance.RegionId, instance.ZoneId, instance.InstanceId)
	// If resourceGroupId is not empty, set to providerId
	if instance.ResourceGroupId != "" {
		providerID = fmt.Sprintf("alibabacloud:///%s/%s/%s/%s", instance.ResourceGroupId, instance.RegionId, instance.ZoneId, instance.InstanceId)
	}

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

	r.machine.Labels[machinecontroller.MachineRegionLabelName] = instance.RegionId

	r.machine.Labels[machinecontroller.MachineAZLabelName] = instance.ZoneId

	if instance.InstanceType != "" {
		r.machine.Labels[machinecontroller.MachineInstanceTypeLabelName] = instance.InstanceType
	}

	if instance.Status != "" {
		r.machine.Annotations[machinecontroller.MachineInstanceStateAnnotationName] = instance.Status
	}

	return nil
}

func (r *Reconciler) getMachineInstances() ([]*ecs.Instance, error) {
	if r.providerStatus.InstanceID != nil && *r.providerStatus.InstanceID != "" {
		i, err := getExistingInstanceByID(*r.providerStatus.InstanceID, r.providerSpec.RegionID, r.alibabacloudClient)
		if err != nil {
			klog.Warningf("%s: Failed to find existing instance by id %s: %v", r.machine.Name, *r.providerStatus.InstanceID, err)
		} else {
			klog.Infof("%s: Found instance by id: %s", r.machine.Name, *r.providerStatus.InstanceID)
			return []*ecs.Instance{i}, nil
		}
	}

	return getExistingInstances(r.machine, r.providerSpec.RegionID, r.alibabacloudClient)
}
