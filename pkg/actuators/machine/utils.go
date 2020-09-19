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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	alibabacloudproviderv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1beta1"
	aliClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"
	"k8s.io/klog"
)

// upstreamMachineClusterIDLabel is the label that a machine must have to identify the cluster to which it belongs
const upstreamMachineClusterIDLabel = "sigs.k8s.io/cluster-api-cluster"

const (
	// InstanceStateNameStarting is a InstanceStateName enum value
	InstanceStateNameStarting = "Starting"

	// InstanceStateNameRunning is a InstanceStateName enum value
	InstanceStateNameRunning = "Running"

	// InstanceStateNameStopping is a InstanceStateName enum value
	InstanceStateNameStopping = "Stopping"

	// InstanceStateNameStopped is a InstanceStateName enum value
	InstanceStateNameStopped = "Stopped"
)

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// StringSlice converts a slice of string values into a slice of
// string pointers
func StringSlice(src []string) []*string {
	dst := make([]*string, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}
	return dst
}

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}

// BoolValue returns the value of the bool pointer passed in or
// false if the pointer is nil.
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

// existingInstanceStates returns the list of states an EC2 instance can be in
// while being considered "existing", i.e. mostly anything but "Terminated".
func existingInstanceStates() []*string {
	return []*string{
		String(InstanceStateNameRunning),
		String(InstanceStateNameStarting),
		String(InstanceStateNameStopped),
		String(InstanceStateNameStopping),
	}
}

// providerConfigFromMachine gets the machine provider config MachineSetSpec from the
// specified cluster-api MachineSpec.
// func providerConfigFromMachine(machine *machinev1.Machine, codec *providerconfigv1.AlicloudProviderConfigCodec) (*providerconfigv1.AlibabaCloudMachineProviderConfig, error) {
// 	if machine.Spec.ProviderSpec.Value == nil {
// 		return nil, fmt.Errorf("unable to find machine provider config: Spec.ProviderSpec.Value is not set")
// 	}

// 	var config providerconfigv1.AlibabaCloudMachineProviderConfig
// 	if err := codec.DecodeProviderSpec(&machine.Spec.ProviderSpec, &config); err != nil {
// 		return nil, err
// 	}

// 	return &config, nil
// }

const (
	machineCreationSucceedReason  = "MachineCreationSucceeded"
	machineCreationSucceedMessage = "machine successfully created"
	machineCreationFailedReason   = "MachineCreationFailed"
)

func shouldUpdateCondition(
	oldCondition alibabacloudproviderv1.AlibabaCloudMachineProviderCondition,
	newCondition alibabacloudproviderv1.AlibabaCloudMachineProviderCondition,
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
func setAliCloudMachineProviderCondition(condition alibabacloudproviderv1.AlibabaCloudMachineProviderCondition, conditions []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition) []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	now := metav1.Now()
	currentCondition := findMachineProviderCondition(conditions, condition.Type)
	if currentCondition == nil {
		klog.Infof("Adding new provider condition %v", condition)
		conditions = append(
			conditions,
			alibabacloudproviderv1.AlibabaCloudMachineProviderCondition{
				Type:               condition.Type,
				Status:             condition.Status,
				Reason:             condition.Reason,
				Message:            condition.Message,
				LastTransitionTime: now,
				LastProbeTime:      now,
			},
		)
	} else {
		if shouldUpdateCondition(
			*currentCondition,
			condition,
		) {
			klog.Infof("Updating provider condition %v", condition)
			if currentCondition.Status != condition.Status {
				currentCondition.LastTransitionTime = now
			}
			currentCondition.Status = condition.Status
			currentCondition.Reason = condition.Reason
			currentCondition.Message = condition.Message
			currentCondition.LastProbeTime = now
		}
	}
	return conditions
}

// findMachineProviderCondition finds in the machine the condition that has the
// specified condition type. If none exists, then returns nil.
func findMachineProviderCondition(conditions []alibabacloudproviderv1.AlibabaCloudMachineProviderCondition, conditionType alibabacloudproviderv1.AlibabaCloudMachineProviderConditionType) *alibabacloudproviderv1.AlibabaCloudMachineProviderCondition {
	for i, condition := range conditions {
		if condition.Type == conditionType {
			return &conditions[i]
		}
	}
	return nil
}

// getRunningInstance returns the ECS instance for a given machine. If multiple instances match our machine,
// the most recently launched will be returned. If no instance exists, an error will be returned.
func getRunningInstance(machine *machinev1.Machine, client aliClient.Client) (*ecs.Instance, error) {
	instances, err := getRunningInstances(machine, client)
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
func getRunningInstances(machine *machinev1.Machine, client aliClient.Client) ([]*ecs.Instance, error) {
	return getInstances(machine, client, StringSlice([]string{InstanceStateNameRunning}))
}

// getInstances returns all instances that have a tag matching our machine name,
// and cluster ID.
func getInstances(machine *machinev1.Machine, client aliClient.Client, instanceStateFilter []*string) ([]*ecs.Instance, error) {

	clusterID, ok := getClusterID(machine)
	if !ok {
		return []*ecs.Instance{}, fmt.Errorf("unable to get cluster ID for machine: %q", machine.Name)
	}

	describeInstancesRequest := ecs.CreateDescribeInstancesRequest()
	tags := clusterTagFilter(clusterID, machine.Name)
	describeInstancesRequest.Tag = &tags
	describeInstancesRequest.Scheme = "https"

	result, err := client.DescribeInstances(describeInstancesRequest)
	if err != nil {
		return []*ecs.Instance{}, err
	}

	instances := make([]*ecs.Instance, 0, len(result.Instances.Instance))
	for _, instance := range result.Instances.Instance {
		err := instanceHasAllowedState(&instance, instanceStateFilter)
		if err != nil {
			klog.Errorf("Excluding instance matching %s: %v", machine.Name, err)
		} else {
			instances = append(instances, &instance)
		}

	}

	return instances, nil
}

// deleteInstances terminates all provided instances with a single ECS request.
func deleteInstances(client aliClient.Client, instances []*ecs.Instance) error {
	// Cleanup all older instances:
	for _, instance := range instances {
		if instance.InstanceChargeType == "PrePaid" {
			modifyInstanceChargeType(client, instance)
		}
		klog.Infof("Cleaning up extraneous instance for machine: %v, state: %v, launchTime: %v", instance.InstanceId, instance.Status, instance.CreationTime)
		deleteInstanceRequest := ecs.CreateDeleteInstanceRequest()
		deleteInstanceRequest.InstanceId = instance.InstanceId
		deleteInstanceRequest.Force = requests.NewBoolean(true)
		deleteInstanceRequest.Scheme = "https"
		_, err := client.DeleteInstance(deleteInstanceRequest)
		if err != nil {
			klog.Errorf("Error delete instances: %v", err)
			return fmt.Errorf("error terminating instances: %v", err)
		}
	}

	return nil
}

//modifyInstanceChargeType convert chargeType to PostPaid
func modifyInstanceChargeType(client aliClient.Client, instance *ecs.Instance) error {

	klog.Infof("Convert instance charge type for machine: %v, state: %v, launchTime: %v", instance.InstanceId, instance.Status, instance.CreationTime)
	modifyrequest := ecs.CreateModifyInstanceChargeTypeRequest()
	instancesIds, _ := json.Marshal([]string{instance.InstanceId})
	modifyrequest.InstanceIds = string(instancesIds)
	modifyrequest.InstanceChargeType = "PostPaid"
	modifyrequest.Scheme = "https"

	_, err := client.ModifyInstanceChargeType(modifyrequest)
	if err != nil {
		klog.Errorf("Error convert instances chargetype: %v", err)
		return fmt.Errorf("error convert instances chargetype: %v", err)
	}

	return nil
}

// getRunningFromInstances returns all running instances from a list of instances.
func getRunningFromInstances(instances []*ecs.Instance) []*ecs.Instance {
	var runningInstances []*ecs.Instance
	for _, instance := range instances {
		if instance.Status == InstanceStateNameRunning {
			runningInstances = append(runningInstances, instance)
		}
	}
	return runningInstances
}

func getExistingInstances(machine *machinev1.Machine, client aliClient.Client) ([]*ecs.Instance, error) {
	return getInstances(machine, client, existingInstanceStates())
}

func getExistingInstanceByID(id string, client aliClient.Client) (*ecs.Instance, error) {
	return getInstanceByID(id, client, existingInstanceStates())
}

//
func createInstance(machine *machinev1.Machine, machineProviderConfig *alibabacloudproviderv1.AlibabaCloudMachineProviderConfig, userData []byte, client aliClient.Client) (*ecs.Instance, error) {
	securityGroupsID, err := checkSecurityGroupsID(machineProviderConfig.VpcId, machineProviderConfig.Placement.Region, machineProviderConfig.SecurityGroupId, client)
	if err != nil {
		return nil, fmt.Errorf("error getting security groups ID: %v", err)
	}

	ImageId, err := checkImageId(machineProviderConfig.Placement.Region, machineProviderConfig.ImageId, client)
	if err != nil {
		return nil, fmt.Errorf("error getting image ID: %v", err)
	}

	createInstanceRequest := ecs.CreateCreateInstanceRequest()
	//securityGroupID
	createInstanceRequest.SecurityGroupId = securityGroupsID
	//imageID
	createInstanceRequest.ImageId = ImageId
	//instanceType
	createInstanceRequest.InstanceType = machineProviderConfig.InstanceType
	//instanceName
	if machineProviderConfig.InstanceName != "" {
		createInstanceRequest.InstanceName = machineProviderConfig.InstanceName
	} else {
		createInstanceRequest.InstanceName = machine.Spec.Name
	}
	//vswitchID
	createInstanceRequest.VSwitchId = machineProviderConfig.VSwitchId
	//systemDisk
	createInstanceRequest.SystemDiskCategory = machineProviderConfig.SystemDiskCategory
	createInstanceRequest.SystemDiskSize = requests.NewInteger64(machineProviderConfig.SystemDiskSize)
	if machineProviderConfig.SystemDiskDiskName != "" {
		createInstanceRequest.SystemDiskDiskName = machineProviderConfig.SystemDiskDiskName
	}
	if machineProviderConfig.SystemDiskDescription != "" {
		createInstanceRequest.SystemDiskDescription = machineProviderConfig.SystemDiskDescription
	}
	//keyPairName
	createInstanceRequest.KeyPairName = machineProviderConfig.KeyPairName
	//publicIP
	if BoolValue(machineProviderConfig.PublicIP) {
		createInstanceRequest.InternetMaxBandwidthOut = requests.NewInteger64(100)
	}
	//ramRoleName
	if machineProviderConfig.RamRoleName != "" {
		createInstanceRequest.RamRoleName = machineProviderConfig.RamRoleName
	}
	//instanceChargeType
	if machineProviderConfig.InstanceChargeType != "" {
		createInstanceRequest.InstanceChargeType = machineProviderConfig.InstanceChargeType
	}

	//No effect when instanceChargeType is not PrePaid
	if machineProviderConfig.Period != requests.NewInteger64(0) {
		createInstanceRequest.Period = machineProviderConfig.Period
	}
	if machineProviderConfig.PeriodUnit != "" {
		createInstanceRequest.PeriodUnit = machineProviderConfig.PeriodUnit
	}
	if machineProviderConfig.AutoRenew == requests.NewBoolean(true) && machineProviderConfig.AutoRenewPeriod != requests.NewInteger64(0) {
		createInstanceRequest.AutoRenew = machineProviderConfig.AutoRenew
		createInstanceRequest.AutoRenewPeriod = machineProviderConfig.AutoRenewPeriod
	}

	//spotStrategy
	if machineProviderConfig.SpotStrategy != "" {
		createInstanceRequest.SpotStrategy = machineProviderConfig.SpotStrategy
	}

	clusterID, ok := getClusterID(machine)
	if !ok {
		klog.Errorf("Unable to get cluster ID for machine: %q", machine.Name)
		return nil, err
	}

	//tags
	createInstanceTags := make([]ecs.CreateInstanceTag, 0)
	if len(machineProviderConfig.Tags) > 0 {
		for _, tag := range machineProviderConfig.Tags {
			createInstanceTags = append(createInstanceTags, ecs.CreateInstanceTag{
				Key:   tag.Key,
				Value: tag.Value,
			})
		}
	}
	createInstanceTags = append(createInstanceTags, []ecs.CreateInstanceTag{
		{Key: fmt.Sprintf("%s%s", clusterFilterKeyPrefix, clusterID), Value: clusterFilterValue},
		{Key: "Name", Value: machine.Name},
	}...)
	tagList := removeDuplicatedTags(createInstanceTags)
	createInstanceRequest.Tag = &tagList

	//dataDisk
	if len(machineProviderConfig.DataDisks) > 0 {
		dataDisks := make([]ecs.CreateInstanceDataDisk, 0)
		for _, dataDisk := range machineProviderConfig.DataDisks {
			dataDisks = append(dataDisks, ecs.CreateInstanceDataDisk{
				Size:     strconv.FormatInt(dataDisk.Size, 10),
				Category: dataDisk.Category,
			})
		}
		createInstanceRequest.DataDisk = &dataDisks
	}
	//userData
	createInstanceRequest.UserData = base64.StdEncoding.EncodeToString(userData)

	createInstanceRequest.Scheme = "https"

	//createInstance
	createInstanceResponse, err := client.CreateInstance(createInstanceRequest)
	if err != nil {
		klog.Errorf("Error creating ECS instance: %v", err)
		return nil, fmt.Errorf("error creating ECS instance: %v", err)
	}

	klog.Infof("The ECS instance %s created", createInstanceResponse.InstanceId)

	//waitForInstance stopped
	klog.Infof("Wait for  ECS instance %s stopped", createInstanceResponse.InstanceId)
	if err := client.WaitForInstance(createInstanceResponse.InstanceId, "Stopped", machineProviderConfig.RegionId, 300); err != nil {
		klog.Errorf("Error waiting ECS instance stopped: %v", err)
		return nil, err
	}
	klog.Infof("The   ECS instance %s stopped", createInstanceResponse.InstanceId)

	klog.Infof("Start  ECS instance %s ", createInstanceResponse.InstanceId)
	//start instance
	startInstanceRequest := ecs.CreateStartInstanceRequest()
	startInstanceRequest.RegionId = machineProviderConfig.RegionId
	startInstanceRequest.InstanceId = createInstanceResponse.InstanceId
	startInstanceRequest.Scheme = "https"

	_, err = client.StartInstance(startInstanceRequest)
	if err != nil {
		klog.Errorf("Error starting ECS instance: %v", err)
		return nil, fmt.Errorf("error starting ECS instance: %v", err)
	}

	//waitForInstanceRunning
	klog.Infof("Wait for  ECS instance %s running", createInstanceResponse.InstanceId)

	if err := client.WaitForInstance(createInstanceResponse.InstanceId, "Running", machineProviderConfig.RegionId, 300); err != nil {
		klog.Errorf("Error waiting ECS instance running: %v", err)
		return nil, err
	}
	klog.Infof("The   ECS instance %s running", createInstanceResponse.InstanceId)

	//describeInstance
	describeInstancesRequest := ecs.CreateDescribeInstancesRequest()
	describeInstancesRequest.RegionId = machineProviderConfig.RegionId
	instancesIds, _ := json.Marshal([]string{createInstanceResponse.InstanceId})
	describeInstancesRequest.InstanceIds = string(instancesIds)
	describeInstancesRequest.Scheme = "https"

	describeInstancesResponse, err := client.DescribeInstances(describeInstancesRequest)
	if err != nil {
		return nil, err
	}

	if len(describeInstancesResponse.Instances.Instance) <= 0 {
		return nil, fmt.Errorf("instance %s not found", createInstanceResponse.InstanceId)
	}

	return &describeInstancesResponse.Instances.Instance[0], nil
}

// extractNodeAddresses maps the instance information from ECS to an array of NodeAddresses
func extractNodeAddresses(instance *ecs.Instance, domainNames []string) ([]corev1.NodeAddress, error) {
	// Not clear if the order matters here, but we might as well indicate a sensible preference order

	if instance == nil {
		return nil, fmt.Errorf("nil instance passed to extractNodeAddresses")
	}

	addresses := []corev1.NodeAddress{}

	// handle internal network interfaces
	for _, networkInterface := range instance.NetworkInterfaces.NetworkInterface {
		// Treating IPv6 addresses as type NodeInternalIP to match what the KNI
		// patch to the AWS cloud-provider code is doing:
		//
		// https://github.com/openshift-kni/origin/commit/7db21c1e26a344e25ae1b825d4f21e7bef5c3650
		// for _, ipv6Address := range networkInterface.Ipv6Addresses {
		// 	if addr := aws.StringValue(ipv6Address.Ipv6Address); addr != "" {
		// 		ip := net.ParseIP(addr)
		// 		if ip == nil {
		// 			return nil, fmt.Errorf("EC2 instance had invalid IPv6 address: %s (%q)", aws.StringValue(instance.InstanceId), addr)
		// 		}
		// 		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalIP, Address: ip.String()})
		// 	}
		// }

		if ipAddress := networkInterface.PrimaryIpAddress; ipAddress != "" {
			ip := net.ParseIP(ipAddress)
			if ip == nil {
				return nil, fmt.Errorf("ECS instance had invalid private address: %s (%q)", instance.InstanceId, ipAddress)
			}
			addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalIP, Address: ip.String()})
		}

	}

	// TODO: Other IP addresses (multiple ips)?
	if len(instance.PublicIpAddress.IpAddress) > 0 {
		publicIPAddress := instance.PublicIpAddress.IpAddress[0]
		if publicIPAddress != "" {
			ip := net.ParseIP(publicIPAddress)
			if ip == nil {
				return nil, fmt.Errorf("EC2 instance had invalid public address: %s (%s)", instance.InstanceId, publicIPAddress)
			}
			addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeExternalIP, Address: ip.String()})
		}
	}
	privateDNSName := instance.HostName
	if privateDNSName != "" {
		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeHostName, Address: privateDNSName})
		// for _, dn := range domainNames {
		// 	customHostName := strings.Join([]string{strings.Split(privateDNSName, ".")[0], dn}, ".")
		// 	if customHostName != privateDNSName {
		// 		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalDNS, Address: customHostName})
		// 	}
		// }
	}
	addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalDNS, Address: strings.Join([]string{instance.RegionId, instance.InstanceId}, ".")})

	return addresses, nil
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

// getStoppedInstances returns all stopped instances that have a tag matching our machine name,
// and cluster ID.
func getStoppedInstances(machine *machinev1.Machine, client aliClient.Client) ([]*ecs.Instance, error) {
	stoppedInstanceStateFilter := []*string{String(InstanceStateNameStopped), String(InstanceStateNameStopping)}
	return getInstances(machine, client, stoppedInstanceStateFilter)
}

// validateMachine check the label that a machine must have to identify the cluster to which it belongs is present.
func validateMachine(machine machinev1.Machine) error {
	if machine.Labels[machinev1.MachineClusterIDLabel] == "" {
		return machinecontroller.InvalidMachineConfiguration("%v: missing %q label", machine.GetName(), machinev1.MachineClusterIDLabel)
	}

	return nil
}

func instanceHasAllowedState(instance *ecs.Instance, instanceStateFilter []*string) error {
	if instance.InstanceId == "" {
		return fmt.Errorf("instance has nil ID")
	}

	if instance.Status == "" {
		return fmt.Errorf("instance %s has nil state", instance.InstanceId)
	}

	if len(instanceStateFilter) == 0 {
		return nil
	}

	actualState := instance.Status
	for _, allowedState := range instanceStateFilter {
		if StringValue(allowedState) == actualState {
			return nil
		}
	}

	allowedStates := make([]string, 0, len(instanceStateFilter))
	for _, allowedState := range instanceStateFilter {
		allowedStates = append(allowedStates, StringValue(allowedState))
	}
	return fmt.Errorf("instance %s state %q is not in %s", instance.InstanceId, actualState, strings.Join(allowedStates, ", "))
}

// getInstanceByID returns the instance with the given ID if it exists.
func getInstanceByID(id string, client aliClient.Client, instanceStateFilter []*string) (*ecs.Instance, error) {
	if id == "" {
		return nil, fmt.Errorf("instance-id not specified")
	}

	describeInstancesRequest := ecs.CreateDescribeInstancesRequest()
	describeInstancesRequest.InstanceIds = fmt.Sprintf("[\"%s\"]", id)
	describeInstancesRequest.Scheme = "https"

	result, err := client.DescribeInstances(describeInstancesRequest)
	if err != nil {
		return nil, err
	}

	if len(result.Instances.Instance) != 1 {
		return nil, fmt.Errorf("found %d instances for instance-id %s", len(result.Instances.Instance), id)
	}

	instance := result.Instances.Instance[0]

	return &instance, instanceHasAllowedState(&instance, instanceStateFilter)
}

// getClusterID get cluster ID by machine.openshift.io/cluster-api-cluster label
func getClusterID(machine *machinev1.Machine) (string, bool) {
	clusterID, ok := machine.Labels[machinev1.MachineClusterIDLabel]
	// TODO: remove 347-350
	// NOTE: This block can be removed after the label renaming transition to machine.openshift.io
	if !ok {
		clusterID, ok = machine.Labels[upstreamMachineClusterIDLabel]
	}
	return clusterID, ok
}
