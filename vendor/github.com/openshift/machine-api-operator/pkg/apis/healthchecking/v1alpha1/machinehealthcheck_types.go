package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapNodeUnhealthyConditions contains the name of the unhealthy conditions config map
const ConfigMapNodeUnhealthyConditions = "node-unhealthy-conditions"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineHealthCheck is the Schema for the machinehealthchecks API
// +k8s:openapi-gen=true
type MachineHealthCheck struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineHealthCheckSpec   `json:"spec,omitempty"`
	Status MachineHealthCheckStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineHealthCheckList contains a list of MachineHealthCheck
type MachineHealthCheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachineHealthCheck `json:"items"`
}

// MachineHealthCheckSpec defines the desired state of MachineHealthCheck
type MachineHealthCheckSpec struct {
	Selector metav1.LabelSelector `json:"selector"`
}

// MachineHealthCheckStatus defines the observed state of MachineHealthCheck
type MachineHealthCheckStatus struct {
	// TODO(alberto)
}
