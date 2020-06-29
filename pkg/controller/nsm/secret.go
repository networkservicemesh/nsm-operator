package nsm

import (
	"bytes"

	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/pkg/apis/nsm/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileNSM) secretForWebhook(nsm *nsmv1alpha1.NSM) *corev1.Secret {
	var k, c bytes.Buffer
	host := "nsm-admission-webhook-svc.nsm.svc"
	generateRSACerts(host, true, &k, &c)
	caCert = c.Bytes()
	key := k.Bytes()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webhookSecretName,
			Namespace: nsm.Namespace,
			Labels:    labelsForNSMAdmissionWebhook(nsm.Name),
		},
		Data: map[string][]byte{
			corev1.TLSCertKey:       caCert,
			corev1.TLSPrivateKeyKey: key,
		},
	}
	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, secret, r.scheme)
	return secret
}
