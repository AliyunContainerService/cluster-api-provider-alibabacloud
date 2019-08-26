<<<<<<< HEAD
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
=======
# Build the manager binary
FROM registry.svc.ci.openshift.org/openshift/release:golang-1.12 as builder

# Copy in the go src
WORKDIR /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY test/  test/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./machine-controller-manager github.com/AliyunContainerService/cluster-api-provider-alibabacloud/cmd/manager
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./manager ./vendor/github.com/openshift/cluster-api/cmd/manager

# Copy the controller-manager into a thin image
FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base
WORKDIR /
<<<<<<< HEAD
<<<<<<< HEAD
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alicloud/manager .
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alicloud/machine-controller-manager .
>>>>>>> ac8840a3 (update image)
=======
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alicloud/manager ./
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alicloud/machine-controller-manager ./
>>>>>>> 62fa6daf (fix config)
=======
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud/manager ./
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud/machine-controller-manager ./
>>>>>>> 8dbd34ff (update project name)
