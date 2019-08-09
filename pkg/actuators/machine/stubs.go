package machine

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machinecontroller "github.com/openshift/machine-api-operator/pkg/controller/machine"

	alibabacloudproviderv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudprovider/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultNamespace                     = "default"
	stubZoneID                           = "cn-beijing"
	stubRegionID                         = "cn-beijing-f"
	alibabaCloudCredentialsSecretName    = "alibabacloud-credentials-secret"
	alibabaCloudMasterUserDataSecretName = "master-user-data-secret"
	alibabaCloudWorkerUserDataSecretName = "worker-user-data-secret"

	stubMasterMachineName = "alibabacloud-actuator-testing-master-machine"
	stubWorkerMachineName = "alibabacloud-actuator-testing-worker-machine"

	stubKeyName                 = "alibabacloud-actuator-key-name"
	stubClusterID               = "alibabacloud-actuator-cluster"
	stubImageID                 = "centos_7_9_x64_20G_alibase_20210318.vhd"
	stubVpcID                   = "vpc-3ze4u29pd4lniym7i1xnp"
	stubVSwitchID               = "vsw-7ze567qrl5das7q8s4rei"
	stubInstanceID              = "i-2ze3hj0qh9d290rpax7w"
	stubSecurityGroupId         = "sg-2zeebk9qd965vc2xqq4w"
	stubSystemDiskCategory      = "cloud_essd"
	stubSystemDiskSize          = 120
	stubInternetMaxBandwidthOut = 100
	stubPassword                = "Hello$1234"
	stubInstanceType            = "ecs.c6.2xlarge"
)

<<<<<<< HEAD
func stubAlibabaCloudCredentialsSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      alibabaCloudCredentialsSecretName,
			Namespace: defaultNamespace,
		},
		Data: map[string][]byte{
			"accessKeyID":     []byte(os.Getenv("ALIBABACLOUD_ACCESS_KEY_ID")),
			"accessKeySecret": []byte(os.Getenv("ALIBABACLOUD_SECRET_ACCESS_KEY")),
		},
	}
}

func stubProviderConfig() *alibabacloudproviderv1.AlibabaCloudMachineProviderConfig {
	return &alibabacloudproviderv1.AlibabaCloudMachineProviderConfig{
		InstanceType:            stubInstanceType,
		ImageID:                 stubImageID,
		RegionID:                stubRegionID,
		ZoneID:                  stubZoneID,
		SecurityGroupID:         stubSecurityGroupId,
		VpcID:                   stubVpcID,
		VSwitchID:               stubVSwitchID,
		SystemDiskCategory:      stubSystemDiskCategory,
		SystemDiskSize:          stubSystemDiskSize,
		InternetMaxBandwidthOut: stubInternetMaxBandwidthOut,
		Password:                stubPassword,
		CredentialsSecret: &corev1.LocalObjectReference{
			Name: alibabaCloudCredentialsSecretName,
		},
		Tags: []alibabacloudproviderv1.Tag{
			{Key: "openshift-node-group-config", Value: "node-config-master"},
			{Key: "host-type", Value: "master"},
			{Key: "sub-host-type", Value: "default"},
		},
	}
}

func stubMasterMachine() (*machinev1.Machine, error) {
	masterMachine, err := stubMachine(stubMasterMachineName, map[string]string{
		"node-role.kubernetes.io/master": "",
		"node-role.kubernetes.io/infra":  "",
	})

	if err != nil {
		return nil, err
=======
const userDataBlob = `#cloud-config
write_files:
- path: /root/node_bootstrap/node_settings.yaml
 owner: 'root:root'
 permissions: '0640'
 content: |
   node_config_name: node-config-master
runcmd:
- [ cat, /root/node_bootstrap/node_settings.yaml]
`

func stubProviderConfig() *providerconfigv1.AlicloudMachineProviderConfig {
	return &providerconfigv1.AlicloudMachineProviderConfig{
		ImageId: "centos_7_06_64_20G_alibase_20190619.vhd",
		CredentialsSecret: &apiv1.LocalObjectReference{
			Name: alicloudCredentialsSecretName,
		},
		InstanceType: "ecs.n4.xlarge",
		RegionId:     "cn-hangzhou",
		VpcId:        "vpc-bp1td11g1i90b1fjnm7jw",
		VSwitchId:    "vsw-bp1ra53n8ban94mbbgb4w",
		Tags: []providerconfigv1.TagSpecification{
			{
				Key:   "openshift-node-group-config",
				Value: "node-config-master",
			},
			{
				Key:   "host-type",
				Value: "master",
			},
			{
				Key:   "sub-host-type",
				Value: "default",
			},
		},
		SecurityGroupId: "sg-bp1iccjoxddumf300okm",
		PublicIP:        true,
>>>>>>> ebdd9bd0 (update test case)
	}

	return masterMachine, nil
}

func stubWorkerMachine() (*machinev1.Machine, error) {
	workerMachine, err := stubMachine(stubWorkerMachineName, map[string]string{
		"node-role.kubernetes.io/infra": "",
	})

	if err != nil {
		return nil, err
	}

	return workerMachine, nil
}

func stubMachine(machineName string, machineLabels map[string]string) (*machinev1.Machine, error) {
	machineSpec := stubProviderConfig()

	providerSpec, err := alibabacloudproviderv1.RawExtensionFromProviderSpec(machineSpec)
	if err != nil {
		return nil, fmt.Errorf("codec.EncodeProviderSpec failed: %v", err)
	}

	machine := &machinev1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machineName,
			Namespace: defaultNamespace,
			Labels: map[string]string{
				machinev1.MachineClusterIDLabel: stubClusterID,
			},
			Annotations: map[string]string{
				// skip node draining since it's not mocked
				machinecontroller.ExcludeNodeDrainingAnnotation: "",
			},
		},
		Spec: machinev1.MachineSpec{
			ObjectMeta: machinev1.ObjectMeta{
				Labels: machineLabels,
			},
<<<<<<< HEAD
			ProviderSpec: machinev1.ProviderSpec{
				Value: providerSpec,
			},
		},
=======
			ProviderSpec: *providerSpec,
		},
	}

	return machine, nil
}

func stubCluster() *clusterv1.Cluster {
	return &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID,
			Namespace: defaultNamespace,
		},
	}
}

func stubUserDataSecret() *corev1.Secret {
	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      userDataSecretName,
			Namespace: defaultNamespace,
		},
		Data: map[string][]byte{
			userDataSecretKey: []byte(userDataBlob),
		},
	}
}

func stubAlicloudCredentialsSecret() *corev1.Secret {
	secret := utils.GenerateAlicloudCredentialsSecretFromEnv(alicloudCredentialsSecretName, defaultNamespace)
	aa, _ := json.Marshal(secret)
	fmt.Printf("=====%s=====", string(aa))
	return secret
}

func stubInstance(imageID, instanceID string) ecs.Instance {
	return ecs.Instance{
		ImageId:      imageID,
		InstanceId:   instanceID,
		Status:       "Running",
		CreationTime: time.Now().String(),
		PublicIpAddress: ecs.PublicIpAddressInDescribeInstances{
			IpAddress: []string{"1.1.1.1"},
		},
		InnerIpAddress: ecs.InnerIpAddressInDescribeInstances{
			IpAddress: []string{"1.1.1.1"},
		},
		Tags: ecs.TagsInDescribeInstances{
			Tag: []ecs.Tag{
				{
					TagKey:   "key1",
					TagValue: "value1",
				},
				{
					TagKey:   "key2",
					TagValue: "value2",
				},
			},
		},
		ZoneId: "cn-hangzhou-a",
		SecurityGroupIds: ecs.SecurityGroupIdsInDescribeInstances{
			SecurityGroupId: []string{"sg-abc"},
		},
	}
}

func stubCreateInstance(instanceID string) *ecs.CreateInstanceResponse {
	return &ecs.CreateInstanceResponse{
		InstanceId: instanceID,
>>>>>>> ebdd9bd0 (update test case)
	}
	return machine, nil
}

<<<<<<< HEAD
func stubRunInstancesRequest() *ecs.RunInstancesRequest {
	request := ecs.CreateRunInstancesRequest()
	request.Scheme = "https"
	request.RegionId = stubRegionID
	request.InstanceType = stubInstanceType
	request.ImageId = stubImageID
	request.VSwitchId = stubVSwitchID
	request.SecurityGroupId = stubSecurityGroupId
	request.Password = stubPassword
	request.MinAmount = requests.NewInteger(1)
	request.Amount = requests.NewInteger(1)

	request.SystemDiskCategory = stubSystemDiskCategory
	request.SystemDiskSize = strconv.Itoa(stubSystemDiskSize)

	return request
}

func stubRunInstancesResponse() *ecs.RunInstancesResponse {
	response := ecs.CreateRunInstancesResponse()
	response.InstanceIdSets = ecs.InstanceIdSets{
		InstanceIdSet: []string{stubInstanceID},
=======
func stubSecurityGroups(sgId string) *ecs.DescribeSecurityGroupsResponse {
	return &ecs.DescribeSecurityGroupsResponse{
		SecurityGroups: ecs.SecurityGroups{
			SecurityGroup: []ecs.SecurityGroup{
				{
					SecurityGroupId: sgId,
				},
			},
		},
	}
}

func stubImages(imageId string) *ecs.DescribeImagesResponse {
	return &ecs.DescribeImagesResponse{
		Images: ecs.Images{
			Image: []ecs.Image{
				{
					ImageId: imageId,
				},
			},
		},
>>>>>>> ebdd9bd0 (update test case)
	}

<<<<<<< HEAD
	return response
=======
func stubDescribeInstances(imageID, instanceID string) *ecs.DescribeInstancesResponse {
	return &ecs.DescribeInstancesResponse{
		Instances: ecs.InstancesInDescribeInstances{
			Instance: []ecs.Instance{
				{
					ImageId:    imageID,
					InstanceId: instanceID,
				},
			},
		},
	}
>>>>>>> ebdd9bd0 (update test case)
}
