package nsm

import "time"

const (
	version                   = "v0.2.0"
	webhookName               = "nsm-admission-webhook"
	webhookSecretName         = webhookName + "-certs"
	webhookServiceName        = webhookName + "-svc"
	webhookServicePort        = 443
	webhookServiceTargetPort  = 443
	webhookMutatingConfigName = webhookName + "-cfg"

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
