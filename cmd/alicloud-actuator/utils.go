package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"k8s.io/client-go/tools/record"
	"os/exec"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/ghodss/yaml"

	machineactuator "github.com/AliyunContainerService/cluster-api-provider-alicloud/pkg/actuators/machine"
	"github.com/AliyunContainerService/cluster-api-provider-alicloud/pkg/apis/alicloudprovider/v1alpha1"
	alicloudclient "github.com/AliyunContainerService/cluster-api-provider-alicloud/pkg/client"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type manifestParams struct {
	ClusterID string
}

func readMachineManifest(manifestParams *manifestParams, manifestLoc string) (*machinev1.Machine, error) {
	machine := &machinev1.Machine{}
	manifestBytes, err := ioutil.ReadFile(manifestLoc)
	if err != nil {
		return nil, fmt.Errorf("unable to read %v: %v", manifestLoc, err)
	}

	t, err := template.New("machineuserdata").Parse(string(manifestBytes))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, *manifestParams)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(buf.Bytes(), &machine); err != nil {
		return nil, fmt.Errorf("unable to unmarshal %v: %v", manifestLoc, err)
	}

	return machine, nil
}

func readClusterResources(manifestParams *manifestParams, machineLoc, alicloudCredentialSecretLoc, userDataLoc string) (*machinev1.Machine, *apiv1.Secret, *apiv1.Secret, error) {
	machine, err := readMachineManifest(manifestParams, machineLoc)
	if err != nil {
		return nil, nil, nil, err
	}

	var aliCloudCredentialsSecret *apiv1.Secret
	if alicloudCredentialSecretLoc != "" {
		aliCloudCredentialsSecret = &apiv1.Secret{}
		alicloudBytes, err := ioutil.ReadFile(alicloudCredentialSecretLoc)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("alicloud credentials manifest %q: %v", alicloudCredentialSecretLoc, err)
		}

		if err = yaml.Unmarshal(alicloudBytes, &aliCloudCredentialsSecret); err != nil {
			return nil, nil, nil, fmt.Errorf("alicloud credentials manifest %q: %v", alicloudCredentialSecretLoc, err)
		}
	}

	var userDataSecret *apiv1.Secret
	if userDataLoc != "" {
		userDataSecret = &apiv1.Secret{}
		userdataBytes, err := ioutil.ReadFile(userDataLoc)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("user data manifest %q: %v", userDataLoc, err)
		}

		if err = yaml.Unmarshal(userdataBytes, &userDataSecret); err != nil {
			return nil, nil, nil, fmt.Errorf("user data manifest %q: %v", userDataLoc, err)
		}
	}

	return machine, aliCloudCredentialsSecret, userDataSecret, nil
}

// CreateActuator creates actuator with fake clientsets
func createActuator(machine *machinev1.Machine, aliCloudCredentials, userData *apiv1.Secret) (*machineactuator.Actuator, error) {
	objList := []runtime.Object{machine}
	if aliCloudCredentials != nil {
		objList = append(objList, aliCloudCredentials)
	}
	if userData != nil {
		objList = append(objList, userData)
	}
	fakeClient := fake.NewFakeClient(objList...)

	codec, err := v1alpha1.NewCodec()
	if err != nil {
		return nil, err
	}

	params := machineactuator.ActuatorParams{
		Client:                fakeClient,
		AliCloudClientBuilder: alicloudclient.NewClient,
		Codec:                 codec,
		// use empty recorder dropping any event recorded
		EventRecorder: &record.FakeRecorder{},
	}

	actuator, err := machineactuator.NewActuator(params)
	if err != nil {
		return nil, err
	}
	return actuator, nil
}

func cmdRun(binaryPath string, args ...string) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)
	return cmd.CombinedOutput()
}
