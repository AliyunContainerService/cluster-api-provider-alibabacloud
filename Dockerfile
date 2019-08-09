<<<<<<< HEAD
FROM registry.ci.openshift.org/openshift/release:golang-1.15 AS builder
WORKDIR /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud
=======
FROM registry.svc.ci.openshift.org/openshift/release:golang-1.12 AS builder
WORKDIR /go/src/github.com/AliyunContainerService/cluster-api-provider-alicloud
>>>>>>> 5a63acd2 (update test case)
COPY . .
# VERSION env gets set in the openshift/release image and refers to the golang version, which interfers with our own
RUN unset VERSION \
  && GOPROXY=off NO_DOCKER=1 make build

<<<<<<< HEAD
FROM registry.ci.openshift.org/openshift/origin-v4.0:base
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud/bin/machine-controller-manager /
=======
FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base
RUN INSTALL_PKGS=" \
      openssh \
      " && \
    yum install -y $INSTALL_PKGS && \
    rpm -V $INSTALL_PKGS && \
    yum clean all
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alicloud/bin/manager /
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alicloud/bin/machine-controller-manager /
>>>>>>> 5a63acd2 (update test case)
