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
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud/manager ./
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud/machine-controller-manager ./