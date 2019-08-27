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
	"time"

	providerconfigv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1alpha1"
	aliClient "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/golang/glog"
	clusterv1 "github.com/openshift/cluster-api/pkg/apis/cluster/v1alpha1"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	clustererror "github.com/openshift/cluster-api/pkg/controller/error"
	apierrors "github.com/openshift/cluster-api/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	userDataSecretKey         = "userData"
	ecsInstanceIDNotFoundCode = "InvalidInstanceID.NotFound"
	requeueAfterSeconds       = 20
	requeueAfterFatalSeconds  = 180

	// MachineCreationSucceeded indicates success for machine creation
	MachineCreationSucceeded = "MachineCreationSucceeded"

	// MachineCreationFailed indicates that machine creation failed
	MachineCreationFailed = "MachineCreationFailed"
)

const (
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

	codec         *providerconfigv1.AlicloudProviderConfigCodec
	eventRecorder record.EventRecorder
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	Client                client.Client
	Config                *rest.Config
	AliCloudClientBuilder aliClient.AliCloudClientBuilderFuncType
	Codec                 *providerconfigv1.AlicloudProviderConfigCodec
	EventRecorder         record.EventRecorder
}

// NewActuator returns a new AliCloud Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {
	actuator := &Actuator{
		client:                params.Client,
		config:                params.Config,
		aliCloudClientBuilder: params.AliCloudClientBuilder,
		codec:                 params.Codec,
		eventRecorder:         params.EventRecorder,
	}
	return actuator, nil
}

// Set corresponding event based on error. It also returns the original error
// for convenience, so callers can do "return handleMachineError(...)".
func (a *Actuator) handleMachineError(machine *machinev1.Machine, err *apierrors.MachineError, eventAction string) error {
	if eventAction != noEventAction {
		a.eventRecorder.Eventf(machine, corev1.EventTypeWarning, "Failed"+eventAction, "%v", err.Reason)
	}

	glog.Errorf("%s: Machine error: %v", machine.Name, err.Message)
	return err
}

//Create runs an ECS instance
func (a *Actuator) Create(context context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
	glog.Infof("%s create machine", machine.Name)
	instance, err := a.CreateMachine(cluster, machine)
	if err != nil {
		glog.Errorf("%s: error creating machine: %v", machine.Name, err)
		updateConditionError := a.updateMachineProviderConditions(machine, providerconfigv1.MachineCreation, MachineCreationFailed, err.Error())
		if updateConditionError != nil {
			glog.Errorf("%s: error updating machine conditions: %v", machine.Name, updateConditionError)
		}
		return err
	}

	return a.updateStatus(machine, instance)
}

// updateStatus calculates the new machine status, checks if anything has changed, and updates if so.
func (a *Actuator) updateStatus(machine *machinev1.Machine, instance *ecs.Instance) error {

	glog.Infof("%s: Updating status", machine.Name)

	// Starting with a fresh status as we assume full control of it here.
	alicloudStatus := &providerconfigv1.AlibabaCloudMachineProviderStatus{}
	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, alicloudStatus); err != nil {
		glog.Errorf("%s: Error decoding machine provider status: %v", machine.Name, err)
		return err
	}

	// Save this, we need to check if it changed later.
	networkAddresses := []corev1.NodeAddress{}

	// Instance may have existed but been deleted outside our control, clear it's status if so:
	if instance == nil {
		alicloudStatus.InstanceID = nil
		alicloudStatus.InstanceStatus = nil
	} else {
		alicloudStatus.InstanceID = &instance.InstanceId
		alicloudStatus.InstanceStatus = &instance.Status
		if len(instance.PublicIpAddress.IpAddress) > 0 {
			networkAddresses = append(networkAddresses, corev1.NodeAddress{
				Type:    corev1.NodeExternalIP,
				Address: instance.PublicIpAddress.IpAddress[0],
			})
		}
		if len(instance.InnerIpAddress.IpAddress) > 0 {
			networkAddresses = append(networkAddresses, corev1.NodeAddress{
				Type:    corev1.NodeInternalIP,
				Address: instance.InnerIpAddress.IpAddress[0],
			})
		}
	}
	glog.Infof("%s: finished calculating Alicloud status", machine.Name)

	alicloudStatus.Conditions = setAliCloudMachineProviderCondition(alicloudStatus.Conditions, providerconfigv1.AlibabaCloudMachineProviderCondition{
		Type:    providerconfigv1.MachineCreation,
		Status:  corev1.ConditionTrue,
		Reason:  MachineCreationSucceeded,
		Message: "machine successfully created",
	})

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
	if err != nil {
		glog.Errorf("%s: unable to obtain AliCloud client: %v", machine.Name, err)
		return nil, a.handleMachineError(machine, apierrors.CreateMachine("error creating alicloud services: %v", err), createEventAction)
	}

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
		}
	}

	instance, err := createInstance(machine, machineProviderConfig, userData, alicloudClient)
	if err != nil {
		return nil, a.handleMachineError(machine, apierrors.CreateMachine("error launching instance: %v", err), createEventAction)
	}

	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Created", "Created Machine %v", machine.Name)
	return instance, nil
}

// updateMachineProviderConditions updates conditions set within machine provider status.
func (a *Actuator) updateMachineProviderConditions(machine *machinev1.Machine, conditionType providerconfigv1.AliCloudMachineProviderConditionType, reason string, msg string) error {

	glog.Infof("%s: updating machine conditions", machine.Name)

	aliCloudStatus := &providerconfigv1.AlibabaCloudMachineProviderStatus{}
	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, aliCloudStatus); err != nil {
		glog.Errorf("%s: error decoding machine provider status: %v", machine.Name, err)
		return err
	}

	aliCloudStatus.Conditions = setAliCloudMachineProviderCondition(aliCloudStatus.Conditions, providerconfigv1.AlibabaCloudMachineProviderCondition{
		Type:    conditionType,
		Status:  corev1.ConditionTrue,
		Reason:  reason,
		Message: msg,
	})

	if err := a.updateMachineStatus(machine, aliCloudStatus, nil); err != nil {
		return err
	}

	return nil
}

//update status
func (a *Actuator) updateMachineStatus(machine *machinev1.Machine, aliCloudStatus *providerconfigv1.AlibabaCloudMachineProviderStatus, networkAddresses []corev1.NodeAddress) error {
	alicloudStatusRaw, err := a.codec.EncodeProviderStatus(aliCloudStatus)
	if err != nil {
		glog.Errorf("%s: error encoding Alicloud provider status: %v", machine.Name, err)
		return err
	}

	machineCopy := machine.DeepCopy()
	machineCopy.Status.ProviderStatus = alicloudStatusRaw
	if networkAddresses != nil {
		machineCopy.Status.Addresses = networkAddresses
	}

	oldAlicloudStatus := &providerconfigv1.AlibabaCloudMachineProviderStatus{}
	if err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, oldAlicloudStatus); err != nil {
		glog.Errorf("%s: error updating machine status: %v", machine.Name, err)
		return err
	}

	if !equality.Semantic.DeepEqual(alicloudStatusRaw, oldAlicloudStatus) || !equality.Semantic.DeepEqual(machine.Status.Addresses, machineCopy.Status.Addresses) {
		glog.Infof("%s: machine status has changed, updating", machine.Name)
		time := metav1.Now()
		machineCopy.Status.LastUpdated = &time

		if err := a.client.Status().Update(context.Background(), machineCopy); err != nil {
			glog.Errorf("%s: error updating machine status: %v", machine.Name, err)
			return err
		}
	} else {
		glog.Infof("%s: status unchanged", machine.Name)
	}

	return nil
}

// Delete deletes a machine and updates its finalizer
func (a *Actuator) Delete(context context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
	glog.Infof("%s: deleting machine", machine.Name)
	if err := a.DeleteMachine(cluster, machine); err != nil {
		glog.Errorf("%s: error deleting machine: %v", machine.Name, err)
		return err
	}
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

	instances, err := getRunningInstances(machine, aliCloudClient, region)
	if err != nil {
		glog.Errorf("%s: error getting running instances: %v", machine.Name, err)
		return err
	}
	if len(instances) == 0 {
		glog.Warningf("%s: no instances found to delete for machine", machine.Name)
		return nil
	}

	err = deleteInstances(aliCloudClient, instances)
	if err != nil {
		return a.handleMachineError(machine, apierrors.DeleteMachine(err.Error()), noEventAction)
	}
	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Deleted", "Deleted machine %v", machine.Name)

	return nil
}

//update
func (a *Actuator) Update(context context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) error {
	glog.Infof("%s: updating machine", machine.Name)

	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		return a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error decoding MachineProviderConfig: %v", err), updateEventAction)
	}

	region := machineProviderConfig.RegionId
	glog.Infof("%s: obtaining ECS client for region", machine.Name)
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
	// Get all instances not deleted.
	existingInstances, err := getExistingInstances(machine, aliCloudClient, region)
	if err != nil {
		glog.Errorf("%s: error getting existing instances: %v", machine.Name, err)
		return err
	}
	existingLen := len(existingInstances)
	glog.Infof("%s: found %d existing instances for machine", machine.Name, existingLen)

	// Parent controller should prevent this from ever happening by calling Exists and then Create,
	// but instance could be deleted between the two calls.
	if existingLen == 0 {
		glog.Warningf("%s: attempted to update machine but no instances found", machine.Name)

		_ = a.handleMachineError(machine, apierrors.UpdateMachine("no instance found, reason unknown"), updateEventAction)

		// Update status to clear out machine details.
		if err := a.updateStatus(machine, nil); err != nil {
			return err
		}
		// This is an unrecoverable error condition.  We should delay to
		// minimize unnecessary API calls.
		return &clustererror.RequeueAfterError{RequeueAfter: requeueAfterFatalSeconds * time.Second}
	}
	runningInstances := getRunningFromInstances(existingInstances)
	runningLen := len(runningInstances)
	var newestInstance *ecs.Instance
	if runningLen > 0 {
		// It would be very unusual to have more than one here, but it is
		// possible if someone manually provisions a machine with same tag name.
		glog.Infof("%s: found %d running instances for machine", machine.Name, runningLen)
		newestInstance = runningInstances[0]
	} else {
		// Didn't find any running instances, just newest existing one.
		// In most cases, there should only be one existing Instance.
		newestInstance = existingInstances[0]
	}

	a.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Updated", "Updated machine %v", machine.Name)

	// We do not support making changes to pre-existing instances, just update status.
	return a.updateStatus(machine, newestInstance)
}

func (a *Actuator) Exists(context context.Context, cluster *clusterv1.Cluster, machine *machinev1.Machine) (bool, error) {
	glog.Infof("%s: Checking if machine exists", machine.Name)

	instances, err := a.getMachineInstances(cluster, machine)
	if err != nil {
		glog.Errorf("%s: Error getting running instances: %v", machine.Name, err)
		return false, err
	}
	if len(instances) == 0 {
		glog.Infof("%s: Instance does not exist", machine.Name)
		return false, nil
	}

	// If more than one result was returned, it will be handled in Update.
	glog.Infof("%s: Instance exists as %q", machine.Name, instances[0].InstanceId)
	return true, nil
}

func (a *Actuator) Describe(cluster *clusterv1.Cluster, machine *machinev1.Machine) (*ecs.Instance, error) {
	glog.Infof("%s: Checking if machine exists", machine.Name)

	instances, err := a.getMachineInstances(cluster, machine)
	if err != nil {
		glog.Errorf("%s: Error getting running instances: %v", machine.Name, err)
		return nil, err
	}
	if len(instances) == 0 {
		glog.Infof("%s: Instance does not exist", machine.Name)
		return nil, nil
	}

	return instances[0], nil
}

func (a *Actuator) getMachineInstances(cluster *clusterv1.Cluster, machine *machinev1.Machine) ([]*ecs.Instance, error) {
	machineProviderConfig, err := providerConfigFromMachine(machine, a.codec)
	if err != nil {
		glog.Errorf("%s: Error decoding MachineProviderConfig: %v", machine.Name, err)
		return nil, err
	}

	region := machineProviderConfig.RegionId
	credentialsSecretName := ""
	if machineProviderConfig.CredentialsSecret != nil {
		credentialsSecretName = machineProviderConfig.CredentialsSecret.Name
	}
	aliCloudClient, err := a.aliCloudClientBuilder(a.client, credentialsSecretName, machine.Namespace, region)
	if err != nil {
		errMsg := fmt.Sprintf("%s: Error getting ECS client: %v", machine.Name, err)
		glog.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	return getExistingInstances(machine, aliCloudClient, region)
}

//get clusterID
func getClusterID(machine *machinev1.Machine) (string, bool) {
	clusterID, ok := machine.Labels[providerconfigv1.ClusterIDLabel]
	// NOTE: This block can be removed after the label renaming transition to machine.openshift.io
	if !ok {
		clusterID, ok = machine.Labels["sigs.k8s.io/cluster-api-cluster"]
	}
	return clusterID, ok
}
