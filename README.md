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

## Requirements

In order to have NSM working check the minimal requirements [here][requirements].

## Install

At this point to install the operator it's enough to apply the manifests:
Once we get ready with OLM, everything will be installed by it.

```
git clone git@github.com:acmenezes/nsm-operator.git
cd nsm-operator
kubectl apply -f deploy/crds/nsm.networkservicemesh.io_nsms_crd.yaml
kubectl apply -f deploy/operator_resources.yaml
```
## Usage 

To create a new NSM custom resource, after deploying the operator itself, use the CR manifest below.
Here we have an example with the Vector Packet Processing as a forwarding plane for nsm.

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