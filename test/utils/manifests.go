package utils

import (
	"fmt"
	machinev1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"os"

	providerconfigv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1alpha1"
	alicloudclient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateAlicloudCredentialsSecretFromEnv generates secret with AliCloud credentials
func GenerateAlicloudCredentialsSecretFromEnv(secretName, namespace string) *apiv1.Secret {
	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			alicloudclient.AliCloudAccessKeyId:     []byte(os.Getenv("ALICLOUD_ACCESS_KEY_ID")),
			alicloudclient.AliCloudAccessKeySecret: []byte(os.Getenv("ALICLOUD_ACCESS_KEY_SECRET")),
		},
	}
}

func testingAlicloudMachineProviderSpec(alicloudCredentialsSecretName string, clusterID string) *providerconfigv1.AlibabaCloudMachineProviderConfig {
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

func TestingMachineProviderSpec(alicloudCredentialsSecretName string, clusterID string) (machinev1beta1.ProviderSpec, error) {
	machinePc := testingAlicloudMachineProviderSpec(alicloudCredentialsSecretName, clusterID)
	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		return machinev1beta1.ProviderSpec{}, fmt.Errorf("failed creating codec: %v", err)
	}
	config, err := codec.EncodeProviderSpec(machinePc)
	if err != nil {
		return machinev1beta1.ProviderSpec{}, fmt.Errorf("EncodeToProviderConfig failed: %v", err)
	}
	return *config, nil
}
