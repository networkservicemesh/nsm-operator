package nsm

import (
	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/pkg/apis/nsm/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileNSM) deploymentForWebhook(nsm *nsmv1alpha1.NSM) *appsv1.Deployment {
	ls := labelsForNSMAdmissionWebhook(nsm.Name)
	replicas := webhookReplicas
	registry := nsmRegistry
	org := nsmOrg
	tag := nsmVersion
	webhookPullPolicy := nsmPullPolicy

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webhookName,
			Namespace: nsm.Namespace,
			Labels:    labelsForNSMAdmissionWebhook(nsm.Name),
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
					ServiceAccountName: "nsm-webhook-acc",
					Containers: []corev1.Container{{
						Name:            webhookName,
						Image:           registry + "/" + org + "/admission-webhook:" + tag,
						ImagePullPolicy: webhookPullPolicy,
						// SecurityContext: &corev1.SecurityContext{
						// 	Capabilities: &corev1.Capabilities{
						// 		Add: []corev1.Capability{"NET_BIND_SERVICE"},
						// 	},
						// },
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

						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/liveness",
									Port: intstr.FromInt(probePort),
								},
							},
							InitialDelaySeconds: probeInitialDelay,
							PeriodSeconds:       probePeriod,
							TimeoutSeconds:      probeTimeout},

						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/readiness",
									Port: intstr.FromInt(probePort),
								},
							},
							InitialDelaySeconds: probeInitialDelay,
							PeriodSeconds:       probePeriod,
							TimeoutSeconds:      probeTimeout},
					}},
					Volumes: []corev1.Volume{{
						Name: "webhook-certs",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: webhookSecretName,
							},
						},
					},
					},
				},
			},
		},
	}
	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, deploy, r.scheme)
	return deploy
}
