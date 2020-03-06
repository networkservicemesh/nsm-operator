# Copyright 2020 networkservicemesh.io
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Operator Variables
OPERATOR_IMAGE = quay.io/acmenezes/nsm-operator
OPERATOR_TAG = v0.0.99

# Base bundle build dir
BUNDLE_DOCKERFILE_PATH = deploy/bundle

# Kubernetes Bundle Variables
BUNDLE_K8S_OUTPUT_PATH = build/_output/kubernetes/bundle
BUNDLE_K8S_MANIFESTS_PATH = deploy/olm-catalog/kubernetes/nsm-operator
BUNDLE_K8S_IMAGE = quay.io/acmenezes/nsm-operator-k8s-registry
BUNDLE_K8S_TAG = v0.0.1

# Openshift Bundle Variables
BUNDLE_OSP_OUTPUT_PATH = build/_output/openshift/bundle
BUNDLE_OSP_MANIFESTS_PATH = deploy/olm-catalog/openshift/nsm-operator
BUNDLE_OSP_IMAGE = quay.io/acmenezes/nsm-operator-osp-registry
BUNDLE_OSP_TAG = v0.0.1


# bundle build
.phony: bundle-build
bundle-build:
	@echo "Cleaning up the build path"
	@echo ""
	rm -rf build/_output/kubernetes
	rm -rf build/_output/openshift

	@echo "building kubernetes bundle image..."
	@echo ""
	mkdir -p ${BUNDLE_K8S_OUTPUT_PATH}/manifests
	cp -r ${BUNDLE_DOCKERFILE_PATH}/* ${BUNDLE_K8S_OUTPUT_PATH}
	cp -r ${BUNDLE_K8S_MANIFESTS_PATH} ${BUNDLE_K8S_OUTPUT_PATH}/manifests
	docker build -t ${BUNDLE_K8S_IMAGE}:${BUNDLE_K8S_TAG} ${BUNDLE_K8S_OUTPUT_PATH}

	@echo "building OpenShift bundle image..."
	@echo ""
	mkdir -p ${BUNDLE_OSP_OUTPUT_PATH}/manifests
	cp -r ${BUNDLE_DOCKERFILE_PATH}/* ${BUNDLE_OSP_OUTPUT_PATH}
	cp -r ${BUNDLE_OSP_MANIFESTS_PATH} ${BUNDLE_OSP_OUTPUT_PATH}/manifests
	docker build -t ${BUNDLE_OSP_IMAGE}:${BUNDLE_OSP_TAG} ${BUNDLE_OSP_OUTPUT_PATH}


# bundle push
.phony: bundle-push
bundle-push:
	@echo "pushing kubernetes bundle to "${BUNDLE_K8S_IMAGE}":"${BUNDLE_K8S_TAG}
	@echo ""
	docker login quay.io/acmenezes
	docker push ${BUNDLE_K8S_IMAGE}":"${BUNDLE_K8S_TAG}

	@echo "pushing OpenShift bundle to "${BUNDLE_OSP_IMAGE}":"${BUNDLE_OSP_TAG}
	@echo ""
	docker login quay.io/acmenezes
	docker push ${BUNDLE_OSP_IMAGE}":"${BUNDLE_OSP_TAG}

# generate k8s crds 
.phony: gen-k8s-crds
gen-k8s-crds:
	@echo "Generating Kubernetes Code for custom resource..."
	@echo ""
	operator-sdk generate k8s
	@echo "Generating/updating CRDs for the APIs"
	@echo ""
	operator-sdk generate crds

# operator build
.phony: operator-build
operator-build:
	@echo "Building operator container image..."
	@echo ""
	operator-sdk build ${OPERATOR_IMAGE}":"${OPERATOR_TAG}

# operator push
.phony: operator-push
operator-push:
	@echo "Pushing operator container image to ${OPERATOR_IMAGE}"
	@echo ""
	docker login quay.io/acmenezes
	docker push ${OPERATOR_IMAGE}":"${OPERATOR_TAG}



#TODO:


# install dependencies for make

# probably from python scripts:
# test openshift manual
# test openshift olm
# test kubernetes manual
# test kubernetes olm
# test all

# openapi validation

# generate CSVs + tags

# update readthedocs
