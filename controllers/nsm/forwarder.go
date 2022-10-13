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

type ForwarderReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	FType  string
}

func NewForwarderReconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme, ftype string) *ForwarderReconciler {
	return &ForwarderReconciler{
		Client: client,
		Log:    log,
		Scheme: scheme,
		FType:  ftype,
	}
}

func (r *ForwarderReconciler) Reconcile(ctx context.Context, nsm *nsmv1alpha1.NSM) error {

	for _, fp := range nsm.Spec.Forwarders {

		ds := &appsv1.DaemonSet{}
		Name := fp.Name
		if Name == "" {
			Name = "forwarder-" + fp.Type
		}
		err := r.Client.Get(ctx, types.NamespacedName{Name: Name, Namespace: nsm.ObjectMeta.Namespace}, ds)
		if err != nil {
			if apierrors.IsNotFound(err) {

				objectMeta := newObjectMeta(Name, "nsm", map[string]string{"app": "nsm"})
				ds = r.daemonSetForForwarder(nsm, objectMeta, r.FType)

				err = r.Client.Create(context.TODO(), ds)
				if err != nil {
					r.Log.Error(err, "failed to create deployment for "+Name)
					return err
				}
				r.Log.Info("nsm " + Name + " daemonset created")
				return nil
			}
			return err
		}
		r.Log.Info("nsm " + Name + " daemonset already exists, skipping creation")
	}
	return nil
}

func (r *ForwarderReconciler) daemonSetForForwarder(nsm *nsmv1alpha1.NSM, objectMeta metav1.ObjectMeta, FType string) *appsv1.DaemonSet {

	privmode := true
	forwarderLabel := map[string]string{"app": "forwarder", "spiffe.io/spiffe-id": "true"}

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
					HostPID:            true,
					HostNetwork:        true,
					DNSPolicy:          corev1.DNSClusterFirstWithHostNet,
					Containers: []corev1.Container{

						// forwarding plane container
						{
							Name:            objectMeta.Name,
							Image:           getForwarderImage(nsm, FType),
							ImagePullPolicy: nsm.Spec.NsmPullPolicy,
							SecurityContext: &corev1.SecurityContext{
								Privileged: &privmode,
							},
							Env:            getEnvVars(nsm, FType),
							ReadinessProbe: getReadinessProbe(FType),
							LivenessProbe:  getLivenessProbe(FType),
							StartupProbe:   getStartupProbe(FType),
							VolumeMounts:   getVolumeMounts(FType),
							Resources:      getForwarderResourceReqs(FType),
						}},
					Volumes: getVolumes(FType),
				},
			},
		},
	}

	// Set NSM instance as the owner and controller
	controllerutil.SetControllerReference(nsm, daemonset, r.Scheme)
	return daemonset
}

func getForwarderImage(nsm *nsmv1alpha1.NSM, ftype string) string {

	for _, pf := range nsm.Spec.Forwarders {
		if pf.Type == ftype {
			return pf.Image
		}
	}
	return ""
}

func getForwarderResourceReqs(FType string) corev1.ResourceRequirements {

	if FType == "vpp" {
		return corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("525m"),
				corev1.ResourceMemory: resource.MustParse("500Mi"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse("150m"),
			},
		}
	} else if FType == "ovs" {
		return corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
		}
	} else if FType == "sriov" {
		return corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("400m"),
				corev1.ResourceMemory: resource.MustParse("40Mi"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse("200m"),
			},
		}
	}
	return corev1.ResourceRequirements{}
}

func getEnvVars(nsm *nsmv1alpha1.NSM, FType string) []corev1.EnvVar {
	EnvVars := []corev1.EnvVar{
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
		{Name: "NSM_LOG_LEVEL", Value: getNsmLogLevel(nsm)},
	}
	if FType == "ovs" {
		EnvVars = append(EnvVars, corev1.EnvVar{Name: "NSM_SRIOV_CONFIG_FILE", Value: "/var/lib/networkservicemesh/smartnic.config"})
	} else if FType == "sriov" {
		EnvVars = append(EnvVars, corev1.EnvVar{Name: "NSM_SRIOV_CONFIG_FILE", Value: "/var/lib/networkservicemesh/sriov.config"})
	} else if FType == "vpp" {
		EnvVars = append(EnvVars, corev1.EnvVar{Name: "NSM_LISTEN_ON", Value: "unix:///listen.on.sock"})
		// For VPP there is no default, but later if we implement its configuration it should be added.
		//	EnvVars = append(EnvVars, corev1.EnvVar{Name: "NSM_SRIOV_CONFIG_FILE", Value: "/var/lib/networkservicemesh/sriov.config"})
	}
	return EnvVars
}

func getVolumeMounts(FType string) []corev1.VolumeMount {
	VolMounts := []corev1.VolumeMount{
		{Name: "nsm-socket",
			MountPath: "/var/lib/networkservicemesh/",
		},
		{Name: "spire-agent-socket",
			MountPath: "/run/spire/sockets",
			ReadOnly:  true,
		},
		{Name: "kubelet-socket",
			MountPath: "/var/lib/kubelet",
		},
		{Name: "cgroup",
			MountPath: "/host/sys/fs/cgroup",
		},
		{Name: "vfio",
			MountPath: "/host/dev/vfio",
		},
	}
	if FType == "vpp" {
		VolMounts = append(VolMounts, corev1.VolumeMount{Name: "vpp", MountPath: "/var/run/vpp/external"})
	}
	return VolMounts
}

func getVolumes(FType string) []corev1.Volume {

	volTypeDirOrCreate := corev1.HostPathDirectoryOrCreate
	volTypeDir := corev1.HostPathDirectory

	Volumes := []corev1.Volume{
		{
			Name: "nsm-socket",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/lib/networkservicemesh",
					Type: &volTypeDirOrCreate,
				}}},
		{
			Name: "spire-agent-socket",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/run/spire/sockets",
					Type: &volTypeDir,
				}}},
		{
			Name: "kubelet-socket",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/lib/kubelet",
					Type: &volTypeDir,
				}}},
		{
			Name: "cgroup",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/sys/fs/cgroup",
					Type: &volTypeDir,
				}}},
		{
			Name: "vfio",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/dev/vfio",
					Type: &volTypeDirOrCreate,
				}}}}

	if FType == "vpp" {
		Volumes = append(Volumes, corev1.Volume{Name: "vpp",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/run/vpp",
					Type: &volTypeDirOrCreate}}})
	}
	return Volumes
}

// Probes are set only for VPP forwarder.
// Put a simple "echo" into OVS and SR-IOV forwarder's probes.
func getReadinessProbe(FType string) *corev1.Probe {

	if FType == "vpp" {
		return &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{
						"/bin/grpc-health-probe",
						"-spiffe",
						"-addr=unix:///listen.on.sock",
					},
				},
			},
			FailureThreshold:    120,
			InitialDelaySeconds: 1,
			PeriodSeconds:       1,
			TimeoutSeconds:      2,
		}
	}
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"echo",
				},
			},
		},
	}
}

func getLivenessProbe(FType string) *corev1.Probe {
	if FType == "vpp" {
		return &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{
						"/bin/grpc-health-probe",
						"-spiffe",
						"-addr=unix:///listen.on.sock",
					},
				},
			},
			FailureThreshold:    25,
			InitialDelaySeconds: 10,
			PeriodSeconds:       5,
			TimeoutSeconds:      2,
		}
	}
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"echo",
				},
			},
		},
	}
}

func getStartupProbe(FType string) *corev1.Probe {
	if FType == "vpp" {
		return &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{
						"/bin/grpc-health-probe",
						"-spiffe",
						"-addr=unix:///listen.on.sock",
					},
				},
			},
			FailureThreshold: 25,
			PeriodSeconds:    5,
		}
	}
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"echo",
				},
			},
		},
	}
}
