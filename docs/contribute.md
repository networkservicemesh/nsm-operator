## Contributor's Guide

#### Software Requirements
In order to contribute on the nsm operator project make sure you have the minimal software requirements below:

- [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [golang version 1.13+](https://golang.org/doc/install)
- [operator-sdk v0.15.1](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md)
- [podman version 1.2.0+](https://podman.io/getting-started/installation.html)
 or [docker version 19+](https://docs.docker.com/install/)
- GNU make (https://www.gnu.org/software/make/)
- access to a Kubernetes 1.21+ cluster (for development it can be [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/) or [kind](https://kind.sigs.k8s.io/docs/user/quick-start/))
- access to an [OpenShift](https://docs.openshift.com/container-platform/4.3/installing/installing_aws/installing-aws-default.html) Cluster (Not tested on [CRC](https://github.com/code-ready/crc))
- have both clients [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) and [oc](https://docs.openshift.com/container-platform/4.10/cli_reference/openshift_cli/getting-started-cli.html)

        NOTE: the environment for Kubernetes or OpenShift need to be setup beforehand. It means that a kubeconfig file and/or proper cluster admin credentials must be in place and the Kubeapi-server endpoint must be accessible through the network. The projects make targets and commands will use that connection to communicate with the api, deploy and run all the resources.

#### Developing and Debugging

This source code can be run basically in two ways: 

1 - Running `make run`, which executes go run behind the scenes and will rely on logs and error messages for debugging. With that the operator will run from the developer's machine and communicate externally with the Kubeapi-server endpoint.

2- By using the debug and run from VSCode. With the main.go file in focus you can hit F5 and trigger the VSCode debugger that has other features such as breakpoints and variable watches. For that you will need the golang extension installed in your environment. We do recommend to use VSCode but it's not mandatory. Any golang environment should work just fine. Check the [VSCode debug page](https://code.visualstudio.com/docs/editor/debugging) for more information.

#### Building the operator

First in the make file we need to setup two variables according to the needs:

`IMG` - for the registry where the operator image will be stored. It needs to have all the path in the string plus the tag.

`BUILDER` - is the container builder that can be podman or docker.

The project has a few make targets available for building:

`make build` -  Will generate a binary file under a bin directory with the controller, manager and all the added reconcilers and types involved. That's the operator itself.

`make container-build` - Will generate the container image that runs the operator inside the cluster as a platform component.

`make container-push` - Will push the operator

With the container image available in a registry and that registry edited in the makefile the `make deploy` command can be run in order to deploy that container with all the necessary manifests to the cluster. From that point you should have your dev version running on a pod in the nsm namespace. Then follow the example on the README file to test with the other resources.

#### Resources to learn more about operators

* [Kubebuilder Book](https://book.kubebuilder.io/)

* [Operator SDK](https://sdk.operatorframework.io/)

* [Operator Lifecycle Manager](https://olm.operatorframework.io/)