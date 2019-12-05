package nsm

import (
	nsmv1alpha1 "github.com/acmenezes/nsm-operator/pkg/apis/nsm/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileNSM) mutatingConfigForWebhook(nsm *nsmv1alpha1.NSM) *admissionregv1beta1.MutatingWebhookConfiguration {

	webhookConfig := &v1beta1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: webhookConfigName,
		},
		Webhooks: []v1beta1.MutatingWebhook{
			{
				Name: "vpa.k8s.io",
				ClientConfig: &admissionregv1beta1.WebhookClientConfig{
					Service: admissionregv1beta1.ServiceReference{
						Name: webhookServiceName,
						Namespace: nsm.Namespace,
						Path: "/mutate",
					},
					caBundle: 
				},

				Rules: []v1beta1.RuleWithOperations{
					{
						Operations: []v1beta1.OperationType{v1beta1.Create},
						Rule: v1beta1.Rule{
							APIGroups:   []string{""},
							APIVersions: []string{"v1"},
							Resources:   []string{"pods"},
						},
					},

					},
				},
			},
		},

}