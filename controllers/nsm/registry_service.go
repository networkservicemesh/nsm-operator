package controllers

import (
	"context"

	"github.com/go-logr/logr"
	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/apis/nsm/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type RegistryServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func NewRegistryServiceReconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme) *RegistryServiceReconciler {
	return &RegistryServiceReconciler{
		Client: client,
		Log:    log,
		Scheme: scheme,
	}
}

func (r *RegistryServiceReconciler) Reconcile(ctx context.Context, nsm *nsmv1alpha1.NSM) error {

	svc := &corev1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: "nsm-registry-svc", Namespace: nsm.ObjectMeta.Namespace}, svc)
	if err != nil {
		if apierrors.IsNotFound(err) {
			svc = r.serviceForNsmRegistry(nsm)
			err = r.Client.Create(context.TODO(), svc)
			if err != nil {
				r.Log.Error(err, "failed to create service for nsm-registry")
				return err
			}
			return nil
		}
		return err
	}
	r.Log.Info("nsm registry service already exists, skipping creation")
	return nil
}

func (r *RegistryServiceReconciler) serviceForNsmRegistry(nsm *nsmv1alpha1.NSM) *corev1.Service {

	registryLabel := map[string]string{"app": "nsm-registry", "spiffe.io/spiffe-id": "true"}

	objectMeta := newObjectMeta("nsm-registry-svc", "nsm", map[string]string{"app": "nsm"})

	service := &corev1.Service{
		ObjectMeta: objectMeta,
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "nsm-registry-svc",
					Protocol:   "TCP",
					Port:       5002,
					TargetPort: intstr.FromInt(5002)},
			},
			Selector: registryLabel,
		},
	}
	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, service, r.Scheme)
	return service
}
