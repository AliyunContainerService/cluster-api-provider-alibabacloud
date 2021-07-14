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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client/mock"
	"github.com/golang/mock/gomock"
	"sigs.k8s.io/controller-runtime/pkg/client"

	machineapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"

	. "github.com/onsi/gomega"
	"k8s.io/client-go/tools/record"

	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	// Add types to scheme
	machinev1.AddToScheme(scheme.Scheme)
	configv1.AddToScheme(scheme.Scheme)
}

var (
	k8sClient     client.Client
	eventRecorder record.EventRecorder
)

// Init
func TestMain(m *testing.M) {
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "config", "crds"),
			filepath.Join("..", "..", "..", "vendor", "github.com", "openshift", "api", "config", "v1"),
		},
	}

	configv1.AddToScheme(scheme.Scheme)

	cfg, err := testEnv.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := testEnv.Stop(); err != nil {
			log.Fatal(err)
		}
	}()

	mgr, err := manager.New(cfg, manager.Options{
		Scheme:             scheme.Scheme,
		MetricsBindAddress: "0",
	})
	if err != nil {
		log.Fatal(err)
	}

	mgrCtx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := mgr.Start(mgrCtx); err != nil {
			log.Fatal(err)
		}
	}()
	defer cancel()

	k8sClient = mgr.GetClient()
	eventRecorder = mgr.GetEventRecorderFor("alibabacloud-controller")

	code := m.Run()
	os.Exit(code)
}

func Test_Client(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAlibabaCloudClient := mock.NewMockClient(mockCtrl)

	mockAlibabaCloudClient.EXPECT().RunInstances(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().RunInstances(stubRunInstancesRequest()).Return(stubRunInstancesResponse(), nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().CreateInstance(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeInstances(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DeleteInstances(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().StartInstance(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().RebootInstance(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().StopInstance(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().StartInstances(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().RebootInstances(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().StopInstances(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DeleteInstance(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().AttachInstanceRAMRole(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DetachInstanceRAMRole(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeInstanceStatus(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ReActivateInstances(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeUserData(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeInstanceTypes(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ModifyInstanceAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ModifyInstanceMetadataOptions(gomock.Any()).Return(nil, nil).AnyTimes()

	mockAlibabaCloudClient.EXPECT().TagResources(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ListTagResources(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().UntagResources(gomock.Any()).Return(nil, nil).AnyTimes()

	mockAlibabaCloudClient.EXPECT().AllocatePublicIPAddress(gomock.Any()).Return(nil, nil).AnyTimes()

	mockAlibabaCloudClient.EXPECT().CreateDisk(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().AttachDisk(gomock.Any()).Return(nil, nil).AnyTimes()

	mockAlibabaCloudClient.EXPECT().DescribeDisks(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ModifyDiskChargeType(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ModifyDiskAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ModifyDiskSpec(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ReplaceSystemDisk(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ReInitDisk(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ResetDisk(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ResizeDisk(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DetachDisk(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DeleteDisk(gomock.Any()).Return(nil, nil).AnyTimes()

	mockAlibabaCloudClient.EXPECT().DescribeRegions(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeZones(gomock.Any()).Return(nil, nil).AnyTimes()

	mockAlibabaCloudClient.EXPECT().DescribeImages(gomock.Any()).Return(nil, nil).AnyTimes()

	mockAlibabaCloudClient.EXPECT().CreateSecurityGroup(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().AuthorizeSecurityGroup(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().AuthorizeSecurityGroupEgress(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().RevokeSecurityGroup(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().RevokeSecurityGroupEgress(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().JoinSecurityGroup(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().LeaveSecurityGroup(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeSecurityGroupAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeSecurityGroups(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeSecurityGroupReferences(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ModifySecurityGroupAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ModifySecurityGroupEgressRule(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ModifySecurityGroupPolicy(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ModifySecurityGroupRule(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DeleteSecurityGroup(gomock.Any()).Return(nil, nil).AnyTimes()

	mockAlibabaCloudClient.EXPECT().CreateVpc(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DeleteVpc(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeVpcs(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().CreateVSwitch(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DeleteVSwitch(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeVSwitches(gomock.Any()).Return(nil, nil).AnyTimes()

	mockAlibabaCloudClient.EXPECT().CreateNatGateway(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeNatGateways(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DeleteNatGateway(gomock.Any()).Return(nil, nil).AnyTimes()

	mockAlibabaCloudClient.EXPECT().AllocateEipAddress(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().AssociateEipAddress(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ModifyEipAddressAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeEipAddresses(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().UnassociateEipAddress(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ReleaseEipAddress(gomock.Any()).Return(nil, nil).AnyTimes()

	mockAlibabaCloudClient.EXPECT().CreateLoadBalancer(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DeleteLoadBalancer(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeLoadBalancers(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().CreateLoadBalancerTCPListener(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().SetLoadBalancerTCPListenerAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeLoadBalancerTCPListenerAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().CreateLoadBalancerUDPListener(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().SetLoadBalancerUDPListenerAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeLoadBalancerUDPListenerAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().CreateLoadBalancerHTTPListener(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().SetLoadBalancerHTTPListenerAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeLoadBalancerHTTPListenerAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().CreateLoadBalancerHTTPSListener(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().SetLoadBalancerHTTPSListenerAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeLoadBalancerHTTPSListenerAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().StartLoadBalancerListener(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().StopLoadBalancerListener(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DeleteLoadBalancerListener(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeLoadBalancerListeners(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().AddBackendServers(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().RemoveBackendServers(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().SetBackendServers(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeHealthStatus(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().CreateVServerGroup(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().SetVServerGroupAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().AddVServerGroupBackendServers(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().RemoveVServerGroupBackendServers(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().ModifyVServerGroupBackendServers(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DeleteVServerGroup(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeVServerGroups(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAlibabaCloudClient.EXPECT().DescribeVServerGroupAttribute(gomock.Any()).Return(nil, nil).AnyTimes()
}

func Test_Machine(t *testing.T) {
	ctx := context.TODO()
	gs := NewWithT(t)

	machine, err := stubWorkerMachine()
	gs.Expect(err).ToNot(HaveOccurred())
	gs.Expect(stubMachine).ToNot(BeNil())

	// Create the machine
	gs.Expect(k8sClient.Create(ctx, machine)).To(Succeed())
	defer func() {
		gs.Expect(k8sClient.Delete(ctx, machine)).To(Succeed())
	}()
}

func Test_HandleMachineErrors(t *testing.T) {
	masterMachine, err := stubMasterMachine()
	if err != nil {
		t.Fatal(err)
	}

	configs := make([]map[string]interface{}, 0)
	//create
	configs = append(configs, map[string]interface{}{
		"Name":        "Create event for create action",
		"EventAction": createEventAction,
		"Error":       machineapierrors.InvalidMachineConfiguration("failed to create machine %q scope: %v", masterMachine.Name, errors.New("failed to get machine config")),
		"Event":       fmt.Sprintf("Warning FailedCreate InvalidConfiguration: failed to create machine \"alibabacloud-actuator-testing-master-machine\" scope: %v", errors.New("failed to get machine config")),
	})

	configs = append(configs, map[string]interface{}{
		"Name":        "Create event for create action",
		"EventAction": createEventAction,
		"Error":       machineapierrors.InvalidMachineConfiguration("failed to reconcile machine %q: %v", masterMachine.Name, errors.New("failed to set machine cloud provider specifics")),
		"Event":       fmt.Sprintf("Warning FailedCreate InvalidConfiguration: failed to reconcile machine \"alibabacloud-actuator-testing-master-machine\": %v", errors.New("failed to set machine cloud provider specifics")),
	})

	for _, config := range configs {
		eventsChannel := make(chan string, 1)

		params := ActuatorParams{
			// use fake recorder and store an event into one item long buffer for subsequent check
			EventRecorder: &record.FakeRecorder{
				Events: eventsChannel,
			},
		}

		actuator := NewActuator(params)

		actuator.handleMachineError(masterMachine, config["Error"].(*machineapierrors.MachineError), config["EventAction"].(string))
		select {
		case event := <-eventsChannel:
			if event != config["Event"] {
				t.Errorf("Expected %q event, got %q", config["Event"], event)
			} else {
				t.Logf("ok")
			}
		}
	}
}
