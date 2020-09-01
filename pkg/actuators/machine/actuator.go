/*
Copyright 2018 The Kubernetes Authors.
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
	"fmt"

	aliClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	"github.com/golang/glog"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ecsInstanceIDNotFoundCode = "InvalidInstanceID.NotFound"

	// MachineCreationSucceeded indicates success for machine creation
	MachineCreationSucceeded = "MachineCreationSucceeded"

	// MachineCreationFailed indicates that machine creation failed
	MachineCreationFailed = "MachineCreationFailed"
)

const (
	scopeFailFmt      = "%s: failed to create scope for machine: %w"
	reconcilerFailFmt = "%s: reconciler failed to %s machine: %w"
	createEventAction = "Create"
	updateEventAction = "Update"
	deleteEventAction = "Delete"
	noEventAction     = ""
)

// Actuator is responsible for performing machine reconciliation.
type Actuator struct {
	aliCloudClientBuilder aliClient.AliCloudClientBuilderFuncType
	client                client.Client
	config                *rest.Config
	
	eventRecorder record.EventRecorder
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	Client                client.Client
	EventRecorder         record.EventRecorder
	AliCloudClientBuilder aliClient.AliCloudClientBuilderFuncType
}

// NewActuator returns a new AliCloud Actuator
func NewActuator(params ActuatorParams) *Actuator {
	return &Actuator{
		client:                params.Client,
		eventRecorder:         params.EventRecorder,
		aliCloudClientBuilder: params.AliCloudClientBuilder,
	}
}

// Set corresponding event based on error. It also returns the original error
// for convenience, so callers can do "return handleMachineError(...)".
func (a *Actuator) handleMachineError(machine *machinev1.Machine, err error, eventAction string) error {
	if eventAction != noEventAction {
		a.eventRecorder.Eventf(machine, corev1.EventTypeWarning, "Failed"+eventAction, "%v", err)
	}

	glog.Errorf("%s: Machine error: %v", machine.Name, err)
	return err
}

//Create runs an ECS instance
func (a *Actuator) Create(ctx context.Context, machine *machinev1.Machine) error {
	klog.Infof("%s create machine", machine.Name)
	scope, err := newMachineScope(machineScopeParams{
		Context:               ctx,
		client:                a.client,
		machine:               machine,
		aliCloudClientBuilder: a.aliCloudClientBuilder,
	})
	if err != nil {
		fmtErr := fmt.Errorf(scopeFailFmt, machine.GetName(), err)
		return a.handleMachineError(machine, fmtErr, createEventAction)
	}
	if err := newReconciler(scope).create(); err != nil {
		if err := scope.patchMachine(); err != nil {
			return err
		}
		fmtErr := fmt.Errorf(reconcilerFailFmt, machine.GetName(), createEventAction, err)
		return a.handleMachineError(machine, fmtErr, createEventAction)
	}
	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, createEventAction, "Created Machine %v", machine.GetName())
	return scope.patchMachine()
}

//Exists ...
func (a *Actuator) Exists(ctx context.Context, machine *machinev1.Machine) (bool, error) {
	klog.Infof("%s: actuator checking if machine exists", machine.GetName())
	scope, err := newMachineScope(machineScopeParams{
		Context:               ctx,
		client:                a.client,
		machine:               machine,
		aliCloudClientBuilder: a.aliCloudClientBuilder,
	})
	if err != nil {
		return false, fmt.Errorf(scopeFailFmt, machine.GetName(), err)
	}
	return newReconciler(scope).exists()
}

//Update ...
func (a *Actuator) Update(ctx context.Context, machine *machinev1.Machine) error {
	klog.Infof("%s: actuator updating machine", machine.GetName())
	scope, err := newMachineScope(machineScopeParams{
		Context:               ctx,
		client:                a.client,
		machine:               machine,
		aliCloudClientBuilder: a.aliCloudClientBuilder,
	})
	if err != nil {
		fmtErr := fmt.Errorf(scopeFailFmt, machine.GetName(), err)
		return a.handleMachineError(machine, fmtErr, updateEventAction)
	}
	if err := newReconciler(scope).update(); err != nil {
		// Update machine and machine status in case it was modified
		if err := scope.patchMachine(); err != nil {
			return err
		}
		fmtErr := fmt.Errorf(reconcilerFailFmt, machine.GetName(), updateEventAction, err)
		return a.handleMachineError(machine, fmtErr, updateEventAction)
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
	klog.Infof("%s: actuator deleting machine", machine.GetName())
	scope, err := newMachineScope(machineScopeParams{
		Context:               ctx,
		client:                a.client,
		machine:               machine,
		aliCloudClientBuilder: a.aliCloudClientBuilder,
	})
	if err != nil {
		fmtErr := fmt.Errorf(scopeFailFmt, machine.GetName(), err)
		return a.handleMachineError(machine, fmtErr, deleteEventAction)
	}
	if err := newReconciler(scope).delete(); err != nil {
		if err := scope.patchMachine(); err != nil {
			return err
		}
		fmtErr := fmt.Errorf(reconcilerFailFmt, machine.GetName(), deleteEventAction, err)
		return a.handleMachineError(machine, fmtErr, deleteEventAction)
	}
	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, deleteEventAction, "Deleted machine %v", machine.GetName())
	return scope.patchMachine()
}

// updateStatus calculates the new machine status, checks if anything has changed, and updates if so.
// func (a *Actuator) updateStatus(machine *machinev1.Machine, instance *ecs.Instance) error {

// 	glog.Infof("%s: Updating status", machine.Name)

// 	// Starting with a fresh status as we assume full control of it here.
// 	alicloudStatus := &providerconfigv1.AlibabaCloudMachineProviderStatus{}
// 	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, alicloudStatus); err != nil {
// 		glog.Errorf("%s: Error decoding machine provider status: %v", machine.Name, err)
// 		return err
// 	}

// 	// Save this, we need to check if it changed later.
// 	networkAddresses := []corev1.NodeAddress{}

// 	// Instance may have existed but been deleted outside our control, clear it's status if so:
// 	if instance == nil {
// 		alicloudStatus.InstanceID = nil
// 		alicloudStatus.InstanceStatus = nil
// 	} else {
// 		alicloudStatus.InstanceID = &instance.InstanceId
// 		alicloudStatus.InstanceStatus = &instance.Status
// 		if len(instance.PublicIpAddress.IpAddress) > 0 {
// 			networkAddresses = append(networkAddresses, corev1.NodeAddress{
// 				Type:    corev1.NodeExternalIP,
// 				Address: instance.PublicIpAddress.IpAddress[0],
// 			})
// 		}
// 		if len(instance.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
// 			networkAddresses = append(networkAddresses, corev1.NodeAddress{
// 				Type:    corev1.NodeInternalIP,
// 				Address: instance.VpcAttributes.PrivateIpAddress.IpAddress[0],
// 			})
// 		}
// 		networkAddresses = append(networkAddresses, corev1.NodeAddress{
// 			Type:    corev1.NodeInternalDNS,
// 			Address: strings.Join([]string{instance.RegionId, instance.InstanceId}, "."),
// 		})
// 	}
// 	glog.Infof("%s: finished calculating Alicloud status", machine.Name)

// 	alicloudStatus.Conditions = setAliCloudMachineProviderCondition(alicloudStatus.Conditions, providerconfigv1.AlibabaCloudMachineProviderCondition{
// 		Type:    providerconfigv1.MachineCreation,
// 		Status:  corev1.ConditionTrue,
// 		Reason:  MachineCreationSucceeded,
// 		Message: "machine successfully created",
// 	})

// 	if err := a.updateMachineStatus(machine, alicloudStatus, networkAddresses); err != nil {
// 		return err
// 	}

// 	// If machine state is still pending, we will return an error to keep the controllers
// 	// attempting to update status until it hits a more permanent state. This will ensure
// 	// we get a public IP populated more quickly.
// 	if alicloudStatus.InstanceStatus != nil && *alicloudStatus.InstanceStatus == "Starting" {
// 		glog.Infof("%s: Instance state still pending, returning an error to requeue", machine.Name)
// 		// return &clustererror.RequeueAfterError{RequeueAfter: requeueAfterSeconds * time.Second}
// 	}
// 	return nil
// }

// CreateMachine starts a new ECS instance as described by the cluster and machine resources
// func (a *Actuator) CreateMachine(machine *machinev1.Machine) (*ecs.Instance, error) {
// 	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
// 	if err != nil {
// 		return nil, a.handleMachineError(machine, fmt.Errorf("error decoding MachineProviderConfig: %v", err), createEventAction)
// 	}

// 	credentialsSecretName := ""
// 	if machineProviderConfig.CredentialsSecret != nil {
// 		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
// 	}
// 	alicloudClient, err := a.aliCloudClientBuilder(a.client, credentialsSecretName, machine.Namespace, machineProviderConfig.RegionId)
// 	if err != nil {
// 		glog.Errorf("%s: unable to obtain AliCloud client: %v", machine.Name, err)
// 		return nil, a.handleMachineError(machine, fmt.Errorf("error creating alicloud services: %v", err), createEventAction)
// 	}

// 	userData := []byte{}
// 	if machineProviderConfig.UserDataSecret != nil {
// 		var userDataSecret corev1.Secret
// 		err := a.client.Get(context.Background(), client.ObjectKey{Namespace: machine.Namespace, Name: machineProviderConfig.UserDataSecret.Name}, &userDataSecret)
// 		if err != nil {
// 			return nil, a.handleMachineError(machine, fmt.Errorf("error getting user data secret %s: %v", machineProviderConfig.UserDataSecret.Name, err), createEventAction)
// 		}
// 		if data, exists := userDataSecret.Data[userDataSecretKey]; exists {
// 			userData = data
// 		} else {
// 			glog.Warningf("%s: Secret %v/%v does not have %q field set. Thus, no user data applied when creating an instance.", machine.Name, machine.Namespace, machineProviderConfig.UserDataSecret.Name, userDataSecretKey)
// 		}
// 	}

// 	instance, err := createInstance(machine, machineProviderConfig, userData, alicloudClient)
// 	if err != nil {
// 		return nil, a.handleMachineError(machine, fmt.Errorf("error launching instance: %v", err), createEventAction)
// 	}

// 	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Created", "Created Machine %v", machine.Name)
// 	return instance, nil
// }

// updateMachineProviderConditions updates conditions set within machine provider status.
// func (a *Actuator) updateMachineProviderConditions(machine *machinev1.Machine, conditionType providerconfigv1.AliCloudMachineProviderConditionType, reason string, msg string) error {

// 	glog.Infof("%s: updating machine conditions", machine.Name)

// 	aliCloudStatus := &providerconfigv1.AlibabaCloudMachineProviderStatus{}
// 	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, aliCloudStatus); err != nil {
// 		glog.Errorf("%s: error decoding machine provider status: %v", machine.Name, err)
// 		return err
// 	}

// 	aliCloudStatus.Conditions = setAliCloudMachineProviderCondition(aliCloudStatus.Conditions, providerconfigv1.AlibabaCloudMachineProviderCondition{
// 		Type:    conditionType,
// 		Status:  corev1.ConditionTrue,
// 		Reason:  reason,
// 		Message: msg,
// 	})

// 	if err := a.updateMachineStatus(machine, aliCloudStatus, nil); err != nil {
// 		return err
// 	}

// 	return nil
// }

// //update status
// func (a *Actuator) updateMachineStatus(machine *machinev1.Machine, aliCloudStatus *providerconfigv1.AlibabaCloudMachineProviderStatus, networkAddresses []corev1.NodeAddress) error {
// 	alicloudStatusRaw, err := a.codec.EncodeProviderStatus(aliCloudStatus)
// 	if err != nil {
// 		glog.Errorf("%s: error encoding Alicloud provider status: %v", machine.Name, err)
// 		return err
// 	}

// 	machineCopy := machine.DeepCopy()
// 	machineCopy.Status.ProviderStatus = alicloudStatusRaw
// 	if networkAddresses != nil {
// 		machineCopy.Status.Addresses = networkAddresses
// 	}

// 	oldAlicloudStatus := &providerconfigv1.AlibabaCloudMachineProviderStatus{}
// 	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, oldAlicloudStatus); err != nil {
// 		glog.Errorf("%s: error updating machine status: %v", machine.Name, err)
// 		return err
// 	}

// 	if !equality.Semantic.DeepEqual(alicloudStatusRaw, oldAlicloudStatus) || !equality.Semantic.DeepEqual(machine.Status.Addresses, machineCopy.Status.Addresses) {
// 		glog.Infof("%s: machine status has changed, updating", machine.Name)
// 		time := metav1.Now()
// 		machineCopy.Status.LastUpdated = &time

// 		if err := a.client.Status().Update(context.Background(), machineCopy); err != nil {
// 			glog.Errorf("%s: error updating machine status: %v", machine.Name, err)
// 			return err
// 		}
// 	} else {
// 		glog.Infof("%s: status unchanged", machine.Name)
// 	}

// 	return nil
// }

// DeleteMachine deletes an ECS instance
// func (a *Actuator) DeleteMachine(machine *machinev1.Machine) error {
// 	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
// 	if err != nil {
// 		return a.handleMachineError(machine, fmt.Errorf("error decoding MachineProviderConfig: %v", err), deleteEventAction)
// 	}

// 	region := machineProviderConfig.RegionId
// 	credentialsSecretName := ""
// 	if machineProviderConfig.CredentialsSecret != nil {
// 		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
// 	}
// 	aliCloudClient, err := a.aliCloudClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
// 	if err != nil {
// 		errMsg := fmt.Errorf("%s: error getting ECS client: %v", machine.Name, err)
// 		glog.Error(errMsg)
// 		return errMsg
// 	}

// 	instances, err := getRunningInstances(machine, aliCloudClient, region)
// 	if err != nil {
// 		glog.Errorf("%s: error getting running instances: %v", machine.Name, err)
// 		return err
// 	}
// 	if len(instances) == 0 {
// 		glog.Warningf("%s: no instances found to delete for machine", machine.Name)
// 		return nil
// 	}

// 	err = deleteInstances(aliCloudClient, instances)
// 	if err != nil {
// 		return a.handleMachineError(machine, machinecontroller.DeleteMachine(err.Error()), noEventAction)
// 	}
// 	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Deleted", "Deleted machine %v", machine.Name)

// 	return nil
// }

//Describe ...
// func (a *Actuator) Describe(machine *machinev1.Machine) (*ecs.Instance, error) {
// 	glog.Infof("%s: Checking if machine exists", machine.Name)

// 	instances, err := a.getMachineInstances(machine)
// 	if err != nil {
// 		glog.Errorf("%s: Error getting running instances: %v", machine.Name, err)
// 		return nil, err
// 	}
// 	if len(instances) == 0 {
// 		glog.Infof("%s: Instance does not exist", machine.Name)
// 		return nil, nil
// 	}

// 	return instances[0], nil
// }

// func (a *Actuator) getMachineInstances(machine *machinev1.Machine) ([]*ecs.Instance, error) {
// 	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
// 	if err != nil {
// 		glog.Errorf("%s: Error decoding MachineProviderConfig: %v", machine.Name, err)
// 		return nil, err
// 	}

// 	region := machineProviderConfig.RegionId
// 	credentialsSecretName := ""
// 	if machineProviderConfig.CredentialsSecret != nil {
// 		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
// 	}
// 	aliCloudClient, err := a.aliCloudClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
// 	if err != nil {
// 		errMsg := fmt.Sprintf("%s: Error getting ECS client: %v", machine.Name, err)
// 		glog.Errorf(errMsg)
// 		return nil, fmt.Errorf(errMsg)
// 	}

// 	return getExistingInstances(machine, aliCloudClient, region)
// }
