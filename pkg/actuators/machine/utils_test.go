package machine
<<<<<<< HEAD
<<<<<<< HEAD
=======

import (
	"reflect"
	"testing"

	providerconfigv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1alpha1"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	// Add types to scheme
	machinev1.AddToScheme(scheme.Scheme)
}

func TestProviderConfigFromMachine(t *testing.T) {

	providerConfig := &providerconfigv1.AlibabaCloudMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "alicloudproviderconfig.openshift.io/v1alpha1",
			Kind:       "AlibabaCloudMachineProviderConfig",
		},
		ImageId: "centos_7_06_64_20G_alibase_20190619.vhd",
		CredentialsSecret: &corev1.LocalObjectReference{
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

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Error(err)
	}
	encodedProviderSpec, err := codec.EncodeProviderSpec(providerConfig)
	if err != nil {
		t.Error(err)
	}

	testCases := []struct {
		machine *machinev1.Machine
	}{
		{
			machine: &machinev1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "configFromSpecProviderConfigValue",
					Namespace: "",
					Labels: map[string]string{
						"foo": "a",
					},
				},
				TypeMeta: metav1.TypeMeta{
					Kind: "Machine",
				},
				Spec: machinev1.MachineSpec{
					ProviderSpec: *encodedProviderSpec,
				},
			},
		},
	}

	for _, tc := range testCases {
		decodedProviderConfig, err := providerConfigFromMachine(tc.machine, codec)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(decodedProviderConfig, providerConfig) {
			t.Errorf("Test case %s. Expected: %v, got: %v", tc.machine.Name, providerConfig, decodedProviderConfig)
		}
	}
}
>>>>>>> ebdd9bd0 (update test case)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
