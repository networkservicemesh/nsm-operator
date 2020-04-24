package nsm

import (
	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/pkg/apis/nsm/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileNSM) deamonSetForNSMGR(nsm *nsmv1alpha1.NSM) *appsv1.DaemonSet {

	registry := nsmRegistry
	org := nsmOrg
	tag := nsmVersion
	pullPolicy := nsmPullPolicy
	volType := corev1.HostPathDirectoryOrCreate
	privmode := true
	insecure := "true"

	if nsm.Spec.Insecure {
		insecure = "true"
	} else {
		insecure = "false"
	}

	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nsmgr",
			Namespace: nsm.Namespace,
			Labels:    map[string]string{"app": "nsmgr-daemonset"},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "nsmgr-daemonset"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "nsmgr-daemonset"},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "nsmgr-acc",
					Containers: []corev1.Container{

						// nsmdp container
						{
							Name:            "nsmdp",
							Image:           registry + "/" + org + "/nsmdp:" + tag,
							ImagePullPolicy: pullPolicy,
							Env: []corev1.EnvVar{
								{Name: "INSECURE", Value: insecure},

								// {Name: "TRACER_ENABLED", Value: nsm.Spec.JaegerTracing},
								// TODO: Jaeger tracing feature
								// {Name: "JAEGER_AGENT_HOST", Value: nsm.Spec.JaegerTracing},
								// {Name: "JAEGER_AGENT_PORT", Value: nsm.Spec.JaegerTracing}
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "kubelet-socket",
									MountPath: "/var/lib/kubelet/device-plugins",
								},
								{Name: "nsm-socket",
									MountPath: "/var/lib/networkservicemesh",
								},
								{Name: "spire-agent-socket",
									MountPath: "/run/spire/sockets",
									ReadOnly:  true,
								},
							},
						},
						// nsmd container
						{
							Name:            "nsmd",
							Image:           registry + "/" + org + "/nsmd:" + tag,
							ImagePullPolicy: pullPolicy,
							SecurityContext: &corev1.SecurityContext{
								Privileged: &privmode,
							},
							Env: []corev1.EnvVar{
								{Name: "INSECURE", Value: insecure},
								// {Name: "TRACER_ENABLED", Value: nsm.Spec.JaegerTracing},

								// TODO: Jaeger tracing feature
								// {Name: "JAEGER_AGENT_HOST", Value: nsm.Spec.JaegerTracing},
								// {Name: "JAEGER_AGENT_PORT", Value: nsm.Spec.JaegerTracing}
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "nsm-socket",
									MountPath: "/var/lib/networkservicemesh",
								},
								{Name: "nsm-plugin-socket",
									MountPath: "/var/lib/networkservicemesh/plugins",
								},
								{Name: "spire-agent-socket",
									MountPath: "/run/spire/sockets",
									ReadOnly:  true,
								},
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
						},
						// nsmd-k8s container
						{
							Name:            "nsmd-k8s",
							Image:           registry + "/" + org + "/nsmd-k8s:" + tag,
							ImagePullPolicy: pullPolicy,
							SecurityContext: &corev1.SecurityContext{
								Privileged: &privmode,
							},
							Env: []corev1.EnvVar{
								{Name: "INSECURE", Value: insecure},
								{Name: "POD_NAME", ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.name",
									}}},
								{Name: "POD_UID", ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.uid",
									}}},
								{Name: "NODE_NAME", ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "spec.nodeName",
									}}},

								// {Name: "TRACER_ENABLED", Value: nsm.Spec.JaegerTracing},
								// TODO: Jaeger tracing feature
								// {Name: "JAEGER_AGENT_HOST", Value: nsm.Spec.JaegerTracing},
								// {Name: "JAEGER_AGENT_PORT", Value: nsm.Spec.JaegerTracing}
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "nsm-plugin-socket",
									MountPath: "/var/lib/networkservicemesh/plugins",
								},
								{Name: "spire-agent-socket",
									MountPath: "/run/spire/sockets",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "kubelet-socket",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/kubelet/device-plugins",
									Type: &volType,
								}}},
						{
							Name: "nsm-socket",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/networkservicemesh",
									Type: &volType,
								}}},
						{
							Name: "nsm-plugin-socket",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/networkservicemesh/plugins",
									Type: &volType,
								}}},
						{
							Name: "spire-agent-socket",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/run/spire/sockets",
									Type: &volType,
								}}},
					},
				},
			},
		},
	}

	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, daemonset, r.scheme)
	return daemonset
}

func (r *ReconcileNSM) deamonSetForForwardingPlane(nsm *nsmv1alpha1.NSM) *appsv1.DaemonSet {

	registry := nsmRegistry
	org := nsmOrg
	tag := nsmVersion
	pullPolicy := nsmPullPolicy
	image := nsm.Spec.ForwardingPlaneImage
	fp := nsm.Spec.ForwardingPlaneName
	volType := corev1.HostPathDirectoryOrCreate
	mountPropagationMode := corev1.MountPropagationBidirectional
	privmode := true

	insecure := "true"

	if nsm.Spec.Insecure {
		insecure = "true"
	} else {
		insecure = "false"
	}

	daemonset := &appsv1.DaemonSet{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "nsm-" + fp + "-forwarder",
			Namespace: nsm.Namespace,
			Labels:    map[string]string{"app": "nsm-" + fp + "-forwarder"},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "nsm-" + fp + "-forwarder"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "nsm-" + fp + "-forwarder"},
				},
				Spec: corev1.PodSpec{
					HostPID:            true,
					HostNetwork:        true,
					ServiceAccountName: "forward-plane-acc",
					Containers: []corev1.Container{

						// forwarding plane container
						{
							Name:            image,
							Image:           registry + "/" + org + "/" + image + ":" + tag,
							ImagePullPolicy: pullPolicy,
							SecurityContext: &corev1.SecurityContext{
								Privileged: &privmode,
							},
							Env: []corev1.EnvVar{
								{Name: "INSECURE", Value: insecure},

								// {Name: "TRACER_ENABLED", Value: nsm.Spec.JaegerTracing},
								// TODO: Jaeger tracing feature
								// {Name: "JAEGER_AGENT_HOST", Value: nsm.Spec.JaegerTracing},
								// {Name: "JAEGER_AGENT_PORT", Value: nsm.Spec.JaegerTracing}
								{Name: "NSM_FORWARDER_SRC_IP", ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "status.podIP",
									}}},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "workspace",
									MountPath:        "/var/lib/networkservicemesh/",
									MountPropagation: &mountPropagationMode,
								},
								{Name: "spire-agent-socket",
									MountPath: "/run/spire/sockets",
									ReadOnly:  true,
								},
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
						},
						// TODO: Resources CPU limits and requests
					},
					Volumes: []corev1.Volume{
						{
							Name: "workspace",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/networkservicemesh",
									Type: &volType,
								}}},
						{
							Name: "spire-agent-socket",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/run/spire/sockets",
									Type: &volType,
								}}},
					},
				},
			},
		},
	}

	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, daemonset, r.scheme)
	return daemonset
}
