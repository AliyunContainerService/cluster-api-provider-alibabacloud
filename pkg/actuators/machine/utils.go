/*
Copyright 2018 The Kubernetes Authors.

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
	"fmt"
	providerconfigv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1alpha1"
	aliCloudClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/golang/glog"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// providerConfigFromMachine gets the machine provider config MachineSetSpec from the
// specified cluster-api MachineSpec.
func providerConfigFromMachine(machine *machinev1.Machine, codec *providerconfigv1.AlicloudProviderConfigCodec) (*providerconfigv1.AlibabaCloudMachineProviderConfig, error) {
	if machine.Spec.ProviderSpec.Value == nil {
		return nil, fmt.Errorf("unable to find machine provider config: Spec.ProviderSpec.Value is not set")
	}

	var config providerconfigv1.AlibabaCloudMachineProviderConfig
	if err := codec.DecodeProviderSpec(&machine.Spec.ProviderSpec, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

const (
	machineCreationSucceedReason  = "MachineCreationSucceeded"
	machineCreationSucceedMessage = "machine successfully created"
	machineCreationFailedReason   = "MachineCreationFailed"
)

func shouldUpdateCondition(
	oldCondition providerconfigv1.AlibabaCloudMachineProviderCondition,
	newCondition providerconfigv1.AlibabaCloudMachineProviderCondition,
) bool {
	if oldCondition.Status != newCondition.Status ||
		oldCondition.Reason != newCondition.Reason ||
		oldCondition.Message != newCondition.Message {
		return true
	}
	return false
}

// setAliCloudMachineProviderCondition sets the condition for the machine and
// returns the new slice of conditions.
// If the machine does not already have a condition with the specified type,
// a condition will be added to the slice.
// If the machine does already have a condition with the specified type,
// the condition will be updated if either of the following are true.
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
	} else {
		if shouldUpdateCondition(
			*currentCondition,
			newCondition,
		) {
			glog.Infof("Updating provider condition %v", newCondition)
			if currentCondition.Status != newCondition.Status {
				currentCondition.LastTransitionTime = now
			}
			currentCondition.Status = newCondition.Status
			currentCondition.Reason = newCondition.Reason
			currentCondition.Message = newCondition.Message
			currentCondition.LastProbeTime = now
		}
	}
	return conditions
}

// findMachineProviderCondition finds in the machine the condition that has the
// specified condition type. If none exists, then returns nil.
func findMachineProviderCondition(conditions []providerconfigv1.AlibabaCloudMachineProviderCondition, conditionType providerconfigv1.AliCloudMachineProviderConditionType) *providerconfigv1.AlibabaCloudMachineProviderCondition {
	for i, condition := range conditions {
		if condition.Type == conditionType {
			return &conditions[i]
		}
	}
	return nil
}

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
	}

	instances := make([]*ecs.Instance, 0, len(result.Instances.Instance))
	for _, instance := range result.Instances.Instance {
		instances = append(instances, &instance)
	}

	return instances, nil
}

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
		if err != nil {
			glog.Errorf("Error delete instances: %v", err)
			return fmt.Errorf("error terminating instances: %v", err)
		}
	}

	return nil
}

// getRunningFromInstances returns all running instances from a list of instances.
func getRunningFromInstances(instances []*ecs.Instance) []*ecs.Instance {
	var runningInstances []*ecs.Instance
	for _, instance := range instances {
		if instance.Status == "Running" {
			runningInstances = append(runningInstances, instance)
		}
	}
	return runningInstances
}

func getExistingInstances(machine *machinev1.Machine, client aliCloudClient.Client, regionId string) ([]*ecs.Instance, error) {
	return getInstances(machine, client, "Running", regionId)
}
