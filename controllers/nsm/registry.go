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

type RegistryReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func NewRegistryReconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme) *RegistryReconciler {
	return &RegistryReconciler{
		Client: client,
		Log:    log,
		Scheme: scheme,
	}
}

func (r *RegistryReconciler) Reconcile(ctx context.Context, nsm *nsmv1alpha1.NSM) error {

	deploy := &appsv1.Deployment{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: "nsm-registry", Namespace: nsm.ObjectMeta.Namespace}, deploy)
	if err != nil {
		if apierrors.IsNotFound(err) {
			deploy = r.DeploymentForRegistry(nsm)
			err = r.Client.Create(context.TODO(), deploy)
			if err != nil {
				r.Log.Error(err, "failed to create deployment for nsm-registry")
				return err
			}
			return nil
		}
		return err
	}
	r.Log.Info("nsm registry deployment already exists, skipping creation")
	return nil
}

func (r *RegistryReconciler) DeploymentForRegistry(nsm *nsmv1alpha1.NSM) *appsv1.Deployment {

	privmode := true

	objectMeta := newObjectMeta("nsm-registry", "nsm", map[string]string{"app": "nsm"})

	registryLabel := map[string]string{"app": "nsm-registry", "spiffe.io/spiffe-id": "true"}
	volTypeDirectory := corev1.HostPathDirectory

	envVar := getEnvVar(nsm)

	deploy := &appsv1.Deployment{
		ObjectMeta: objectMeta,
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: registryLabel,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: registryLabel,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountName,
					Containers: []corev1.Container{{
						Name:            "nsm-registry",
						Image:           nsm.Spec.Registry.Image,
						ImagePullPolicy: nsm.Spec.NsmPullPolicy,
						SecurityContext: &corev1.SecurityContext{
							Privileged: &privmode,
						},
						Env: *envVar,
						Ports: []corev1.ContainerPort{{
							ContainerPort: 5002,
							HostPort:      5002}},
						VolumeMounts: []corev1.VolumeMount{
							{Name: "spire-agent-socket",
								MountPath: "/run/spire/sockets",
							}}}},
					Volumes: []corev1.Volume{{
						Name: "spire-agent-socket",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/run/spire/sockets",
								Type: &volTypeDirectory,
							}}},
					},
				},
			},
		},
	}
	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, deploy, r.Scheme)
	return deploy
}

func getEnvVar(nsm *nsmv1alpha1.NSM) *[]corev1.EnvVar {

	switch nsm.Spec.Registry.Type {

	case "memory":
		return &[]corev1.EnvVar{
			{Name: "SPIFFE_ENDPOINT_SOCKET", Value: "unix:///run/spire/sockets/agent.sock"},
			{Name: "REGISTRY_MEMORY_LISTEN_ON", Value: "tcp://:5002"},
			{Name: "REGISTRY_MEMORY_PROXY_REGISTRY_URL", Value: "nsmgr-proxy:5004"},
			{Name: "REGISTRY_MEMORY_LOG_LEVEL", Value: "TRACE"},
		}
	case "k8s":
		return &[]corev1.EnvVar{
			{Name: "SPIFFE_ENDPOINT_SOCKET", Value: "unix:///run/spire/sockets/agent.sock"},
			{Name: "REGISTRY_K8S_LISTEN_ON", Value: "tcp://:5002"},
			{Name: "REGISTRY_K8S_PROXY_REGISTRY_URL", Value: "nsmgr-proxy:5004"},
			{Name: "REGISTRY_K8S_LOG_LEVEL", Value: "TRACE"},
			{Name: "REGISTRY_K8S_NAMESPACE", ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				},
			}},
		}
	}

	return nil
}
