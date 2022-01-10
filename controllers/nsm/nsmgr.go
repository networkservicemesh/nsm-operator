package controllers

import (
	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/apis/nsm/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *NSMReconciler) deamonSetForNSMGR(nsm *nsmv1alpha1.NSM, objectMeta metav1.ObjectMeta) runtime.Object {

	volType := corev1.HostPathDirectoryOrCreate
	volTypeSpire := corev1.HostPathDirectory
	privmode := true

	nsmgrLabel := map[string]string{"app": "nsmgr"}

	daemonset := &appsv1.DaemonSet{
		ObjectMeta: objectMeta,
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: nsmgrLabel,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: nsmgrLabel,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountName,
					Containers: []corev1.Container{

						// nsmdp container
						{
							Name:            "nsmgr",
							Image:           nsm.Spec.NsmgrImage + ":" + nsm.Spec.Version,
							ImagePullPolicy: nsm.Spec.NsmPullPolicy,
							SecurityContext: &corev1.SecurityContext{
								Privileged: &privmode,
							},
							Ports: []corev1.ContainerPort{{
								ContainerPort: 5001,
								HostPort:      5001}},

							Env: []corev1.EnvVar{
								{Name: "SPIFFE_ENDPOINT_SOCKET", Value: "unix:///run/spire/sockets/agent.sock"},
								{Name: "NSM_NAME", ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.name",
									}}},
								{Name: "NSM_REGISTRY_URL", Value: "nsm-registry-svc:5002"},
								{Name: "POD_IP", ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "status.podIP",
									}}},
								{Name: "NSM_LISTEN_ON", Value: "unix:///var/lib/networkservicemesh/nsm.io.sock,tcp://:5001"},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "nsm-socket",
									MountPath: "/var/lib/networkservicemesh",
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
