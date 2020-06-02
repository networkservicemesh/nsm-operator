package nsm

import (
	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/pkg/apis/nsm/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileNSM) serviceForWebhook(nsm *nsmv1alpha1.NSM) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webhookServiceName,
			Namespace: nsm.Namespace,
			Labels:    labelsForNSMAdmissionWebhook(nsm.Name),
			// TODO: Solve TLS Certs for OCP - This annotation below is specific to OpenShift and needs to be addressed other way
			Annotations: map[string]string{"service.beta.openshift.io/serving-cert-secret-name": "nsm-admission-webhook-certs"},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: webhookName + "-port", Port: webhookServicePort, TargetPort: intstr.FromInt(webhookServiceTargetPort)},
			},
			Selector: map[string]string{"app": webhookName},
		},
	}
	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, service, r.scheme)
	return service
}
