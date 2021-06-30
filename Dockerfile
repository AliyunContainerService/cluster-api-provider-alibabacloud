FROM registry.ci.openshift.org/openshift/release:golang-1.15 AS builder
WORKDIR /go/src/sigs.k8s.io/cluster-api-provider-alibabacloud
COPY . .
# VERSION env gets set in the openshift/release image and refers to the golang version, which interfers with our own
RUN unset VERSION \
  && GOPROXY=off NO_DOCKER=1 make build

FROM registry.ci.openshift.org/openshift/origin-v4.0:base
COPY --from=builder /go/src/sigs.k8s.io/cluster-api-provider-alibabacloud/bin/machine-controller-manager /
#COPY --from=builder /go/src/sigs.k8s.io/cluster-api-provider-alibabacloud/bin/termination-handler /
