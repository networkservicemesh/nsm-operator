package controllers

import (
	"context"

	"github.com/go-logr/logr"
	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/apis/nsm/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
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
			r.Log.Info("nsm registry deployment created")
			return nil
		}
		return err
	}
	r.Log.Info("nsm registry deployment already exists, skipping creation")
	return nil
}

func (r *RegistryReconciler) DeploymentForRegistry(nsm *nsmv1alpha1.NSM) *appsv1.Deployment {

	objectMeta := newObjectMeta("nsm-registry", "nsm", map[string]string{"app": "nsm"})

	registryLabel := map[string]string{"app": "nsm-registry", "spiffe.io/spiffe-id": "true"}
	volTypeDirectory := corev1.HostPathDirectory

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
						Env:             *getEnvVar(nsm),
						Ports: []corev1.ContainerPort{{
							ContainerPort: 5002,
							HostPort:      5002}},
						VolumeMounts: []corev1.VolumeMount{
							{Name: "spire-agent-socket",
								MountPath: "/run/spire/sockets",
							}},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("200m"),
								corev1.ResourceMemory: resource.MustParse("40Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("100m"),
							},
						},
					}},
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
	if nsm.Spec.Registry.EnvVars != nil {
		return &nsm.Spec.Registry.EnvVars
	}
	prefix := "NSM_"
	switch nsm.Spec.Registry.Type {
	case "memory":
		prefix = "REGISTRY_MEMORY_"
	case "k8s":
		// From version 1.7.0 the prefix of the environment variables changed to NSM, instead of REGISTRY_K8S
		if nsm.Spec.Version < "v1.7.0" {
			prefix = "REGISTRY_K8S_"
		}
	}

	return &[]corev1.EnvVar{{Name: "SPIFFE_ENDPOINT_SOCKET", Value: "unix:///run/spire/sockets/agent.sock"},
		{Name: prefix + "LISTEN_ON", Value: "tcp://:5002"},
		{Name: prefix + "PROXY_REGISTRY_URL", Value: "nsmgr-proxy:5004"},
		{Name: prefix + "LOG_LEVEL", Value: getNsmLogLevel(nsm)},
		{Name: prefix + "NAMESPACE", ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		}},
	}
}
