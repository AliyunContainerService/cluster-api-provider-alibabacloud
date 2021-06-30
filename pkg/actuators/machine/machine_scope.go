/*
Copyright 2021 The Kubernetes Authors.

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

package machine

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

<<<<<<< HEAD
	alibabacloudproviderv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudprovider/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"
=======
	"k8s.io/klog/v2"

	alibabacloudproviderv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudprovider/v1beta1"
	corev1 "k8s.io/api/core/v1"
>>>>>>> e879a141 (alibabacloud machine-api provider)

	v1beta1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudprovider/v1beta1"
	alibabacloudClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machineapierros "github.com/openshift/machine-api-operator/pkg/controller/machine"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineScope defines a scope defined around a machine and its cluster.
type machineScope struct {
	context.Context

	// client for interacting with alibabacloud
	alibabacloudClient alibabacloudClient.Client
	// api server controller runtime client
	client runtimeclient.Client
	// machine resource
	machine            *machinev1.Machine
	machineToBePatched runtimeclient.Patch
	providerSpec       *v1beta1.AlibabaCloudMachineProviderConfig
	providerStatus     *v1beta1.AlibabaCloudMachineProviderStatus
}

// machineScopeParams defines the input parameters used to create a new MachineScope.
type machineScopeParams struct {
	context.Context

<<<<<<< HEAD
	alibabacloudClientBuilder alibabacloudClient.AlibabaCloudClientBuilderFunc
=======
	alibabacloudClientBuilder alibabacloudClient.AlibabaCloudClientBuilderFuncType
>>>>>>> e879a141 (alibabacloud machine-api provider)
	// api server controller runtime client
	client runtimeclient.Client
	// machine resource
	machine *machinev1.Machine
	// api server controller runtime client for the openshift-config-managed namespace
	configManagedClient runtimeclient.Client
}

// newMachineScope init machineScope instance
func newMachineScope(params machineScopeParams) (*machineScope, error) {
	providerSpec, err := v1beta1.ProviderSpecFromRawExtension(params.machine.Spec.ProviderSpec.Value)
	if err != nil {
		return nil, machineapierros.InvalidMachineConfiguration("failed to get machine config: %v", err)
	}

	providerStatus, err := v1beta1.ProviderStatusFromRawExtension(params.machine.Status.ProviderStatus)
	if err != nil {
		return nil, machineapierros.InvalidMachineConfiguration("failed to get machine provider status: %v", err.Error())
	}

	credentialsSecretName := ""
	if providerSpec.CredentialsSecret != nil {
		credentialsSecretName = providerSpec.CredentialsSecret.Name
	}

	aliClient, err := params.alibabacloudClientBuilder(params.client, credentialsSecretName, params.machine.Namespace, providerSpec.RegionID, params.configManagedClient)
	if err != nil {
		return nil, machineapierros.InvalidMachineConfiguration("failed to create alibabacloud client: %v", err)
	}

	return &machineScope{
		Context:            params.Context,
		alibabacloudClient: aliClient,
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
		return machineapierros.InvalidMachineConfiguration("failed to get machine provider status: %v", err)
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
func (s *machineScope) getUserData() (string, error) {
	if s.providerSpec == nil || s.providerSpec.UserDataSecret == nil {
		return "", nil
	}

	userDataSecret := &corev1.Secret{}

	objKey := runtimeclient.ObjectKey{
		Namespace: s.machine.Namespace,
		Name:      s.providerSpec.UserDataSecret.Name,
	}

	if err := s.client.Get(s.Context, objKey, userDataSecret); err != nil {
		return "", fmt.Errorf("error getting user data secret %s in namespace %s: %w", s.providerSpec.UserDataSecret.Name, s.providerSpec.GetNamespace(), err)
	}

	userData, exists := userDataSecret.Data[userDataSecretKey]
	if !exists {
		return "", fmt.Errorf("secret %v/%v does not have userData field set. thus, no user data applied when creating an instance", s.providerSpec.GetNamespace(), s.providerSpec.UserDataSecret.Name)
	}

	return base64.StdEncoding.EncodeToString(userData), nil
}

func (s *machineScope) setProviderStatus(instance *ecs.Instance, condition alibabacloudproviderv1.AlibabaCloudMachineProviderCondition) error {
	klog.Infof("%s: Updating status", s.machine.Name)

<<<<<<< HEAD
	// assign value to providerStatus
=======
	networkAddresses := []corev1.NodeAddress{}

>>>>>>> e879a141 (alibabacloud machine-api provider)
	if instance == nil {
		s.providerStatus.InstanceID = nil
		s.providerStatus.InstanceState = nil
	} else {
		s.providerStatus.InstanceID = &instance.InstanceId
		s.providerStatus.InstanceState = &instance.Status
<<<<<<< HEAD
	}

	networkAddresses, err := s.getNetworkAddress(instance)
	if err != nil {
		klog.Errorf("%s : Error get Network Address :%v", s.machine.Name, err)
		return nil
	}
	s.machine.Status.Addresses = networkAddresses
	s.providerStatus.Conditions = setMachineProviderCondition(condition, s.providerStatus.Conditions)

	return nil
}

func (s *machineScope) getNetworkAddress(instance *ecs.Instance) ([]corev1.NodeAddress, error) {
	klog.Infof("%s : Setting network address", s.machine.Name)

	networkAddresses := make([]corev1.NodeAddress, 0)

	addresses, err := extractNodeAddressesFromInstance(instance)
	if err != nil {
		klog.Errorf("%s: Error extracting instance IP addresses: %v", s.machine.Name, err)
		return nil, err
	}

	networkAddresses = append(networkAddresses, addresses...)

	klog.Infof("%s: finished calculating alibabacloud status", s.machine.Name)

	return networkAddresses, nil
}

// extractNodeAddressesFromInstance maps the instance information from ECS to an array of NodeAddresses
func extractNodeAddressesFromInstance(instance *ecs.Instance) ([]corev1.NodeAddress, error) {

	if instance == nil {
		return nil, fmt.Errorf("the ecs instance is nil")
=======

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
	klog.Infof("%s: finished calculating alibabacloud status", s.machine.Name)

	s.machine.Status.Addresses = networkAddresses
	s.providerStatus.Conditions = setAlibabaCloudMachineProviderCondition(condition, s.providerStatus.Conditions)

	return nil
}

// extractNodeAddresses maps the instance information from ECS to an array of NodeAddresses
func extractNodeAddresses(instance *ecs.Instance, domainNames []string) ([]corev1.NodeAddress, error) {
	// Not clear if the order matters here, but we might as well indicate a sensible preference order

	if instance == nil {
		return nil, fmt.Errorf("nil instance passed to extractNodeAddresses")
>>>>>>> e879a141 (alibabacloud machine-api provider)
	}

	addresses := make([]corev1.NodeAddress, 0)

	// handle internal network interfaces
	for _, networkInterface := range instance.NetworkInterfaces.NetworkInterface {
<<<<<<< HEAD
=======
		// skip network interfaces that are not currently in use
		//if networkInterface.= ecs.NetworkInterfaceStatusInUse {
		//	continue
		//}

		// Treating IPv6 addresses as type NodeInternalIP to match what the KNI
		// patch to the alibabacloud cloud-provider code is doing:
		//
>>>>>>> e879a141 (alibabacloud machine-api provider)
		// https://github.com/openshift-kni/origin/commit/7db21c1e26a344e25ae1b825d4f21e7bef5c3650
		for _, ipv6Address := range networkInterface.Ipv6Sets.Ipv6Set {
			if addr := ipv6Address.Ipv6Address; addr != "" {
				ip := net.ParseIP(addr)
				if ip == nil {
					return nil, fmt.Errorf("ECS instance had invalid IPv6 address: %s (%q)", instance.InstanceId, addr)
				}
				addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalIP, Address: ip.String()})
			}
		}

		for _, internalIP := range networkInterface.PrivateIpSets.PrivateIpSet {
			if ipAddress := internalIP.PrivateIpAddress; ipAddress != "" {
				ip := net.ParseIP(ipAddress)
				if ip == nil {
					return nil, fmt.Errorf("ECS instance had invalid private address: %s (%q)", instance.InstanceId, ipAddress)
				}
				addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalIP, Address: ip.String()})
			}
		}
	}

	//// TODO: Other IP addresses (multiple ips)?
	for _, publicIPAddress := range instance.PublicIpAddress.IpAddress {
		if publicIPAddress != "" {
			ip := net.ParseIP(publicIPAddress)
			if ip == nil {
				return nil, fmt.Errorf("ECS instance had invalid public address: %s (%s)", instance.InstanceId, publicIPAddress)
			}
			addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeExternalIP, Address: ip.String()})
		}
	}

	return addresses, nil
}
<<<<<<< HEAD
=======

func (s *machineScope) getCustomDomainFromDHCP(vpcID string) ([]string, error) {
	//vpc, err := s.awsClient.DescribeVpcs(&ec2.DescribeVpcsInput{
	//	VpcIds: []*string{vpcID},
	//})
	//if err != nil {
	//	klog.Errorf("%s: error describing vpc: %v", s.machine.Name, err)
	//	return nil, err
	//}
	//
	//if len(vpc.Vpcs) == 0 || vpc.Vpcs[0] == nil || vpc.Vpcs[0].DhcpOptionsId == nil {
	//	return nil, nil
	//}
	//
	//dhcp, err := s.awsClient.DescribeDHCPOptions(&ec2.DescribeDhcpOptionsInput{
	//	DhcpOptionsIds: []*string{vpc.Vpcs[0].DhcpOptionsId},
	//})
	//if err != nil {
	//	klog.Errorf("%s: error describing dhcp: %v", s.machine.Name, err)
	//	return nil, err
	//}
	//
	//if dhcp == nil || len(dhcp.DhcpOptions) == 0 || dhcp.DhcpOptions[0] == nil {
	//	return nil, nil
	//}
	//
	//for _, i := range dhcp.DhcpOptions[0].DhcpConfigurations {
	//	if i.Key != nil && *i.Key == dhcpDomainKeyName && len(i.Values) > 0 && i.Values[0] != nil && i.Values[0].Value != nil {
	//		return strings.Split(*i.Values[0].Value, " "), nil
	//	}
	//}
	return nil, nil
}
>>>>>>> e879a141 (alibabacloud machine-api provider)
