package integration

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	//"github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	machineactuator "github.com/AliyunContainerService/cluster-api-provider-alicloud/pkg/actuators/machine"
	alicloudclient "github.com/AliyunContainerService/cluster-api-provider-alicloud/pkg/client"
	clusterv1 "github.com/openshift/cluster-api/pkg/apis/cluster/v1alpha1"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"

	"github.com/ghodss/yaml"
)

const (
	controllerLogName = "alicloudMachine"
	defaultLogLevel   = "info"

	defaultNamespace              = "default"
	alicloudCredentialsSecretName = "alicloud-credentials-secret"
	userDataSecretName            = "alicloud-actuator-user-data-secret"

	clusterID = "tb-asg-35"
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

func testMachineAPIResources(clusterID string) (*machinev1.Machine, *clusterv1.Cluster, *apiv1.Secret, *apiv1.Secret, error) {
	machine := &machinev1.Machine{}

	bytes, err := ioutil.ReadFile("/Users/xiaobing/Documents/workspace/aone/go/src/github.com/AliyunContainerService/cluster-api-provider-alicloud/examples/machines.yaml")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if err = yaml.Unmarshal(bytes, &machine); err != nil {
		return nil, nil, nil, nil, err
	}

	alicloudCredentialsSecret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      alicloudCredentialsSecretName,
			Namespace: defaultNamespace,
		},
		Data: map[string][]byte{
			alicloudclient.AliCloudAccessKeyId:     []byte(os.Getenv("ALICLOUD_ACCESS_KEY_ID")),
			alicloudclient.AliCloudAccessKeySecret: []byte(os.Getenv("ALICLOUD_ACCESS_KEY_SECRET")),
		},
	}

	userDataSecret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      userDataSecretName,
			Namespace: defaultNamespace,
		},
		Data: map[string][]byte{
			"user-data": []byte(userDataBlob),
		},
	}

	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID,
			Namespace: defaultNamespace,
		},
	}

	return machine, cluster, alicloudCredentialsSecret, userDataSecret, nil
}

func TestCreateAndDeleteMachine(t *testing.T) {

	if os.Getenv("GOPATH") == "" {
		t.Fatalf("GOPATH not set")
	}

	machine, cluster, alicloudCredentialsSecret, userDataSecret, err := testMachineAPIResources(clusterID)
	if err != nil {
		t.Fatal(err)
	}

	fakeClient := fake.NewFakeClient(machine, alicloudCredentialsSecret, userDataSecret)

	params := machineactuator.ActuatorParams{
		Client:                fakeClient,
		AliCloudClientBuilder: alicloudclient.NewClient,
	}

	actuator, err := machineactuator.NewActuator(params)
	if err != nil {
		t.Fatalf("Could not create Openstack machine actuator: %v", err)
	}

	// Create a machine
	if err := actuator.Create(context.TODO(), cluster, machine); err != nil {
		t.Errorf("Unable to create instance for machine: %v", err)
	}

	// Get the machine
	if exists, err := actuator.Exists(context.TODO(), cluster, machine); err != nil || !exists {
		t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
	} else {
		t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
	}

	// TODO(jchaloup): Wait until the machine is ready

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
}
