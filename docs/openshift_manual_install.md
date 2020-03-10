## Manual Openshift Install:

First let's install the Network Service Mesh requirements. They have been consolidated in a single yaml file for convinience. It will apply all initial RBAC and service account settings for NSM itself not the operator, some important CRDs used by NSM and create the nsm namespace as well.

```
oc apply -f deploy/nsm_requirements/openshift.yaml
```

Once the requirements were applied, we can deploy the operator. First let's get into the nsm namespace:

```
oc project nsm
```

Now let's apply the RBAC permissions for the operator itself:
```
oc apply -f deploy/role.yaml
oc apply -f deploy/service_account.yaml
oc apply -f deploy/role_binding.yaml
```
Once we have the RBAC permissions we need to deploy the NSM custom resource definition:
```
oc apply -f deploy/crds/nsm.networkservicemesh.io_nsms_crd.yaml
```

Finally let's deploy the operator:
```
oc apply -f deploy/operator.yaml
```

Check that you have a running operator with:
```
oc get pods
```
And you should see something like below:
```
NAME                                     READY   STATUS    RESTARTS   AGE
nsm-operator-6dbfdc8dc5-4gq2d            1/1     Running   0          98s
```
Once the operator is running we can create a Network Service Mesh instance. To see how to do it and some network services example check the [usage page](usage.md).