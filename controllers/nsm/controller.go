/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/apis/nsm/v1alpha1"
)

// NSMReconciler reconciles a NSM object
type NSMReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=nsm.networkservicemesh.io,resources=nsms,verbs=get;list;watch;create;update;patch;delete,namespace=nsm
// +kubebuilder:rbac:groups=nsm.networkservicemesh.io,resources=nsms/status,verbs=get;update;patch,namespace=nsm
// +kubebuilder:rbac:groups=apps,resources=daemonsets;deployments;replicasets,verbs=get;list;watch;create;update;patch;delete,namespace=nsm
// +kubebuilder:rbac:groups=apps,resourceNames=nsm-operator,resources=deployments/finalizers,verbs=update,namespace=nsm
// +kubebuilder:rbac:groups=core,resources=secrets;services;services/finalizers;configmaps;events;persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete,namespace=nsm
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;create,namespace=nsm
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get,namespace=nsm
// +kubebuilder:rbac:groups=admissionregistration,resources=mutatingwebhookconfigurations;mutatingwebhookconfigurations/finalizers,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups="*",resources="*",verbs="*"

const (
	serviceAccountName  string = "nsm-operator"
	registryMemoryImage string = "ghcr.io/networkservicemesh/cmd-registry-memory"
	registryK8sImage    string = "ghcr.io/networkservicemesh/cmd-registry-k8s"
	nsmgrImage          string = "ghcr.io/networkservicemesh/cmd-nsmgr"
	exclPrefImage       string = "ghcr.io/networkservicemesh/cmd-exclude-prefixes-k8s"
	forwarderImage      string = "ghcr.io/networkservicemesh/cmd-forwarder-"
)

// Reconcile for NSMs
func (r *NSMReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	Log := log.FromContext(ctx).WithValues("nsm", req.NamespacedName)

	// Fetch the NSM instance
	nsm := &nsmv1alpha1.NSM{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, nsm)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Update the status field to creating
	if nsm.Status.Phase == nsmv1alpha1.NSMPhaseInitial {
		nsm.Status.Phase = nsmv1alpha1.NSMPhaseCreating
		if updateErr := r.Client.Status().Update(context.TODO(), nsm); updateErr != nil {
			Log.Info("Failed to update status", "Error", err.Error())
		}
	}

	// setting up default images for registry
	if nsm.Spec.Registry.Image == "" {
		switch nsm.Spec.Registry.Type {
		case "memory":
			nsm.Spec.Registry.Image = registryMemoryImage + ":" + nsm.Spec.Version
		case "k8s":
			nsm.Spec.Registry.Image = registryK8sImage + ":" + nsm.Spec.Version
		}
	}

	if nsm.Spec.Nsmgr.Image == "" {
		nsm.Spec.Nsmgr.Image = nsmgrImage + ":" + nsm.Spec.Version
	}

	if nsm.Spec.ExclPref.Image == "" {
		nsm.Spec.ExclPref.Image = exclPrefImage + ":" + nsm.Spec.Version
	}

	reconcilers := []Reconciler{
		NewRegistryReconciler(r.Client, Log, r.Scheme),
		NewRegistryServiceReconciler(r.Client, Log, r.Scheme),
		NewNsmgrReconciler(r.Client, Log, r.Scheme),
	}

	// Add admission-webhook-k8s reconciler on demand
	if nsm.Spec.Webhook.Image == "" {
		reconcilers = append(reconcilers,
			NewWebhookReconciler(r.Client, Log, r.Scheme),
			NewWebhookServiceReconciler(r.Client, Log, r.Scheme))
	}

	// Add forwarder reconcilers
	for _, pf := range nsm.Spec.Forwarders {
		reconcilers = append(reconcilers,
			NewForwarderReconciler(r.Client, Log, r.Scheme, pf.Type))
	}

	// Call all reconcilers
	for _, r := range reconcilers {
		err := r.Reconcile(ctx, nsm)
		if err != nil {
			Log.Error(err, "error while reconciling")
			return ctrl.Result{}, err
		}
	}

	// Update Status field after creating all resources
	if nsm.Status.Phase != nsmv1alpha1.NSMPhaseRunning {
		nsm.Status.Phase = nsmv1alpha1.NSMPhaseRunning
		if updateErr := r.Client.Status().Update(context.TODO(), nsm); updateErr != nil {
			Log.Info("Failed to update status", "Error", updateErr.Error())
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager registers the controlller with the manager and adds the owned resource types
func (r *NSMReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nsmv1alpha1.NSM{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.DaemonSet{}).
		Complete(r)
}

// Get value for NSM_LOG_LEVEL environment variable (defaul: "INFO")
func getNsmLogLevel(nsm *nsmv1alpha1.NSM) string {
	if nsm.Spec.NsmLogLevel != "" {
		return nsm.Spec.NsmLogLevel
	}
	return "INFO"
}
