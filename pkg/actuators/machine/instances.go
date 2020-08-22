package machine

import (
	"fmt"
	"sort"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/golang/glog"
	"k8s.io/klog"

	aliClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
)

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
func checkImageId(regionId string, ImageId string, client aliClient.Client) (string, error) {
	klog.Infof("check imageId in region %s", regionId)
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

// removeStoppedMachine removes all instances of a specific machine that are in a stopped state.
func removeStoppedMachine(machine *machinev1.Machine, client aliClient.Client) error {
	instances, err := getStoppedInstances(machine, client)
	if err != nil {
		klog.Errorf("Error getting stopped instances: %v", err)
		return fmt.Errorf("error getting stopped instances: %v", err)
	}

	if len(instances) == 0 {
		klog.Infof("No stopped instances found for machine %v", machine.Name)
		return nil
	}

	err = deleteInstances(client, instances)
	return err
}
