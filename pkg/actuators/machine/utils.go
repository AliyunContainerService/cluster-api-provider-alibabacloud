package machine

import (
	"fmt"
<<<<<<< HEAD
<<<<<<< HEAD
	"time"

=======
	"time"

	"k8s.io/klog/v2"

>>>>>>> e879a141 (alibabacloud machine-api provider)
	alibabacloudproviderv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudprovider/v1beta1"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"
	corev1 "k8s.io/api/core/v1"
<<<<<<< HEAD
=======
	providerconfigv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1alpha1"
	aliCloudClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/golang/glog"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
>>>>>>> 8dbd34ff (update project name)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

<<<<<<< HEAD
<<<<<<< HEAD
// upstreamMachineClusterIDLabel is the label that a machine must have to identify the cluster to which it belongs
const upstreamMachineClusterIDLabel = "sigs.k8s.io/cluster-api-cluster"

// supportedInstanceStates returns the list of states an ECS instance
func supportedInstanceStates() []string {
=======
// upstreamMachineClusterIDLabel is the label that a machine must have to identify the cluster to which it belongs
const upstreamMachineClusterIDLabel = "sigs.k8s.io/cluster-api-cluster"

// existingInstanceStates returns the list of states an ECS instance can be in
// while being considered "existing", i.e. mostly anything but "Terminated".
func existingInstanceStates() []string {
>>>>>>> e879a141 (alibabacloud machine-api provider)
	return []string{
		ECSInstanceStatusPending,
		ECSInstanceStatusStarting,
		ECSInstanceStatusRunning,
		ECSInstanceStatusStopping,
		ECSInstanceStatusStopped,
<<<<<<< HEAD
=======
// providerConfigFromMachine gets the machine provider config MachineSetSpec from the
// specified cluster-api MachineSpec.
func providerConfigFromMachine(machine *machinev1.Machine, codec *providerconfigv1.AlicloudProviderConfigCodec) (*providerconfigv1.AlibabaCloudMachineProviderConfig, error) {
	if machine.Spec.ProviderSpec.Value == nil {
		return nil, fmt.Errorf("unable to find machine provider config: Spec.ProviderSpec.Value is not set")
>>>>>>> c7e62b88 (fix testcase)
	}
}

<<<<<<< HEAD
=======
	}
}

>>>>>>> e879a141 (alibabacloud machine-api provider)
// validateMachine check the label that a machine must have to identify the cluster to which it belongs is present.
func validateMachine(machine machinev1.Machine) error {
	if machine.Labels[machinev1.MachineClusterIDLabel] == "" {
		return machinecontroller.InvalidMachineConfiguration("%v: missing %q label", machine.GetName(), machinev1.MachineClusterIDLabel)
<<<<<<< HEAD
=======
	var config providerconfigv1.AlibabaCloudMachineProviderConfig
	if err := codec.DecodeProviderSpec(&machine.Spec.ProviderSpec, &config); err != nil {
		return nil, err
>>>>>>> c7e62b88 (fix testcase)
	}

<<<<<<< HEAD
	return nil
=======
	return &config, nil
>>>>>>> ebdd9bd0 (update test case)
=======
	}

	return nil
>>>>>>> e879a141 (alibabacloud machine-api provider)
}

// getClusterID get cluster ID by machine.openshift.io/cluster-api-cluster label
func getClusterID(machine *machinev1.Machine) (string, bool) {
	clusterID, ok := machine.Labels[machinev1.MachineClusterIDLabel]

<<<<<<< HEAD
<<<<<<< HEAD
	if !ok {
		clusterID, ok = machine.Labels[upstreamMachineClusterIDLabel]
	}
	return clusterID, ok
}

func conditionSuccess() alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	return alibabacloudproviderv1.AlibabaCloudMachineProviderCondition{
		Type:    alibabacloudproviderv1.MachineCreation,
		Status:  corev1.ConditionTrue,
		Reason:  alibabacloudproviderv1.MachineCreationSucceeded,
		Message: "Machine successfully created",
=======
func shouldUpdateCondition(
	oldCondition providerconfigv1.AlibabaCloudMachineProviderCondition,
	newCondition providerconfigv1.AlibabaCloudMachineProviderCondition,
) bool {
	if oldCondition.Status != newCondition.Status ||
		oldCondition.Reason != newCondition.Reason ||
		oldCondition.Message != newCondition.Message {
		return true
>>>>>>> c7e62b88 (fix testcase)
	}
}

=======
	if !ok {
		clusterID, ok = machine.Labels[upstreamMachineClusterIDLabel]
	}
	return clusterID, ok
}

func conditionSuccess() alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	return alibabacloudproviderv1.AlibabaCloudMachineProviderCondition{
		Type:    alibabacloudproviderv1.MachineCreation,
		Status:  corev1.ConditionTrue,
		Reason:  alibabacloudproviderv1.MachineCreationSucceeded,
		Message: "Machine successfully created",
	}
}

>>>>>>> e879a141 (alibabacloud machine-api provider)
func conditionFailed() alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	return alibabacloudproviderv1.AlibabaCloudMachineProviderCondition{
		Type:   alibabacloudproviderv1.MachineCreation,
		Status: corev1.ConditionFalse,
		Reason: alibabacloudproviderv1.MachineCreationFailed,
	}
}

<<<<<<< HEAD
// setMachineProviderCondition sets the condition for the machine and
=======
// setAlibabaCloudMachineProviderCondition sets the condition for the machine and
>>>>>>> e879a141 (alibabacloud machine-api provider)
// returns the new slice of conditions.
// If the machine does not already have a condition with the specified type,
// a condition will be added to the slice
// If the machine does already have a condition with the specified type,
// the condition will be updated if either of the following are true.
<<<<<<< HEAD
<<<<<<< HEAD
func setMachineProviderCondition(condition alibabacloudproviderv1.AlibabaCloudMachineProviderCondition, conditions []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition) []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	now := metav1.Now()

	if existingCondition := findProviderCondition(conditions, condition.Type); existingCondition == nil {
		condition.LastProbeTime = now
		condition.LastTransitionTime = now
		conditions = append(conditions, condition)
=======
// 1) Requested Status is different than existing status.
// 2) requested Reason is different that existing one.
// 3) requested Message is different that existing one.
func setAliCloudMachineProviderCondition(conditions []providerconfigv1.AlibabaCloudMachineProviderCondition, newCondition providerconfigv1.AlibabaCloudMachineProviderCondition) []providerconfigv1.AlibabaCloudMachineProviderCondition {
	now := metav1.Now()
	currentCondition := findMachineProviderCondition(conditions, newCondition.Type)
	if currentCondition == nil {
		glog.Infof("Adding new provider condition %v", newCondition)
		conditions = append(
			conditions,
			providerconfigv1.AlibabaCloudMachineProviderCondition{
				Type:               newCondition.Type,
				Status:             newCondition.Status,
				Reason:             newCondition.Reason,
				Message:            newCondition.Message,
				LastTransitionTime: now,
				LastProbeTime:      now,
			},
		)
>>>>>>> c7e62b88 (fix testcase)
=======
func setAlibabaCloudMachineProviderCondition(condition alibabacloudproviderv1.AlibabaCloudMachineProviderCondition, conditions []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition) []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	now := metav1.Now()

	if existingCondition := findProviderCondition(conditions, condition.Type); existingCondition == nil {
		condition.LastProbeTime = now
		condition.LastTransitionTime = now
		conditions = append(conditions, condition)
>>>>>>> e879a141 (alibabacloud machine-api provider)
	} else {
		updateExistingCondition(&condition, existingCondition)
	}

	return conditions
}

<<<<<<< HEAD
<<<<<<< HEAD
func findProviderCondition(conditions []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition, conditionType alibabacloudproviderv1.AlibabaCloudMachineProviderConditionType) *alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	for i := range conditions {
		if conditions[i].Type == conditionType {
=======
// findMachineProviderCondition finds in the machine the condition that has the
// specified condition type. If none exists, then returns nil.
func findMachineProviderCondition(conditions []providerconfigv1.AlibabaCloudMachineProviderCondition, conditionType providerconfigv1.AliCloudMachineProviderConditionType) *providerconfigv1.AlibabaCloudMachineProviderCondition {
	for i, condition := range conditions {
		if condition.Type == conditionType {
>>>>>>> c7e62b88 (fix testcase)
=======
func findProviderCondition(conditions []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition, conditionType alibabacloudproviderv1.AlibabaCloudMachineProviderConditionType) *alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	for i := range conditions {
		if conditions[i].Type == conditionType {
>>>>>>> e879a141 (alibabacloud machine-api provider)
			return &conditions[i]
		}
	}
	return nil
}

<<<<<<< HEAD
<<<<<<< HEAD
func updateExistingCondition(newCondition, existingCondition *alibabacloudproviderv1.AlibabaCloudMachineProviderCondition) {
	if !shouldUpdateCondition(newCondition, existingCondition) {
		return
=======
// getRunningInstance returns the ECS instance for a given machine. If multiple instances match our machine,
// the most recently launched will be returned. If no instance exists, an error will be returned.
func getRunningInstance(machine *machinev1.Machine, client aliCloudClient.Client, regionId string) (*ecs.Instance, error) {
	instances, err := getRunningInstances(machine, client, regionId)
	if err != nil {
		return nil, err
	}
	if len(instances) == 0 {
		return nil, fmt.Errorf("no instance found for machine: %s", machine.Name)
	}

	sortInstances(instances)
	return instances[0], nil
}

// getRunningInstances returns all running instances that have a tag matching our machine name,
// and cluster ID.
func getRunningInstances(machine *machinev1.Machine, client aliCloudClient.Client, regionId string) ([]*ecs.Instance, error) {
	return getInstances(machine, client, "Running", regionId)
}

// getInstances returns all instances that have a tag matching our machine name,
// and cluster ID.
func getInstances(machine *machinev1.Machine, client aliCloudClient.Client, instanceStatus string, regionId string) ([]*ecs.Instance, error) {

	clusterID, ok := getClusterID(machine)
	if !ok {
		return []*ecs.Instance{}, fmt.Errorf("unable to get cluster ID for machine: %q", machine.Name)
	}

	describeInstancesRequest := ecs.CreateDescribeInstancesRequest()
	describeInstancesRequest.RegionId = regionId
	tags := clusterTagFilter(clusterID, machine.Name)
	describeInstancesRequest.Tag = &tags
	describeInstancesRequest.Status = instanceStatus
	describeInstancesRequest.Scheme = "https"

	result, err := client.DescribeInstances(describeInstancesRequest)
	if err != nil {
		return []*ecs.Instance{}, err
>>>>>>> ebdd9bd0 (update test case)
=======
func updateExistingCondition(newCondition, existingCondition *alibabacloudproviderv1.AlibabaCloudMachineProviderCondition) {
	if !shouldUpdateCondition(newCondition, existingCondition) {
		return
>>>>>>> e879a141 (alibabacloud machine-api provider)
	}

	if existingCondition.Status != newCondition.Status {
		existingCondition.LastTransitionTime = metav1.Now()
	}
	existingCondition.Status = newCondition.Status
	existingCondition.Reason = newCondition.Reason
	existingCondition.Message = newCondition.Message
	existingCondition.LastProbeTime = newCondition.LastProbeTime
}

func shouldUpdateCondition(newCondition, existingCondition *alibabacloudproviderv1.AlibabaCloudMachineProviderCondition) bool {
	return newCondition.Reason != existingCondition.Reason || newCondition.Message != existingCondition.Message
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
// WaitForResult wait func
func WaitForResult(name string, predicate func() (bool, interface{}, error), returnWhenError bool, delay int, timeout int) (interface{}, error) {
	endTime := time.Now().Add(time.Duration(timeout) * time.Second)
	delaySecond := time.Duration(delay) * time.Second
	for {
		// Execute the function
		satisfied, result, err := predicate()
<<<<<<< HEAD
=======
// deleteInstances terminates all provided instances with a single ECS request.
func deleteInstances(client aliCloudClient.Client, instances []*ecs.Instance) error {
	// Cleanup all older instances:
	for _, instance := range instances {
		glog.Infof("Cleaning up extraneous instance for machine: %v, state: %v, launchTime: %v", instance.InstanceId, instance.Status, instance.CreationTime)
		deleteInstanceRequest := ecs.CreateDeleteInstanceRequest()
		deleteInstanceRequest.InstanceId = instance.InstanceId
		deleteInstanceRequest.Force = requests.NewBoolean(true)
		deleteInstanceRequest.Scheme = "https"
		_, err := client.DeleteInstance(deleteInstanceRequest)
>>>>>>> ebdd9bd0 (update test case)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
		if err != nil {
			klog.Errorf("%s Invoke func %++s error %++v", name, "predicate func() (bool, error)", err)
			if returnWhenError {
				return result, err
			}
		}
		if satisfied {
			return result, nil
		}
		// Sleep
		time.Sleep(delaySecond)
		// If a timeout is set, and that's been exceeded, shut it down
		if timeout >= 0 && time.Now().After(endTime) {
			return nil, fmt.Errorf("wait for %s timeout", name)
		}
	}
}
