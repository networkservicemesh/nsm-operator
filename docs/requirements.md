## Service Accounts and Other Objects
```
cat <<EOF | kubectl apply -f -
---
apiVersion: v1
kind: Namespace
metadata:
  name: spire
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: nsm-role
  labels:
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
rules:
  - apiGroups: ["networkservicemesh.io"]
    resources:
      - "networkservices"
      - "networkserviceendpoints"
      - "networkservicemanagers"
    verbs: ["*"]
  - apiGroups: ["apiextensions.k8s.io"]
    resources: ["customresourcedefinitions"]
    verbs: ["*"]
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get"]
  - apiGroups: [""]
    resources: ["nodes", "services", "namespaces"]
    verbs: ["get", "list", "watch"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: nsm-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: nsm-role
subjects:
  - kind: ServiceAccount
    name: nsmgr-acc
    namespace: default
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: aggregate-network-services-view
  labels:
    # Add these permissions to the "view" default role.
    rbac.authorization.k8s.io/aggregate-to-view: "true"
rules:
  - apiGroups: ["networkservicemesh.io"]
    resources: ["networkservices"]
    verbs: ["get", "list", "watch"]
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: networkserviceendpoints.networkservicemesh.io
spec:
  conversion:
    strategy: None
  group: networkservicemesh.io
  names:
    kind: NetworkServiceEndpoint
    listKind: NetworkServiceEndpointList
    plural: networkserviceendpoints
    shortNames:
      - nse
      - nses
    singular: networkserviceendpoint
  scope: Namespaced
  version: v1alpha1
  versions:
    - name: v1alpha1
      served: true
      storage: true

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: networkservicemanagers.networkservicemesh.io
spec:
  conversion:
    strategy: None
  group: networkservicemesh.io
  names:
    kind: NetworkServiceManager
    listKind: NetworkServiceManagerList
    plural: networkservicemanagers
    shortNames:
      - nsm
      - nsms
    singular: networkservicemanager
  scope: Namespaced
  version: v1alpha1
  versions:
    - name: v1alpha1
      served: true
      storage: true
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: networkservices.networkservicemesh.io
spec:
  conversion:
    strategy: None
  group: networkservicemesh.io
  names:
    kind: NetworkService
    listKind: NetworkServiceList
    plural: networkservices
    shortNames:
      - netsvc
      - netsvcs
    singular: networkservice
  scope: Namespaced
  version: v1alpha1
  versions:
    - name: v1alpha1
      served: true
      storage: true
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nse-acc
  namespace: default
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nsc-acc
  namespace: default
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nsmgr-acc
  namespace: default
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: forward-plane-acc
  namespace: default
```
## Spire

In order to have nsm working spire must be already installed in the cluster. 

You can do that using helm like below:

```
git clone git@github.com:acmenezes/nsm-operator.git
cd nsm-operator
helm install spire deploy/helm/spire
```