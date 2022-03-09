package controllers

import (
	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/apis/nsm/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *NSMReconciler) deamonSetForForwardingPlane(nsm *nsmv1alpha1.NSM, objectMeta metav1.ObjectMeta) runtime.Object {

	volType := corev1.HostPathDirectoryOrCreate
	// mountPropagationMode := corev1.MountPropagationBidirectional
	privmode := true

	forwarderLabel := map[string]string{"app": "forwarder"}
	volTypeSpire := corev1.HostPathDirectory

	daemonset := &appsv1.DaemonSet{

		ObjectMeta: objectMeta,
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: forwarderLabel,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: forwarderLabel,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountName,
					// HostPID:            true,
					HostNetwork: true,

					Containers: []corev1.Container{

						// forwarding plane container
						{
							Name:            objectMeta.Name,
							Image:           getForwarderImage(nsm, objectMeta.Name),
							ImagePullPolicy: nsm.Spec.NsmPullPolicy,
							SecurityContext: &corev1.SecurityContext{
								Privileged: &privmode,
							},
							Env: []corev1.EnvVar{
								{Name: "SPIFFE_ENDPOINT_SOCKET", Value: "unix:///run/spire/sockets/agent.sock"},
								{Name: "NSM_TUNNEL_IP", ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "status.podIP",
									}}},
								{Name: "NSM_CONNECT_TO", Value: "unix:///var/lib/networkservicemesh/nsm.io.sock"},
								{Name: "NSM_NAME", ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.name",
									}}},

								// {Name: "JAEGER_AGENT_PORT", Value: nsm.Spec.JaegerTracing}
							},

							VolumeMounts: []corev1.VolumeMount{
								{Name: "nsm-socket",
									MountPath: "/var/lib/networkservicemesh/",
									// MountPropagation: &mountPropagationMode,
								},
								{Name: "spire-agent-socket",
									MountPath: "/run/spire/sockets",
									ReadOnly:  true,
								},
							},
						}},
					Volumes: []corev1.Volume{
						{
							Name: "nsm-socket",
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
									Type: &volTypeSpire,
								}}},
					},
				},
			},
		},
	}

	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, daemonset, r.Scheme)
	return daemonset
}

func getForwarderImage(nsm *nsmv1alpha1.NSM, name string) string {

	for _, pf := range nsm.Spec.Forwarders {
		if pf.Name == name {
			return pf.Image
		}
	}
	return ""
}
