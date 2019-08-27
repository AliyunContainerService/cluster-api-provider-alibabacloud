package machine

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/record"
	"strings"
	"testing"

	providerconfigv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1alpha1"
	aliCloudClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	mockaliCloud "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client/mock"
	clusterv1 "github.com/openshift/cluster-api/pkg/apis/cluster/v1alpha1"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func init() {
	// Add types to scheme
	machinev1.AddToScheme(scheme.Scheme)
}

const (
	noError              = ""
	aliCloudServiceError = "error creating alicloud service"
	launchInstanceError  = "error launching instance"
)

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

func TestActuator(t *testing.T) {
	machine, err := stubMachine()
	if err != nil {
		t.Fatal(err)
	}

	cluster := stubCluster()
	alicloudCredentialsSecret := stubAlicloudCredentialsSecret()
	userDataSecret := stubUserDataSecret()

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Fatalf("unable to build codec: %v", err)
	}

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
		}

		machineStatus := &providerconfigv1.AlibabaCloudMachineProviderStatus{}
		if err := codec.DecodeProviderStatus(updatedMachine.Status.ProviderStatus, machineStatus); err != nil {
			return nil, fmt.Errorf("error decoding machine provider status: %v", err)
		}
		return machineStatus, nil
	}

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

			describeInstancesResponse: &ecs.DescribeInstancesResponse{
				Instances: ecs.InstancesInDescribeInstances{
					Instance: []ecs.Instance{},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), cluster, machine)
			},
		},
		{
			name: "Delete machine failed (error terminating instances)",

			deleteInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *clusterv1.Cluster, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), cluster, machine)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fakeClient := fake.NewFakeClient(machine, alicloudCredentialsSecret, userDataSecret)
			mockCtrl := gomock.NewController(t)
			mockAliCloudClient := mockaliCloud.NewMockClient(mockCtrl)

			params := ActuatorParams{
				Client: fakeClient,
				AliCloudClientBuilder: func(client client.Client, secretName, namespace, region string) (aliCloudClient.Client, error) {
					if tc.error == aliCloudServiceError {
						return nil, fmt.Errorf(aliCloudServiceError)
					}
					return mockAliCloudClient, nil
				},
				Codec: codec,
				// use empty recorder dropping any event recorded
				EventRecorder: &record.FakeRecorder{},
			}

			actuator, err := NewActuator(params)
			if err != nil {
				t.Fatalf("Could not create Alicloud machine actuator: %v", err)
			}

			mockAliCloudClient.EXPECT().StartInstance(gomock.Any()).Return(&ecs.StartInstanceResponse{}, tc.runInstancesErr).AnyTimes()

			if tc.describeInstancesResponse == nil {
				mockAliCloudClient.EXPECT().DescribeInstances(gomock.Any()).Return(stubDescribeInstances("centos_7_06_64_20G_alibase_20190619.vhd", "i-bp1bsuzspukvo4t56a4f"), tc.describeInstancesErr).AnyTimes()
			} else {
				mockAliCloudClient.EXPECT().DescribeInstances(gomock.Any()).Return(&ecs.DescribeInstancesResponse{}, tc.describeInstancesErr).AnyTimes()
			}

			mockAliCloudClient.EXPECT().DescribeSecurityGroups(gomock.Any()).Return(stubSecurityGroups("sg-bp1iccjoxddumf300okm"), nil)
			mockAliCloudClient.EXPECT().DescribeImages(gomock.Any()).Return(stubImages("centos_7_06_64_20G_alibase_20190619.vhd"), nil)

			mockAliCloudClient.EXPECT().CreateInstance(gomock.Any()).Return(stubCreateInstance("i-bp1bsuzspukvo4t56a4f"), tc.runInstancesErr)
			mockAliCloudClient.EXPECT().WaitForInstance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			mockAliCloudClient.EXPECT().WaitForInstance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			mockAliCloudClient.EXPECT().DeleteInstance(gomock.Any()).Return(&ecs.DeleteInstanceResponse{}, tc.deleteInstancesErr).AnyTimes()

			if tc.machine == nil {
				tc.operation(fakeClient, actuator, cluster, machine)
			} else {
				tc.operation(fakeClient, actuator, cluster, tc.machine)
			}
		})
	}
}
