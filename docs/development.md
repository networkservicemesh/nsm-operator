## Contributor's Guide

#### Software Requirements
In order to contribute on the nsm operator project make sure you have the minimal software requirements below:

- [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [golang version 1.13+](https://golang.org/doc/install)
- [operator-sdk v0.15.1](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md)
- [podman version 1.2.0+](https://podman.io/getting-started/installation.html), [buildah version 1.7+](https://github.com/containers/buildah/blob/master/install.md) or [docker version 19+](https://docs.docker.com/install/)
- GNU make (https://www.gnu.org/software/make/)
- access to a Kubernetes 1.14+ cluster (can be [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/))
- access to an [OpenShift](https://docs.openshift.com/container-platform/4.3/installing/installing_aws/installing-aws-default.html) Cluster (version 4+, can not be minishift)
- have both clients [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) and [oc](https://docs.openshift.com/container-platform/4.3/cli_reference/openshift_cli/getting-started-cli.html)

#### Project's Page and Contribution Guidelines

At this point we have a great opportunity for people wanting to join the project. You can check the project's board [here](https://github.com/networkservicemesh/nsm-operator/projects/1). We have quite a lot yet to put in. This is the opportunity to shape the project from the beginning. 

If you want to develop new features you can submit an Issue and fork the project making a new branch from the development branch and commit changes against that new branch. After that you can submit a PR against your copy of the development branch. We'll review it and merge it if is aligned with the project.

#### Operator Framework Important References

* [Getting Started With Operator SDK](https://github.com/operator-framework/getting-started/blob/master/README.md)

* [Submitting your operator](https://github.com/operator-framework/community-operators/blob/master/docs/contributing.md)

* [Create a Cluster Service Version](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/doc/design/building-your-csv.md)

* [Operator SDK Scorecard Tests](https://github.com/operator-framework/operator-sdk/blob/master/doc/test-framework/scorecard.md)

* [Operator Lifecycle Manager](https://github.com/operator-framework/operator-lifecycle-manager)
