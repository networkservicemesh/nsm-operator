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

type WebhookReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func NewWebhookReconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme) *WebhookReconciler {
	return &WebhookReconciler{
		Client: client,
		Log:    log,
		Scheme: scheme,
	}
}

func (r *WebhookReconciler) Reconcile(ctx context.Context, nsm *nsmv1alpha1.NSM) error {

	deploy := &appsv1.Deployment{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: "admission-webhook-k8s", Namespace: nsm.ObjectMeta.Namespace}, deploy)
	if err != nil {
		if apierrors.IsNotFound(err) {
			deploy = r.DeploymentForWebhook(nsm)
			err = r.Client.Create(context.TODO(), deploy)
			if err != nil {
				r.Log.Error(err, "failed to create deployment for admission-webhook-k8s")
				return err
			}
			r.Log.Info("admission-webhook-k8s deployment created")
			return nil
		}
		return err
	}
	r.Log.Info("admission-webhook-k8s deployment already exists, skipping creation")
	return nil
}

func (r *WebhookReconciler) DeploymentForWebhook(nsm *nsmv1alpha1.NSM) *appsv1.Deployment {

	privmode := true

	objectMeta := newObjectMeta("admission-webhook-k8s", "nsm", map[string]string{"app": "nsm"})
	webhookLabel := map[string]string{"app": "admission-webhook-k8s"}

	envVars := []corev1.EnvVar{}
	if nsm.Spec.Webhook.EnvVars != nil {
		envVars = nsm.Spec.Nsmgr.EnvVars
	} else {
		envVars = []corev1.EnvVar{
			{Name: "SPIFFE_ENDPOINT_SOCKET", Value: "unix:///run/spire/sockets/agent.sock"},
			{Name: "NSM_SERVICE_NAME", Value: "admission-webhook-svc"},
			{Name: "NSM_NAME", ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				}}},
			{Name: "NSM_NAMESPACE", ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				}}},
			{Name: "NSM_ANNOTATION", Value: "networkservicemesh.io"},
			{Name: "NSM_CONTAINER_IMAGES", Value: "ghcr.io/networkservicemesh/cmd-nsc:" + nsm.Spec.Version},
			{Name: "NSM_INIT_CONTAINER_IMAGES", Value: "ghcr.io/networkservicemesh/cmd-nsc-init:" + nsm.Spec.Version},
			{Name: "NSM_LABELS", Value: "spiffe.io/spiffe-id:true"},
			{Name: "NSM_ENVS", Value: "NSM_LOG_LEVEL=" + getNsmLogLevel(nsm)},
		}
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: objectMeta,
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: webhookLabel,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: webhookLabel,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountName,
					Containers: []corev1.Container{{
						Name:            "admission-webhook-k8s",
						Image:           nsm.Spec.Webhook.Image,
						ImagePullPolicy: nsm.Spec.NsmPullPolicy,
						SecurityContext: &corev1.SecurityContext{
							Privileged: &privmode,
						},
						Env: envVars,
					}},
				},
			},
		},
	}
	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, deploy, r.Scheme)
	return deploy
}
