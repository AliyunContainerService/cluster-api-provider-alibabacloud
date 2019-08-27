# OpenShift cluster-api-provider-alibabacloud

This repository hosts an implementation of a provider for Alibabacloud for the
OpenShift [machine-api](https://github.com/openshift/cluster-api).

This provider runs as a machine-controller deployed by the
[machine-api-operator](https://github.com/openshift/machine-api-operator)

### How to build the images in the RH infrastructure
The Dockerfiles use `as builder` in the `FROM` instruction which is not currently supported
by the RH's docker fork (see [https://github.com/kubernetes-sigs/kubebuilder/issues/268](https://github.com/kubernetes-sigs/kubebuilder/issues/268)).
One needs to run the `imagebuilder` command instead of the `docker build`.

Note: this info is RH only, it needs to be backported every time the `README.md` is synced with the upstream one.

## Test locally built alibabacloud actuator

1. **Tear down machine-controller**

   Deployed machine API plane (`machine-api-controllers` deployment) is (among other
   controllers) running `machine-controller`. In order to run locally built one,
   simply edit `machine-api-controllers` deployment and remove `machine-controller` container from it.

1. **Build and run alibabacloud actuator outside of the cluster**

   ```sh
   $ go build -o bin/manager github.com/AliyunContainerService/cluster-api-provider-alibabacloud/cmd/manager
   ```

   ```sh
   $ ./bin/manager --kubeconfig ~/.kube/config --logtostderr -v 5 -alsologtostderr
   ```

2. **Build and run alibabacloud actuator outside of the cluster**

   ```sh
   $ go build -o bin/alicloud-actuator github.com/AliyunContainerService/cluster-api-provider-alibabacloud/cmd/alicloud-actuator
   ```



