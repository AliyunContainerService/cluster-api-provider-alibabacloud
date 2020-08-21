package machine

import (
	"context"
	"fmt"

	alibabacloudproviderv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1beta1"
	aliClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machineapierros "github.com/openshift/machine-api-operator/pkg/controller/machine"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	userDataSecretKey = "userData"
)

// dhcpDomainKeyName is a variable so we can reference it in unit tests.
// var dhcpDomainKeyName = "domain-name"

// machineScopeParams defines the input parameters used to create a new MachineScope.
type machineScopeParams struct {
	context.Context

	aliCloudClientBuilder aliClient.AliCloudClientBuilderFuncType
	// api server controller runtime client
	client runtimeclient.Client
	// machine resource
	machine *machinev1.Machine
}

type machineScope struct {
	context.Context

	// client for interacting with AliCloud
	aliClient aliClient.Client
	// api server controller runtime client
	client runtimeclient.Client
	// machine resource
	machine            *machinev1.Machine
	machineToBePatched runtimeclient.Patch
	providerSpec       *alibabacloudproviderv1.AlibabaCloudMachineProviderConfig
	providerStatus     *alibabacloudproviderv1.AlibabaCloudMachineProviderStatus
}

func newMachineScope(params machineScopeParams) (*machineScope, error) {
	providerSpec, err := alibabacloudproviderv1.ProviderSpecFromRawExtension(params.machine.Spec.ProviderSpec.Value)
	if err != nil {
		return nil, machineapierros.InvalidMachineConfiguration("failed to get machine config: %v", err)
	}

	providerStatus, err := alibabacloudproviderv1.ProviderStatusFromRawExtension(params.machine.Status.ProviderStatus)
	if err != nil {
		return nil, machineapierros.InvalidMachineConfiguration("failed to get machine provider status: %v", err.Error())
	}

	credentialsSecretName := ""
	if providerSpec.CredentialsSecret != nil {
		credentialsSecretName = providerSpec.CredentialsSecret.Name
	}
	fmt.Printf("secret: %s, namespaces: %s, region: %s", credentialsSecretName, params.machine.Namespace, providerSpec.Placement.Region) //Remove when Pre Release
	aliClient, err := params.aliCloudClientBuilder(params.client, credentialsSecretName, params.machine.Namespace, providerSpec.Placement.Region)
	if err != nil {
		return nil, machineapierros.InvalidMachineConfiguration("failed to create alicloud client: %v", err.Error())
	}

	return &machineScope{
		Context:            params.Context,
		aliClient:          aliClient,
		client:             params.client,
		machine:            params.machine,
		machineToBePatched: runtimeclient.MergeFrom(params.machine.DeepCopy()),
		providerSpec:       providerSpec,
		providerStatus:     providerStatus,
	}, nil
}

// Patch patches the machine spec and machine status after reconciling.
func (s *machineScope) patchMachine() error {
	klog.V(3).Infof("%v: patching machine", s.machine.GetName())

	providerStatus, err := alibabacloudproviderv1.RawExtensionFromProviderStatus(s.providerStatus)
	if err != nil {
		return machineapierros.InvalidMachineConfiguration("failed to get machine provider status: %v", err.Error())
	}
	s.machine.Status.ProviderStatus = providerStatus

	statusCopy := *s.machine.Status.DeepCopy()

	// patch machine
	if err := s.client.Patch(context.Background(), s.machine, s.machineToBePatched); err != nil {
		klog.Errorf("Failed to patch machine %q: %v", s.machine.GetName(), err)
		return err
	}

	s.machine.Status = statusCopy

	// patch status
	if err := s.client.Status().Patch(context.Background(), s.machine, s.machineToBePatched); err != nil {
		klog.Errorf("Failed to patch machine status %q: %v", s.machine.GetName(), err)
		return err
	}

	return nil
}

// getUserData fetches the user-data from the secret referenced in the Machine's
// provider spec, if one is set.
func (s *machineScope) getUserData() ([]byte, error) {
	if s.providerSpec == nil || s.providerSpec.UserDataSecret == nil {
		return nil, nil
	}

	userDataSecret := &corev1.Secret{}

	objKey := runtimeclient.ObjectKey{
		Namespace: s.machine.Namespace,
		Name:      s.providerSpec.UserDataSecret.Name,
	}

	if err := s.client.Get(s.Context, objKey, userDataSecret); err != nil {
		return nil, err
	}

	userData, exists := userDataSecret.Data[userDataSecretKey]
	if !exists {
		return nil, fmt.Errorf("secret %s missing %s key", objKey, userDataSecretKey)
	}

	return userData, nil
}

func (s *machineScope) setProviderStatus(instance *ecs.Instance, condition alibabacloudproviderv1.AlibabaCloudMachineProviderCondition) error {
	klog.Infof("%s: Updating status", s.machine.Name)

	networkAddresses := []corev1.NodeAddress{}

	// TODO: remove 139-141, no need to clean instance id ands state
	// Instance may have existed but been deleted outside our control, clear it's status if so:
	if instance == nil {
		s.providerStatus.InstanceID = nil
		s.providerStatus.InstanceState = nil
	} else {
		s.providerStatus.InstanceID = String(instance.InstanceId)
		s.providerStatus.InstanceState = String(instance.Status)

		domainNames, err := s.getCustomDomainFromDHCP(instance.VpcAttributes.VpcId)

		if err != nil {
			return err
		}

		addresses, err := extractNodeAddresses(instance, domainNames)
		if err != nil {
			klog.Errorf("%s: Error extracting instance IP addresses: %v", s.machine.Name, err)
			return err
		}

		networkAddresses = append(networkAddresses, addresses...)
	}
	klog.Infof("%s: finished calculating AWS status", s.machine.Name)

	s.machine.Status.Addresses = networkAddresses
	s.providerStatus.Conditions = setAliCloudMachineProviderCondition(condition, s.providerStatus.Conditions)

	return nil
}

func (s *machineScope) getCustomDomainFromDHCP(vpcID string) ([]string, error) {
	// vpc, err := s.awsClient.DescribeVpcs(&ec2.DescribeVpcsInput{
	// 	VpcIds: []*string{vpcID},
	// })
	// if err != nil {
	// 	klog.Errorf("%s: error describing vpc: %v", s.machine.Name, err)
	// 	return nil, err
	// }

	// if len(vpc.Vpcs) == 0 || vpc.Vpcs[0] == nil || vpc.Vpcs[0].DhcpOptionsId == nil {
	// 	return nil, nil
	// }

	// dhcp, err := s.awsClient.DescribeDHCPOptions(&ec2.DescribeDhcpOptionsInput{
	// 	DhcpOptionsIds: []*string{vpc.Vpcs[0].DhcpOptionsId},
	// })
	// if err != nil {
	// 	klog.Errorf("%s: error describing dhcp: %v", s.machine.Name, err)
	// 	return nil, err
	// }

	// if dhcp == nil || len(dhcp.DhcpOptions) == 0 || dhcp.DhcpOptions[0] == nil {
	// 	return nil, nil
	// }

	// for _, i := range dhcp.DhcpOptions[0].DhcpConfigurations {
	// 	if i.Key != nil && *i.Key == dhcpDomainKeyName && len(i.Values) > 0 && i.Values[0] != nil && i.Values[0].Value != nil {
	// 		return strings.Split(*i.Values[0].Value, " "), nil
	// 	}
	// }
	return nil, nil
}
