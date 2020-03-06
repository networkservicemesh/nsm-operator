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

```

```

Now to run network services check the examples on the [usage page](usage.md).