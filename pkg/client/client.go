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

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/providers"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	corev1 "k8s.io/api/core/v1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	AliCloudAccessKeyId       = "alicloud_access_key_id"
	AliCloudAccessKeySecret   = "alicloud_access_key_secret"
	AliCloudAccessKeyStsToken = "alicloud_access_key_sts_token"

	InstanceDefaultTimeout = 900
	DefaultWaitForInterval = 5
)

type AliCloudClientBuilderFuncType func(client client.Client, secretName, namespace, region string) (Client, error)

// Client is a wrapper object for actual AliCloud SDK clients to allow for easier testing.
type Client interface {
	//ecs
	CreateInstance(*ecs.CreateInstanceRequest) (*ecs.CreateInstanceResponse, error)
	RunInstances(*ecs.RunInstancesRequest) (*ecs.RunInstancesResponse, error)
	DescribeInstances(*ecs.DescribeInstancesRequest) (*ecs.DescribeInstancesResponse, error)
	StartInstance(*ecs.StartInstanceRequest) (*ecs.StartInstanceResponse, error)
	StopInstance(*ecs.StopInstanceRequest) (*ecs.StopInstanceResponse, error)
	DeleteInstance(*ecs.DeleteInstanceRequest) (*ecs.DeleteInstanceResponse, error)
	DescribeImages(*ecs.DescribeImagesRequest) (*ecs.DescribeImagesResponse, error)
	DescribeSecurityGroups(*ecs.DescribeSecurityGroupsRequest) (*ecs.DescribeSecurityGroupsResponse, error)
	DescribeRegions(*ecs.DescribeRegionsRequest) (*ecs.DescribeRegionsResponse, error)
	DescribeZones(*ecs.DescribeZonesRequest) (*ecs.DescribeZonesResponse, error)
	DescribeDisks(*ecs.DescribeDisksRequest) (*ecs.DescribeDisksResponse, error)
	//waitForInstance
	WaitForInstance(instanceId, instanceStatus, regionId string, timeout int) error

	//vpc
	DescribeVpcs(*vpc.DescribeVpcsRequest) (*vpc.DescribeVpcsResponse, error)
	DescribeVSwitches(*vpc.DescribeVSwitchesRequest) (*vpc.DescribeVSwitchesResponse, error)
}

//create instance
func (c *aliCloudClient) CreateInstance(request *ecs.CreateInstanceRequest) (*ecs.CreateInstanceResponse, error) {
	return c.ecs2Client.CreateInstance(request)
}

//run instances
func (c *aliCloudClient) RunInstances(request *ecs.RunInstancesRequest) (*ecs.RunInstancesResponse, error) {
	return c.ecs2Client.RunInstances(request)
}

//describe instances
func (c *aliCloudClient) DescribeInstances(request *ecs.DescribeInstancesRequest) (*ecs.DescribeInstancesResponse, error) {
	return c.ecs2Client.DescribeInstances(request)
}

//start instances
func (c *aliCloudClient) StartInstance(request *ecs.StartInstanceRequest) (*ecs.StartInstanceResponse, error) {
	return c.ecs2Client.StartInstance(request)
}

//stop instances
func (c *aliCloudClient) StopInstance(request *ecs.StopInstanceRequest) (*ecs.StopInstanceResponse, error) {
	return c.ecs2Client.StopInstance(request)
}

//deleteInstance
func (c *aliCloudClient) DeleteInstance(request *ecs.DeleteInstanceRequest) (*ecs.DeleteInstanceResponse, error) {
	return c.ecs2Client.DeleteInstance(request)
}

//describe images
func (c *aliCloudClient) DescribeImages(request *ecs.DescribeImagesRequest) (*ecs.DescribeImagesResponse, error) {
	return c.ecs2Client.DescribeImages(request)
}

//describe securityGroups
func (c *aliCloudClient) DescribeSecurityGroups(request *ecs.DescribeSecurityGroupsRequest) (*ecs.DescribeSecurityGroupsResponse, error) {
	return c.ecs2Client.DescribeSecurityGroups(request)
}

//describe regions
func (c *aliCloudClient) DescribeRegions(request *ecs.DescribeRegionsRequest) (*ecs.DescribeRegionsResponse, error) {
	return c.ecs2Client.DescribeRegions(request)
}

//describe zones
func (c *aliCloudClient) DescribeZones(request *ecs.DescribeZonesRequest) (*ecs.DescribeZonesResponse, error) {
	return c.ecs2Client.DescribeZones(request)
}

//describe disks
func (c *aliCloudClient) DescribeDisks(request *ecs.DescribeDisksRequest) (*ecs.DescribeDisksResponse, error) {
	return c.ecs2Client.DescribeDisks(request)
}

//wait4Instances
func (c *aliCloudClient) WaitForInstance(instanceId, instanceStatus, regionId string, timeout int) error {
	// WaitForInstance waits for instance to given status
	// when instance.NotFound wait until timeout
	if timeout <= 0 {
		timeout = InstanceDefaultTimeout
	}

	for {
		describeInstancesRequest := ecs.CreateDescribeInstancesRequest()
		describeInstancesRequest.RegionId = regionId
		instancesIds, _ := json.Marshal([]string{instanceId})
		describeInstancesRequest.InstanceIds = string(instancesIds)
		describeInstancesRequest.Scheme = "https"
		describeInstancesResponse, err := c.ecs2Client.DescribeInstances(describeInstancesRequest)
		if err == nil {
			if len(describeInstancesResponse.Instances.Instance) > 0 && describeInstancesResponse.Instances.Instance[0].Status == instanceStatus {
				break
			}
		}
		timeout = timeout - DefaultWaitForInterval
		if timeout <= 0 {
			return fmt.Errorf("Timeout")
		}
		time.Sleep(DefaultWaitForInterval * time.Second)
	}
	return nil
}

//describe vpcs
func (c *aliCloudClient) DescribeVpcs(request *vpc.DescribeVpcsRequest) (*vpc.DescribeVpcsResponse, error) {
	return c.vpc2Client.DescribeVpcs(request)
}

//describe vswitches
func (c *aliCloudClient) DescribeVSwitches(request *vpc.DescribeVSwitchesRequest) (*vpc.DescribeVSwitchesResponse, error) {
	return c.vpc2Client.DescribeVSwitches(request)
}

type aliCloudClient struct {
	ecs2Client *ecs.Client
	vpc2Client *vpc.Client
}

func NewClient(ctrlRuntimeClient client.Client, secretName, namespace, region string) (Client, error) {
	aliCloudConfig := &providers.Configuration{
		AccessKeyID:     os.Getenv("ALICLOUD_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("ALICLOUD_ACCESS_KEY_SECRET"),
	}

	if secretName != "" {
		var secret corev1.Secret
		if err := ctrlRuntimeClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: secretName}, &secret); err != nil {
			return nil, err
		}
		accessKeyID, ok := secret.Data[AliCloudAccessKeyId]
		if !ok {
			return nil, fmt.Errorf("AliCloud credentials secret %v did not contain key %v",
				secretName, AliCloudAccessKeyId)
		}
		aliCloudConfig.AccessKeyID = string(accessKeyID)

		accessKeySecret, ok := secret.Data[AliCloudAccessKeySecret]
		if !ok {
			return nil, fmt.Errorf("AliCloud credentials secret %v did not contain key %v",
				secretName, AliCloudAccessKeySecret)
		}
		aliCloudConfig.AccessKeySecret = string(accessKeySecret)
		accessKeyStsToken, _ := secret.Data[AliCloudAccessKeyStsToken]
		aliCloudConfig.AccessKeyStsToken = string(accessKeyStsToken)
	}

	return initAliCloudClient(aliCloudConfig, region)
}

//init client from ak directly
func NewClientFromKeys(accessKeyId, accessKeySecret, region string) (Client, error) {
	return initAliCloudClient(&providers.Configuration{AccessKeyID: accessKeyId, AccessKeySecret: accessKeySecret}, region)
}

func initAliCloudClient(configuration *providers.Configuration, region string) (Client, error) {
	p := &providers.ConfigurationProvider{Configuration: configuration}
	credential, err := p.Retrieve()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrive credential %++v ", err)
	}

	ecsClient, err := ecs.NewClientWithOptions(region, &sdk.Config{}, credential)
	if err != nil {
		return nil, fmt.Errorf("Failed to init ecsClient %++v ", err)
	}

	vpcClient, err := vpc.NewClientWithOptions(region, &sdk.Config{}, credential)
	if err != nil {
		return nil, fmt.Errorf("Failed to init vpcClient %++v ", err)
	}

	return &aliCloudClient{
		ecs2Client: ecsClient,
		vpc2Client: vpcClient,
	}, nil
}
