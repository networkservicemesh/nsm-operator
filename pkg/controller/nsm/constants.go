package nsm

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

const (
	nsmRegistry   = "docker.io"
	nsmOrg        = "networkservicemesh"
	nsmVersion    = "latest"
	nsmPullPolicy = corev1.PullIfNotPresent

	webhookName                     = "nsm-admission-webhook"
	webhookSecretName               = webhookName + "-certs"
	webhookServiceName              = webhookName + "-svc"
	webhookServicePort              = 443
	webhookServiceTargetPort        = 443
	webhookMutatingConfigName       = webhookName + "-cfg"
	webhookReplicas           int32 = 1
	// Deployment inputs for liveness and readiness probes to pods
	probePort         = 5555
	probeInitialDelay = 10
	probePeriod       = 10
	probeTimeout      = 3

	// TLS Certs configuration for webhook
	rsaBits  = 2048
	validFor = 365 * 24 * time.Hour
)

func labelsForNSMAdmissionWebhook(crName string) map[string]string {
	return map[string]string{"app": webhookName, "nsm-cr": crName}
}
