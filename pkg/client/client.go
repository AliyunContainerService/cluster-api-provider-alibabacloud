package client

import (
	"context"
	"fmt"

	"github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/utils"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"

	"k8s.io/klog/v2"

	"github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/version"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/providers"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	machineapiapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"
	corev1 "k8s.io/api/core/v1"
<<<<<<< HEAD
<<<<<<< HEAD
	apimachineryerrors "k8s.io/apimachinery/pkg/api/errors"
=======
	"os"
>>>>>>> ebdd9bd0 (update test case)
=======
	apimachineryerrors "k8s.io/apimachinery/pkg/api/errors"
>>>>>>> e879a141 (alibabacloud machine-api provider)
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate go run ../../vendor/github.com/golang/mock/mockgen -source=./client.go -destination=./mock/client_mock.go -package=mock
<<<<<<< HEAD

// AlibabaCloudClientBuilderFunc is function type for building alibabacloud client
type AlibabaCloudClientBuilderFunc func(client client.Client, secretName, namespace, region string, configManagedClient client.Client) (Client, error)

// machineProviderUserAgent is a named handler that will add cluster-api-provider-alibabacloud
var machineProviderUserAgent = fmt.Sprintf("openshift.io cluster-api-provider-alibabacloud:%s", version.Version.String())

const (
	kubeAccessKeyID           = "accessKeyID"
	kubeAccessKeySecret       = "accessKeySecret"
	kubeAccessKeyStsToken     = "accessKeyStsToken"
	kubeRoleArn               = "roleArn"
	kubeRoleSessionName       = "roleSessionName"
	kubeRoleSessionExpiration = "roleSessionExpiration"
	kubeRoleName              = "roleName"

	// KubeCloudConfigNamespace is the namespace where the kube cloud config ConfigMap is located
	KubeCloudConfigNamespace = "openshift-config-managed"

	kubeCloudConfigName = "kube-cloud-config"
)

=======

// AlibabaCloudClientBuilderFunc is function type for building alibabacloud client
type AlibabaCloudClientBuilderFunc func(client client.Client, secretName, namespace, region string, configManagedClient client.Client) (Client, error)

// machineProviderUserAgent is a named handler that will add cluster-api-provider-alibabacloud
var machineProviderUserAgent = fmt.Sprintf("openshift.io cluster-api-provider-alibabacloud:%s", version.Version.String())

const (
	kubeAccessKeyID           = "accessKeyID"
	kubeAccessKeySecret       = "accessKeySecret"
	kubeAccessKeyStsToken     = "accessKeyStsToken"
	kubeRoleArn               = "roleArn"
	kubeRoleSessionName       = "roleSessionName"
	kubeRoleSessionExpiration = "roleSessionExpiration"
	kubeRoleName              = "roleName"

	// KubeCloudConfigNamespace is the namespace where the kube cloud config ConfigMap is located
	KubeCloudConfigNamespace = "openshift-config-managed"

	kubeCloudConfigName = "kube-cloud-config"
)

>>>>>>> e879a141 (alibabacloud machine-api provider)
// Client is a wrapper object for actual alibabacloud SDK clients to allow for easier testing.
type Client interface {
	//Ecs
	RunInstances(*ecs.RunInstancesRequest) (*ecs.RunInstancesResponse, error)
	CreateInstance(*ecs.CreateInstanceRequest) (*ecs.CreateInstanceResponse, error)
	DescribeInstances(*ecs.DescribeInstancesRequest) (*ecs.DescribeInstancesResponse, error)
	DeleteInstances(*ecs.DeleteInstancesRequest) (*ecs.DeleteInstancesResponse, error)
	StartInstance(*ecs.StartInstanceRequest) (*ecs.StartInstanceResponse, error)
	RebootInstance(request *ecs.RebootInstanceRequest) (*ecs.RebootInstanceResponse, error)
	StopInstance(*ecs.StopInstanceRequest) (*ecs.StopInstanceResponse, error)
	StartInstances(*ecs.StartInstancesRequest) (*ecs.StartInstancesResponse, error)
	RebootInstances(request *ecs.RebootInstancesRequest) (*ecs.RebootInstancesResponse, error)
	StopInstances(*ecs.StopInstancesRequest) (*ecs.StopInstancesResponse, error)
	DeleteInstance(*ecs.DeleteInstanceRequest) (*ecs.DeleteInstanceResponse, error)
	AttachInstanceRAMRole(*ecs.AttachInstanceRamRoleRequest) (*ecs.AttachInstanceRamRoleResponse, error)
	DetachInstanceRAMRole(*ecs.DetachInstanceRamRoleRequest) (*ecs.DetachInstanceRamRoleResponse, error)
	DescribeInstanceStatus(*ecs.DescribeInstanceStatusRequest) (*ecs.DescribeInstanceStatusResponse, error)
	ReActivateInstances(*ecs.ReActivateInstancesRequest) (*ecs.ReActivateInstancesResponse, error)
	DescribeUserData(*ecs.DescribeUserDataRequest) (*ecs.DescribeUserDataResponse, error)
	DescribeInstanceTypes(*ecs.DescribeInstanceTypesRequest) (*ecs.DescribeInstanceTypesResponse, error)
	ModifyInstanceAttribute(*ecs.ModifyInstanceAttributeRequest) (*ecs.ModifyInstanceAttributeResponse, error)
	ModifyInstanceMetadataOptions(*ecs.ModifyInstanceMetadataOptionsRequest) (*ecs.ModifyInstanceMetadataOptionsResponse, error)

	TagResources(*ecs.TagResourcesRequest) (*ecs.TagResourcesResponse, error)
	ListTagResources(*ecs.ListTagResourcesRequest) (*ecs.ListTagResourcesResponse, error)
	UntagResources(*ecs.UntagResourcesRequest) (*ecs.UntagResourcesResponse, error)

	//Network
	AllocatePublicIPAddress(*ecs.AllocatePublicIpAddressRequest) (*ecs.AllocatePublicIpAddressResponse, error)

	//Disk
	CreateDisk(*ecs.CreateDiskRequest) (*ecs.CreateDiskResponse, error)
	AttachDisk(*ecs.AttachDiskRequest) (*ecs.AttachDiskResponse, error)
	DescribeDisks(*ecs.DescribeDisksRequest) (*ecs.DescribeDisksResponse, error)
	ModifyDiskChargeType(*ecs.ModifyDiskChargeTypeRequest) (*ecs.ModifyDiskChargeTypeResponse, error)
	ModifyDiskAttribute(*ecs.ModifyDiskAttributeRequest) (*ecs.ModifyDiskAttributeResponse, error)
	ModifyDiskSpec(*ecs.ModifyDiskSpecRequest) (*ecs.ModifyDiskSpecResponse, error)
	ReplaceSystemDisk(*ecs.ReplaceSystemDiskRequest) (*ecs.ReplaceSystemDiskResponse, error)
	ReInitDisk(*ecs.ReInitDiskRequest) (*ecs.ReInitDiskResponse, error)
	ResetDisk(*ecs.ResetDiskRequest) (*ecs.ResetDiskResponse, error)
	ResizeDisk(*ecs.ResizeDiskRequest) (*ecs.ResizeDiskResponse, error)
	DetachDisk(*ecs.DetachDiskRequest) (*ecs.DetachDiskResponse, error)
	DeleteDisk(*ecs.DeleteDiskRequest) (*ecs.DeleteDiskResponse, error)

	//Region & Zone
	DescribeRegions(*ecs.DescribeRegionsRequest) (*ecs.DescribeRegionsResponse, error)
	DescribeZones(*ecs.DescribeZonesRequest) (*ecs.DescribeZonesResponse, error)

	//Images
	DescribeImages(*ecs.DescribeImagesRequest) (*ecs.DescribeImagesResponse, error)

	//SecurityGroup
	CreateSecurityGroup(*ecs.CreateSecurityGroupRequest) (*ecs.CreateSecurityGroupResponse, error)
	AuthorizeSecurityGroup(*ecs.AuthorizeSecurityGroupRequest) (*ecs.AuthorizeSecurityGroupResponse, error)
	AuthorizeSecurityGroupEgress(*ecs.AuthorizeSecurityGroupEgressRequest) (*ecs.AuthorizeSecurityGroupEgressResponse, error)
	RevokeSecurityGroup(*ecs.RevokeSecurityGroupRequest) (*ecs.RevokeSecurityGroupResponse, error)
	RevokeSecurityGroupEgress(*ecs.RevokeSecurityGroupEgressRequest) (*ecs.RevokeSecurityGroupEgressResponse, error)
	JoinSecurityGroup(*ecs.JoinSecurityGroupRequest) (*ecs.JoinSecurityGroupResponse, error)
	LeaveSecurityGroup(*ecs.LeaveSecurityGroupRequest) (*ecs.LeaveSecurityGroupResponse, error)
	DescribeSecurityGroupAttribute(*ecs.DescribeSecurityGroupAttributeRequest) (*ecs.DescribeSecurityGroupAttributeResponse, error)
	DescribeSecurityGroups(*ecs.DescribeSecurityGroupsRequest) (*ecs.DescribeSecurityGroupsResponse, error)
	DescribeSecurityGroupReferences(*ecs.DescribeSecurityGroupReferencesRequest) (*ecs.DescribeSecurityGroupReferencesResponse, error)
	ModifySecurityGroupAttribute(*ecs.ModifySecurityGroupAttributeRequest) (*ecs.ModifySecurityGroupAttributeResponse, error)
	ModifySecurityGroupEgressRule(*ecs.ModifySecurityGroupEgressRuleRequest) (*ecs.ModifySecurityGroupEgressRuleResponse, error)
	ModifySecurityGroupPolicy(*ecs.ModifySecurityGroupPolicyRequest) (*ecs.ModifySecurityGroupPolicyResponse, error)
	ModifySecurityGroupRule(*ecs.ModifySecurityGroupRuleRequest) (*ecs.ModifySecurityGroupRuleResponse, error)
	DeleteSecurityGroup(*ecs.DeleteSecurityGroupRequest) (*ecs.DeleteSecurityGroupResponse, error)

	//VPC
	CreateVpc(*vpc.CreateVpcRequest) (*vpc.CreateVpcResponse, error)
	DeleteVpc(*vpc.DeleteVpcRequest) (*vpc.DeleteVpcResponse, error)
	DescribeVpcs(*vpc.DescribeVpcsRequest) (*vpc.DescribeVpcsResponse, error)
	CreateVSwitch(*vpc.CreateVSwitchRequest) (*vpc.CreateVSwitchResponse, error)
	DeleteVSwitch(*vpc.DeleteVSwitchRequest) (*vpc.DeleteVSwitchResponse, error)
	DescribeVSwitches(*vpc.DescribeVSwitchesRequest) (*vpc.DescribeVSwitchesResponse, error)

	//Natgateway
	CreateNatGateway(*vpc.CreateNatGatewayRequest) (*vpc.CreateNatGatewayResponse, error)
	DescribeNatGateways(*vpc.DescribeNatGatewaysRequest) (*vpc.DescribeNatGatewaysResponse, error)
	DeleteNatGateway(*vpc.DeleteNatGatewayRequest) (*vpc.DeleteNatGatewayResponse, error)

	//EIP
	AllocateEipAddress(*vpc.AllocateEipAddressRequest) (*vpc.AllocateEipAddressResponse, error)
	AssociateEipAddress(*vpc.AssociateEipAddressRequest) (*vpc.AssociateEipAddressResponse, error)
	ModifyEipAddressAttribute(*vpc.ModifyEipAddressAttributeRequest) (*vpc.ModifyEipAddressAttributeResponse, error)
	DescribeEipAddresses(*vpc.DescribeEipAddressesRequest) (*vpc.DescribeEipAddressesResponse, error)
	UnassociateEipAddress(*vpc.UnassociateEipAddressRequest) (*vpc.UnassociateEipAddressResponse, error)
	ReleaseEipAddress(*vpc.ReleaseEipAddressRequest) (*vpc.ReleaseEipAddressResponse, error)

	//SLB
	CreateLoadBalancer(*slb.CreateLoadBalancerRequest) (*slb.CreateLoadBalancerResponse, error)
	DeleteLoadBalancer(*slb.DeleteLoadBalancerRequest) (*slb.DeleteLoadBalancerResponse, error)
	DescribeLoadBalancers(*slb.DescribeLoadBalancersRequest) (*slb.DescribeLoadBalancersResponse, error)
	CreateLoadBalancerTCPListener(*slb.CreateLoadBalancerTCPListenerRequest) (*slb.CreateLoadBalancerTCPListenerResponse, error)
	SetLoadBalancerTCPListenerAttribute(*slb.SetLoadBalancerTCPListenerAttributeRequest) (*slb.SetLoadBalancerTCPListenerAttributeResponse, error)
	DescribeLoadBalancerTCPListenerAttribute(*slb.DescribeLoadBalancerTCPListenerAttributeRequest) (*slb.DescribeLoadBalancerTCPListenerAttributeResponse, error)
	CreateLoadBalancerUDPListener(*slb.CreateLoadBalancerUDPListenerRequest) (*slb.CreateLoadBalancerUDPListenerResponse, error)
	SetLoadBalancerUDPListenerAttribute(*slb.SetLoadBalancerUDPListenerAttributeRequest) (*slb.SetLoadBalancerUDPListenerAttributeResponse, error)
	DescribeLoadBalancerUDPListenerAttribute(*slb.DescribeLoadBalancerUDPListenerAttributeRequest) (*slb.DescribeLoadBalancerUDPListenerAttributeResponse, error)
	CreateLoadBalancerHTTPListener(*slb.CreateLoadBalancerHTTPListenerRequest) (*slb.CreateLoadBalancerHTTPListenerResponse, error)
	SetLoadBalancerHTTPListenerAttribute(*slb.SetLoadBalancerHTTPListenerAttributeRequest) (*slb.SetLoadBalancerHTTPListenerAttributeResponse, error)
	DescribeLoadBalancerHTTPListenerAttribute(*slb.DescribeLoadBalancerHTTPListenerAttributeRequest) (*slb.DescribeLoadBalancerHTTPListenerAttributeResponse, error)
	CreateLoadBalancerHTTPSListener(*slb.CreateLoadBalancerHTTPSListenerRequest) (*slb.CreateLoadBalancerHTTPSListenerResponse, error)
	SetLoadBalancerHTTPSListenerAttribute(*slb.SetLoadBalancerHTTPSListenerAttributeRequest) (*slb.SetLoadBalancerHTTPSListenerAttributeResponse, error)
	DescribeLoadBalancerHTTPSListenerAttribute(*slb.DescribeLoadBalancerHTTPSListenerAttributeRequest) (*slb.DescribeLoadBalancerHTTPSListenerAttributeResponse, error)
	StartLoadBalancerListener(*slb.StartLoadBalancerListenerRequest) (*slb.StartLoadBalancerListenerResponse, error)
	StopLoadBalancerListener(*slb.StopLoadBalancerListenerRequest) (*slb.StopLoadBalancerListenerResponse, error)
	DeleteLoadBalancerListener(*slb.DeleteLoadBalancerListenerRequest) (*slb.DeleteLoadBalancerListenerResponse, error)
	DescribeLoadBalancerListeners(*slb.DescribeLoadBalancerListenersRequest) (*slb.DescribeLoadBalancerListenersResponse, error)
	AddBackendServers(*slb.AddBackendServersRequest) (*slb.AddBackendServersResponse, error)
	RemoveBackendServers(*slb.RemoveBackendServersRequest) (*slb.RemoveBackendServersResponse, error)
	SetBackendServers(*slb.SetBackendServersRequest) (*slb.SetBackendServersResponse, error)
	DescribeHealthStatus(*slb.DescribeHealthStatusRequest) (*slb.DescribeHealthStatusResponse, error)
	CreateVServerGroup(*slb.CreateVServerGroupRequest) (*slb.CreateVServerGroupResponse, error)
	SetVServerGroupAttribute(*slb.SetVServerGroupAttributeRequest) (*slb.SetVServerGroupAttributeResponse, error)
	AddVServerGroupBackendServers(*slb.AddVServerGroupBackendServersRequest) (*slb.AddVServerGroupBackendServersResponse, error)
	RemoveVServerGroupBackendServers(*slb.RemoveVServerGroupBackendServersRequest) (*slb.RemoveVServerGroupBackendServersResponse, error)
	ModifyVServerGroupBackendServers(*slb.ModifyVServerGroupBackendServersRequest) (*slb.ModifyVServerGroupBackendServersResponse, error)
	DeleteVServerGroup(*slb.DeleteVServerGroupRequest) (*slb.DeleteVServerGroupResponse, error)
	DescribeVServerGroups(*slb.DescribeVServerGroupsRequest) (*slb.DescribeVServerGroupsResponse, error)
	DescribeVServerGroupAttribute(*slb.DescribeVServerGroupAttributeRequest) (*slb.DescribeVServerGroupAttributeResponse, error)
}

type alibabacloudClient struct {
	ecsClient *ecs.Client
	vpcClient *vpc.Client
	slbClient *slb.Client
}

func (client *alibabacloudClient) RunInstances(request *ecs.RunInstancesRequest) (*ecs.RunInstancesResponse, error) {
	return client.ecsClient.RunInstances(request)
}

func (client *alibabacloudClient) CreateInstance(request *ecs.CreateInstanceRequest) (*ecs.CreateInstanceResponse, error) {
	return client.ecsClient.CreateInstance(request)
}

func (client *alibabacloudClient) DescribeInstances(request *ecs.DescribeInstancesRequest) (*ecs.DescribeInstancesResponse, error) {
	return client.ecsClient.DescribeInstances(request)
}

func (client *alibabacloudClient) DeleteInstances(request *ecs.DeleteInstancesRequest) (*ecs.DeleteInstancesResponse, error) {
	return client.ecsClient.DeleteInstances(request)
}

func (client *alibabacloudClient) StartInstance(request *ecs.StartInstanceRequest) (*ecs.StartInstanceResponse, error) {
	return client.ecsClient.StartInstance(request)
}

func (client *alibabacloudClient) RebootInstance(request *ecs.RebootInstanceRequest) (*ecs.RebootInstanceResponse, error) {
	return client.ecsClient.RebootInstance(request)
}

func (client *alibabacloudClient) StopInstance(request *ecs.StopInstanceRequest) (*ecs.StopInstanceResponse, error) {
	return client.ecsClient.StopInstance(request)
}

func (client *alibabacloudClient) StartInstances(request *ecs.StartInstancesRequest) (*ecs.StartInstancesResponse, error) {
	return client.ecsClient.StartInstances(request)
}

func (client *alibabacloudClient) RebootInstances(request *ecs.RebootInstancesRequest) (*ecs.RebootInstancesResponse, error) {
	return client.ecsClient.RebootInstances(request)
}

func (client *alibabacloudClient) StopInstances(request *ecs.StopInstancesRequest) (*ecs.StopInstancesResponse, error) {
	return client.ecsClient.StopInstances(request)
}

func (client *alibabacloudClient) DeleteInstance(request *ecs.DeleteInstanceRequest) (*ecs.DeleteInstanceResponse, error) {
	return client.ecsClient.DeleteInstance(request)
}

<<<<<<< HEAD
<<<<<<< HEAD
func (client *alibabacloudClient) AttachInstanceRAMRole(request *ecs.AttachInstanceRamRoleRequest) (*ecs.AttachInstanceRamRoleResponse, error) {
	return client.ecsClient.AttachInstanceRamRole(request)
=======
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
>>>>>>> ebdd9bd0 (update test case)
=======
func (client *alibabacloudClient) AttachInstanceRAMRole(request *ecs.AttachInstanceRamRoleRequest) (*ecs.AttachInstanceRamRoleResponse, error) {
	return client.ecsClient.AttachInstanceRamRole(request)
>>>>>>> e879a141 (alibabacloud machine-api provider)
}

func (client *alibabacloudClient) DetachInstanceRAMRole(request *ecs.DetachInstanceRamRoleRequest) (*ecs.DetachInstanceRamRoleResponse, error) {
	return client.ecsClient.DetachInstanceRamRole(request)
}

func (client *alibabacloudClient) DescribeInstanceStatus(request *ecs.DescribeInstanceStatusRequest) (*ecs.DescribeInstanceStatusResponse, error) {
	return client.ecsClient.DescribeInstanceStatus(request)
}

func (client *alibabacloudClient) ReActivateInstances(request *ecs.ReActivateInstancesRequest) (*ecs.ReActivateInstancesResponse, error) {
	return client.ecsClient.ReActivateInstances(request)
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
func (client *alibabacloudClient) DescribeUserData(request *ecs.DescribeUserDataRequest) (*ecs.DescribeUserDataResponse, error) {
	return client.ecsClient.DescribeUserData(request)
}

func (client *alibabacloudClient) DescribeInstanceTypes(request *ecs.DescribeInstanceTypesRequest) (*ecs.DescribeInstanceTypesResponse, error) {
	return client.ecsClient.DescribeInstanceTypes(request)
}

func (client *alibabacloudClient) ModifyInstanceAttribute(request *ecs.ModifyInstanceAttributeRequest) (*ecs.ModifyInstanceAttributeResponse, error) {
	return client.ecsClient.ModifyInstanceAttribute(request)
}

func (client *alibabacloudClient) ModifyInstanceMetadataOptions(request *ecs.ModifyInstanceMetadataOptionsRequest) (*ecs.ModifyInstanceMetadataOptionsResponse, error) {
	return client.ecsClient.ModifyInstanceMetadataOptions(request)
}

func (client *alibabacloudClient) AllocatePublicIPAddress(request *ecs.AllocatePublicIpAddressRequest) (*ecs.AllocatePublicIpAddressResponse, error) {
	return client.ecsClient.AllocatePublicIpAddress(request)
}

func (client *alibabacloudClient) CreateDisk(request *ecs.CreateDiskRequest) (*ecs.CreateDiskResponse, error) {
	return client.ecsClient.CreateDisk(request)
}

func (client *alibabacloudClient) AttachDisk(request *ecs.AttachDiskRequest) (*ecs.AttachDiskResponse, error) {
	return client.ecsClient.AttachDisk(request)
}

func (client *alibabacloudClient) DescribeDisks(request *ecs.DescribeDisksRequest) (*ecs.DescribeDisksResponse, error) {
	return client.ecsClient.DescribeDisks(request)
}

func (client *alibabacloudClient) ModifyDiskChargeType(request *ecs.ModifyDiskChargeTypeRequest) (*ecs.ModifyDiskChargeTypeResponse, error) {
	return client.ecsClient.ModifyDiskChargeType(request)
}

func (client *alibabacloudClient) ModifyDiskAttribute(request *ecs.ModifyDiskAttributeRequest) (*ecs.ModifyDiskAttributeResponse, error) {
	return client.ecsClient.ModifyDiskAttribute(request)
}

func (client *alibabacloudClient) ModifyDiskSpec(request *ecs.ModifyDiskSpecRequest) (*ecs.ModifyDiskSpecResponse, error) {
	return client.ecsClient.ModifyDiskSpec(request)
}

func (client *alibabacloudClient) ReplaceSystemDisk(request *ecs.ReplaceSystemDiskRequest) (*ecs.ReplaceSystemDiskResponse, error) {
	return client.ecsClient.ReplaceSystemDisk(request)
<<<<<<< HEAD
}

func (client *alibabacloudClient) ReInitDisk(request *ecs.ReInitDiskRequest) (*ecs.ReInitDiskResponse, error) {
	return client.ecsClient.ReInitDisk(request)
}

func (client *alibabacloudClient) ResetDisk(request *ecs.ResetDiskRequest) (*ecs.ResetDiskResponse, error) {
	return client.ecsClient.ResetDisk(request)
}

func (client *alibabacloudClient) ResizeDisk(request *ecs.ResizeDiskRequest) (*ecs.ResizeDiskResponse, error) {
	return client.ecsClient.ResizeDisk(request)
}

func (client *alibabacloudClient) DetachDisk(request *ecs.DetachDiskRequest) (*ecs.DetachDiskResponse, error) {
	return client.ecsClient.DetachDisk(request)
}

func (client *alibabacloudClient) DeleteDisk(request *ecs.DeleteDiskRequest) (*ecs.DeleteDiskResponse, error) {
	return client.ecsClient.DeleteDisk(request)
}

func (client *alibabacloudClient) DescribeRegions(request *ecs.DescribeRegionsRequest) (*ecs.DescribeRegionsResponse, error) {
	return client.ecsClient.DescribeRegions(request)
}

func (client *alibabacloudClient) DescribeZones(request *ecs.DescribeZonesRequest) (*ecs.DescribeZonesResponse, error) {
	return client.ecsClient.DescribeZones(request)
}

func (client *alibabacloudClient) DescribeImages(request *ecs.DescribeImagesRequest) (*ecs.DescribeImagesResponse, error) {
	return client.ecsClient.DescribeImages(request)
}

func (client *alibabacloudClient) CreateSecurityGroup(request *ecs.CreateSecurityGroupRequest) (*ecs.CreateSecurityGroupResponse, error) {
	return client.ecsClient.CreateSecurityGroup(request)
}

func (client *alibabacloudClient) AuthorizeSecurityGroup(request *ecs.AuthorizeSecurityGroupRequest) (*ecs.AuthorizeSecurityGroupResponse, error) {
	return client.ecsClient.AuthorizeSecurityGroup(request)
}

func (client *alibabacloudClient) AuthorizeSecurityGroupEgress(request *ecs.AuthorizeSecurityGroupEgressRequest) (*ecs.AuthorizeSecurityGroupEgressResponse, error) {
	return client.ecsClient.AuthorizeSecurityGroupEgress(request)
}

func (client *alibabacloudClient) RevokeSecurityGroup(request *ecs.RevokeSecurityGroupRequest) (*ecs.RevokeSecurityGroupResponse, error) {
	return client.ecsClient.RevokeSecurityGroup(request)
}

func (client *alibabacloudClient) RevokeSecurityGroupEgress(request *ecs.RevokeSecurityGroupEgressRequest) (*ecs.RevokeSecurityGroupEgressResponse, error) {
	return client.ecsClient.RevokeSecurityGroupEgress(request)
}

func (client *alibabacloudClient) JoinSecurityGroup(request *ecs.JoinSecurityGroupRequest) (*ecs.JoinSecurityGroupResponse, error) {
	return client.ecsClient.JoinSecurityGroup(request)
}

func (client *alibabacloudClient) LeaveSecurityGroup(request *ecs.LeaveSecurityGroupRequest) (*ecs.LeaveSecurityGroupResponse, error) {
	return client.ecsClient.LeaveSecurityGroup(request)
}

func (client *alibabacloudClient) DescribeSecurityGroupAttribute(request *ecs.DescribeSecurityGroupAttributeRequest) (*ecs.DescribeSecurityGroupAttributeResponse, error) {
	return client.ecsClient.DescribeSecurityGroupAttribute(request)
}

func (client *alibabacloudClient) DescribeSecurityGroups(request *ecs.DescribeSecurityGroupsRequest) (*ecs.DescribeSecurityGroupsResponse, error) {
	return client.ecsClient.DescribeSecurityGroups(request)
}

func (client *alibabacloudClient) DescribeSecurityGroupReferences(request *ecs.DescribeSecurityGroupReferencesRequest) (*ecs.DescribeSecurityGroupReferencesResponse, error) {
	return client.ecsClient.DescribeSecurityGroupReferences(request)
}

func (client *alibabacloudClient) ModifySecurityGroupAttribute(request *ecs.ModifySecurityGroupAttributeRequest) (*ecs.ModifySecurityGroupAttributeResponse, error) {
	return client.ecsClient.ModifySecurityGroupAttribute(request)
}

func (client *alibabacloudClient) ModifySecurityGroupEgressRule(request *ecs.ModifySecurityGroupEgressRuleRequest) (*ecs.ModifySecurityGroupEgressRuleResponse, error) {
	return client.ecsClient.ModifySecurityGroupEgressRule(request)
}

func (client *alibabacloudClient) ModifySecurityGroupPolicy(request *ecs.ModifySecurityGroupPolicyRequest) (*ecs.ModifySecurityGroupPolicyResponse, error) {
	return client.ecsClient.ModifySecurityGroupPolicy(request)
}

func (client *alibabacloudClient) ModifySecurityGroupRule(request *ecs.ModifySecurityGroupRuleRequest) (*ecs.ModifySecurityGroupRuleResponse, error) {
	return client.ecsClient.ModifySecurityGroupRule(request)
}

func (client *alibabacloudClient) DeleteSecurityGroup(request *ecs.DeleteSecurityGroupRequest) (*ecs.DeleteSecurityGroupResponse, error) {
	return client.ecsClient.DeleteSecurityGroup(request)
}

func (client *alibabacloudClient) TagResources(request *ecs.TagResourcesRequest) (*ecs.TagResourcesResponse, error) {
	return client.ecsClient.TagResources(request)
}

func (client *alibabacloudClient) ListTagResources(request *ecs.ListTagResourcesRequest) (*ecs.ListTagResourcesResponse, error) {
	return client.ecsClient.ListTagResources(request)
}

func (client *alibabacloudClient) UntagResources(request *ecs.UntagResourcesRequest) (*ecs.UntagResourcesResponse, error) {
	return client.ecsClient.UntagResources(request)
}

func (client *alibabacloudClient) CreateVpc(request *vpc.CreateVpcRequest) (*vpc.CreateVpcResponse, error) {
	return client.vpcClient.CreateVpc(request)
}

func (client *alibabacloudClient) DeleteVpc(request *vpc.DeleteVpcRequest) (*vpc.DeleteVpcResponse, error) {
	return client.vpcClient.DeleteVpc(request)
}

func (client *alibabacloudClient) DescribeVpcs(request *vpc.DescribeVpcsRequest) (*vpc.DescribeVpcsResponse, error) {
	return client.vpcClient.DescribeVpcs(request)
}

func (client *alibabacloudClient) CreateVSwitch(request *vpc.CreateVSwitchRequest) (*vpc.CreateVSwitchResponse, error) {
	return client.vpcClient.CreateVSwitch(request)
}

func (client *alibabacloudClient) DeleteVSwitch(request *vpc.DeleteVSwitchRequest) (*vpc.DeleteVSwitchResponse, error) {
	return client.vpcClient.DeleteVSwitch(request)
}

func (client *alibabacloudClient) DescribeVSwitches(request *vpc.DescribeVSwitchesRequest) (*vpc.DescribeVSwitchesResponse, error) {
	return client.vpcClient.DescribeVSwitches(request)
}

func (client *alibabacloudClient) CreateNatGateway(request *vpc.CreateNatGatewayRequest) (*vpc.CreateNatGatewayResponse, error) {
	return client.vpcClient.CreateNatGateway(request)
}

func (client *alibabacloudClient) DescribeNatGateways(request *vpc.DescribeNatGatewaysRequest) (*vpc.DescribeNatGatewaysResponse, error) {
	return client.vpcClient.DescribeNatGateways(request)
}

func (client *alibabacloudClient) DeleteNatGateway(request *vpc.DeleteNatGatewayRequest) (*vpc.DeleteNatGatewayResponse, error) {
	return client.vpcClient.DeleteNatGateway(request)
}

=======
}

func (client *alibabacloudClient) ReInitDisk(request *ecs.ReInitDiskRequest) (*ecs.ReInitDiskResponse, error) {
	return client.ecsClient.ReInitDisk(request)
}

func (client *alibabacloudClient) ResetDisk(request *ecs.ResetDiskRequest) (*ecs.ResetDiskResponse, error) {
	return client.ecsClient.ResetDisk(request)
}

func (client *alibabacloudClient) ResizeDisk(request *ecs.ResizeDiskRequest) (*ecs.ResizeDiskResponse, error) {
	return client.ecsClient.ResizeDisk(request)
}

func (client *alibabacloudClient) DetachDisk(request *ecs.DetachDiskRequest) (*ecs.DetachDiskResponse, error) {
	return client.ecsClient.DetachDisk(request)
}

func (client *alibabacloudClient) DeleteDisk(request *ecs.DeleteDiskRequest) (*ecs.DeleteDiskResponse, error) {
	return client.ecsClient.DeleteDisk(request)
}

func (client *alibabacloudClient) DescribeRegions(request *ecs.DescribeRegionsRequest) (*ecs.DescribeRegionsResponse, error) {
	return client.ecsClient.DescribeRegions(request)
}

func (client *alibabacloudClient) DescribeZones(request *ecs.DescribeZonesRequest) (*ecs.DescribeZonesResponse, error) {
	return client.ecsClient.DescribeZones(request)
}

func (client *alibabacloudClient) DescribeImages(request *ecs.DescribeImagesRequest) (*ecs.DescribeImagesResponse, error) {
	return client.ecsClient.DescribeImages(request)
}

func (client *alibabacloudClient) CreateSecurityGroup(request *ecs.CreateSecurityGroupRequest) (*ecs.CreateSecurityGroupResponse, error) {
	return client.ecsClient.CreateSecurityGroup(request)
}

func (client *alibabacloudClient) AuthorizeSecurityGroup(request *ecs.AuthorizeSecurityGroupRequest) (*ecs.AuthorizeSecurityGroupResponse, error) {
	return client.ecsClient.AuthorizeSecurityGroup(request)
}

func (client *alibabacloudClient) AuthorizeSecurityGroupEgress(request *ecs.AuthorizeSecurityGroupEgressRequest) (*ecs.AuthorizeSecurityGroupEgressResponse, error) {
	return client.ecsClient.AuthorizeSecurityGroupEgress(request)
}

func (client *alibabacloudClient) RevokeSecurityGroup(request *ecs.RevokeSecurityGroupRequest) (*ecs.RevokeSecurityGroupResponse, error) {
	return client.ecsClient.RevokeSecurityGroup(request)
}

func (client *alibabacloudClient) RevokeSecurityGroupEgress(request *ecs.RevokeSecurityGroupEgressRequest) (*ecs.RevokeSecurityGroupEgressResponse, error) {
	return client.ecsClient.RevokeSecurityGroupEgress(request)
}

func (client *alibabacloudClient) JoinSecurityGroup(request *ecs.JoinSecurityGroupRequest) (*ecs.JoinSecurityGroupResponse, error) {
	return client.ecsClient.JoinSecurityGroup(request)
}

func (client *alibabacloudClient) LeaveSecurityGroup(request *ecs.LeaveSecurityGroupRequest) (*ecs.LeaveSecurityGroupResponse, error) {
	return client.ecsClient.LeaveSecurityGroup(request)
}

func (client *alibabacloudClient) DescribeSecurityGroupAttribute(request *ecs.DescribeSecurityGroupAttributeRequest) (*ecs.DescribeSecurityGroupAttributeResponse, error) {
	return client.ecsClient.DescribeSecurityGroupAttribute(request)
}

func (client *alibabacloudClient) DescribeSecurityGroups(request *ecs.DescribeSecurityGroupsRequest) (*ecs.DescribeSecurityGroupsResponse, error) {
	return client.ecsClient.DescribeSecurityGroups(request)
}

func (client *alibabacloudClient) DescribeSecurityGroupReferences(request *ecs.DescribeSecurityGroupReferencesRequest) (*ecs.DescribeSecurityGroupReferencesResponse, error) {
	return client.ecsClient.DescribeSecurityGroupReferences(request)
}

func (client *alibabacloudClient) ModifySecurityGroupAttribute(request *ecs.ModifySecurityGroupAttributeRequest) (*ecs.ModifySecurityGroupAttributeResponse, error) {
	return client.ecsClient.ModifySecurityGroupAttribute(request)
}

func (client *alibabacloudClient) ModifySecurityGroupEgressRule(request *ecs.ModifySecurityGroupEgressRuleRequest) (*ecs.ModifySecurityGroupEgressRuleResponse, error) {
	return client.ecsClient.ModifySecurityGroupEgressRule(request)
}

func (client *alibabacloudClient) ModifySecurityGroupPolicy(request *ecs.ModifySecurityGroupPolicyRequest) (*ecs.ModifySecurityGroupPolicyResponse, error) {
	return client.ecsClient.ModifySecurityGroupPolicy(request)
}

func (client *alibabacloudClient) ModifySecurityGroupRule(request *ecs.ModifySecurityGroupRuleRequest) (*ecs.ModifySecurityGroupRuleResponse, error) {
	return client.ecsClient.ModifySecurityGroupRule(request)
}

func (client *alibabacloudClient) DeleteSecurityGroup(request *ecs.DeleteSecurityGroupRequest) (*ecs.DeleteSecurityGroupResponse, error) {
	return client.ecsClient.DeleteSecurityGroup(request)
}

func (client *alibabacloudClient) TagResources(request *ecs.TagResourcesRequest) (*ecs.TagResourcesResponse, error) {
	return client.ecsClient.TagResources(request)
}

func (client *alibabacloudClient) ListTagResources(request *ecs.ListTagResourcesRequest) (*ecs.ListTagResourcesResponse, error) {
	return client.ecsClient.ListTagResources(request)
}

func (client *alibabacloudClient) UntagResources(request *ecs.UntagResourcesRequest) (*ecs.UntagResourcesResponse, error) {
	return client.ecsClient.UntagResources(request)
}

func (client *alibabacloudClient) CreateVpc(request *vpc.CreateVpcRequest) (*vpc.CreateVpcResponse, error) {
	return client.vpcClient.CreateVpc(request)
}

func (client *alibabacloudClient) DeleteVpc(request *vpc.DeleteVpcRequest) (*vpc.DeleteVpcResponse, error) {
	return client.vpcClient.DeleteVpc(request)
}

func (client *alibabacloudClient) DescribeVpcs(request *vpc.DescribeVpcsRequest) (*vpc.DescribeVpcsResponse, error) {
	return client.vpcClient.DescribeVpcs(request)
}

func (client *alibabacloudClient) CreateVSwitch(request *vpc.CreateVSwitchRequest) (*vpc.CreateVSwitchResponse, error) {
	return client.vpcClient.CreateVSwitch(request)
}

func (client *alibabacloudClient) DeleteVSwitch(request *vpc.DeleteVSwitchRequest) (*vpc.DeleteVSwitchResponse, error) {
	return client.vpcClient.DeleteVSwitch(request)
}

func (client *alibabacloudClient) DescribeVSwitches(request *vpc.DescribeVSwitchesRequest) (*vpc.DescribeVSwitchesResponse, error) {
	return client.vpcClient.DescribeVSwitches(request)
}

func (client *alibabacloudClient) CreateNatGateway(request *vpc.CreateNatGatewayRequest) (*vpc.CreateNatGatewayResponse, error) {
	return client.vpcClient.CreateNatGateway(request)
}

func (client *alibabacloudClient) DescribeNatGateways(request *vpc.DescribeNatGatewaysRequest) (*vpc.DescribeNatGatewaysResponse, error) {
	return client.vpcClient.DescribeNatGateways(request)
}

func (client *alibabacloudClient) DeleteNatGateway(request *vpc.DeleteNatGatewayRequest) (*vpc.DeleteNatGatewayResponse, error) {
	return client.vpcClient.DeleteNatGateway(request)
}

>>>>>>> e879a141 (alibabacloud machine-api provider)
func (client *alibabacloudClient) AllocateEipAddress(request *vpc.AllocateEipAddressRequest) (*vpc.AllocateEipAddressResponse, error) {
	return client.vpcClient.AllocateEipAddress(request)
}

func (client *alibabacloudClient) AssociateEipAddress(request *vpc.AssociateEipAddressRequest) (*vpc.AssociateEipAddressResponse, error) {
	return client.vpcClient.AssociateEipAddress(request)
}

func (client *alibabacloudClient) ModifyEipAddressAttribute(request *vpc.ModifyEipAddressAttributeRequest) (*vpc.ModifyEipAddressAttributeResponse, error) {
	return client.vpcClient.ModifyEipAddressAttribute(request)
}

func (client *alibabacloudClient) DescribeEipAddresses(request *vpc.DescribeEipAddressesRequest) (*vpc.DescribeEipAddressesResponse, error) {
	return client.vpcClient.DescribeEipAddresses(request)
}

func (client *alibabacloudClient) UnassociateEipAddress(request *vpc.UnassociateEipAddressRequest) (*vpc.UnassociateEipAddressResponse, error) {
	return client.vpcClient.UnassociateEipAddress(request)
}

func (client *alibabacloudClient) ReleaseEipAddress(request *vpc.ReleaseEipAddressRequest) (*vpc.ReleaseEipAddressResponse, error) {
	return client.vpcClient.ReleaseEipAddress(request)
}

func (client *alibabacloudClient) CreateLoadBalancer(request *slb.CreateLoadBalancerRequest) (*slb.CreateLoadBalancerResponse, error) {
	return client.slbClient.CreateLoadBalancer(request)
}

func (client *alibabacloudClient) DeleteLoadBalancer(request *slb.DeleteLoadBalancerRequest) (*slb.DeleteLoadBalancerResponse, error) {
	return client.slbClient.DeleteLoadBalancer(request)
}

func (client *alibabacloudClient) DescribeLoadBalancers(request *slb.DescribeLoadBalancersRequest) (*slb.DescribeLoadBalancersResponse, error) {
	return client.slbClient.DescribeLoadBalancers(request)
}

func (client *alibabacloudClient) CreateLoadBalancerTCPListener(request *slb.CreateLoadBalancerTCPListenerRequest) (*slb.CreateLoadBalancerTCPListenerResponse, error) {
	return client.slbClient.CreateLoadBalancerTCPListener(request)
}

func (client *alibabacloudClient) SetLoadBalancerTCPListenerAttribute(request *slb.SetLoadBalancerTCPListenerAttributeRequest) (*slb.SetLoadBalancerTCPListenerAttributeResponse, error) {
	return client.slbClient.SetLoadBalancerTCPListenerAttribute(request)
}

func (client *alibabacloudClient) DescribeLoadBalancerTCPListenerAttribute(request *slb.DescribeLoadBalancerTCPListenerAttributeRequest) (*slb.DescribeLoadBalancerTCPListenerAttributeResponse, error) {
	return client.slbClient.DescribeLoadBalancerTCPListenerAttribute(request)
}

func (client *alibabacloudClient) CreateLoadBalancerUDPListener(request *slb.CreateLoadBalancerUDPListenerRequest) (*slb.CreateLoadBalancerUDPListenerResponse, error) {
	return client.slbClient.CreateLoadBalancerUDPListener(request)
}

func (client *alibabacloudClient) SetLoadBalancerUDPListenerAttribute(request *slb.SetLoadBalancerUDPListenerAttributeRequest) (*slb.SetLoadBalancerUDPListenerAttributeResponse, error) {
	return client.slbClient.SetLoadBalancerUDPListenerAttribute(request)
}

func (client *alibabacloudClient) DescribeLoadBalancerUDPListenerAttribute(request *slb.DescribeLoadBalancerUDPListenerAttributeRequest) (*slb.DescribeLoadBalancerUDPListenerAttributeResponse, error) {
	return client.slbClient.DescribeLoadBalancerUDPListenerAttribute(request)
}

func (client *alibabacloudClient) CreateLoadBalancerHTTPListener(request *slb.CreateLoadBalancerHTTPListenerRequest) (*slb.CreateLoadBalancerHTTPListenerResponse, error) {
	return client.slbClient.CreateLoadBalancerHTTPListener(request)
}

func (client *alibabacloudClient) SetLoadBalancerHTTPListenerAttribute(request *slb.SetLoadBalancerHTTPListenerAttributeRequest) (*slb.SetLoadBalancerHTTPListenerAttributeResponse, error) {
	return client.slbClient.SetLoadBalancerHTTPListenerAttribute(request)
}

func (client *alibabacloudClient) DescribeLoadBalancerHTTPListenerAttribute(request *slb.DescribeLoadBalancerHTTPListenerAttributeRequest) (*slb.DescribeLoadBalancerHTTPListenerAttributeResponse, error) {
	return client.slbClient.DescribeLoadBalancerHTTPListenerAttribute(request)
}

func (client *alibabacloudClient) CreateLoadBalancerHTTPSListener(request *slb.CreateLoadBalancerHTTPSListenerRequest) (*slb.CreateLoadBalancerHTTPSListenerResponse, error) {
	return client.slbClient.CreateLoadBalancerHTTPSListener(request)
}

func (client *alibabacloudClient) SetLoadBalancerHTTPSListenerAttribute(request *slb.SetLoadBalancerHTTPSListenerAttributeRequest) (*slb.SetLoadBalancerHTTPSListenerAttributeResponse, error) {
	return client.slbClient.SetLoadBalancerHTTPSListenerAttribute(request)
}

func (client *alibabacloudClient) DescribeLoadBalancerHTTPSListenerAttribute(request *slb.DescribeLoadBalancerHTTPSListenerAttributeRequest) (*slb.DescribeLoadBalancerHTTPSListenerAttributeResponse, error) {
	return client.slbClient.DescribeLoadBalancerHTTPSListenerAttribute(request)
}

func (client *alibabacloudClient) StartLoadBalancerListener(request *slb.StartLoadBalancerListenerRequest) (*slb.StartLoadBalancerListenerResponse, error) {
	return client.slbClient.StartLoadBalancerListener(request)
}

func (client *alibabacloudClient) StopLoadBalancerListener(request *slb.StopLoadBalancerListenerRequest) (*slb.StopLoadBalancerListenerResponse, error) {
	return client.slbClient.StopLoadBalancerListener(request)
}

func (client *alibabacloudClient) DeleteLoadBalancerListener(request *slb.DeleteLoadBalancerListenerRequest) (*slb.DeleteLoadBalancerListenerResponse, error) {
	return client.slbClient.DeleteLoadBalancerListener(request)
}

func (client *alibabacloudClient) DescribeLoadBalancerListeners(request *slb.DescribeLoadBalancerListenersRequest) (*slb.DescribeLoadBalancerListenersResponse, error) {
	return client.slbClient.DescribeLoadBalancerListeners(request)
}

func (client *alibabacloudClient) AddBackendServers(request *slb.AddBackendServersRequest) (*slb.AddBackendServersResponse, error) {
	return client.slbClient.AddBackendServers(request)
}

func (client *alibabacloudClient) RemoveBackendServers(request *slb.RemoveBackendServersRequest) (*slb.RemoveBackendServersResponse, error) {
	return client.slbClient.RemoveBackendServers(request)
}

func (client *alibabacloudClient) SetBackendServers(request *slb.SetBackendServersRequest) (*slb.SetBackendServersResponse, error) {
	return client.slbClient.SetBackendServers(request)
}

func (client *alibabacloudClient) DescribeHealthStatus(request *slb.DescribeHealthStatusRequest) (*slb.DescribeHealthStatusResponse, error) {
	return client.slbClient.DescribeHealthStatus(request)
}

func (client *alibabacloudClient) CreateVServerGroup(request *slb.CreateVServerGroupRequest) (*slb.CreateVServerGroupResponse, error) {
	return client.slbClient.CreateVServerGroup(request)
}

func (client *alibabacloudClient) SetVServerGroupAttribute(request *slb.SetVServerGroupAttributeRequest) (*slb.SetVServerGroupAttributeResponse, error) {
	return client.slbClient.SetVServerGroupAttribute(request)
}

func (client *alibabacloudClient) AddVServerGroupBackendServers(request *slb.AddVServerGroupBackendServersRequest) (*slb.AddVServerGroupBackendServersResponse, error) {
	return client.slbClient.AddVServerGroupBackendServers(request)
}

func (client *alibabacloudClient) RemoveVServerGroupBackendServers(request *slb.RemoveVServerGroupBackendServersRequest) (*slb.RemoveVServerGroupBackendServersResponse, error) {
	return client.slbClient.RemoveVServerGroupBackendServers(request)
}

func (client *alibabacloudClient) ModifyVServerGroupBackendServers(request *slb.ModifyVServerGroupBackendServersRequest) (*slb.ModifyVServerGroupBackendServersResponse, error) {
	return client.slbClient.ModifyVServerGroupBackendServers(request)
}

func (client *alibabacloudClient) DeleteVServerGroup(request *slb.DeleteVServerGroupRequest) (*slb.DeleteVServerGroupResponse, error) {
	return client.slbClient.DeleteVServerGroup(request)
}

func (client *alibabacloudClient) DescribeVServerGroups(request *slb.DescribeVServerGroupsRequest) (*slb.DescribeVServerGroupsResponse, error) {
	return client.slbClient.DescribeVServerGroups(request)
}

func (client *alibabacloudClient) DescribeVServerGroupAttribute(request *slb.DescribeVServerGroupAttributeRequest) (*slb.DescribeVServerGroupAttributeResponse, error) {
	return client.slbClient.DescribeVServerGroupAttribute(request)
}

// NewClient creates our client wrapper object for the actual alibabacloud clients we use.
<<<<<<< HEAD
<<<<<<< HEAD
func NewClient(ctrlRuntimeClient client.Client, secretName, namespace, regionID string, configManagedClient client.Client) (Client, error) {
	config, err := newConfiguration(ctrlRuntimeClient, secretName, namespace, configManagedClient)
=======
// For authentication the underlying clients will use either the cluster alibabacloud credentials
// secret if defined (i.e. in the root cluster),
// otherwise the IAM profile of the master where the actuator will run. (target clusters)
<<<<<<< HEAD
func NewClient(ctrlRuntimeClient client.Client, secretName, namespace, region string, configManagedClient client.Client) (Client, error) {
	config, err := newConfiguration(ctrlRuntimeClient, secretName, namespace, region, configManagedClient)
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
=======
>>>>>>> 24c35849 (fix stop ecs instance func)
func NewClient(ctrlRuntimeClient client.Client, secretName, namespace, regionID string, configManagedClient client.Client) (Client, error) {
	config, err := newConfiguration(ctrlRuntimeClient, secretName, namespace, configManagedClient)
>>>>>>> 60dde8f7 (update Makefile)
	if err != nil {
		return nil, err
	}

	provider := providers.NewConfigurationCredentialProvider(config)

	credential, err := provider.Retrieve()
	if err != nil {
		klog.Errorf("error %v ", err)
		return nil, err
<<<<<<< HEAD
	}

	sdkConfig := &sdk.Config{
		UserAgent: machineProviderUserAgent,
		Scheme:    "HTTPS",
	}
	//init ecsClient
	ecsClient, err := ecs.NewClientWithOptions(regionID, sdkConfig, credential)
	if err != nil {
		klog.Errorf("failed to init ecs client %v", err)
		return nil, err
	}

	//init vpcClient
	vpcClient, err := vpc.NewClientWithOptions(regionID, sdkConfig, credential)
	if err != nil {
		klog.Errorf("failed to init vpc client %v", err)
		return nil, err
	}

	//init slbClient
	slbClient, err := slb.NewClientWithOptions(regionID, sdkConfig, credential)
	if err != nil {
		klog.Errorf("failed to init slb client %v", err)
		return nil, err
	}

	return &alibabacloudClient{
		ecsClient: ecsClient,
		vpcClient: vpcClient,
		slbClient: slbClient,
	}, nil
}

//Init alibabacloud configuration
//https://github.com/aliyun/alibaba-cloud-sdk-go/blob/master/sdk/auth/credentials/providers/configuration.go
func newConfiguration(ctrlRuntimeClient client.Client, secretName, namespace string, configManagedClient client.Client) (*providers.Configuration, error) {
	config := &providers.Configuration{}
=======
func NewClient(ctrlRuntimeClient client.Client, secretName, namespace, region string) (Client, error) {
	aliCloudConfig := &providers.Configuration{
		AccessKeyID:     os.Getenv("ALICLOUD_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("ALICLOUD_ACCESS_KEY_SECRET"),
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
	}
>>>>>>> ebdd9bd0 (update test case)

	sdkConfig := &sdk.Config{
		UserAgent: machineProviderUserAgent,
		Scheme:    "HTTPS",
	}
	//init ecsClient
	ecsClient, err := ecs.NewClientWithOptions(regionID, sdkConfig, credential)
	if err != nil {
		klog.Errorf("failed to init ecs client %v", err)
		return nil, err
	}

	//init vpcClient
	vpcClient, err := vpc.NewClientWithOptions(regionID, sdkConfig, credential)
	if err != nil {
		klog.Errorf("failed to init vpc client %v", err)
		return nil, err
	}

	//init slbClient
	slbClient, err := slb.NewClientWithOptions(regionID, sdkConfig, credential)
	if err != nil {
		klog.Errorf("failed to init slb client %v", err)
		return nil, err
	}

	return &alibabacloudClient{
		ecsClient: ecsClient,
		vpcClient: vpcClient,
		slbClient: slbClient,
	}, nil
}

//Init alibabacloud configuration
//https://github.com/aliyun/alibaba-cloud-sdk-go/blob/master/sdk/auth/credentials/providers/configuration.go
func newConfiguration(ctrlRuntimeClient client.Client, secretName, namespace string, configManagedClient client.Client) (*providers.Configuration, error) {
	config := &providers.Configuration{}

	if secretName != "" {
		var secret corev1.Secret
		if err := ctrlRuntimeClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: secretName}, &secret); err != nil {
			if apimachineryerrors.IsNotFound(err) {
				return nil, machineapiapierrors.InvalidMachineConfiguration("alibabacloud credentials secret %s/%s: %v not found", namespace, secretName, err)
			}
			return nil, err
		}
		err := fetchCredentialsFileFromSecret(&secret, config)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch  credentials file from Secret: %v", err)
		}
<<<<<<< HEAD
<<<<<<< HEAD
	}

	return config, nil
}

// fetchCredentialsFileFromSecret fetch credentials from screct
// alibabacloud accessKeyID accessKeySecret stsToken
//roleArn, roleName ,roleSessionName ,roleSessionExpiration
func fetchCredentialsFileFromSecret(secret *corev1.Secret, config *providers.Configuration) error {
	//accessKeyID/accessKeySecret
	if len(secret.Data[kubeAccessKeyID]) > 0 && len(secret.Data[kubeAccessKeySecret]) > 0 {
		config.AccessKeyID = utils.ByteArray2String(secret.Data[kubeAccessKeyID])
		config.AccessKeySecret = utils.ByteArray2String(secret.Data[kubeAccessKeySecret])
	}

	//ststoken
	if len(secret.Data[kubeAccessKeyStsToken]) > 0 {
		config.AccessKeyStsToken = utils.ByteArray2String(secret.Data[kubeAccessKeyStsToken])
=======
		aliCloudConfig.AccessKeyID = string(accessKeyID)

		accessKeySecret, ok := secret.Data[AliCloudAccessKeySecret]
		if !ok {
			return nil, fmt.Errorf("AliCloud credentials secret %v did not contain key %v",
				secretName, AliCloudAccessKeySecret)
		}
		aliCloudConfig.AccessKeySecret = string(accessKeySecret)
		accessKeyStsToken, _ := secret.Data[AliCloudAccessKeyStsToken]
		aliCloudConfig.AccessKeyStsToken = string(accessKeyStsToken)
>>>>>>> ebdd9bd0 (update test case)
	}

	// roleArn ,roleSessionName ,roleSessionExpiration
	if len(secret.Data[kubeRoleArn]) > 0 && len(secret.Data[kubeRoleSessionName]) > 0 && len(secret.Data[kubeRoleSessionExpiration]) > 0 {
		if roleSessionExpiration, err := utils.String2IntPointer(string(secret.Data[kubeRoleSessionExpiration])); roleSessionExpiration != nil && err == nil {
			config.RoleArn = string(secret.Data[kubeRoleArn])
			config.RoleSessionName = string(secret.Data[kubeRoleSessionName])
			config.RoleSessionExpiration = roleSessionExpiration
		}
	}

	// roleName
	if len(secret.Data[kubeRoleName]) > 0 {
		config.RoleName = string(secret.Data[kubeRoleName])
	}

	return nil
}

// fetchCredentialsFileFromConfigMap fetch oleArn, roleName ,roleSessionName ,roleSessionExpiration from configmap
=======
	}

	if err := fetchCredentialsFileFromConfigMap(namespace, configManagedClient, config); err != nil {
		klog.Errorf("failed to fetch credentials from ConfigMap : %v", err)
	}

	return config, nil
}

// fetchCredentialsFileFromSecret fetch credentials from screct
// alibabacloud accessKeyID accessKeySecret
func fetchCredentialsFileFromSecret(secret *corev1.Secret, config *providers.Configuration) error {
	//accessKeyID/accessKeySecret
	if len(secret.Data[kubeAccessKeyID]) > 0 && len(secret.Data[kubeAccessKeySecret]) > 0 {
		config.AccessKeyID = utils.ByteArray2String(secret.Data[kubeAccessKeyID])
		config.AccessKeySecret = utils.ByteArray2String(secret.Data[kubeAccessKeySecret])
	}

	//ststoken
	if len(secret.Data[kubeAccessKeyStsToken]) > 0 {
		config.AccessKeyStsToken = utils.ByteArray2String(secret.Data[kubeAccessKeyStsToken])
	}

	return nil
}

<<<<<<< HEAD
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
// fetchCredentialsFileFromConfigMap
>>>>>>> 24c35849 (fix stop ecs instance func)
// roleArn, roleName ,roleSessionName ,roleSessionExpiration
func fetchCredentialsFileFromConfigMap(namespace string, configManagedClient client.Client, config *providers.Configuration) error {
	cm := &corev1.ConfigMap{}
	switch err := configManagedClient.Get(
		context.Background(),
		client.ObjectKey{Namespace: namespace, Name: kubeCloudConfigName},
		cm,
	); {
	case apimachineryerrors.IsNotFound(err):
		// no cloud config ConfigMap
		return nil
	case err != nil:
		return fmt.Errorf("failed to get kube-cloud-config ConfigMap: %w", err)
	}

	if len(cm.Data[kubeRoleArn]) > 0 && len(cm.Data[kubeRoleSessionName]) > 0 && len(cm.Data[kubeRoleSessionExpiration]) > 0 {
		if roleSessionExpiration, err := utils.String2IntPointer(cm.Data[kubeRoleSessionExpiration]); roleSessionExpiration != nil && err == nil {
			config.RoleArn = cm.Data[kubeRoleArn]
			config.RoleSessionName = cm.Data[kubeRoleSessionName]
			config.RoleSessionExpiration = roleSessionExpiration
		}
	}

	if len(cm.Data[kubeRoleName]) > 0 {
		config.RoleName = cm.Data[kubeRoleName]
	}

	return nil
}
