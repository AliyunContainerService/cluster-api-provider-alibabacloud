module github.com/AliyunContainerService/cluster-api-provider-alibabacloud

go 1.15

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1153
	github.com/blang/semver v3.5.1+incompatible
	github.com/go-logr/logr v0.4.0
	github.com/golang/mock v1.6.0
	github.com/onsi/gomega v1.10.5
	github.com/openshift/api v0.0.0-20210416115537-a60c0dc032fd
	github.com/openshift/machine-api-operator v0.2.1-0.20210504014029-a132ec00f7dd
	github.com/stretchr/testify v1.7.0

	// kube 1.18
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/controller-runtime v0.9.0-beta.1.0.20210512131817-ce2f0c92d77e
	sigs.k8s.io/controller-tools v0.3.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20210420175812-638f9f3fbb42
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20210408182022-987bc3d6a107
)
