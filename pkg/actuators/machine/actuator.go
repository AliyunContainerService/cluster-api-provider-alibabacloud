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

	alibabacloudClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"

	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machineapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	controllerclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	createEventAction = "Create"
	updateEventAction = "Update"
	deleteEventAction = "Delete"
	noEventAction     = ""

	userDataSecretKey = "userData"
)

// Actuator is responsible for performing machine reconciliation.
type Actuator struct {
	client        controllerclient.Client
	eventRecorder record.EventRecorder

	alibabacloudClientBuilder alibabacloudClient.AlibabaCloudClientBuilderFunc
	configManagedClient       controllerclient.Client

	reconcilerBuilder func(scope *machineScope) *Reconciler
}

// ActuatorParams holds parameter information for Actuator.
type ActuatorParams struct {
	Client        controllerclient.Client
	EventRecorder record.EventRecorder

	AlibabaCloudClientBuilder alibabacloudClient.AlibabaCloudClientBuilderFunc
	ConfigManagedClient       controllerclient.Client

	ReconcilerBuilder func(scope *machineScope) *Reconciler
}

// NewActuator returns an actuator.
func NewActuator(params ActuatorParams) *Actuator {
	return &Actuator{
		client:                    params.Client,
		eventRecorder:             params.EventRecorder,
		alibabacloudClientBuilder: params.AlibabaCloudClientBuilder,
		configManagedClient:       params.ConfigManagedClient,
		reconcilerBuilder:         params.ReconcilerBuilder,
	}
}

// Set corresponding event based on error. It also returns the original error
// for convenience, so callers can do "return handleMachineError(...)".
func (a *Actuator) handleMachineError(machine *machinev1.Machine, err *machineapierrors.MachineError, eventAction string) error {
	klog.Errorf("%v error: %v", machine.GetName(), err)
	if eventAction != noEventAction {
		a.eventRecorder.Eventf(machine, corev1.EventTypeWarning, "Failed"+eventAction, "%v: %v", err.Reason, err.Message)
	}
	klog.Errorf("Machine error: %v", err.Message)
	return err
}

// Create creates a machine and is invoked by the machine controller.
func (a *Actuator) Create(ctx context.Context, machine *machinev1.Machine) error {
	klog.Infof("%s actuator creating machine to namespace %s", machine.Name, machine.Namespace)

	scope, err := newMachineScope(machineScopeParams{
		Context:                   ctx,
		client:                    a.client,
		machine:                   machine,
		alibabacloudClientBuilder: a.alibabacloudClientBuilder,
		configManagedClient:       a.configManagedClient,
	})

	if err != nil {
		return a.handleMachineError(machine, machineapierrors.InvalidMachineConfiguration("failed to create machine %q scope: %v", machine.Name, err), createEventAction)
	}

	if err = a.reconcilerBuilder(scope).Create(context.Background()); err != nil {
		if err := scope.patchMachine(); err != nil {
			return err
		}
		return a.handleMachineError(machine, machineapierrors.InvalidMachineConfiguration("failed to reconcile machine %q: %v", machine.Name, err), createEventAction)
	}
	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, createEventAction, "Created Machine %v", machine.GetName())
	return scope.patchMachine()
}

// Update attempts to sync machine state with an existing instance.
func (a *Actuator) Update(ctx context.Context, machine *machinev1.Machine) error {
	klog.Infof("%s actuator updating machine to namespace %s", machine.Name, machine.Namespace)

	scope, err := newMachineScope(machineScopeParams{
		Context:                   ctx,
		client:                    a.client,
		machine:                   machine,
		alibabacloudClientBuilder: a.alibabacloudClientBuilder,
		configManagedClient:       a.configManagedClient,
	})

	if err != nil {
		return a.handleMachineError(machine, machineapierrors.InvalidMachineConfiguration("failed to create machine %q scope: %v", machine.Name, err), updateEventAction)
	}

	if err = a.reconcilerBuilder(scope).Update(context.Background()); err != nil {
		if err := scope.patchMachine(); err != nil {
			return err
		}
		return a.handleMachineError(machine, machineapierrors.InvalidMachineConfiguration("failed to reconcile machine %q: %v", machine.Name, err), updateEventAction)
	}

	previousResourceVersion := scope.machine.ResourceVersion

	if err := scope.patchMachine(); err != nil {
		return err
	}

	currentResourceVersion := scope.machine.ResourceVersion

	// Create event only if machine object was modified
	if previousResourceVersion != currentResourceVersion {
		a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, updateEventAction, "Updated Machine %v", machine.GetName())
	}

	return nil
}

// Delete deletes a machine and updates its finalizer
func (a *Actuator) Delete(ctx context.Context, machine *machinev1.Machine) error {
	klog.Infof("%s actuator deleting machine from namespace %s", machine.Name, machine.Namespace)

	scope, err := newMachineScope(machineScopeParams{
		Context:                   ctx,
		client:                    a.client,
		machine:                   machine,
		alibabacloudClientBuilder: a.alibabacloudClientBuilder,
		configManagedClient:       a.configManagedClient,
	})

	if err != nil {
		return a.handleMachineError(machine, machineapierrors.DeleteMachine("failed to create machine %q scope: %v", machine.Name, err), deleteEventAction)
	}

	if err = a.reconcilerBuilder(scope).Delete(context.Background()); err != nil {
		if err := scope.patchMachine(); err != nil {
			return err
		}
		return a.handleMachineError(machine, machineapierrors.InvalidMachineConfiguration("failed to reconcile machine %q: %v", machine.Name, err), deleteEventAction)
	}

	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Deleted", "Deleted machine %q", machine.Name)
	return scope.patchMachine()
}

// Exists determines if the given machine currently exists.
// A machine which is not terminated is considered as existing.
func (a *Actuator) Exists(ctx context.Context, machine *machinev1.Machine) (bool, error) {
	klog.Infof("%s: actuator checking if machine exists in namespace %s", machine.Name, machine.Namespace)

	scope, err := newMachineScope(machineScopeParams{
		Context:                   ctx,
		client:                    a.client,
		machine:                   machine,
		alibabacloudClientBuilder: a.alibabacloudClientBuilder,
		configManagedClient:       a.configManagedClient,
	})

	if err != nil {
		return false, a.handleMachineError(machine, machineapierrors.InvalidMachineConfiguration("failed to create machine %q scope: %v", machine.Name, err), createEventAction)
	}

	isExists, err := a.reconcilerBuilder(scope).Exists(context.Background())
	if err != nil {
		klog.Errorf("failed to check machine %s exists: %v", machine.Name, err)
	}

	return isExists, err
}
