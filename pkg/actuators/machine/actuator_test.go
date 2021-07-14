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
	"errors"
	"fmt"
	"testing"

	machineapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"

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
