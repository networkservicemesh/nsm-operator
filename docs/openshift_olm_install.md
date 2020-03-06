## Openshift Embedded Operator Hub

First let's create a namespace named nsm:

```
oc create namespace nsm
```

Now after logging in to the Openshift console, clicking on OperatorHub we can filter by nsm on the text box and we will find the nsm operator badge.

<img src="./img/1-nsm-operator-hub.png">

After clicking on it we have the presentation page with the install button. Let's run it. Click install!

<img src="./img/2-nsm-operator-install.png">

We land on the subscription home page. Here is were you subscribe to the nsm operator offered by the Red Hat's Openshift OperatorHub. Pick the nsm namespace and let the defaults on the other options. Finally click subscribe.

<img src="./img/3-nsm-operator-subscription.png">

After a while we should be able to see the operator with the Status `InstallSuceeded`. Then click on the Network Service Mesh among the Provided APIs.

<img src="./img/4-nsm-operator-installed.png">

You will see a button saying Create NSM. Let's click it!

<img src="./img/5-nsm-servicemesh.png">

Below we present 2 methods that can be used with the embedded openshift OperatorHub. The edit form method and the yaml. Both are good. Choose one and go forward. At this point in time the default choices are the only supported. Let's click create.

Edit form method:
<img src="./img/6-nsm-servicemesh-create-form.png">
YAML create method:
<img src="./img/7-nsm-servicemesh-create-yaml.png">

After clicking create, a few seconds after we should have a NSM enabled cluster ready to be configure with new network services. Check if it got to the `running` Status.

<img src="./img/8-nsm-servicemesh-running.png">

And finally clicking on the created nsm object and on the Resources tab we can check it's components running on Openshift:

<img src="./img/9-nsm-resources.png">

Now to run network services check the examples on the [usage page](usage.md).