<<<<<<< HEAD
<<<<<<< HEAD
# OpenShift cluster-api-provider-alibabacloud

<<<<<<< HEAD
This repository hosts an implementation of a provider for AlibabaCloud for the
=======
This repository hosts an implementation of a provider for Alibabacloud for the
>>>>>>> 8dbd34ff (update project name)
=======
# OpenShift cluster-api-provider-alibabacloud

This repository hosts an implementation of a provider for AlibabaCloud for the
>>>>>>> 3c01667c (ignore vendor)
OpenShift [machine-api](https://github.com/openshift/cluster-api).

This provider runs as a machine-controller deployed by the
[machine-api-operator](https://github.com/openshift/machine-api-operator)

### How to build the images in the RH infrastructure
The Dockerfiles use `as builder` in the `FROM` instruction which is not currently supported
by the RH's docker fork (see [https://github.com/kubernetes-sigs/kubebuilder/issues/268](https://github.com/kubernetes-sigs/kubebuilder/issues/268)).
One needs to run the `imagebuilder` command instead of the `docker build`.

Note: this info is RH only, it needs to be backported every time the `README.md` is synced with the upstream one.

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 3c01667c (ignore vendor)
## Deploy machine API plane with minikube

1. **Install kvm**

   Depending on your virtualization manager you can choose a different [driver](https://github.com/kubernetes/minikube/blob/master/docs/drivers.md).
   In order to install kvm, you can run (as described in the [drivers](https://github.com/kubernetes/minikube/blob/master/docs/drivers.md#kvm2-driver) documentation):

    ```sh
    $ sudo yum install libvirt-daemon-kvm qemu-kvm libvirt-daemon-config-network
    $ systemctl start libvirtd
    $ sudo usermod -a -G libvirt $(whoami)
    $ newgrp libvirt
    ```

   To install to kvm2 driver:

    ```sh
    curl -Lo docker-machine-driver-kvm2 https://storage.googleapis.com/minikube/releases/latest/docker-machine-driver-kvm2 \
    && chmod +x docker-machine-driver-kvm2 \
    && sudo cp docker-machine-driver-kvm2 /usr/local/bin/ \
    && rm docker-machine-driver-kvm2
    ```

2. **Deploying the cluster**

   To install minikube `v1.1.0`, you can run:

    ```sg
    $ curl -Lo minikube https://storage.googleapis.com/minikube/releases/v1.1.0/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
    ```

   To deploy the cluster:

    ```
    $ minikube start --vm-driver kvm2 --kubernetes-version v1.13.1 --v 5
    $ eval $(minikube docker-env)
    ```

3. **Deploying machine API controllers**

<<<<<<< HEAD
   For development purposes the AlibabaCloud machine controller itself will run out of the machine API stack.
=======
   For development purposes the aws machine controller itself will run out of the machine API stack.
>>>>>>> 3c01667c (ignore vendor)
   Otherwise, docker images needs to be built, pushed into a docker registry and deployed within the stack.

   To deploy the stack:
    ```
    kustomize build config | kubectl apply -f -
    ```

<<<<<<< HEAD
4. **Deploy secret with AlibabaCloud credentials**

   AlibabaCloud actuator assumes existence of a secret file (references in machine object) with base64 encoded credentials:
=======
4. **Deploy secret with AWS credentials**

   AWS actuator assumes existence of a secret file (references in machine object) with base64 encoded credentials:
>>>>>>> 3c01667c (ignore vendor)

   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
<<<<<<< HEAD
     name: alibabacloud-credentials-secret
     namespace: default
   type: Opaque
   data:
     accessKeyID: FILLIN
     accessKeySecret: FILLIN
   ```

   Save the above resource as **_secret.yaml_** and then apply it:
   ```sh
   kubectl apply -f secret.yaml
   ```

## Test locally built AlibabaCloud actuator
=======
## Test locally built alibabacloud actuator
>>>>>>> 8dbd34ff (update project name)
=======
     name: aws-credentials-secret
     namespace: default
   type: Opaque
   data:
     aws_access_key_id: FILLIN
     aws_secret_access_key: FILLIN
   ```

   You can use `examples/render-aws-secrets.sh` script to generate the secret:
   ```sh
   ./examples/render-aws-secrets.sh examples/addons.yaml | kubectl apply -f -
   ```

5. **Provision AWS resource**

   The actuator expects existence of certain resource in AWS such as:
    - vpc
    - subnets
    - security groups
    - etc.

   To create them, you can run:

   ```sh
   $ ENVIRONMENT_ID=aws-actuator-k8s ./hack/aws-provision.sh install
   ```

   To delete the resources, you can run:

   ```sh
   $ ENVIRONMENT_ID=aws-actuator-k8s ./hack/aws-provision.sh destroy
   ```

   All machine manifests expect `ENVIRONMENT_ID` to be set to `aws-actuator-k8s`.

## Test locally built alibabacloud actuator
>>>>>>> 3c01667c (ignore vendor)

1. **Tear down machine-controller**

   Deployed machine API plane (`machine-api-controllers` deployment) is (among other
   controllers) running `machine-controller`. In order to run locally built one,
   simply edit `machine-api-controllers` deployment and remove `machine-controller` container from it.

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
2. **Build and run AlibabaCloud actuator outside of the cluster**
=======
1. **Build and run alibabacloud actuator outside of the cluster**
>>>>>>> 3c01667c (ignore vendor)

   ```sh
   $ go build -o bin/machine-controller-manager github.com/AliyunContainerService/cluster-api-provider-alibabacloud/cmd/manager
   ```

   ```sh
<<<<<<< HEAD
   $ ./bin/machine-controller-manager --kubeconfig ~/.kube/config --logtostderr -v 5 -alsologtostderr
   ```
   If running in container with `podman`, or locally without `docker` installed, and encountering issues, see [hacking-guide](https://github.com/openshift/machine-api-operator/blob/master/docs/dev/hacking-guide.md#troubleshooting-make-targets).

=======
   $ .bin/machine-controller-manager --kubeconfig ~/.kube/config --logtostderr -v 5 -alsologtostderr
   ```
   If running in container with `podman`, or locally without `docker` installed, and encountering issues, see [hacking-guide](https://github.com/openshift/machine-api-operator/blob/master/docs/dev/hacking-guide.md#troubleshooting-make-targets).


>>>>>>> 3c01667c (ignore vendor)
1. **Deploy k8s apiserver through machine manifest**:

   To deploy user data secret with kubernetes apiserver initialization (under [config/master-user-data-secret.yaml](config/master-user-data-secret.yaml)):

<<<<<<< HEAD
   ```sh
=======
   ```yaml
>>>>>>> 3c01667c (ignore vendor)
   $ kubectl apply -f config/master-user-data-secret.yaml
   ```

   To deploy kubernetes master machine (under [config/master-machine.yaml](config/master-machine.yaml)):

<<<<<<< HEAD
   ```sh
   $ kubectl apply -f config/master-machine.yaml
   ```

1. **Join worker node through machine manifest**:

   To deploy user data secret with kubernetes apiserver initialization (under [config/worker-user-data-secret.yaml](config/worker-user-data-secret.yaml)):

   ```sh
   $ kubectl apply -f config/worker-user-data-secret.yaml
   ```

   To deploy kubernetes worker machine (under [config/worker-machine.yaml](config/worker-machine.yaml)):

   ```sh
   $ kubectl apply -f config/worker-machine.yaml
   ```

=======
   ```yaml
   $ kubectl apply -f config/master-machine.yaml
   ```

>>>>>>> 3c01667c (ignore vendor)
1. **Pull kubeconfig from created master machine**

   The master public IP can be accessed from AlibabaCloud Portal. Once done, you
   can collect the kube config by running:

   ```
<<<<<<< HEAD
   $ ssh -i SSHPMKEY root@PUBLICIP 'sudo cat /root/.kube/config' > kubeconfig
   $ kubectl --kubeconfig=kubeconfig config set-cluster kubernetes --server=https://PUBLICIP:6443
=======
   $ ssh -i SSHPMKEY ec2-user@PUBLICIP 'sudo cat /root/.kube/config' > kubeconfig
   $ kubectl --kubeconfig=kubeconfig config set-cluster kubernetes --server=https://PUBLICIP:8443
>>>>>>> 3c01667c (ignore vendor)
   ```

   Once done, you can access the cluster via `kubectl`. E.g.

   ```sh
   $ kubectl --kubeconfig=kubeconfig get nodes
   ```

<<<<<<< HEAD

## Deploy machine API plane with AlibabaCloud ACK Cluster

1. **Creating ACK Cluster**

    You can create a Kubernetes cluster using the CLI, TerraForm, or ACK console

    CLI Document:
    ```
   https://www.alibabacloud.com/help/doc-detail/198808.htm
    ```

   TerraForm Document:
    ```
   https://www.alibabacloud.com/help/doc-detail/252824.htm
    ```

   ACK Console Document:
    ```
   https://www.alibabacloud.com/help/doc-detail/86488.htm
    ```


2. **Deploying machine API controllers**

   For development purposes the AlibabaCloud machine controller itself will run out of the machine API stack.
   Otherwise, docker images needs to be built, pushed into a docker registry and deployed within the stack.

   To deploy the machine crds:
    ```sh
    $ kubectl apply -f config/crds/
    ```

   To deploy the machine rbac:
    ```sh
    $ kubectl apply -f config/rbac/
    ```

   To deploy the machine controller:
    ```sh
    $ kubectl apply -f config/controllers/
    ```

3. **Deploy secret with AlibabaCloud credentials**

   AlibabaCloud actuator assumes existence of a secret file (references in machine object) with base64 encoded credentials:

   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: alibabacloud-credentials-secret
     namespace: default
   type: Opaque
   data:
     accessKeyID: FILLIN
     accessKeySecret: FILLIN
   ```

   Save the above resource as **_secret.yaml_** and then apply it:
   ```sh
   $ kubectl apply -f secret.yaml
   ``` 

1. **Deploy k8s apiserver through machine manifest**:

   To deploy user data secret with kubernetes apiserver initialization (under [config/master-user-data-secret.yaml](config/master-user-data-secret.yaml)):

   ```sh
   $ kubectl apply -f config/master-user-data-secret.yaml
   ```

   To deploy kubernetes master machine (under [config/master-machine.yaml](config/master-machine.yaml)):

   ```sh
   $ kubectl apply -f config/master-machine.yaml
   ```

1. **Join worker node through machine manifest**:

   To deploy user data secret with kubernetes apiserver initialization (under [config/worker-user-data-secret.yaml](config/worker-user-data-secret.yaml)):

   ```sh
   $ kubectl apply -f config/worker-user-data-secret.yaml
   ```

   To deploy kubernetes worker machine (under [config/worker-machine.yaml](config/worker-machine.yaml)):

   ```sh
   $ kubectl apply -f config/worker-machine.yaml
=======
## Deploy k8s cluster in AlibabaCloud with machine API plane deployed

1. **Generate bootstrap user data**

   To generate bootstrap script for machine api plane, simply run:

   ```sh
   $ ./config/generate-bootstrap.sh
   ```

   The script requires `ALIBABACLOUD_ACCESS_KEY_ID` and `ALIBABACLOUD_SECRET_ACCESS_KEY` environment variables to be set.
   It generates `config/bootstrap.yaml` secret for master machine
   under `config/master-machine.yaml`.

   The generated bootstrap secret contains user data responsible for:
    - deployment of kube-apiserver
    - deployment of machine API plane with alibabacloud machine controllers
    - generating worker machine user data script secret deploying a node
    - deployment of worker machineset

1. **Deploy machine API plane through machine manifest**:

   First, deploy generated bootstrap secret:

   ```yaml
   $ kubectl apply -f config/bootstrap.yaml
   ```

   Then, deploy master machine (under [config/master-machine.yaml](config/master-machine.yaml)):

   ```yaml
   $ kubectl apply -f config/master-machine.yaml
>>>>>>> 3c01667c (ignore vendor)
   ```

1. **Pull kubeconfig from created master machine**

   The master public IP can be accessed from AlibabaCloud Portal. Once done, you
   can collect the kube config by running:

   ```
<<<<<<< HEAD
   $ ssh -i SSHPMKEY root@PUBLICIP 'sudo cat /root/.kube/config' > kubeconfig
=======
   $ ssh -i SSHPMKEY ecs-user@PUBLICIP 'sudo cat /root/.kube/config' > kubeconfig
>>>>>>> 3c01667c (ignore vendor)
   $ kubectl --kubeconfig=kubeconfig config set-cluster kubernetes --server=https://PUBLICIP:6443
   ```

   Once done, you can access the cluster via `kubectl`. E.g.

   ```sh
   $ kubectl --kubeconfig=kubeconfig get nodes
   ```

<<<<<<< HEAD
### Add worker nodes to the ACK cluster via Machine-API

1. **Deploy secret with AlibabaCloud worker nodes userdata**

   AlibabaCloud actuator assumes existence of a secret file (references in machine object) with base64 encoded userdata:

   How do I get the script to add worker nodes? You can refer to the documentation
   
   ```
   https://www.alibabacloud.com/help/doc-detail/86919.htm
   ```

   And then generate the userdata:

   ```sh
   $ echo '#!/bin/bash  <Your worker node script>' | base64
   ```

   Replace FILLIN with userdata:
   
   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: worker-user-data-secret
     namespace: default
   type: Opaque
   data:
    userData: FILLIN
   ```

   Save the above resource as **_worker-user-data-secret.yaml_** and then apply it:
   ```sh
   kubectl apply -f worker-user-data-secret.yaml
   ``` 

 1. **Add worker machine to ACK Cluster**
    
   ```yaml
   apiVersion: machine.openshift.io/v1beta1
   kind: Machine
   metadata:
     name: alibabacloud-actuator-testing-machine
     namespace: default
     labels:
       machine.openshift.io/cluster-api-cluster: alibabacloud-actuator-k8s
   spec:
     metadata:
       labels:
         node-role.kubernetes.io/infra: ""
     providerSpec:
       value:
         apiVersion: alibabacloudproviderconfig.openshift.io/v1alpha1
         kind: AlibabaCloudMachineProviderConfig
         instanceType: FILLIN
         imageId: FILLIN
         regionId: FILLIN
         zoneId: FILLIN
         securityGroupId: FILLIN
         vpcId: FILLIN
         vSwitchId: FILLIN
         systemDiskCategory: FILLIN
         systemDiskSize: FILLIN
         internetMaxBandwidthOut: FILLIN
         password: FILLIN
         tags:
           - key: openshift-node-group-config
             value: node-config-node
           - key: host-type
             value: node
           - key: sub-host-type
             value: default
         userDataSecret:
           name: alibabacloud-worker-user-data-secret
         credentialsSecret:
           name: alibabacloud-credentials-secret
   ```
   
     Save the above resource as **_worker-machine-with-user-data.yaml_** and then apply it:

   ```sh
   kubectl apply -f worker-machine-with-user-data.yaml
   ``` 

   Once done, you can describe the machine via `kubectl`. E.g.
   
   ```sh
   $ kubectl  get machine
   ```
 
# Upstream Implementation
Other branches of this repository may choose to track the upstream
Kubernetes [Cluster-API AlibabaCloud provider](https://github.com/AliyunContainerService/cluster-api-provider-alibabacloud)
=======
1. **Build and run alibabacloud actuator outside of the cluster**
=======
2. **Build and run alibabacloud actuator outside of the cluster**
>>>>>>> fc426375 (fix manager run error)

   ```sh
   $ go build -o bin/manager github.com/AliyunContainerService/cluster-api-provider-alibabacloud/cmd/manager
   ```

   ```sh
   $ ./bin/manager --kubeconfig ~/.kube/config --logtostderr -v 5 -alsologtostderr
   ```

3. **Build and run alibabacloud actuator outside of the cluster**

   ```sh
   $ go build -o bin/alicloud-actuator github.com/AliyunContainerService/cluster-api-provider-alibabacloud/cmd/alicloud-actuator
   ```

4. **Run machine controller in a kubernetes cluser.**
   
   Deploy crds
   
   ```sh
   $ kubectl apply -f ./config/crds/
   ```
   
   Deploy rbac
   
   ```sh
   $ kubectl apply -f ./config/rbac/
   ```   
   
   Before you deploy machine-controller, you edit the file ./config/configmap/user_config.yaml first, and fill your info ,and then 
   
   ```sh
   $ kubectl apply  -f ./config/configmap/
   ```   
   
   
   ```sh
   $ kubectl apply -f ./config/controllers/
   ```   
   




>>>>>>> 8dbd34ff (update project name)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
# Upstream Implementation
Other branches of this repository may choose to track the upstream
Kubernetes [Cluster-API AlibabaCloud provider](https://github.com/AliyunContainerService/cluster-api-provider-alibabacloud)

In the future, we may align the master branch with the upstream project as it
stabilizes within the community.
>>>>>>> 3c01667c (ignore vendor)
