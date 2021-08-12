package machines

import (
	. "github.com/onsi/gomega"
	MachineV1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"

	providerconfigv1 "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudprovider/v1alpha1"
)

func createSecretAndWait(f *framework.Framework, secret *apiv1.Secret) {
	_, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
	Expect(err).NotTo(HaveOccurred())

	err = wait.Poll(framework.PollInterval, framework.PoolTimeout, func() (bool, error) {
		if _, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{}); err != nil {
			return false, nil
		}
		return true, nil
	})
	Expect(err).NotTo(HaveOccurred())
}

func getMachineProviderStatus(f *framework.Framework, machine *MachineV1beta1.Machine) *providerconfigv1.AlibabaCloudMachineProviderStatus {
	machine, err := f.CAPIClient.MachineV1beta1().Machines(machine.Namespace).Get(machine.Name, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	codec, err := providerconfigv1.NewCodec()
	Expect(err).NotTo(HaveOccurred())

	machineProviderStatus := &providerconfigv1.AlibabaCloudMachineProviderStatus{}
	err = codec.DecodeProviderStatus(machine.Status.ProviderStatus, machineProviderStatus)
	Expect(err).NotTo(HaveOccurred())

	return machineProviderStatus
}

func getMachineCondition(f *framework.Framework, machine *MachineV1beta1.Machine) providerconfigv1.AlibabaCloudMachineProviderCondition {
	conditions := getMachineProviderStatus(f, machine).Conditions
	Expect(len(conditions)).To(Equal(1), "ambiguous conditions: %#v", conditions)
	return conditions[0]
}
