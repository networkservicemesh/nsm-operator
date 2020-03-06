## Kubernetes Operator Lifecycle Manager Install

Here we use Minikube to demonstrate the installation process. Please follow the instructions on your kubernetes distribution of choice.

For Minikube we can start a small setup such as below:

```
minikube start -p nsm --cpus=2 --disk-size='15g' --memory='4g'
```
Check if the node is ready:
`kubectl get nodes`
```
NAME   STATUS   ROLES    AGE     VERSION
nsm    Ready    master   6m58s   v1.17.2
```
After you get the node ready, let's install OLM on minikube.
```
kubectl apply -f https://github.com/operator-framework/operator-lifecycle-manager/releases/download/0.14.1/crds.yaml
kubectl apply -f https://github.com/operator-framework/operator-lifecycle-manager/releases/download/0.14.1/olm.yaml
```

Let's check if the OLM namespaces were created:
`kubectl get ns`

```
NAME              STATUS   AGE
default           Active   93s
kube-node-lease   Active   94s
kube-public       Active   94s
kube-system       Active   94s
olm               Active   25s
operators         Active   25s
```

Let's check if all OLM pods are running:
`kubectl get pods -n olm`
```
NAME                                READY   STATUS    RESTARTS   AGE
catalog-operator-64b6b59c4f-kp7w9   1/1     Running   0          53s
olm-operator-844fb69f58-nk2zd       1/1     Running   0          53s
operatorhubio-catalog-qz57n         1/1     Running   1          42s
packageserver-7db898d5d5-495k2      1/1     Running   0          39s
packageserver-7db898d5d5-f8hph      1/1     Running   0          26s
```

## Installing the operator

First we create the nsm namespace:
```
kubectl create namespace nsm
```
Then we apply the nsm catalog source to the OLM namespace:
```
kubectl apply -f deploy/catalog_source.yaml -n olm
```
Check if it was deployed correctly:

`kubectl get catalogsources -n olm`
```
NAME                    DISPLAY                          TYPE   PUBLISHER               AGE
nsm-catalog             Network Service Mesh Operators   grpc   networkservicemesh.io   127m
operatorhubio-catalog   Community Operators              grpc   OperatorHub.io          132m
```
And should now have also a pod serving that custom catalog in the OLM namespace:

`kubectl get pods -n olm`
```
NAME                                READY   STATUS    RESTARTS   AGE
catalog-operator-64b6b59c4f-2hdt4   1/1     Running   0          134m
nsm-catalog-p2cfc                   1/1     Running   0          129m
olm-operator-844fb69f58-z4rz2       1/1     Running   0          134m
operatorhubio-catalog-7txwm         1/1     Running   0          133m
packageserver-75965f5c68-6gqh9      1/1     Running   0          133m
packageserver-75965f5c68-msjd8      1/1     Running   0          133m
```
Now we apply the operator group. It basically signals to OLM that it can deploy operators to that namespace were the operator group resides. And after that we can subscribe to our nsm operator using the subscription.
```
kubectl apply -f deploy/operator_group.yaml -n nsm
kubectl apply -f deploy/subscription.yaml -n nsm
```
Now we can check what Install Plan and Cluster Service Version have been deployed and get more information about the operator server by the catalog source.

`kubectl get installplan -n nsm`
```
NAME            CSV                   APPROVAL    APPROVED
install-l8wvr   nsm-operator.v0.0.1   Automatic   true
```

`kubectl get csv -n nsm`
```
NAME                  DISPLAY                         VERSION   REPLACES   PHASE
nsm-operator.v0.0.1   Network Service Mesh Operator   0.0.1                Succeeded
```
Finally let's check our operator pod

`kubectl get pods -n nsm`
```
NAME                            READY   STATUS    RESTARTS   AGE
nsm-operator-8588c8dd6c-tsnx7   1/1     Running   0          8m56s
```

Now to run network services check the examples on the [usage page](usage.md).