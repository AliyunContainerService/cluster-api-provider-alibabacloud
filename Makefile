# Copyright 2021 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO111MODULE = on
export GO111MODULE
GOFLAGS ?= -mod=vendor
export GOFLAGS
GOPROXY ?=
export GOPROXY

DBG ?= 0

ifeq ($(DBG),1)
GOGCFLAGS ?= -gcflags=all="-N -l"
endif

<<<<<<< HEAD
<<<<<<< HEAD
GOARCH  ?= $(shell go env GOARCH)
GOOS    ?= $(shell go env GOOS)

VERSION     ?= $(shell git describe --tags --abbrev=7)
REPO_PATH   ?= github.com/AliyunContainerService/cluster-api-provider-alibabacloud
LD_FLAGS    ?= -X $(REPO_PATH)/pkg/version.Raw=$(VERSION) $(shell hack/version.sh)  -extldflags "-static"
=======
VERSION     ?= $(shell git describe --always --abbrev=7)
REPO_PATH   ?= github.com/AliyunContainerService/cluster-api-provider-alibabacloud
LD_FLAGS    ?= -X $(REPO_PATH)/pkg/version.Version=$(VERSION) -extldflags "-static"
>>>>>>> 8dbd34ff (update project name)
=======
GOARCH  ?= $(shell go env GOARCH)
GOOS    ?= $(shell go env GOOS)

VERSION     ?= $(shell git describe --always --abbrev=7)
REPO_PATH   ?= github.com/AliyunContainerService/cluster-api-provider-alibabacloud
LD_FLAGS    ?= -X $(REPO_PATH)/pkg/version.Raw=$(VERSION) -extldflags "-static"
>>>>>>> e879a141 (alibabacloud machine-api provider)
MUTABLE_TAG ?= latest
IMAGE        = origin-alibabacloud-machine-controllers

# race tests need CGO_ENABLED, everything else should have it disabled
CGO_ENABLED = 0
unit : CGO_ENABLED = 1

.PHONY: all
all: generate build images check

NO_DOCKER ?= 0

ifeq ($(shell command -v podman > /dev/null 2>&1 ; echo $$? ), 0)
	ENGINE=podman
else ifeq ($(shell command -v docker > /dev/null 2>&1 ; echo $$? ), 0)
	ENGINE=docker
else
	NO_DOCKER=1
endif

USE_DOCKER ?= 0
ifeq ($(USE_DOCKER), 1)
	ENGINE=docker
endif

ifeq ($(NO_DOCKER), 1)
  DOCKER_CMD =
  IMAGE_BUILD_CMD = imagebuilder
else
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
  DOCKER_CMD = $(ENGINE) run --rm -e CGO_ENABLED=$(CGO_ENABLED) -e GOARCH=$(GOARCH) -e GOOS=$(GOOS) -v "$(PWD)":/go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud:Z -w /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud openshift/origin-release:golang-1.15
  IMAGE_BUILD_CMD = $(ENGINE) build
=======
  DOCKER_CMD := docker run --rm -e CGO_ENABLED=1 -v "$(PWD)":/go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud:Z -w /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud openshift/origin-release:golang-1.12
  IMAGE_BUILD_CMD = docker build
>>>>>>> 8dbd34ff (update project name)
=======
  DOCKER_CMD = $(ENGINE) run --rm -e CGO_ENABLED=$(CGO_ENABLED) -e GOARCH=$(GOARCH) -e GOOS=$(GOOS) -v "$(PWD)":/go/src/sigs.k8s.io/cluster-api-provider-alibabacloud:Z -w /go/src/sigs.k8s.io/cluster-api-provider-alibabacloud openshift/origin-release:golang-1.15
=======
  DOCKER_CMD = $(ENGINE) run --rm -e CGO_ENABLED=$(CGO_ENABLED) -e GOARCH=$(GOARCH) -e GOOS=$(GOOS) -v "$(PWD)":/go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud:Z -w /go/src/github.com/AliyunContainerService/cluster-api-provider-alibabacloud openshift/origin-release:golang-1.15
>>>>>>> 7e2c5241 (remove test case)
  IMAGE_BUILD_CMD = $(ENGINE) build
>>>>>>> e879a141 (alibabacloud machine-api provider)
endif

.PHONY: vendor
vendor:
	$(DOCKER_CMD) hack/go-mod.sh
.PHONY: generate
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
generate: gogen goimports

gogen:
	$(DOCKER_CMD) go generate ./pkg/... ./cmd/...
<<<<<<< HEAD
=======
generate:
	go install $(GOGCFLAGS) -ldflags '-extldflags "-static"' github.com/AliyunContainerService/cluster-api-provider-alibabacloud/vendor/github.com/golang/mock/mockgen
	go generate ./pkg/... ./cmd/...
>>>>>>> 8dbd34ff (update project name)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)

.PHONY: test
test: ## Run tests
	@echo -e "\033[32mTesting...\033[0m"
	$(DOCKER_CMD) hack/ci-test.sh

bin:
	@mkdir $@



##@ Development
CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)


deepcopy: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt",year=2021 paths="./pkg/apis/alibabacloudprovider/v1beta1/."


.PHONY: build
build: ## build binaries
	$(DOCKER_CMD) CGO_ENABLED=0 go build $(GOGCFLAGS) -o "bin/machine-controller-manager" \
               -ldflags "$(LD_FLAGS)" "$(REPO_PATH)/cmd/manager"
<<<<<<< HEAD
<<<<<<< HEAD
=======
	$(DOCKER_CMD) go build $(GOGCFLAGS) -o bin/manager -ldflags '-extldflags "-static"' \
               "$(REPO_PATH)/vendor/github.com/openshift/cluster-api/cmd/manager"

alicloud-actuator:
	$(DOCKER_CMD) go build $(GOGCFLAGS) -o bin/alicloud-actuator github.com/AliyunContainerService/cluster-api-provider-alibabacloud/cmd/alicloud-actuator
>>>>>>> 8dbd34ff (update project name)
=======
	# $(DOCKER_CMD) CGO_ENABLED=0 go build  $(GOGCFLAGS) -o "bin/termination-handler" \
	#             -ldflags "$(LD_FLAGS)" "$(REPO_PATH)/cmd/termination-handler"
>>>>>>> e879a141 (alibabacloud machine-api provider)

.PHONY: images
images: ## Create images
ifeq ($(NO_DOCKER), 1)
	./hack/imagebuilder.sh
endif
	$(IMAGE_BUILD_CMD) -t "$(IMAGE):$(VERSION)" -t "$(IMAGE):$(MUTABLE_TAG)" ./

.PHONY: push
push:
	$(ENGINE) push "$(IMAGE):$(VERSION)"
	$(ENGINE) push "$(IMAGE):$(MUTABLE_TAG)"

.PHONY: check
check: fmt vet lint test # Check your code

.PHONY: unit
unit: # Run unit test
	$(DOCKER_CMD) go test -race -cover ./cmd/... ./pkg/...

<<<<<<< HEAD
<<<<<<< HEAD
.PHONY: test-e2e
test-e2e: ## Run e2e tests
	 hack/e2e.sh

.PHONY: lint
lint: ## Go lint your code
	$(DOCKER_CMD) hack/go-lint.sh -min_confidence 0.3 $$(go list -f '{{ .ImportPath }}' ./... | grep -v -e 'github.com/AliyunContainerService/cluster-api-provider-alibabacloud/test' -e 'github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/cloud/alibabacloud/client/mock')
=======
#.PHONY: integration
#integration: ## Run integration test
#	$(DOCKER_CMD) go test -v github.com/AliyunContainerService/cluster-api-provider-alibabacloud/test/integration

#.PHONY: build-e2e
#build-e2e:
#	go test -c -o bin/e2e.test github.com/AliyunContainerService/cluster-api-provider-alibabacloud/test/machines

#.PHONY: test-e2e
#test-e2e: ## Run e2e tests
#	hack/e2e.sh


#.PHONY: lint
#lint: ## Go lint your code
<<<<<<< HEAD
#	hack/go-lint.sh -min_confidence 0.3 $$(go list -f '{{ .ImportPath }}' ./... | grep -v -e 'sigs.k8s.io/cluster-api-provider-alicloud/test' -e 'sigs.k8s.io/cluster-api-provider-alicloud/pkg/cloud/alicloud/client/mock')
>>>>>>> 5ed2bd4c (format)
=======
#	hack/go-lint.sh -min_confidence 0.3 $$(go list -f '{{ .ImportPath }}' ./... | grep -v -e 'sigs.k8s.io/cluster-api-provider-alibabacloud/test' -e 'sigs.k8s.io/cluster-api-provider-alibabacloud/pkg/cloud/alicloud/client/mock')
>>>>>>> 8dbd34ff (update project name)
=======
.PHONY: test-e2e
test-e2e: ## Run e2e tests
	 hack/e2e.sh

.PHONY: lint
lint: ## Go lint your code
<<<<<<< HEAD
	$(DOCKER_CMD) hack/go-lint.sh -min_confidence 0.3 $$(go list -f '{{ .ImportPath }}' ./... | grep -v -e 'sigs.k8s.io/cluster-api-provider-alibabacloud/test' -e 'sigs.k8s.io/cluster-api-provider-alibabacloud/pkg/cloud/alibabacloud/client/mock')
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	$(DOCKER_CMD) hack/go-lint.sh -min_confidence 0.3 $$(go list -f '{{ .ImportPath }}' ./... | grep -v -e 'github.com/AliyunContainerService/cluster-api-provider-alibabacloud/test' -e 'github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/cloud/alibabacloud/client/mock')
>>>>>>> 7e2c5241 (remove test case)

.PHONY: fmt
fmt: ## Go fmt your code
	$(DOCKER_CMD) hack/go-fmt.sh .

.PHONY: goimports
goimports:
	$(DOCKER_CMD) hack/goimports.sh .
	hack/verify-diff.sh

<<<<<<< HEAD
<<<<<<< HEAD
.PHONY: vet
vet: ## Apply go vet to all go files
	$(DOCKER_CMD) hack/go-vet.sh ./...
=======
#.PHONY: vet
#vet: ## Apply go vet to all go files
#	hack/go-vet.sh ./...
>>>>>>> 5ed2bd4c (format)
=======
.PHONY: vet
vet: ## Apply go vet to all go files
	$(DOCKER_CMD) hack/go-vet.sh ./...
>>>>>>> e879a141 (alibabacloud machine-api provider)

.PHONY: help
help:
	@grep -E '^[a-zA-Z/0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'