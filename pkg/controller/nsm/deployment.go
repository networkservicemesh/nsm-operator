package nsm

import (
	nsmv1alpha1 "github.com/acmenezes/nsm-operator/pkg/apis/nsm/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileNSM) deploymentForWebhook(nsm *nsmv1alpha1.NSM) *appsv1.Deployment {
	ls := labelsForNSMAdmissionWebhook(nsm.Name)
	replicas := nsm.Spec.Replicas
	registry := nsm.Spec.Registry
	org := nsm.Spec.Org
	tag := nsm.Spec.Tag
	webhookPullPolicy := nsm.Spec.PullPolicy

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nsm-admission-webhook",
			Namespace: nsm.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{

						Name:            nsm.Spec.WebhookName,
						Image:           registry + "/" + org + "/admission-webhook:" + tag,
						ImagePullPolicy: webhookPullPolicy,
						Env: []corev1.EnvVar{
							{Name: "REPO", Value: org},
							{Name: "TAG", Value: tag},
						},
						// TODO: check Jaeger tracing option and insert other envs

						VolumeMounts: []corev1.VolumeMount{
							{Name: "webhook-certs",
								MountPath: "/etc/webhook/certs",
								ReadOnly:  true},
						},

						// 	LivenessProbe: *v1.Probe{
						// 		Handler:             HttpGetAvction{Path: "/liveness", Port: "5555"},
						// 		InitialDelaySeconds: 10,
						// 		PeriodSeconds:       3,
						// 		TimeoutSeconds:      3},

						// 	ReadinessProbe: *v1.Probe{
						// 		{Handler: HttpGetAction{Path: "/readiness", Port: "5555"},
						// 			InitialDelaySeconds: 10,
						// 			PeriodSeconds:       3,
						// 			TimeoutSeconds:      3}},
						// }},
					}},
					Volumes: []corev1.Volume{{
						Name: "webhook-certs",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: nsm.Spec.WebhookSecretName,
							},
						},
					},
					},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	controllerutil.SetControllerReference(nsm, deploy, r.scheme)
	return deploy
}

func labelsForNSMAdmissionWebhook(name string) map[string]string {
	return map[string]string{"app": "nsm-admission-webhook", "nsm-admission-webhook-cr": name}
}
