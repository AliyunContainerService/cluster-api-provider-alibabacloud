package machine

import (
	"fmt"
	"time"

	alibabacloudproviderv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudprovider/v1beta1"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// upstreamMachineClusterIDLabel is the label that a machine must have to identify the cluster to which it belongs
const upstreamMachineClusterIDLabel = "sigs.k8s.io/cluster-api-cluster"

// supportedInstanceStates returns the list of states an ECS instance
func supportedInstanceStates() []string {
	return []string{
		ECSInstanceStatusPending,
		ECSInstanceStatusStarting,
		ECSInstanceStatusRunning,
		ECSInstanceStatusStopping,
		ECSInstanceStatusStopped,
	}
}

// validateMachine check the label that a machine must have to identify the cluster to which it belongs is present.
func validateMachine(machine machinev1.Machine) error {
	if machine.Labels[machinev1.MachineClusterIDLabel] == "" {
		return machinecontroller.InvalidMachineConfiguration("%v: missing %q label", machine.GetName(), machinev1.MachineClusterIDLabel)
	}

	return nil
}

// getClusterID get cluster ID by machine.openshift.io/cluster-api-cluster label
func getClusterID(machine *machinev1.Machine) (string, bool) {
	clusterID, ok := machine.Labels[machinev1.MachineClusterIDLabel]

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

func conditionFailed() alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	return alibabacloudproviderv1.AlibabaCloudMachineProviderCondition{
		Type:   alibabacloudproviderv1.MachineCreation,
		Status: corev1.ConditionFalse,
		Reason: alibabacloudproviderv1.MachineCreationFailed,
	}
}

// setMachineProviderCondition sets the condition for the machine and
// returns the new slice of conditions.
// If the machine does not already have a condition with the specified type,
// a condition will be added to the slice
// If the machine does already have a condition with the specified type,
// the condition will be updated if either of the following are true.
func setMachineProviderCondition(condition alibabacloudproviderv1.AlibabaCloudMachineProviderCondition, conditions []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition) []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	now := metav1.Now()

	if existingCondition := findProviderCondition(conditions, condition.Type); existingCondition == nil {
		condition.LastProbeTime = now
		condition.LastTransitionTime = now
		conditions = append(conditions, condition)
	} else {
		updateExistingCondition(&condition, existingCondition)
	}

	return conditions
}

func findProviderCondition(conditions []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition, conditionType alibabacloudproviderv1.AlibabaCloudMachineProviderConditionType) *alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return &conditions[i]
		}
	}
	return nil
}

func updateExistingCondition(newCondition, existingCondition *alibabacloudproviderv1.AlibabaCloudMachineProviderCondition) {
	if !shouldUpdateCondition(newCondition, existingCondition) {
		return
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

// WaitForResult wait func
func WaitForResult(name string, predicate func() (bool, interface{}, error), returnWhenError bool, delay int, timeout int) (interface{}, error) {
	endTime := time.Now().Add(time.Duration(timeout) * time.Second)
	delaySecond := time.Duration(delay) * time.Second
	for {
		// Execute the function
		satisfied, result, err := predicate()
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
