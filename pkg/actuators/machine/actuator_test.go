/*
Copyright 2021 The Kubernetes Authors.
<<<<<<< HEAD

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

<<<<<<< HEAD
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
=======
	providerconfigv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1alpha1"
	aliCloudClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	mockaliCloud "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client/mock"
	clusterv1 "github.com/openshift/cluster-api/pkg/apis/cluster/v1alpha1"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
>>>>>>> 8dbd34ff (update project name)

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

<<<<<<< HEAD
// Init
func TestMain(m *testing.M) {
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "config", "crds"),
			filepath.Join("..", "..", "..", "vendor", "github.com", "openshift", "api", "config", "v1"),
		},
	}

	configv1.AddToScheme(scheme.Scheme)
=======
func TestMachineEvents(t *testing.T) {
	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Fatalf("unable to build codec: %v", err)
	}

	machine, err := stubMachine()
	if err != nil {
		t.Fatal(err)
	}

	cluster := stubCluster()
	aliCloudCredentialsSecret := stubAlicloudCredentialsSecret()
	userDataSecret := stubUserDataSecret()

	machineInvalidProviderConfig := machine.DeepCopy()
	machineInvalidProviderConfig.Spec.ProviderSpec.Value = nil

	workerMachine := machine.DeepCopy()
	workerMachine.Spec.Labels["node-role.kubernetes.io/worker"] = ""

	cases := []struct {
		name                     string
		machine                  *machinev1.Machine
		error                    string
		operation                func(actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine)
		event                    string
		describeInstancesReponse *ecs.DescribeInstancesResponse
		describeInstancesErr     error
		runInstancesErr          error
		deleteInstancesErr       error
		lbErr                    error
		regInstancesWithLbErr    error
	}{
		{
			name:    "Create machine event failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Warning FailedCreate InvalidConfiguration",
		},
		{
			name:    "Create machine event failed (error creating alicloud service)",
			machine: machine,
			error:   aliCloudServiceError,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Warning FailedCreate CreateError",
		},
		{
			name:                     "Create machine event failed (error launching instance)",
			machine:                  machine,
			runInstancesErr:          fmt.Errorf("error"),
			describeInstancesReponse: &ecs.DescribeInstancesResponse{},
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Warning FailedCreate CreateError",
		},
		{
			name:    "Create machine event succeed",
			machine: machine,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Normal Created Created Machine alicloud-actuator-testing-machine",
		},
		{
			name:    "Create worker machine event succeed",
			machine: workerMachine,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Normal Created Created Machine alicloud-actuator-testing-machine",
		},
		{
			name:    "Delete machine event failed",
			machine: machineInvalidProviderConfig,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.DeleteMachine(cluster, machine)
			},
			event: "Warning FailedDelete InvalidConfiguration",
		},
		{
			name:    "Delete machine event succeed",
			machine: machine,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.DeleteMachine(cluster, machine)
			},
			event: "Normal Deleted Deleted machine alicloud-actuator-testing-machine",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			mockAlicloudClient := mockaliCloud.NewMockClient(mockCtrl)

			eventsChannel := make(chan string, 1)

			params := ActuatorParams{
				Client: fake.NewFakeClient(tc.machine, aliCloudCredentialsSecret, userDataSecret),
				AliCloudClientBuilder: func(client client.Client, secretName, namespace, region string) (aliCloudClient.Client, error) {
					if tc.error == aliCloudServiceError {
						return nil, fmt.Errorf(aliCloudServiceError)
					}
					return mockAlicloudClient, nil
				},
				Codec: codec,
				// use fake recorder and store an event into one item long buffer for subsequent check
				EventRecorder: &record.FakeRecorder{
					Events: eventsChannel,
				},
			}

			mockAlicloudClient.EXPECT().StartInstance(gomock.Any()).Return(&ecs.StartInstanceResponse{}, tc.runInstancesErr).AnyTimes()

			if tc.describeInstancesReponse == nil {
				mockAlicloudClient.EXPECT().DescribeInstances(gomock.Any()).Return(stubDescribeInstances("centos_7_06_64_20G_alibase_20190619.vhd", "i-bp1bsuzspukvo4t56a4f"), tc.describeInstancesErr).AnyTimes()
			} else {
				mockAlicloudClient.EXPECT().DescribeInstances(gomock.Any()).Return(&ecs.DescribeInstancesResponse{}, tc.describeInstancesErr).AnyTimes()
			}

			mockAlicloudClient.EXPECT().DescribeSecurityGroups(gomock.Any()).Return(stubSecurityGroups("sg-bp1iccjoxddumf300okm"), nil)
			mockAlicloudClient.EXPECT().DescribeImages(gomock.Any()).Return(stubImages("centos_7_06_64_20G_alibase_20190619.vhd"), nil)

			mockAlicloudClient.EXPECT().CreateInstance(gomock.Any()).Return(stubCreateInstance("i-bp1bsuzspukvo4t56a4f"), nil)
			mockAlicloudClient.EXPECT().WaitForInstance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			mockAlicloudClient.EXPECT().WaitForInstance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			mockAlicloudClient.EXPECT().DeleteInstance(gomock.Any()).Return(&ecs.DeleteInstanceResponse{}, tc.deleteInstancesErr).AnyTimes()
			mockAlicloudClient.EXPECT().DeleteInstance(gomock.Any()).Return(&ecs.DeleteInstanceResponse{}, nil)

			actuator, err := NewActuator(params)
			if err != nil {
				t.Fatalf("Could not create AliCloud machine actuator: %v", err)
			}

			tc.operation(actuator, cluster, tc.machine)
			select {
			case event := <-eventsChannel:
				if event != tc.event {
					t.Errorf("Expected %q event, got %q", tc.event, event)
				}
			default:
				t.Errorf("Expected %q event, got none", tc.event)
			}
		})
	}
}
>>>>>>> ebdd9bd0 (update test case)

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

<<<<<<< HEAD
	mgrCtx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := mgr.Start(mgrCtx); err != nil {
			log.Fatal(err)
=======
	getMachineStatus := func(objectClient client.Client, machine *machinev1.Machine) (*providerconfigv1.AlibabaCloudMachineProviderStatus, error) {
		// Get updated machine object from the cluster client
		key := types.NamespacedName{
			Namespace: machine.Namespace,
			Name:      machine.Name,
		}
		updatedMachine := machinev1.Machine{}
		err := objectClient.Get(context.Background(), client.ObjectKey(key), &updatedMachine)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve machine: %v", err)
>>>>>>> c7e62b88 (fix testcase)
		}
	}()
	defer cancel()

<<<<<<< HEAD
<<<<<<< HEAD
	k8sClient = mgr.GetClient()
	eventRecorder = mgr.GetEventRecorderFor("alibabacloud-controller")
=======
=======
		machineStatus := &providerconfigv1.AlibabaCloudMachineProviderStatus{}
		if err := codec.DecodeProviderStatus(updatedMachine.Status.ProviderStatus, machineStatus); err != nil {
			return nil, fmt.Errorf("error decoding machine provider status: %v", err)
		}
		return machineStatus, nil
	}

>>>>>>> c7e62b88 (fix testcase)
	machineInvalidProviderConfig := machine.DeepCopy()
	machineInvalidProviderConfig.Spec.ProviderSpec.Value = nil

	machineNoClusterID := machine.DeepCopy()
	delete(machineNoClusterID.Labels, providerconfigv1.ClusterIDLabel)

	pendingInstance := stubInstance("centos_7_06_64_20G_alibase_20190711.vhd", "i-bp1bsuzspukvo4t56a4f")
	pendingInstance.Status = "Starting"

	cases := []struct {
		name                      string
		machine                   *machinev1.Machine
		error                     string
		operation                 func(client client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine)
		describeInstancesResponse *ecs.DescribeInstancesResponse
		runInstancesErr           error
		describeInstancesErr      error
		deleteInstancesErr        error
		lbErr                     error
	}{
		{
			name:    "Create machine with success",
			machine: machine,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				createErr := actuator.Create(context.TODO(), cluster, machine)
				assert.NoError(t, createErr)

				machineStatus, err := getMachineStatus(objectClient, machine)
				if err != nil {
					t.Fatalf("Unable to get machine status: %v", err)
				}

				assert.Equal(t, machineStatus.Conditions[0].Reason, MachineCreationSucceeded)

				// Get the machine
				if exists, err := actuator.Exists(context.TODO(), cluster, machine); err != nil || !exists {
					t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
				} else {
					t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
				}

				// Update a machine
				if err := actuator.Update(context.TODO(), cluster, machine); err != nil {
					t.Errorf("Unable to create instance for machine: %v", err)
				}

				// Get the machine
				if exists, err := actuator.Exists(context.TODO(), cluster, machine); err != nil || !exists {
					t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
				} else {
					t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
				}

				// Delete a machine
				if err := actuator.Delete(context.TODO(), cluster, machine); err != nil {
					t.Errorf("Unable to delete instance for machine: %v", err)
				}
			},
		},
		{
			name:            "Create machine with failure",
			machine:         machine,
			runInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				createErr := actuator.Create(context.TODO(), cluster, machine)
				assert.Error(t, createErr)

				machineStatus, err := getMachineStatus(objectClient, machine)
				if err != nil {
					t.Fatalf("Unable to get machine status: %v", err)
				}

				assert.Equal(t, machineStatus.Conditions[0].Reason, MachineCreationFailed)
			},
		},
		{
			name:    "Update machine with success",
			machine: machine,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name:    "Update machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name:  "Update machine failed (error creating alicloud service)",
			error: aliCloudServiceError,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name:                 "Update machine failed (error getting running instances)",
			describeInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name: "Update machine failed (no running instances)",
			describeInstancesResponse: &ecs.DescribeInstancesResponse{
				Instances: ecs.InstancesInDescribeInstances{
					Instance: []ecs.Instance{},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name: "Update machine succeeds (two running instances)",
			describeInstancesResponse: &ecs.DescribeInstancesResponse{
				Instances: ecs.InstancesInDescribeInstances{
					Instance: []ecs.Instance{
						{InstanceId: "i-bp1bsuzspukvo4t56a4f", ImageId: "centos_7_06_64_20G_alibase_20190711.vhd"},
						{InstanceId: "i-bp1bsuzspukvo4t56a4g", ImageId: "centos_7_06_64_20G_alibase_20190711.vhd"},
					},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name: "Update machine status fails (instance pending)",
			describeInstancesResponse: &ecs.DescribeInstancesResponse{
				Instances: ecs.InstancesInDescribeInstances{
					Instance: []ecs.Instance{
						pendingInstance,
					},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name: "Update machine failed (two running instances, error terminating one)",
			describeInstancesResponse: &ecs.DescribeInstancesResponse{
				Instances: ecs.InstancesInDescribeInstances{
					Instance: []ecs.Instance{
						{InstanceId: "i-bp1bsuzspukvo4t56a4f", ImageId: "centos_7_06_64_20G_alibase_20190711.vhd"},
						{InstanceId: "i-bp1bsuzspukvo4t56a4g", ImageId: "centos_7_06_64_20G_alibase_20190711.vhd"},
					},
				},
			},
			deleteInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name:    "Update machine with failure (cluster ID missing)",
			machine: machineNoClusterID,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name:                 "Describe machine fails (error getting running instance)",
			describeInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Describe(cluster, machine)
			},
		},
		{
			name:    "Describe machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Describe(cluster, machine)
			},
		},
		{
			name:  "Describe machine failed (error creating alicloud service)",
			error: aliCloudServiceError,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Describe(cluster, machine)
			},
		},
		{
			name: "Describe machine fails (no running instance)",
			describeInstancesResponse: &ecs.DescribeInstancesResponse{
				Instances: ecs.InstancesInDescribeInstances{
					Instance: []ecs.Instance{},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Describe(cluster, machine)
			},
		},
		{
			name: "Describe machine succeeds",
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Describe(cluster, machine)
			},
		},
		{
			name:    "Exists machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Exists(context.TODO(), cluster, machine)
			},
		},
		{
			name: "Exists machine fails (no running instance)",
			describeInstancesResponse: &ecs.DescribeInstancesResponse{
				Instances: ecs.InstancesInDescribeInstances{
					Instance: []ecs.Instance{},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Exists(context.TODO(), cluster, machine)
			},
		},
		{
			name:    "Delete machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), cluster, machine)
			},
		},
		{
			name:  "Delete machine failed (error creating alicloud service)",
			error: aliCloudServiceError,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), cluster, machine)
			},
		},
		{
			name:                 "Delete machine failed (error getting running instances)",
			describeInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), cluster, machine)
			},
		},
		{
			name: "Delete machine failed (no running instances)",
>>>>>>> ebdd9bd0 (update test case)

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

<<<<<<< HEAD
		actuator := NewActuator(params)
=======
			mockAliCloudClient.EXPECT().CreateInstance(gomock.Any()).Return(stubCreateInstance("i-bp1bsuzspukvo4t56a4f"), tc.runInstancesErr)
			mockAliCloudClient.EXPECT().WaitForInstance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			mockAliCloudClient.EXPECT().WaitForInstance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			mockAliCloudClient.EXPECT().DeleteInstance(gomock.Any()).Return(&ecs.DeleteInstanceResponse{}, tc.deleteInstancesErr).AnyTimes()
>>>>>>> ebdd9bd0 (update test case)

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
=======

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
>>>>>>> e879a141 (alibabacloud machine-api provider)
