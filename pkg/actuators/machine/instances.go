package machine

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/golang/glog"
	"sort"
	"strconv"
	"time"

	providerconfigv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1alpha1"
	aliClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
)

//
func createInstance(machine *machinev1.Machine, machineProviderConfig *providerconfigv1.AlibabaCloudMachineProviderConfig, userData []byte, client aliClient.Client) (*ecs.Instance, error) {
	securityGroupsID, err := checkSecurityGroupsID(machineProviderConfig.VpcId, machineProviderConfig.RegionId, machineProviderConfig.SecurityGroupId, client)
	if err != nil {
		return nil, fmt.Errorf("error getting security groups ID: %v", err)
	}

	ImageId, err := checkImageId(machineProviderConfig.RegionId, machineProviderConfig.ImageId, client)
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
	if machineProviderConfig.PublicIP {
		createInstanceRequest.InternetMaxBandwidthOut = requests.NewInteger64(100)
	}
	//ramRoleName
	if machineProviderConfig.RamRoleName != "" {
		createInstanceRequest.RamRoleName = machineProviderConfig.RamRoleName
	}

	clusterID, ok := getClusterID(machine)
	if !ok {
		glog.Errorf("Unable to get cluster ID for machine: %q", machine.Name)
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
		glog.Errorf("Error creating ECS instance: %v", err)
		return nil, fmt.Errorf("error creating ECS instance: %v", err)
	}

	glog.Infof("The ECS instance %s created", createInstanceResponse.InstanceId)

	//waitForInstance stopped
	glog.Infof("Wait for  ECS instance %s stopped", createInstanceResponse.InstanceId)
	if err := client.WaitForInstance(createInstanceResponse.InstanceId, "Stopped", machineProviderConfig.RegionId, 300); err != nil {
		glog.Errorf("Error waiting ECS instance stopped: %v", err)
		return nil, err
	}
	glog.Infof("The   ECS instance %s stopped", createInstanceResponse.InstanceId)

	glog.Infof("Start  ECS instance %s ", createInstanceResponse.InstanceId)
	//start instance
	startInstanceRequest := ecs.CreateStartInstanceRequest()
	startInstanceRequest.RegionId = machineProviderConfig.RegionId
	startInstanceRequest.InstanceId = createInstanceResponse.InstanceId
	startInstanceRequest.Scheme = "https"

	_, err = client.StartInstance(startInstanceRequest)
	if err != nil {
		glog.Errorf("Error starting ECS instance: %v", err)
		return nil, fmt.Errorf("error starting ECS instance: %v", err)
	}

	//waitForInstanceRunning
	glog.Infof("Wait for  ECS instance %s running", createInstanceResponse.InstanceId)

	if err := client.WaitForInstance(createInstanceResponse.InstanceId, "Running", machineProviderConfig.RegionId, 300); err != nil {
		glog.Errorf("Error waiting ECS instance running: %v", err)
		return nil, err
	}
	glog.Infof("The   ECS instance %s running", createInstanceResponse.InstanceId)

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

// Scan machine tags, and return a deduped tags list
func removeDuplicatedTags(tags []ecs.CreateInstanceTag) []ecs.CreateInstanceTag {
	m := make(map[string]bool)
	result := make([]ecs.CreateInstanceTag, 0)

	// look for duplicates
	for _, entry := range tags {
		if _, value := m[entry.Key]; !value {
			m[entry.Key] = true
			result = append(result, entry)
		}
	}
	return result
}

//check securityGroupId
func checkSecurityGroupsID(vpcId, regionId, securityGroupId string, client aliClient.Client) (string, error) {
	glog.Infof("check security group ID based in vpc %s", vpcId)
	describeSecurityGroupsRequest := ecs.CreateDescribeSecurityGroupsRequest()
	describeSecurityGroupsRequest.RegionId = regionId
	describeSecurityGroupsRequest.SecurityGroupId = securityGroupId
	describeSecurityGroupsRequest.VpcId = vpcId
	describeSecurityGroupsRequest.Scheme = "https"

	describeSecurityGroupsResponse, err := client.DescribeSecurityGroups(describeSecurityGroupsRequest)
	if err != nil {
		return "", fmt.Errorf("error describing security groups: %v", err)
	}

	if len(describeSecurityGroupsResponse.SecurityGroups.SecurityGroup) <= 0 {
		return "", fmt.Errorf("No security group found")
	}

	return securityGroupId, nil
}

//check ImageId
func checkImageId(regionId, ImageId string, client aliClient.Client) (string, error) {
	glog.Infof("check imageId in region %s", regionId)
	describeImagesRequest := ecs.CreateDescribeImagesRequest()
	describeImagesRequest.RegionId = regionId
	describeImagesRequest.ImageId = ImageId
	describeImagesRequest.Scheme = "https"

	describeImagesResponse, err := client.DescribeImages(describeImagesRequest)
	if err != nil {
		return "", fmt.Errorf("error describing images: %v", err)
	}

	if len(describeImagesResponse.Images.Image) <= 0 {
		return "", fmt.Errorf("No image found ")
	}

	return ImageId, nil
}

type instanceList []*ecs.Instance

func (il instanceList) Len() int {
	return len(il)
}

func (il instanceList) Swap(i, j int) {
	il[i], il[j] = il[j], il[i]
}

var (
	timeTemplate1 = "2006-01-02 15:04:05"
)

func (il instanceList) Less(i, j int) bool {
	if il[i].CreationTime == "" && il[j].CreationTime == "" {
		return false
	}
	if il[i].CreationTime != "" && il[j].CreationTime == "" {
		return false
	}
	if il[i].CreationTime == "" && il[j].CreationTime != "" {
		return true
	}

	t1, _ := time.ParseInLocation(timeTemplate1, il[i].CreationTime, time.Local)
	t2, _ := time.ParseInLocation(timeTemplate1, il[j].CreationTime, time.Local)

	return t1.After(t2)
}

// sortInstances will sort a list of instance based on an instace launch time
// from the newest to the oldest.
// This function should only be called with running instances, not those which are stopped or
// terminated.
func sortInstances(instances []*ecs.Instance) {
	sort.Sort(instanceList(instances))
}
