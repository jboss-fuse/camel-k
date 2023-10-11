# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

VERSIONFILE := pkg/util/defaults/defaults.go
VERSION ?= 1.10.3
LAST_RELEASED_IMAGE_NAME := red-hat-camel-k-operator
LAST_RELEASED_VERSION := 1.10.2
RUNTIME_VERSION ?= 1.15.2
BUILDAH_VERSION := 1.14.0
KANIKO_VERSION := 0.17.1
INSTALL_DEFAULT_KAMELETS := true
CONTROLLER_GEN_VERSION := v0.4.1
OPERATOR_SDK_VERSION := v1.28.0
KUSTOMIZE_VERSION := v4.1.2
BASE_IMAGE ?= adoptopenjdk/openjdk11:slim
LOCAL_REPOSITORY := /tmp/artifacts/m2
IMAGE_NAME := registry-proxy.engineering.redhat.com/rh-osbs/integration-camel-k-rhel8-operator

#
# Situations when user wants to override
# the image name and version
# - used in kustomize install
# - used in making bundle
# - need to preserve original image and version as used in other files
#
CUSTOM_IMAGE ?= $(IMAGE_NAME)
CUSTOM_VERSION ?= 1.10.3
OPERATOR_VERSION := $(subst -SNAPSHOT,,$(CUSTOM_VERSION))
CUSTOM_TAG ?= $(OPERATOR_VERSION)

METADATA_IMAGE_NAME := $(CUSTOM_IMAGE)-metadata
BUNDLE_IMAGE_NAME ?= $(CUSTOM_IMAGE)-bundle
GIT_COMMIT := $(shell if [ -d .git ]; then git rev-list -1 HEAD; else echo "$(CUSTOM_VERSION)"; fi)
RELEASE_GIT_REMOTE := upstream
LINT_GOGC := 10
LINT_DEADLINE := 10m

# olm bundle vars
MANAGER := config/manager
MANIFESTS := config/manifests
BUNDLE_DIR := bundle
CHANNELS ?= $(shell v=$(OPERATOR_VERSION) && echo "$${v%\.[0-9]}.x"),candidate,latest
DEFAULT_CHANNEL ?= $(shell v=$(OPERATOR_VERSION) && echo "$${v%\.[0-9]}.x")
PACKAGE := red-hat-camel-k
CSV_VERSION := $(CUSTOM_VERSION:-SNAPSHOT=)
CSV_NAME := $(PACKAGE).v$(CSV_VERSION)
CSV_PRODUCT_NAME := $(PACKAGE)-operator.v$(CSV_VERSION)
CSV_DISPLAY_NAME := Red Hat Integration - Camel K
CSV_SUPPORT := Camel
CSV_REPLACES := $(LAST_RELEASED_IMAGE_NAME).v$(LAST_RELEASED_VERSION)
CSV_SKIP_RANGE := ">=1.6.7 <$(CSV_VERSION)"
CSV_FILENAME := $(PACKAGE).clusterserviceversion.yaml
CSV_PATH := $(MANIFESTS)/bases/$(CSV_FILENAME)
DEFAULT_CSV := $(MANIFESTS)/bases/camel-k.clusterserviceversion.yaml
ADD_REDHAT_LABELS_PATCH := $(MANAGER)/patch-add-redhat-labels.yaml

# Used to push pre-release artifacts
STAGING_IMAGE_NAME := docker.io/camelk/camel-k

STAGING_RUNTIME_REPO :=

# Kamelets options
INSTALL_DEFAULT_KAMELETS ?= true
KAMELET_CATALOG_REPO := https://github.com/apache/camel-kamelets.git
# Optional branch for the default Kamelet catalog (change this to a tag before release)
KAMELET_CATALOG_REPO_TAG := v0.7.1

# When packaging artifacts into the docker image, you can "copy" them from local maven
# or "download" them from Apache Snapshots and Maven Central
PACKAGE_ARTIFACTS_STRATEGY := copy

# When performing integration tests, it is not necessary to always execute build, especially
# in e2e tests when lots of tests are being executed sequentially & the build has already taken place.
# In this case, if DO_TEST_PREBUILD is set to false then TEST_PREBUILD will be set to blank
# and those rules with a dependency on the build rule will no longer call it (see do-build)
DO_TEST_PREBUILD ?= true
TEST_PREBUILD = build

# OLM (Operator Lifecycle Manager and Operator Hub): uncomment to override operator settings at build time
GOLDFLAGS += -X 'github.com/apache/camel-k/pkg/cmd.VersionVariant=Red Hat'
GOLDFLAGS += -X 'github.com/apache/camel-k/pkg/install.KamelCLIDownloadURLTemplate=https://mirror.openshift.com/pub/openshift-v4/clients/camel-k/%s/camel-k-client-%s-%s-64bit.tar.gz'
GOLDFLAGS += -X 'github.com/apache/camel-k/pkg/util/olm.DefaultOperatorName=$(LAST_RELEASED_IMAGE_NAME)'
GOLDFLAGS += -X 'github.com/apache/camel-k/pkg/util/olm.DefaultPackage=$(PACKAGE)'
GOLDFLAGS += -X 'github.com/apache/camel-k/pkg/util/olm.DefaultChannel=$(DEFAULT_CHANNEL)'
GOLDFLAGS += -X 'github.com/apache/camel-k/pkg/util/olm.DefaultSource=redhat-operators'
GOLDFLAGS += -X 'github.com/apache/camel-k/pkg/util/maven.DefaultMavenRepositories=https://maven.repository.redhat.com/ga@id=redhat.ga,https://repository.jboss.org/nexus/content/groups/ea@id=jboss.ea,https://maven.repository.redhat.com/earlyaccess/all@id=redhat.ea,https://repo.maven.apache.org/maven2@id=central'
GOLDFLAGS += -X 'github.com/apache/camel-k/pkg/kamelet/repository.DefaultRemoteRepository=none'
GOLDFLAGS += -X 'github.com/apache/camel-k/pkg/install.KamelCLIDownloadName=redhat-kamel-cli'
GOLDFLAGS += -X 'github.com/apache/camel-k/pkg/install.KamelCLIDownloadDisplayName=kamel - Red Hat Integration - Camel K - Command Line Interface'
GOLDFLAGS += -X 'github.com/apache/camel-k/pkg/install.KamelCLIDownloadDescription=CLI for Red Hat Integration - Camel K. The \"kamel\" binary can be used to both configure the cluster and run integrations. Once you have downloaded the \"kamel\" binary, log into the cluster using the \"oc\" client tool and start using the \"kamel\" CLI.'

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Build
ifdef GIT_COMMIT
GOLDFLAGS += -X github.com/apache/camel-k/pkg/util/defaults.GitCommit=$(GIT_COMMIT)
else
$(warning Could not retrieve a valid Git Commit)
endif

GOFLAGS = -ldflags "$(GOLDFLAGS)" -trimpath

define LICENSE_HEADER
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
endef

export LICENSE_HEADER
default: build

# Generates the version file
codegen:
	$(info Regenerating $(VERSIONFILE))
	@echo "/*" > $(VERSIONFILE)
	@echo "$$LICENSE_HEADER" >> $(VERSIONFILE)
	@echo "*/" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)
	@echo "package defaults" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)
	@echo "// ***********************" >> $(VERSIONFILE)
	@echo "//  DO NOT EDIT THIS FILE"  >> $(VERSIONFILE)
	@echo "// ***********************" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)
	@echo "const (" >> $(VERSIONFILE)
	@echo "  // Version -- " >> $(VERSIONFILE)
	@echo "  Version = \"$(CUSTOM_VERSION)\"" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)
	@echo "  // DefaultRuntimeVersion -- " >> $(VERSIONFILE)
	@echo "  DefaultRuntimeVersion = \"$(RUNTIME_VERSION)\"" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)
	@echo "  // BuildahVersion -- " >> $(VERSIONFILE)
	@echo "  BuildahVersion = \"$(BUILDAH_VERSION)\"" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)
	@echo "  // KanikoVersion -- " >> $(VERSIONFILE)
	@echo "  KanikoVersion = \"$(KANIKO_VERSION)\"" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)
	@echo "  // baseImage -- " >> $(VERSIONFILE)
	@echo "  baseImage = \"$(BASE_IMAGE)\"" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)
	@echo "  // LocalRepository -- " >> $(VERSIONFILE)
	@echo "  LocalRepository = \"$(LOCAL_REPOSITORY)\"" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)
	@echo "  // ImageName -- " >> $(VERSIONFILE)
	@echo "  ImageName = \"$(CUSTOM_IMAGE)\"" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)
	@echo "  // installDefaultKamelets -- " >> $(VERSIONFILE)
	@echo "  installDefaultKamelets = $(INSTALL_DEFAULT_KAMELETS)" >> $(VERSIONFILE)
	@echo ")" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)
	@echo "//GitCommit must be provided during application build" >> $(VERSIONFILE)
	@echo "var GitCommit string" >> $(VERSIONFILE)
	@echo "" >> $(VERSIONFILE)

	gofmt -w pkg/util/defaults/defaults.go

generate: generate-deepcopy generate-crd generate-client generate-doc generate-keda generate-strimzi

generate-client:
	./script/gen_client.sh

generate-crd: controller-gen
	CONTROLLER_GEN=$(CONTROLLER_GEN) ./script/gen_crd.sh

generate-doc:
	./script/gen_doc.sh

generate-deepcopy: controller-gen
	cd pkg/apis/camel && $(CONTROLLER_GEN) paths="./..." object

generate-keda:
	cd addons/keda/duck && $(CONTROLLER_GEN) paths="./..." object

generate-strimzi:
	cd addons/strimzi/duck && $(CONTROLLER_GEN) paths="./..." object
	./script/gen_client_strimzi.sh

build: codegen build-resources build-kamel build-compile-integration-tests build-submodules

ci-build: clean codegen set-module-version set-version check-licenses dir-licenses build-kamel cross-compile

do-build: gotestfmt-install
ifeq ($(DO_TEST_PREBUILD),true)
TEST_PREBUILD = build
else
TEST_PREBUILD =
endif

gotestfmt-install:
ifeq (, $(shell command -v gotestfmt 2> /dev/null))
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
endif

test: do-build
	@echo "####### Running unit test..."
	go test ./...
	cd pkg/apis/camel && go test ./...
	cd pkg/client/camel && go test ./...
	cd pkg/kamelet/repository && go test ./...

test-integration: do-build
	FAILED=0; STAGING_RUNTIME_REPO="$(STAGING_RUNTIME_REPO)"; \
	go test -timeout 60m -v ./e2e/global/common -tags=integration $(TEST_INTEGRATION_COMMON_RUN) -json 2>&1 | gotestfmt || FAILED=1; \
	go test -timeout 60m -v ./e2e/global/common/build -tags=integration $(TEST_INTEGRATION_COMMON_BUILD_RUN) -json 2>&1 | gotestfmt || FAILED=1; \
	go test -timeout 60m -v ./e2e/global/common/config -tags=integration $(TEST_INTEGRATION_COMMON_CONFIG_RUN) -json 2>&1 | gotestfmt || FAILED=1; \
	go test -timeout 60m -v ./e2e/global/common/languages -tags=integration $(TEST_INTEGRATION_COMMON_LANG_RUN) -json 2>&1 | gotestfmt || FAILED=1; \
	go test -timeout 60m -v ./e2e/global/common/traits -tags=integration $(TEST_INTEGRATION_COMMON_TRAITS_RUN) -json 2>&1 | gotestfmt || FAILED=1; \
	exit $${FAILED}

test-knative: do-build
	STAGING_RUNTIME_REPO="$(STAGING_RUNTIME_REPO)" \
	go test -timeout 60m -v ./e2e/global/knative -tags=integration $(TEST_KNATIVE_RUN) -json 2>&1 | gotestfmt

test-builder: do-build
	STAGING_RUNTIME_REPO="$(STAGING_RUNTIME_REPO)" \
	go test -timeout 60m -v ./e2e/global/builder -tags=integration $(TEST_BUILDER_RUN) -json 2>&1 | gotestfmt

test-service-binding: do-build
	STAGING_RUNTIME_REPO="$(STAGING_RUNTIME_REPO)" \
	go test -timeout 60m -v ./e2e/global/service-binding -tags=integration $(TEST_SERVICE_RUN) -json 2>&1 | gotestfmt

test-local: do-build
	STAGING_RUNTIME_REPO="$(STAGING_RUNTIME_REPO)" \
	RUNTIME_VERSION="$(RUNTIME_VERSION)" \
	go test -timeout 30m -v ./e2e/local -tags=integration $(TEST_LOCAL_RUN) -json 2>&1 | gotestfmt

test-registry-maven-wagon: do-build
	STAGING_RUNTIME_REPO="$(STAGING_RUNTIME_REPO)" \
	go test -timeout 60m -v ./e2e/global/registry -tags=integration $(TEST_REGISTRY_MAVEN_WAGON_RUN) -json 2>&1 | gotestfmt

###############################
#
# NAMESPACED INSTALL TESTS
#
###############################

test-install: do-build
	FAILED=0; STAGING_RUNTIME_REPO="$(STAGING_RUNTIME_REPO)"; \
	go test -timeout 60m -v ./e2e/namespace/install/ -tags=integration $(TEST_INSTALL_RUN) -json 2>&1 | gotestfmt || FAILED=1; \
	go test -timeout 60m -v ./e2e/namespace/install/cli -tags=integration $(TEST_INSTALL_CLI_RUN) -json 2>&1 | gotestfmt || FAILED=1; \
	go test -timeout 60m -v ./e2e/namespace/install/kustomize -tags=integration $(TEST_INSTALL_KUSTOMIZE_RUN) -json 2>&1 | gotestfmt || FAILED=1; \
	exit $${FAILED}

test-quarkus-native: do-build
	STAGING_RUNTIME_REPO="$(STAGING_RUNTIME_REPO)" \
	go test -timeout 130m -v ./e2e/namespace/native -tags=integration $(TEST_QUARKUS_RUN) -json 2>&1 | gotestfmt

test-upgrade: do-build
	STAGING_RUNTIME_REPO="$(STAGING_RUNTIME_REPO)" \
	go test -timeout 60m -v ./e2e/namespace/upgrade -tags=integration $(TEST_UPGRADE_RUN) -json 2>&1 | gotestfmt
build-kamel:
# Ensure the binary is statically linked when building on Linux due to ABI changes in newer glibc 2.32, otherwise
# it would not run on older versions. See https://github.com/apache/camel-k/pull/2141#issuecomment-800990117
ifeq ($(shell uname -s 2>/dev/null || echo Unknown),Linux)
	CGO_ENABLED=0 go build $(GOFLAGS) -o kamel ./cmd/kamel/*.go
else
	go build $(GOFLAGS) -o kamel ./cmd/kamel/*.go
endif

build-resources:
	./script/get_catalog.sh $(RUNTIME_VERSION) $(STAGING_RUNTIME_REPO)
	go generate ./pkg/...

bundle-kamelets:
ifneq (,$(findstring release,$(MAKECMDGOALS)))
ifneq (,$(findstring $(KAMELET_CATALOG_REPO_TAG), main))
	@echo "You cannot set KAMELET_CATALOG_REPO_TAG=$(KAMELET_CATALOG_REPO_TAG) when doing a release"
	@exit 1
endif
endif
	./script/bundle_kamelets.sh $(KAMELET_CATALOG_REPO) $(KAMELET_CATALOG_REPO_TAG)

build-submodules:
	./script/build_submodules.sh

build-compile-integration-tests:
	go test -run nope -tags="integration" ./e2e/global/... > /dev/null
	go test -run nope -tags="integration" ./e2e/local/... > /dev/null
	go test -run nope -tags="integration" ./e2e/namespace/... > /dev/null

clean:
	# disable gomodules when executing go clean:
	#
	#    https://github.com/golang/go/issues/31002
	#
	GO111MODULE=off go clean
	rm -f camel-k
	rm -f kamel
	rm -f *.test
	rm -rf build/_maven_output
	rm -rf build/_maven_overlay
	rm -rf build/_output
	rm -rf camel-k-client-*.tar.gz
	rm -rf camel-k-examples-*.tar.gz
	rm -rf bundle
ifneq ($(CSV_PATH), $(DEFAULT_CSV))
ifneq ("$(wildcard $(CSV_PATH))","")
	mv $(CSV_PATH) $(DEFAULT_CSV)
endif
endif

dep:
	go mod tidy

lint:
	GOGC=$(LINT_GOGC) golangci-lint run --out-format tab --deadline $(LINT_DEADLINE)

lint-fix:
	GOGC=$(LINT_GOGC) golangci-lint run --out-format tab --deadline $(LINT_DEADLINE) --fix

dir-licenses:
	./script/vendor-license-directory.sh

check-licenses:
	./script/check_licenses.sh

check-platform:
	./script/check_platform.sh

# The below are dependencies needed for maven structured logs. We must bundle into the final container image.
maven-overlay:
	@echo "####### Preparing maven dependencies bundle..."
	mkdir -p build/_maven_overlay
	./script/maven_overlay.sh -s "$(STAGING_RUNTIME_REPO)" -d "$(CAMEL_K_RUNTIME_DIR)" $(RUNTIME_VERSION) build/_maven_overlay

kamel-overlay:
	@echo "####### Copying kamel CLI to output build directory..."
	mkdir -p build/_output/bin
ifeq ($(shell uname -s 2>/dev/null || echo Unknown),Linux)
	cp kamel build/_output/bin
else ifeq ($(shell uname -m), arm64)
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o build/_output/bin/kamel ./cmd/kamel/*.go
else
	GOOS=linux go build $(GOFLAGS) -o build/_output/bin/kamel ./cmd/kamel/*.go
endif

images: build kamel-overlay maven-overlay bundle-kamelets
ifneq (,$(findstring SNAPSHOT,$(RUNTIME_VERSION)))
	./script/package_maven_artifacts.sh -s "$(STAGING_RUNTIME_REPO)" -d "$(CAMEL_K_RUNTIME_DIR)" $(RUNTIME_VERSION)
endif
	@echo "####### Building Camel K operator container image..."
	mkdir -p build/_maven_output
	docker build -t $(CUSTOM_IMAGE):$(CUSTOM_VERSION) -f build/Dockerfile .

images-arch: build kamel-overlay maven-overlay bundle-kamelets
	@echo "####### Building Camel K operator container image for multi architectures..."
	mkdir -p build/_maven_output
	docker buildx rm --all-inactive --force
	docker buildx create --append --name builder
ifeq ($(shell uname -m), x86_x64)
	docker buildx build --platform=linux/amd64 -t $(CUSTOM_IMAGE):$(CUSTOM_VERSION) -f build/Dockerfile.arch .
endif
ifeq ($(shell uname -m), aarch64)
	docker buildx build --platform=linux/arm64 -t $(CUSTOM_IMAGE):$(CUSTOM_VERSION) -f build/Dockerfile.arch .
endif
ifeq ($(shell uname -m), arm64)
	docker buildx build --platform=linux/amd64 -t $(CUSTOM_IMAGE):$(CUSTOM_VERSION) -f build/Dockerfile.arch .
endif

images-dev: build kamel-overlay maven-overlay bundle-kamelets package-artifacts
	@echo "####### Building Camel K operator container development image..."
	docker build -t $(CUSTOM_IMAGE):$(CUSTOM_VERSION) -f build/Dockerfile .

images-push:
	docker push $(CUSTOM_IMAGE):$(CUSTOM_VERSION)

images-push-staging:
	docker tag $(CUSTOM_IMAGE):$(CUSTOM_VERSION) $(STAGING_IMAGE_NAME):$(CUSTOM_VERSION)
	docker push $(STAGING_IMAGE_NAME):$(CUSTOM_VERSION)

get-version:
	@echo $(CUSTOM_VERSION)

get-last-released-version:
	@echo $(LAST_RELEASED_VERSION)

set-version:
	./script/set_version.sh $(CUSTOM_VERSION) $(CUSTOM_IMAGE)

set-module-version:
	./script/set_go_modules_version.sh $(CUSTOM_VERSION)

git-tag:
	./script/git_tag.sh $(CUSTOM_VERSION) $(RELEASE_GIT_REMOTE)

cross-compile:
	# we must wrap the goldflags parameters with quotes as they will need to
	# be processed as a single argument by the cross compile script
	./script/cross_compile.sh $(CUSTOM_VERSION) $(subst ","\",$(GOFLAGS))

package-examples:
	./script/package_examples.sh $(CUSTOM_VERSION)

package-artifacts:
	# useful for development purpose only in order to include any local development runtime dependency
	./script/package_maven_artifacts.sh -s "$(STAGING_RUNTIME_REPO)" -d "$(CAMEL_K_RUNTIME_DIR)" $(RUNTIME_VERSION)

release: clean codegen set-module-version set-version build-resources check-licenses build images images-push cross-compile package-examples git-tag

release-staging: clean codegen set-module-version set-version build-resources check-licenses build images images-push-staging cross-compile package-examples git-tag

release-nightly: clean codegen set-module-version set-version build-resources build images images-push cross-compile package-examples

release-helm:
	./script/release_helm.sh

install-crc:
	./script/install_crc.sh $(CUSTOM_VERSION)
install-minikube:
	./script/install_minikube.sh

get-staging-repo:
	@echo $(or ${STAGING_RUNTIME_REPO},https://repository.apache.org/content/repositories/snapshots@id=apache-snapshots@snapshots)

.PHONY: build build-kamel build-resources dep codegen images images-dev images-push images-push-staging test check test-integration clean release cross-compile package-examples set-version git-tag release-notes check-licenses generate-deepcopy generate-client generate-doc build-resources release-helm release-staging release-nightly get-staging-repo get-version build-submodules set-module-version bundle-kamelets generate-keda generate-strimzi

# find or download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION) ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

kubectl:
ifeq (, $(shell which kubectl))
	$(error "No kubectl found in PATH. Please install and re-run")
endif

kustomize:
ifeq (, $(shell which kustomize))
	@{ \
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v4@$(KUSTOMIZE_VERSION) ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
KUSTOMIZE=$(GOBIN)/kustomize
else
KUSTOMIZE=$(shell which kustomize)
endif

operator-sdk:
ifeq (, $(shell which operator-sdk))
	@{ \
	set -e ;\
	if [ "$(shell uname -s 2>/dev/null || echo Unknown)" == "Darwin" ] ; then \
		curl \
			-L https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_VERSION)/operator-sdk_darwin_amd64 \
			-o operator-sdk ; \
	else \
		curl \
			-L https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_VERSION)/operator-sdk_linux_amd64 \
			-o operator-sdk ; \
	fi ;\
	chmod +x operator-sdk ;\
	mv operator-sdk $(GOBIN)/ ;\
	}
OPERATOR_SDK=$(GOBIN)/operator-sdk
else
	@{ \
	echo -n "operator-sdk already installed: "; \
  operator-sdk version | sed -n 's/.*"v\([^"]*\)".*/\1/p'; \
	echo " If this is less than $(OPERATOR_SDK_VERSION) then please consider moving it aside and allowing the approved version to be downloaded."; \
	}
OPERATOR_SDK=$(shell which operator-sdk)
endif

.PHONY: generate-crd $(BUNDLE_CAMEL_APIS) bundle bundle-build

# - Have to copy pkg/apis since it is a module in its own right
# - The dependency of v1alpha1 api on v1 results in duplication
#   of the v1 schema in those CRDs so need to split them up
# - No way to add a licence header to the CRD other then post-script
BUNDLE_CAMEL_APIS = v1 v1alpha1

$(BUNDLE_CAMEL_APIS): operator-sdk
	@# Clean up api directory, copy over fresh version
	rm -rf api_$@ && cp -rf pkg/apis api_$@
	@# Remove the camel directory and re-copy only the required api
	rm -rf api_$@/camel/* && cp -rf pkg/apis/camel/$@ api_$@/camel/
	@# operator-sdk generate ... cannot execute across separate modules so need to temporarily move api
	$(OPERATOR_SDK) generate kustomize manifests --apis-dir $(addprefix api_, $@) -q
	@# Adds the licence header to the csv file.
	./script/add_license.sh config/manifests/bases ./script/headers/yaml.txt
	./script/add_createdAt.sh config/manifests/bases
	@# Clean up temporary working api directories
	rm -rf api_*

# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
ifneq ($(origin PACKAGE), undefined)
BUNDLE_PACKAGE := --package=$(PACKAGE)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL) $(BUNDLE_PACKAGE)

#
# Tailor the manifest according to default values for this project
# Note. to successfully make the bundle the name must match that specified in the PROJECT file
#
pre-bundle:
# bundle name must match that which appears in PROJECT file
	@sed -i 's/projectName: .*/projectName: $(PACKAGE)/' PROJECT
# finds the single CSV file and renames it
	@find $(MANIFESTS)/bases -type f -name "*.clusterserviceversion.yaml" -execdir mv '{}' $(CSV_FILENAME) ';'
	@sed -i 's~^    containerImage: .*~    containerImage: $(CUSTOM_IMAGE):$(CUSTOM_TAG)~' $(CSV_PATH)
	@sed -i 's/^  name: .*.\(v.*\)/  name: $(CSV_NAME)/' $(CSV_PATH)
	@sed -i 's/^  displayName: .*/  displayName: $(CSV_DISPLAY_NAME)/' $(CSV_PATH)
	@sed -i 's/^  replaces: .*/  replaces: $(CSV_REPLACES)/' $(CSV_PATH)
	@sed -i 's/^  version: .*/  version: $(CSV_VERSION)/' $(CSV_PATH)

bundle: set-version generate-crd kustomize operator-sdk pre-bundle $(BUNDLE_CAMEL_APIS)
	@# Sets the operator image to the preferred image:tag
	@cd config/manifests && $(KUSTOMIZE) edit set image $(IMAGE_NAME)=$(CUSTOM_IMAGE):$(CUSTOM_VERSION)
	@# Build kustomize manifests
	@$(KUSTOMIZE) build config/manifests | \
		$(OPERATOR_SDK) generate bundle \
			-q --overwrite --version $(OPERATOR_VERSION) \
			--kustomize-dir config/manifests $(BUNDLE_METADATA_OPTS)
	@# Move the dockerfile into the bundle directory
ifeq ($(shell uname -s 2>/dev/null || echo Unknown),Darwin)
	@mv bundle.Dockerfile bundle/Dockerfile && sed -i '' 's/bundle\///g' bundle/Dockerfile
else
	@mv bundle.Dockerfile bundle/Dockerfile && sed -i 's/bundle\///g' bundle/Dockerfile
endif
	@# Adds the licence headers to the csv file
	./script/add_license.sh bundle/manifests ./script/headers/yaml.txt
	$(OPERATOR_SDK) bundle validate ./bundle
# operator-sdk requires the name of the operator to be the PACKAGE
# However, the historical name of the operator has the suffix 'operator' so this should
# be added once the validation has been completed
	@sed -i 's/  name: $(CSV_NAME)/  name: $(CSV_PRODUCT_NAME)/' $(BUNDLE_DIR)/manifests/$(CSV_FILENAME)
	@# Execute post bundle-generation modifications
	@CSV_VERSION=$(CSV_VERSION) OPERATOR_VERSION=$(CUSTOM_TAG) ./script/post_bundle_gen.sh

# Build the bundle image.
bundle-build: bundle
	cd bundle && docker build -f Dockerfile -t $(BUNDLE_IMAGE_NAME) .
