## Manual Kubernetes Install:

First let's install the Network Service Mesh requirements. They have been consolidated in a single yaml file for convinience. It will apply all initial RBAC and service account settings for NSM itself not the operator, some important CRDs used by NSM and create the nsm namespace as well.

```
kubectl apply -f deploy/nsm_requirements/k8s.yaml
```

Once the requirements were applied, we can deploy the operator. First let's set the local context to use the nsm namespace. Below we can find an example extracted from a minikube environment:

```
kubectl config set-context nsm/minikube --cluster minikube --user minikube --namespace nsm
kubectl config use-context nsm/minikube
```


Now let's apply the RBAC permissions for the operator itself:
```
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/role_binding.yaml
```
Once we have the RBAC permissions we need to deploy the NSM custom resource definition:
```
kubectl apply -f deploy/crds/nsm.networkservicemesh.io_nsms_crd.yaml
```

Finally let's deploy the operator:
```
kubectl apply -f deploy/operator.yaml
```

Check that you have a running operator with:
```
kubectl get pods
```
And you should see something like below:
```
NAME                            READY   STATUS    RESTARTS   AGE
nsm-operator-6dbfdc8dc5-6fn85   1/1     Running   0          6m5s
```
Once the operator is running we can create a Network Service Mesh instance. To see how to do it and some network services example check the [usage page](usage.md).