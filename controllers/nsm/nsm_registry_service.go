package controllers

import (
	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/apis/nsm/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *NSMReconciler) serviceForNsmRegistry(nsm *nsmv1alpha1.NSM, objectMeta metav1.ObjectMeta) runtime.Object {

	registryLabel := map[string]string{"app": "nsm-registry"}

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
