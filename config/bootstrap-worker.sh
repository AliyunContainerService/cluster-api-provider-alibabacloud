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
docker pull registry.aliyuncs.com/google_containers/coredns:1.8.0
docker tag registry.aliyuncs.com/google_containers/coredns:1.8.0 k8s.gcr.io/coredns/coredns:v1.8.0
docker rmi registry.aliyuncs.com/google_containers/coredns:1.8.0



################################################
######## Deploy kubernetes worker
################################################

kubeadm join FILLIN:6443 --token FILLIN --discovery-token-unsafe-skip-ca-verification


HEREDOC

bash /root/user-data.sh > /root/user-data.logs
