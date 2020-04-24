package nsm

import (
	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/pkg/apis/nsm/v1alpha1"
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileNSM) mutatingConfigForWebhook(nsm *nsmv1alpha1.NSM) *admissionregv1beta1.MutatingWebhookConfiguration {

	var path string
	path = "/mutate"
	mutatingConfig := &admissionregv1beta1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: webhookMutatingConfigName,
		},
		Webhooks: []admissionregv1beta1.MutatingWebhook{
			{
				Name: "admission-webhook.networkservicemesh.io",
				ClientConfig: admissionregv1beta1.WebhookClientConfig{
					Service: &admissionregv1beta1.ServiceReference{
						Name:      webhookServiceName,
						Namespace: nsm.Namespace,
						Path:      &path,
					},
					CABundle: caCert,
				},

				Rules: []admissionregv1beta1.RuleWithOperations{
					{
						Operations: []admissionregv1beta1.OperationType{admissionregv1beta1.Create},
						Rule: admissionregv1beta1.Rule{
							APIGroups:   []string{"apps", "extensions", ""},
							APIVersions: []string{"v1", "v1beta1"},
							Resources:   []string{"pods", "deployments", "services"},
						},
					},
				},
			},
		},
	}
	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, mutatingConfig, r.scheme)
	return mutatingConfig
}
