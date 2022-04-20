package controllers

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/apis/nsm/v1alpha1"
)

type createResourceFunc func(nsm *nsmv1alpha1.NSM, objectMeta metav1.ObjectMeta) client.Object

func setObjectMeta(name string, namespace string, labels map[string]string) metav1.ObjectMeta {
	objectMeta := metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
	}
	return objectMeta
}

func newObjectMeta(name string, namespace string, labels map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
	}
}

func (r *NSMReconciler) reconcileResource(
	createResource createResourceFunc,
	nsm *nsmv1alpha1.NSM,
	resource client.Object,
	objectMeta metav1.ObjectMeta) error {

	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: objectMeta.Name, Namespace: objectMeta.Namespace}, resource)
	if err != nil {
		if errors.IsNotFound(err) {

			fmt.Printf("%s - reconcileResource: creating a new resource for NSM", time.Now())

			resource := createResource(nsm, objectMeta)
			err = r.Client.Create(context.TODO(), resource)

			if err != nil {
				fmt.Printf("%s - reconcileResource: failed to create new resource, err: %s", time.Now(), err)
				return err
			}

			return nil
		}

		fmt.Printf("%s - reconcileResource: error reading resource, err: %s", time.Now(), err)
		return err
	}

	fmt.Printf("%s - resource %s already exists in namespace %s", time.Now(), objectMeta.Name, objectMeta.Namespace)
	fmt.Printf("Please delete any instance of NSM before creating a new one.")

	return nil
}
