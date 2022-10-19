package controllers

import (
	"context"

	"github.com/go-logr/logr"
	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/apis/nsm/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type NsmgrReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func NewNsmgrReconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme) *NsmgrReconciler {
	return &NsmgrReconciler{
		Client: client,
		Log:    log,
		Scheme: scheme,
	}
}

func (r *NsmgrReconciler) Reconcile(ctx context.Context, nsm *nsmv1alpha1.NSM) error {

	ds := &appsv1.DaemonSet{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: "nsmgr", Namespace: nsm.ObjectMeta.Namespace}, ds)
	if err != nil {
		if apierrors.IsNotFound(err) {
			ds = r.daemonSetForNSMGR(nsm)
			err = r.Client.Create(context.TODO(), ds)
			if err != nil {
				r.Log.Error(err, "failed to create daemonset for nsm-registry")
				return err
			}
			return nil
		}
		return err
	}
	r.Log.Info("nsm nsmgr daemonset already exists, skipping creation")
	return nil
}

func (r *NsmgrReconciler) daemonSetForNSMGR(nsm *nsmv1alpha1.NSM) *appsv1.DaemonSet {

	objectMeta := newObjectMeta("nsmgr", "nsm", map[string]string{"app": "nsm"})

	volType := corev1.HostPathDirectoryOrCreate
	volTypeSpire := corev1.HostPathDirectory
	privmode := true

	nsmgrLabel := map[string]string{"app": "nsmgr", "spiffe.io/spiffe-id": "true"}

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
