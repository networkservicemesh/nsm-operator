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

type AdmissionWHServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func NewAdmissionWHServiceReconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme) *AdmissionWHServiceReconciler {
	return &AdmissionWHServiceReconciler{
		Client: client,
		Log:    log,
		Scheme: scheme,
	}
}

func (r *AdmissionWHServiceReconciler) Reconcile(ctx context.Context, nsm *nsmv1alpha1.NSM) error {

	svc := &corev1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: "admission-webhook-svc", Namespace: nsm.ObjectMeta.Namespace}, svc)
	if err != nil {
		if apierrors.IsNotFound(err) {
			svc = r.serviceForAdmissionWH(nsm)
			err = r.Client.Create(context.TODO(), svc)
			if err != nil {
				r.Log.Error(err, "failed to create service for admission-webhook")
				return err
			}
			r.Log.Info("admission-webhook service created")
			return nil
		}
		return err
	}
	r.Log.Info("admission-webhook service already exists, skipping creation")
	return nil
}

func (r *AdmissionWHServiceReconciler) serviceForAdmissionWH(nsm *nsmv1alpha1.NSM) *corev1.Service {

	admissionWHLabel := map[string]string{"app": "admission-webhook-k8s"}

	objectMeta := newObjectMeta("admission-webhook-svc", "nsm", map[string]string{"app": "nsm"})

	service := &corev1.Service{
		ObjectMeta: objectMeta,
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "admission-webhook-svc",
					Protocol:   "TCP",
					Port:       443,
					TargetPort: intstr.FromInt(443)},
			},
			Selector: admissionWHLabel,
		},
	}
	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, service, r.Scheme)
	return service
}
