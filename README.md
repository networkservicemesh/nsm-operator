# Network Service Mesh Operator

[![Go Report Card](https://goreportcard.com/badge/github.com/acmenezes/nsm-operator "Go Report Card")](https://goreportcard.com/report/github.com/acmenezes/nsm-operator)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](code-of-conduct.md) 

A Kubernetes Operator for Installing and Managing Network Service Meshes

## Overview

The Network Service Mesh Operator is a tool to install and manage the [Network Service Mesh][nsm_home] application which <em> is a novel approach solving complicated L2/L3 use cases in Kubernetes that are tricky to address with the existing Kubernetes Network Model. Inspired by Istio, Network Service Mesh maps the concept of a Service Mesh to L2/L3 payloads as part of an attempt to re-imagine NFV in a Cloud-native way! </em>. To  better understand the network service meshes take a look at [what is nsm](nsm_whatis).

The operator is a single pod workload that automates operational human knowledge behind the scenes to create the service mesh infrastructure components deploying a webhook and daemonsets with the network service managers and forwarding plane workloads taking the configuration from the Custom Resource manifest created specifically to be used with the operator. It aims to be platform independed and for such should run well in any kubernetes distribution.

Some of the features intended to be embedded with the operator are

* Installing, configuring and making cleanups if needed;

* Upgrading, backing up and restoring any important state information

* Expose and aggregate metrics from all components through prometheus/grafana

* Operations analytics based on metrics exposed   

* Auto-pilot functions such as distributing NSM registry into multiple pods according to the size of the cluster among other functions that may be addressed as well via automation.

## Install

At this point the installation is via the manifest below

```
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nsm-operator
spec:
  replicas: 1
  selector:
    matchLabels:
      name: nsm-operator
  template:
    metadata:
      labels:
        name: nsm-operator
    spec:
      serviceAccountName: nsm-operator
      containers:
        - name: nsm-operator
          # Replace this with the built image name
          image: quay.io/acmenezes/nsm-operator:v0.0.1
          command:
          - nsm-operator
          imagePullPolicy: Always
          env:
            - name: WATCH_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: OPERATOR_NAME
              value: "nsm-operator"
EOF
```

## Usage 

To create a new NSM custom resource, after deploying the operator itself, use the CR manifest below.
Here we have an example with the Vector Packet Processing as a forwarding plane for nsm.

```
cat <<EOF | kubectl apply -f -
apiVersion: nsm.networkservicemesh.io/v1alpha1
kind: NSM
metadata:
  name: nsm
  namespace: default
spec:
  # repo configs
  registry: docker.io
  org: networkservicemesh
  tag: master
  pullPolicy: IfNotPresent

  # admission webhook configs
  webhookName: nsm-admission-webhook
  replicas: 1

  # INSECURE env var
  insecure: "true"

  # Forwarding Plane Configs
  forwardingPlaneName: vpp
  forwardingPlaneImage: vppagent-forwarder

  # Enable Spire
  spire: "false"

  # Enable Jaeger Tracing
  jaegerTracing: "false"
EOF
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
