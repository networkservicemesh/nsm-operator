# Network Service Mesh Operator

[![Go Report Card](https://goreportcard.com/badge/github.com/acmenezes/nsm-operator "Go Report Card")](https://goreportcard.com/report/github.com/acmenezes/nsm-operator)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](code-of-conduct.md) 

A Kubernetes Operator for Installing and Managing Network Service Meshes

## Overview

The Network Service Mesh Operator is a tool to install and manage the [Network Service Mesh][nsm_home] application which <em> is a novel approach solving complicated L2/L3 use cases in Kubernetes that are tricky to address with the existing Kubernetes Network Model. Inspired by Istio, Network Service Mesh maps the concept of a Service Mesh to L2/L3 payloads as part of an attempt to re-imagine NFV in a Cloud-native way! </em>. To  better understand the network service meshes take a look at [what is nsm][nsm_whatis].

The operator is a single pod workload that automates operational human knowledge behind the scenes to create the service mesh infrastructure components deploying a webhook and daemonsets with the network service managers and forwarding plane workloads taking the configuration from the Custom Resource manifest created specifically to be used with the operator. It aims to be platform independed and for such should run well in any kubernetes distribution.

Some of the features intended to be embedded with the operator are

* Installing, configuring and making cleanups if needed

* Upgrading, backing up and restoring any important state information

* Expose and aggregate metrics from all components through prometheus/grafana

* Operations analytics based on metrics exposed   

* Auto-pilot functions such as distributing NSM registry into multiple pods according to the size of the cluster among other functions that may be addressed as well via automation.

## Install

#### Getting Started

The network service mesh operator is supported for use with kubernetes 1.14 or above and Openshift 4 or above. It can be installed automatically using the Operator Lifecycle Manager or it can be installed manually through standard kubernetes yaml manifests.

#### Operator Lifecycle Manager

The Operator Lifecycle Manager, or OLM, is the preferred method. It's simple, faster and cleaner. <em> OLM extends Kubernetes to provide a declarative way to install, manage, and upgrade Operators and their dependencies in a cluster. </em> It manages the available operators using catalogs, it takes care of automatic updates and dependencies, has a discover mechanism to advertise services provided by available operators, prevents conflicts if two operators try to use the same API, for example, and provides a nice way to build declarative UI controls to configure operator services. OLM comes installed by default in Openshift 4.

If want to install OLM in a kubernetes cluster you can try the [install guide][olm_install_guide].

#### Install Methods:

[Openshift Embedded Operator Hub][openshift_olm_install]
[Openshift Manual Install][openshift_manual_install]
[Kubernetes OLM Install][k8s_olm_install]
[Kubernetes Manual Install][k8s_manual_Install]


## Usage 

To create a new NSM custom resource, after deploying the operator itself, use the CR manifest below. Here we have an example with the Vector Packet Processing as a forwarding plane for nsm.

```
kubectl apply -f deploy/crds/nsm.networkservicemesh.io_v1alpha1_nsm_cr.yaml
```

## Contributing

Anyone interested in contributing to the nsm-operator is welcomed and 
should start by reviewing the [development][docs_dev] documentation.

## License

nsm-operator is released under the Apache 2.0 license. Please check the [LICENSE][license_file] file for details.

[nsm_home]:https://networkservicemesh.io
[nsm_whatis]:https://github.com/networkservicemesh/networkservicemesh/blob/master/docs/what-is-nsm.md
[docs_dev]:./docs/development.md
[license_file]:./LICENSE
[requirements]:./docs/requirements.md
[olm_install_guide]:https://github.com/operator-framework/operator-lifecycle-manager/blob/master/doc/install/install.md
[openshift_olm_install]:./docs/openshift_olm_install.md
[openshift_manual_install]:./docs/openshift_manual_install.md
[k8s_olm_install]:./docs/k8s_olm_install.yaml
[k8s_manual_install]:./docs/k8s_manual_install.yaml
