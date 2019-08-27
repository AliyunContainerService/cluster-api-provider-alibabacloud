package v1alpha1

// AliCloudMachineProviderConditionType is a valid value for AlibabaCloudMachineProviderCondition.Type
type AliCloudMachineProviderConditionType string

// Valid conditions for an AliCloud machine instance
const (
	// MachineCreation indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreation AliCloudMachineProviderConditionType = "MachineCreation"

	MachineCreated AliCloudMachineProviderConditionType = "MachineCreated"
)
