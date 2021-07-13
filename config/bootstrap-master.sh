#!/bin/bash

cat <<HEREDOC > /root/user-data.sh
#!/bin/bash

cat <<EOF | tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF

cat <<EOF | tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sysctl --system


cat <<EOF | tee /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
repo_gpgcheck=0
gpgkey=http://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg
       http://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF

setenforce 0
sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

swapoff -a


yum install -y docker
systemctl enable docker
systemctl start docker

yum install -y kubelet-1.21.0 kubeadm-1.21.0 kubectl-1.21.0 --disableexcludes=kubernetes

systemctl enable --now kubelet


systemctl daemon-reload
systemctl enable kubelet && systemctl start kubelet

docker --version
kubelet --version

cat <<EOF > /etc/default/kubelet
KUBELET_KUBEADM_EXTRA_ARGS=--cgroup-driver=systemd
EOF

echo '1' > /proc/sys/net/bridge/bridge-nf-call-iptables

################################################
######## Config kubeadm images
################################################

# List images
kubeadm config images list

# Replace kube-apiserver
docker pull registry.aliyuncs.com/google_containers/kube-apiserver:v1.21.0
docker tag registry.aliyuncs.com/google_containers/kube-apiserver:v1.21.0 k8s.gcr.io/kube-apiserver:v1.21.0
docker rmi registry.aliyuncs.com/google_containers/kube-apiserver:v1.21.0

# Replace kube-controller-manager
docker pull registry.aliyuncs.com/google_containers/kube-controller-manager:v1.21.0
docker tag registry.aliyuncs.com/google_containers/kube-controller-manager:v1.21.0 k8s.gcr.io/kube-controller-manager:v1.21.0
docker rmi registry.aliyuncs.com/google_containers/kube-controller-manager:v1.21.0

# Replace kube-scheduler
docker pull registry.aliyuncs.com/google_containers/kube-scheduler:v1.21.0
docker tag registry.aliyuncs.com/google_containers/kube-scheduler:v1.21.0 k8s.gcr.io/kube-scheduler:v1.21.0
docker rmi registry.aliyuncs.com/google_containers/kube-scheduler:v1.21.0

# Replace kube-proxy
docker pull registry.aliyuncs.com/google_containers/kube-proxy:v1.21.0
docker tag registry.aliyuncs.com/google_containers/kube-proxy:v1.21.0 k8s.gcr.io/kube-proxy:v1.21.0
docker rmi registry.aliyuncs.com/google_containers/kube-proxy:v1.21.0

# Replace pause
docker pull registry.aliyuncs.com/google_containers/pause:3.4.1
docker tag registry.aliyuncs.com/google_containers/pause:3.4.1 k8s.gcr.io/pause:3.4.1
docker rmi registry.aliyuncs.com/google_containers/pause:3.4.1

# Replace etcd
docker pull registry.aliyuncs.com/google_containers/etcd:3.4.13-0
docker tag registry.aliyuncs.com/google_containers/etcd:3.4.13-0 k8s.gcr.io/etcd:3.4.13-0
docker rmi registry.aliyuncs.com/google_containers/etcd:3.4.13-0

# Replace coredns
<<<<<<< HEAD
docker pull registry.aliyuncs.com/google_containers/coredns:1.8.0
docker tag registry.aliyuncs.com/google_containers/coredns:1.8.0 k8s.gcr.io/coredns/coredns:v1.8.0
docker rmi registry.aliyuncs.com/google_containers/coredns:1.8.0
=======
docker pull coredns/coredns:1.8.0
docker tag coredns/coredns:1.8.0 k8s.gcr.io/coredns/coredns:v1.8.0
docker rmi coredns/coredns:1.8.0
>>>>>>> 56ed82a5 (add master and worker userdata for kubeadm)



################################################
######## Deploy kubernetes master
################################################

<<<<<<< HEAD
kubeadm init  --apiserver-bind-port 6443 --token bign04.m6l27w1gu6vqbloq  --token-ttl 0  --kubernetes-version=v1.21.0 --apiserver-cert-extra-sans=$(curl -s http://100.100.100.200/latest/meta-data/private-ipv4) --apiserver-cert-extra-sans=$(curl -s http://100.100.100.200/latest/meta-data/eipv4)   --pod-network-cidr=10.244.0.0/16
=======
kubeadm init  --apiserver-bind-port 6443 --token bign04.m6l27w1gu6vqbloq  --kubernetes-version=v1.21.0 --apiserver-cert-extra-sans=$(curl -s http://100.100.100.200/latest/meta-data/private-ipv4) --apiserver-cert-extra-sans=$(curl -s http://100.100.100.200/latest/meta-data/eipv4)   --pod-network-cidr=10.244.0.0/16
>>>>>>> 56ed82a5 (add master and worker userdata for kubeadm)

mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config

################################################
######## Deploy kube-flannel
################################################

cat <<EOF > kube-flannel.yaml
---
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: psp.flannel.unprivileged
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: docker/default
    seccomp.security.alpha.kubernetes.io/defaultProfileName: docker/default
    apparmor.security.beta.kubernetes.io/allowedProfileNames: runtime/default
    apparmor.security.beta.kubernetes.io/defaultProfileName: runtime/default
spec:
  privileged: false
  volumes:
  - configMap
  - secret
  - emptyDir
  - hostPath
  allowedHostPaths:
  - pathPrefix: "/etc/cni/net.d"
  - pathPrefix: "/etc/kube-flannel"
  - pathPrefix: "/run/flannel"
  readOnlyRootFilesystem: false
  # Users and groups
  runAsUser:
    rule: RunAsAny
  supplementalGroups:
    rule: RunAsAny
  fsGroup:
    rule: RunAsAny
  # Privilege Escalation
  allowPrivilegeEscalation: false
  defaultAllowPrivilegeEscalation: false
  # Capabilities
  allowedCapabilities: ['NET_ADMIN', 'NET_RAW']
  defaultAddCapabilities: []
  requiredDropCapabilities: []
  # Host namespaces
  hostPID: false
  hostIPC: false
  hostNetwork: true
  hostPorts:
  - min: 0
    max: 65535
  # SELinux
  seLinux:
    # SELinux is unused in CaaSP
    rule: 'RunAsAny'
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: flannel
rules:
- apiGroups: ['extensions']
  resources: ['podsecuritypolicies']
  verbs: ['use']
  resourceNames: ['psp.flannel.unprivileged']
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - nodes/status
  verbs:
  - patch
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: flannel
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: flannel
subjects:
- kind: ServiceAccount
  name: flannel
  namespace: kube-system
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: flannel
  namespace: kube-system
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: kube-flannel-cfg
  namespace: kube-system
  labels:
    tier: node
    app: flannel
data:
  cni-conf.json: |
    {
      "name": "cbr0",
      "cniVersion": "0.3.1",
      "plugins": [
        {
          "type": "flannel",
          "delegate": {
            "hairpinMode": true,
            "isDefaultGateway": true
          }
        },
        {
          "type": "portmap",
          "capabilities": {
            "portMappings": true
          }
        }
      ]
    }
  net-conf.json: |
    {
      "Network": "10.244.0.0/16",
      "Backend": {
        "Type": "vxlan"
      }
    }
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-flannel-ds
  namespace: kube-system
  labels:
    tier: node
    app: flannel
spec:
  selector:
    matchLabels:
      app: flannel
  template:
    metadata:
      labels:
        tier: node
        app: flannel
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/os
                operator: In
                values:
                - linux
      hostNetwork: true
      priorityClassName: system-node-critical
      tolerations:
      - operator: Exists
        effect: NoSchedule
      serviceAccountName: flannel
      initContainers:
      - name: install-cni
<<<<<<< HEAD
        image: registry.cn-hangzhou.aliyuncs.com/k8sos/flannel:v0.14.0
=======
        image: quay.io/coreos/flannel:v0.14.0
>>>>>>> 56ed82a5 (add master and worker userdata for kubeadm)
        command:
        - cp
        args:
        - -f
        - /etc/kube-flannel/cni-conf.json
        - /etc/cni/net.d/10-flannel.conflist
        volumeMounts:
        - name: cni
          mountPath: /etc/cni/net.d
        - name: flannel-cfg
          mountPath: /etc/kube-flannel/
      containers:
      - name: kube-flannel
<<<<<<< HEAD
        image: registry.cn-hangzhou.aliyuncs.com/k8sos/flannel:v0.14.0
=======
        image: quay.io/coreos/flannel:v0.14.0
>>>>>>> 56ed82a5 (add master and worker userdata for kubeadm)
        command:
        - /opt/bin/flanneld
        args:
        - --ip-masq
        - --kube-subnet-mgr
        resources:
          requests:
            cpu: "100m"
            memory: "50Mi"
          limits:
            cpu: "100m"
            memory: "50Mi"
        securityContext:
          privileged: false
          capabilities:
            add: ["NET_ADMIN", "NET_RAW"]
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        volumeMounts:
        - name: run
          mountPath: /run/flannel
        - name: flannel-cfg
          mountPath: /etc/kube-flannel/
      volumes:
      - name: run
        hostPath:
          path: /run/flannel
      - name: cni
        hostPath:
          path: /etc/cni/net.d
      - name: flannel-cfg
        configMap:
          name: kube-flannel-cfg
EOF

kubectl apply -f kube-flannel.yaml

HEREDOC

bash /root/user-data.sh > /root/user-data.logs
