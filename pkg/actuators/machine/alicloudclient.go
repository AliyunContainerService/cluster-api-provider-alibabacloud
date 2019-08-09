package machine

import (
	"fmt"
	aliCloud "github.com/AliyunContainerService/cluster-api-provider-alicloud/pkg/client"
	"github.com/openshift/cluster-api-actuator-pkg/pkg/types"
	machinev1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
)

// AliCloudClientWrapper implements CloudProviderClient for alicloud e2e framework
type AliCloudClientWrapper struct {
	client   aliCloud.Client
	regionId string
}

func (client *AliCloudClientWrapper) GetRunningInstances(machine *machinev1beta1.Machine) ([]interface{}, error) {
	runningInstances, err := getRunningInstances(machine, client.client, client.regionId)
	if err != nil {
		return nil, err
	}

	var instances []interface{}
	for _, instance := range runningInstances {
		instances = append(instances, instance)
	}

	return instances, nil
}

func (client *AliCloudClientWrapper) GetPublicDNSName(machine *machinev1beta1.Machine) (string, error) {
	instance, err := getRunningInstance(machine, client.client, client.regionId)
	if err != nil {
		return "", err
	}

	//todo replace dns using publicIP
	if len(instance.PublicIpAddress.IpAddress) <= 0 {
		return "", fmt.Errorf("machine instance public IP name not set")
	}

	return instance.PublicIpAddress.IpAddress[0], nil
}

func (client *AliCloudClientWrapper) GetPrivateIP(machine *machinev1beta1.Machine) (string, error) {
	instance, err := getRunningInstance(machine, client.client, client.regionId)
	if err != nil {
		return "", err
	}

	if len(instance.InnerIpAddress.IpAddress) <= 0 {
		return "", fmt.Errorf("machine instance private ip not set")
	}

	return instance.InnerIpAddress.IpAddress[0], nil
}

var _ types.CloudProviderClient = &AliCloudClientWrapper{}

func NewAliCloudClientWrapper(client aliCloud.Client, regionId string) *AliCloudClientWrapper {
	return &AliCloudClientWrapper{client: client, regionId: regionId}
}

// GetSecurityGroups gets security groups
func (client *AliCloudClientWrapper) GetSecurityGroups(machine *machinev1beta1.Machine) ([]string, error) {
	instance, err := getRunningInstance(machine, client.client, client.regionId)
	if err != nil {
		return nil, err
	}
	var groups []string
	for _, sgId := range instance.SecurityGroupIds.SecurityGroupId {
		groups = append(groups, sgId)
	}
	return groups, nil
}

// GetRamRole gets RAM role
func (client *AliCloudClientWrapper) GetRamRole(machine *machinev1beta1.Machine) (string, error) {
	_, err := getRunningInstance(machine, client.client, client.regionId)
	if err != nil {
		return "", err
	}

	//TODO The response not include roleName
	return "", nil
}

// GetTags gets tags
func (client *AliCloudClientWrapper) GetTags(machine *machinev1beta1.Machine) (map[string]string, error) {
	instance, err := getRunningInstance(machine, client.client, client.regionId)
	if err != nil {
		return nil, err
	}
	tags := make(map[string]string, len(instance.Tags.Tag))
	for _, tag := range instance.Tags.Tag {
		tags[tag.TagKey] = tag.TagValue
	}
	return tags, nil
}

// GetAvailabilityZoneId gets availability zone
func (client *AliCloudClientWrapper) GetAvailabilityZoneId(machine *machinev1beta1.Machine) (string, error) {
	instance, err := getRunningInstance(machine, client.client, client.regionId)
	if err != nil {
		return "", err
	}
	if instance.ZoneId == "" {
		return "", err
	}
	return instance.ZoneId, nil
}

// GetDataDisks gets volumes attached to instance
func (client *AliCloudClientWrapper) GetDataDisks(machine *machinev1beta1.Machine) (map[string]map[string]interface{}, error) {
	//todo API NOT IMPL
	return nil, nil
}
