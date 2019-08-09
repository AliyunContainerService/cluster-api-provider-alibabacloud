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

<<<<<<< HEAD
// NewActuator returns an actuator.
func NewActuator(params ActuatorParams) *Actuator {
	return &Actuator{
		client:                    params.Client,
		eventRecorder:             params.EventRecorder,
		alibabacloudClientBuilder: params.AlibabaCloudClientBuilder,
		configManagedClient:       params.ConfigManagedClient,
		reconcilerBuilder:         params.ReconcilerBuilder,
=======
// NewActuator returns a new AliCloud Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {
	actuator := &Actuator{
		client:                params.Client,
		config:                params.Config,
		aliCloudClientBuilder: params.AliCloudClientBuilder,
		codec:                 params.Codec,
		eventRecorder:         params.EventRecorder,
>>>>>>> ebdd9bd0 (update test case)
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

<<<<<<< HEAD
=======
	if err := a.updateMachineStatus(machine, alicloudStatus, networkAddresses); err != nil {
		return err
	}

	// If machine state is still pending, we will return an error to keep the controllers
	// attempting to update status until it hits a more permanent state. This will ensure
	// we get a public IP populated more quickly.
	if alicloudStatus.InstanceStatus != nil && *alicloudStatus.InstanceStatus == "Starting" {
		glog.Infof("%s: Instance state still pending, returning an error to requeue", machine.Name)
		return &clustererror.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
	}
	return nil
}

// CreateMachine starts a new ECS instance as described by the cluster and machine resources
func (a *Actuator) CreateMachine(cluster *clusterv1.Cluster, machine *machinev1.Machine) (*ecs.Instance, error) {
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return nil, a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error decoding MachineProviderConfig: %v", err), createEventAction)
	}

	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	alicloudClient, err := a.aliCloudClientBuilder(a.client, credentialsSecretName, machine.Namespace, machineProviderConfig.RegionId)
>>>>>>> ebdd9bd0 (update test case)
	if err != nil {
		return a.handleMachineError(machine, machineapierrors.InvalidMachineConfiguration("failed to create machine %q scope: %v", machine.Name, err), createEventAction)
	}

<<<<<<< HEAD
	if err = a.reconcilerBuilder(scope).Create(context.Background()); err != nil {
		if err := scope.patchMachine(); err != nil {
			return err
=======
	userData := []byte{}
	if machineProviderConfig.UserDataSecret != nil {
		var userDataSecret corev1.Secret
		err := a.client.Get(context.Background(), client.ObjectKey{Namespace: machine.Namespace, Name: machineProviderConfig.UserDataSecret.Name}, &userDataSecret)
		if err != nil {
			return nil, a.handleMachineError(machine, apierrors.CreateMachine("error getting user data secret %s: %v", machineProviderConfig.UserDataSecret.Name, err), createEventAction)
		}
		if data, exists := userDataSecret.Data[userDataSecretKey]; exists {
			userData = data
		} else {
			glog.Warningf("%s: Secret %v/%v does not have %q field set. Thus, no user data applied when creating an instance.", machine.Name, machine.Namespace, machineProviderConfig.UserDataSecret.Name, userDataSecretKey)
>>>>>>> ebdd9bd0 (update test case)
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
<<<<<<< HEAD
=======
	return nil
}

// DeleteMachine deletes an ECS instance
func (a *Actuator) DeleteMachine(cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error decoding MachineProviderConfig: %v", err), deleteEventAction)
	}

	region := machineProviderConfig.RegionId
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	aliCloudClient, err := a.aliCloudClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		errMsg := fmt.Errorf("%s: error getting ECS client: %v", machine.Name, err)
		glog.Error(errMsg)
		return errMsg
	}
>>>>>>> ebdd9bd0 (update test case)

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

<<<<<<< HEAD
=======
	region := machineProviderConfig.RegionId
	glog.Infof("%s: obtaining ECS client for region", machine.Name)
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	aliCloudClient, err := a.aliCloudClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
>>>>>>> ebdd9bd0 (update test case)
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

<<<<<<< HEAD
	isExists, err := a.reconcilerBuilder(scope).Exists(context.Background())
=======
	region := machineProviderConfig.RegionId
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	aliCloudClient, err := a.aliCloudClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
>>>>>>> ebdd9bd0 (update test case)
	if err != nil {
		klog.Errorf("failed to check machine %s exists: %v", machine.Name, err)
	}

	return isExists, err
}
