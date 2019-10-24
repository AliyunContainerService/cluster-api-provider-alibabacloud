/*
Copyright 2019 The Kubernetes Authors.
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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ClusterIDLabel = "machine.openshift.io/cluster-api-cluster"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AlicloudMachineProviderSpec defines the desired state of AlibabaCloudMachineProviderConfig
// +k8s:openapi-gen=true
//type AlicloudMachineProviderSpec struct {
//	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
//	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
//	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
//
//	//The ID of the image file that you specified when you create the instance. Example :centos_7_06_64_20G_alibase_20190619.vhd
//	ImageId string `json:"imageId"`
//
//	// The type of the instance. Example: ecs.n4.large
//	InstanceType string `json:"instanceType"`
//
//	//The region ID of the instance. Example: cn-hangzhou
//	RegionId string `json:"regionId"`
//
//	/**
//	The name of the instance. It must be 2 to 128 characters in length and can contain uppercase and lowercase letters,
//	digits, colons (:), underscores (_), and hyphens (-). It must start with a letter. It cannot start with http:// or https://.
//	If this parameter is not specified, InstanceId is used by default.
//	*/
//	InstanceName string `json:"instanceName"`
//
//	// PublicIP specifies whether the instance should get a public IP. If not present,
//	// it should use the default of its subnet.
//	PublicIP bool `json:"publicIp"`
//
//	//The ID of the VPC
//	VpcId string `json:"vpcId"`
//
//	//The ID of the VSwitch. It must be specified when you create VPC-type instances.
//	VSwitchId string `json:"vSwitchId"`
//
//	// KeyPairName is the name of the KeyPair to use for SSH
//	KeyPairName string `json:"keyPairName"`
//
//	// UserDataSecret contains a local reference to a secret that contains the
//	// UserData to apply to the instance
//	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`
//
//	// CredentialsSecret is a reference to the secret with AliCloud credentials. Otherwise, defaults to permissions
//	// provided by attached RAM role where the actuator is running.
//	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`
//
//	//The RAM role name of the instance
//	RamRoleName string `json:"ramRoleName"`
//
//	//The ID of the security group to which the instance belongs
//	SecurityGroupId string `json:"securityGroupId"`
//
//	//The tags  of the instance
//	Tags []TagSpecification `json:"tags,omitempty"`
//
//	/**
//	The category of the system disk. The default value of the non-optimized instance for phased-out instance types for which I/O optimization is not performed is cloud. The default value for other instances is cloud_efficiency. Valid values:
//
//	cloud: basic disk.
//	cloud_efficiency: ultra disk.
//	cloud_ssd: SSD.
//	ephemeral_ssd: local SSD.
//	cloud_essd: ESSD. ESSDs are still in the public beta test phase and only available in some regions. For more information, see ESSD FAQ.
//	*/
//	SystemDiskCategory string `json:"systemDiskCategory"`
//
//	/**
//	The size of the system disk. Unit: GiB. Valid values: 20 to 500.
//
//	This value must be equal to or greater than max {20, ImageSize}. Default value: max {40, ImageSize}.
//	*/
//	SystemDiskSize int64 `json:"systemDiskSize"`
//
//	/**
//	The name of the system disk. It must be 2 to 128 characters in length and can contain uppercase and lowercase letters, digits, colons (:),
//	underscores (_), and hyphens (-). It must start with a letter. It cannot start with http:// or https://. Default value: null.
//	*/
//	SystemDiskDiskName string `json:"systemDiskDiskName"`
//
//	/**
//	The description of the system disk. It must be 2 to 256 characters in length. It cannot start with http:// or https://.
//	*/
//	SystemDiskDescription string `json:"systemDiskDescription"`
//
//	//DataDisks of the instance
//	DataDisks []DataDiskSpecification `json:"dataDisks,omitempty"`
//
//	/**
//	The release protection attribute of the instance. It indicates whether you can use the ECS console or
//	call the DeleteInstance action to release the instance. Default value: false. Valid values:
//	true: enables release protection.
//	false: disables release protection.
//	*/
//	DeletionProtection bool `json:"deletionProtection"`
//}

// TagSpecification is the name/value pair for a tag
type TagSpecification struct {
	// Name of the tag
	Key string `json:"key"`

	// Value of the tag
	Value string `json:"value"`
}

// DataDiskSpecification is the definition of datadisk
type DataDiskSpecification struct {
	/**
	The category of the nth data disk. Default value: cloud. Valid values:

	cloud: basic disk.
	cloud_efficiency: ultra disk.
	cloud_ssd: SSD.
	ephemeral_ssd: local SSD.
	cloud_essd: ESSD. ESSDs are still in the public beta test phase and only available in some regions
	*/
	Category string `json:"category"`

	/**
	Indicates whether the data disk is released along with the instance. Default value: true.
	*/
	DeleteWithInstance *bool `json:"deleteWithInstance,omitempty"`

	/**
	The description of the data disk. It must be 2 to 256 characters in length. It cannot start with http:// or https://. Default value: null.
	*/
	Description string `json:"description"`

	/**
	The mount point of the nth data disk.Default value : /dev/xvda
	*/
	Device string `json:"device"`

	/**
	The name of the nth data disk. It must be 2 to 128 characters in length and can contain uppercase and lowercase letters,
	digits, colons (:), underscores (_), and hyphens (-). It must start with a letter. It cannot start with http://
	or https://. Default value: null.

	*/
	DiskName string `json:"diskName"`

	/**
	Indicates whether the nth data disk is encrypted. Default value: false.
	*/
	Encrypted bool `json:"encrypted"`

	/**
	The size of the nth data disk. n ranges from 1 to 16. Unit: GiB. Valid values:

	cloud: 5 to 2,000.
	cloud_efficiency: 20 to 32,768.
	cloud_ssd: 20 to 32,768.
	cloud_essd: 20 to 32,768.
	ephemeral_ssd: 5 to 800.
	This value must be equal to or greater than the size of SnapshotId.
	*/
	Size int64 `json:"size"`

	/**
	The ID of the snapshot used to create the nth data disk
	*/
	SnapshotId string `json:"snapshotId"`
}

// AlibabaCloudMachineProviderStatus defines the observed state of AlibabaCloudMachineProviderConfig
// +k8s:openapi-gen=true
type AlibabaCloudMachineProviderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// InstanceID is the instance ID of the machine created in AliCloud
	// +optional
	InstanceID *string `json:"instanceId,omitempty"`

	// InstanceState is the state of the ECS instance for this machine
	// +optional
	InstanceStatus *string `json:"instanceStatus,omitempty"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []AlibabaCloudMachineProviderCondition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlibabaCloudMachineProviderConfig is the Schema for the alicloudmachineproviderconfigs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type AlibabaCloudMachineProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	//The ID of the image file that you specified when you create the instance. Example :centos_7_06_64_20G_alibase_20190619.vhd
	ImageId string `json:"imageId"`

	// The type of the instance. Example: ecs.n4.large
	InstanceType string `json:"instanceType"`

	//The region ID of the instance. Example: cn-hangzhou
	RegionId string `json:"regionId"`

	/**
	The name of the instance. It must be 2 to 128 characters in length and can contain uppercase and lowercase letters,
	digits, colons (:), underscores (_), and hyphens (-). It must start with a letter. It cannot start with http:// or https://.
	If this parameter is not specified, InstanceId is used by default.
	*/
	InstanceName string `json:"instanceName"`

	// PublicIP specifies whether the instance should get a public IP. If not present,
	// it should use the default of its subnet.
	PublicIP bool `json:"publicIp"`

	//The ID of the VPC
	VpcId string `json:"vpcId"`

	//The ID of the VSwitch. It must be specified when you create VPC-type instances.
	VSwitchId string `json:"vSwitchId"`

	// KeyPairName is the name of the KeyPair to use for SSH
	KeyPairName string `json:"keyPairName"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	// CredentialsSecret is a reference to the secret with AliCloud credentials. Otherwise, defaults to permissions
	// provided by attached RAM role where the actuator is running.
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`

	//The RAM role name of the instance
	RamRoleName string `json:"ramRoleName"`

	//The ID of the security group to which the instance belongs
	SecurityGroupId string `json:"securityGroupId"`

	//The tags  of the instance
	Tags []TagSpecification `json:"tags,omitempty"`

	/**
	The category of the system disk. The default value of the non-optimized instance for phased-out instance types for which I/O optimization is not performed is cloud. The default value for other instances is cloud_efficiency. Valid values:

	cloud: basic disk.
	cloud_efficiency: ultra disk.
	cloud_ssd: SSD.
	ephemeral_ssd: local SSD.
	cloud_essd: ESSD. ESSDs are still in the public beta test phase and only available in some regions. For more information, see ESSD FAQ.
	*/
	SystemDiskCategory string `json:"systemDiskCategory"`

	/**
	The size of the system disk. Unit: GiB. Valid values: 20 to 500.

	This value must be equal to or greater than max {20, ImageSize}. Default value: max {40, ImageSize}.
	*/
	SystemDiskSize int64 `json:"systemDiskSize"`

	/**
	The name of the system disk. It must be 2 to 128 characters in length and can contain uppercase and lowercase letters, digits, colons (:),
	underscores (_), and hyphens (-). It must start with a letter. It cannot start with http:// or https://. Default value: null.
	*/
	SystemDiskDiskName string `json:"systemDiskDiskName"`

	/**
	The description of the system disk. It must be 2 to 256 characters in length. It cannot start with http:// or https://.
	*/
	SystemDiskDescription string `json:"systemDiskDescription"`

	//DataDisks of the instance
	DataDisks []DataDiskSpecification `json:"dataDisks,omitempty"`

	/**
	The release protection attribute of the instance. It indicates whether you can use the ECS console or
	call the DeleteInstance action to release the instance. Default value: false. Valid values:
	true: enables release protection.
	false: disables release protection.
	*/
	DeletionProtection bool `json:"deletionProtection"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlibabaCloudMachineProviderList contains a list of AlibabaCloudMachineProviderConfig
type AlibabaCloudMachineProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlibabaCloudMachineProviderConfig `json:"items"`
}

// AlibabaCloudMachineProviderCondition is a condition in a AlibabaCloudMachineProviderStatus
type AlibabaCloudMachineProviderCondition struct {
	// Type is the type of the condition.
	Type AliCloudMachineProviderConditionType `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

func init() {
	SchemeBuilder.Register(&AlibabaCloudMachineProviderConfig{}, &AlibabaCloudMachineProviderList{}, &AlibabaCloudMachineProviderStatus{})
}
