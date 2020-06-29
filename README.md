## Network Service Mesh Operator

[![Go Report Card](https://goreportcard.com/badge/github.com/networkservicemesh/nsm-operator "Go Report Card")](https://goreportcard.com/report/github.com/networkservicemesh/nsm-operator)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](code-of-conduct.md) 

A Kubernetes Operator for Installing and Managing Network Service Meshes

#### Overview

The Network Service Mesh Operator is a tool to install and manage the [Network Service Mesh][nsm_home] application which <em> is a novel approach solving complicated L2/L3 use cases in Kubernetes that are tricky to address with the existing Kubernetes Network Model. Inspired by Istio, Network Service Mesh maps the concept of a Service Mesh to L2/L3 payloads as part of an attempt to re-imagine NFV in a Cloud-native way! </em>. To  better understand the network service meshes take a look at [what is nsm][nsm_whatis].

The operator is a single pod workload that automates operational human knowledge behind the scenes to create the service mesh infrastructure components deploying a webhook and daemonsets with the network service managers and forwarding plane workloads taking the configuration from the Custom Resource manifest created specifically to be used with the operator. It aims to be platform independed and for such should run well in any kubernetes distribution.

Some of the features intended to be embedded with the operator are

* Installing, configuring and making cleanups if needed

* Upgrading, backing up and restoring any important state information

* Expose and aggregate metrics from all components through prometheus/grafana

* Operations analytics based on metrics exposed   

* Auto-pilot functions such as distributing NSM registry into multiple pods according to the size of the cluster among other functions that may be addressed as well via automation.

### Install

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


### Usage 

To create a new NSM custom resource, after deploying the operator itself, use the CR manifest below. Here we have an example with the Vector Packet Processing as a forwarding plane for nsm.

```
kubectl apply -f deploy/crds/nsm.networkservicemesh.io_v1alpha1_nsm_cr.yaml
```
Now we should see something like this:
```
oc get pods -n nsm

NAME                                    READY   STATUS    RESTARTS   AGE
nsm-admission-webhook-b54cb4545-spt42   1/1     Running   0          6m58s
nsm-operator-794d58d4f8-n7tgj           1/1     Running   0          10m
nsm-vpp-forwarder-bh8nv                 1/1     Running   0          6m57s
nsm-vpp-forwarder-blrfb                 1/1     Running   0          6m57s
nsm-vpp-forwarder-mrmgl                 1/1     Running   0          6m57s
nsmgr-jwrbt                             3/3     Running   0          6m57s
nsmgr-s5tg2                             3/3     Running   0          6m57s
nsmgr-xfdf8                             3/3     Running   0          6m57s
```

| NOTE: nsm operator doesn't support the use of secure mode, aka spire, on Openshift yet. Planned for release 0.0.2. Jaegger tracing option is currently under development for both kubernetes and openshift as well. |
| --- |

After that, we have the basic network service mesh control plane and forwarding plan in place. That's the moment we can start to play with network services examples. 

First install the nsm helm repo to the cluster: `helm repo add nsm https://helm.nsm.dev/`

Then check the nsm [official examples](https://github.com/networkservicemesh/networkservicemesh/blob/master/docs/guide-quickstart.md#run) and the [community examples](https://github.com/networkservicemesh/examples) and the official instructions as well.


### Contributing

Anyone interested in contributing to the nsm-operator is welcomed and 
should start by reviewing the [development][docs_dev] documentation.

For any questions you can reach out to us on the CNCF slack cloud-native.slack.com. Take a look at the channels below and feel free to send me messages using `@alexandre menezes` there.

[![Slack Channel](https://img.shields.io/badge/Slack:-%23nsm%20on%20CNCF%20Slack-blue.svg?style=plastic&logo=slack)](https://cloud-native.slack.com/messages/CHQNNUPN1/)

[![Slack Channel](https://img.shields.io/badge/Slack:-%23nsm--dev%20on%20CNCF%20Slack-blue.svg?style=plastic&logo=slack)](https://cloud-native.slack.com/messages/CHSKJ4849/)

[![Slack Invite](https://img.shields.io/badge/Slack-CNCF%20Slack%20Invite-blue.svg?style=plastic&logo=slack)](https://slack.cncf.io/)

#### CNCF Code of Conduct:
  * The Network Service Mesh community follows the [CNCF Community Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).

### License

nsm-operator is released under the Apache 2.0 license. Please check the [LICENSE][license_file] file for details.

[nsm_home]:https://networkservicemesh.io
[nsm_whatis]:https://github.com/networkservicemesh/networkservicemesh/blob/master/docs/what-is-nsm.md
[docs_dev]:./docs/development.md
[license_file]:./LICENSE
[requirements]:./docs/requirements.md
[olm_install_guide]:https://github.com/operator-framework/operator-lifecycle-manager/blob/master/doc/install/install.md
[openshift_olm_install]:./docs/openshift_olm_install.md
[openshift_manual_install]:./docs/openshift_manual_install.md
[k8s_olm_install]:./docs/k8s_olm_install.md
[k8s_manual_install]:./docs/k8s_manual_install.md
