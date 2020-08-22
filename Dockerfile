# Build the manager binary
FROM quay.xiaodiankeji.net/openshift/release:golang-1.14 as builder

# Copy in the go src
WORKDIR /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud
COPY . . 

RUN unset VERSION \
    && GOPROXY=https://goproxy-dian-app.apps.dian-int.com,direct NO_DOCKER=1 make build

# Copy the controller-manager into a thin image
FROM quay.xiaodiankeji.net/openshift/origin-v4.0:base
WORKDIR /
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud/bin/machine-controller-manager ./
COPY --from=builder /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud/bin/termination-handler ./