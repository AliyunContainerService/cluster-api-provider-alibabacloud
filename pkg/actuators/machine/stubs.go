package machine

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"time"

	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	providerconfigv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1alpha1"
	"github.com/AliyunContainerService/cluster-api-provider-alibabacloud/test/utils"
	clusterv1 "github.com/openshift/cluster-api/pkg/apis/cluster/v1alpha1"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	machinecontroller "github.com/openshift/cluster-api/pkg/controller/machine"
)

const (
	defaultNamespace              = "default"
	defaultAvailabilityZone       = "cn-hangzhou-a"
	region                        = "cn-hangzhou"
	alicloudCredentialsSecretName = "alicloud-credentials-secret"
	userDataSecretName            = "alicloud-actuator-user-data-secret"

	keyName   = "alicloud-actuator-key-name"
	clusterID = "alicloud-actuator-cluster"
)

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

func stubProviderConfig() *providerconfigv1.AlibabaCloudMachineProviderConfig {
	return &providerconfigv1.AlibabaCloudMachineProviderConfig{
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
	}
}

func stubMachine() (*machinev1.Machine, error) {
	machinePc := stubProviderConfig()

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		return nil, fmt.Errorf("failed creating codec: %v", err)
	}
	providerSpec, err := codec.EncodeProviderSpec(machinePc)
	if err != nil {
		return nil, fmt.Errorf("codec.EncodeProviderSpec failed: %v", err)
	}

	machine := &machinev1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "alicloud-actuator-testing-machine",
			Namespace: defaultNamespace,
			Labels: map[string]string{
				providerconfigv1.ClusterIDLabel: clusterID,
			},
			Annotations: map[string]string{
				// skip node draining since it's not mocked
				machinecontroller.ExcludeNodeDrainingAnnotation: "",
			},
		},

		Spec: machinev1.MachineSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"node-role.kubernetes.io/master": "",
					"node-role.kubernetes.io/infra":  "",
				},
			},
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
	}
}

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
	}
}

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
}
