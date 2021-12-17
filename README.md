## Network Service Mesh Operator

[![Go Report Card](https://goreportcard.com/badge/github.com/networkservicemesh/nsm-operator "Go Report Card")](https://goreportcard.com/report/github.com/networkservicemesh/nsm-operator)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](code-of-conduct.md) 

A Kubernetes Operator for Installing and Managing Network Service Meshes

#### Overview

The Network Service Mesh Operator is a tool to install and manage the [Network Service Mesh][nsm_home] application which <em> is a novel approach solving complicated L2/L3 use cases in Kubernetes that are tricky to address with the existing Kubernetes Network Model. Inspired by Istio, Network Service Mesh maps the concept of a Service Mesh to L2/L3 payloads as part of an attempt to re-imagine NFV in a Cloud-native way! </em>. To  better understand the network service meshes take a look at [what is nsm][https://networkservicemesh.io/].

The operator is a single pod workload that automates operational human knowledge behind the scenes to create the service mesh infrastructure components deploying a webhook and daemonsets with the network service managers and forwarding plane workloads taking the configuration from the Custom Resource manifest created specifically to be used with the operator. It aims to be platform independed and for such should run well in any kubernetes distribution.

#### Installation Steps:

Step 1 - To install the nsm-operator with all its denpendencies such as spire run:

```
make deploy
```

That command will create the NSM namespace, install spire using the helm chart present on scripts/spire, configure spire and register the nsm-operator service account and namespace on spire and finally install all the necessary RBAC manifests with the nsm-operator deployment.

*** Please remark that for OpenShift both nsm-operator and client applications need priviledged security contexts and security context constraints to make it work. ***

Step 2 - Install an NSM sample instance:

You can find an NSM custom resource example under config/samples/nsm_v1alpha1_nsm.yaml

Here is how it looks like:
```
apiVersion: nsm.networkservicemesh.io/v1alpha1
kind: NSM
metadata:
  name: nsm-sample
  namespace: nsm
spec:
  tag: v1.0.0
  nsmPullPolicy: IfNotPresent

  registry: ghcr.io
  organization: networkservicemesh

  registryMemoryImage: cmd-registry-memory
  nsmgrImage: cmd-nsmgr
  
  # Forwarding Plane Configs
  forwardingPlaneName: vpp
  forwardingPlaneImage: cmd-forwarder-vpp

```

You can deploy it after the operator is deployed by running:

```
kubectl apply -f config/samples/nsm_v1alpha1_nsm.yaml
```

With the CR in place it's possible to check the nsm namespace for it's workloads:

```
kubectl get pods -n nsm

NAME                            READY   STATUS    RESTARTS   AGE
nsm-operator-866f4ff5c8-j6gc8   1/1     Running   0          53s
nsm-registry-c84c97c4c-2mdzs    1/1     Running   0          8s
nsmgr-cqgtw                     1/1     Running   0          8s
nsmgr-fss6h                     1/1     Running   0          8s
nsmgr-ns5tz                     1/1     Running   0          8s
vpp-2wx56                       1/1     Running   0          8s
vpp-8qtjx                       1/1     Running   0          8s
vpp-gv56g                       1/1     Running   0          8s
```

Step 3 - There is a sample ICMP responder in a helm chart format that can be run as below:

```
helm install icmp-responder config/samples/application -n nsm

NAME: icmp-responder
LAST DEPLOYED: Fri Dec 17 10:39:51 2021
NAMESPACE: nsm
STATUS: deployed
REVISION: 1
TEST SUITE: None
```
You should see 2 other Pods running in the nsm namespace:

```
NAME                            READY   STATUS    RESTARTS   AGE
nsc-kernel-6b5d76f6bc-rk74g     1/1     Running   0          46s
nse-kernel-5579898565-6h8zh     1/1     Running   0          46s
nsm-operator-866f4ff5c8-j6gc8   1/1     Running   0          4m55s
nsm-registry-c84c97c4c-2mdzs    1/1     Running   0          4m10s
nsmgr-cqgtw                     1/1     Running   0          4m10s
nsmgr-fss6h                     1/1     Running   0          4m10s
nsmgr-ns5tz                     1/1     Running   0          4m10s
vpp-2wx56                       1/1     Running   0          4m10s
vpp-8qtjx                       1/1     Running   0          4m10s
vpp-gv56g                       1/1     Running   0          4m10s
```

You can check if they succeeded by entering the pods as below:

```
kubectl exec -it nsc-kernel-6b5d76f6bc-rk74g -n nsm -- /bin/sh

/ # 
```

Check the presence of a secondary network provided by NSM:
```
/ # ip addr
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
3: eth0@if34: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 8951 qdisc noqueue state UP 
    link/ether 0a:58:0a:81:02:17 brd ff:ff:ff:ff:ff:ff
    inet 10.129.2.23/23 brd 10.129.3.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::7018:ebff:fe64:85c6/64 scope link 
       valid_lft forever preferred_lft forever
36: nsm@if35: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 8951 qdisc noqueue state UP qlen 1000
    link/ether 6e:57:10:17:f7:8a brd ff:ff:ff:ff:ff:ff
    inet 169.254.0.1/32 scope global nsm
       valid_lft forever preferred_lft forever
    inet6 fe80::6c57:10ff:fe17:f78a/64 scope link 
       valid_lft forever preferred_lft forever
```

And finally ping the endpoint using the endpoint ip address:
```
/ # ping 169.254.0.0
PING 169.254.0.0 (169.254.0.0): 56 data bytes
64 bytes from 169.254.0.0: seq=0 ttl=64 time=1.279 ms
64 bytes from 169.254.0.0: seq=1 ttl=64 time=0.579 ms
64 bytes from 169.254.0.0: seq=2 ttl=64 time=0.543 ms
64 bytes from 169.254.0.0: seq=3 ttl=64 time=0.480 ms
64 bytes from 169.254.0.0: seq=4 ttl=64 time=0.542 ms
^C
```

#### Cleanup

Delete the client application:
```
helm delete icmp-responder -n nsm
```

Delete the nsm CR:
```
kubectl delete nsm nsm-sample -n nsm
```

Delete nsm-operator and all its dependencies:
```
make undeploy
```

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
