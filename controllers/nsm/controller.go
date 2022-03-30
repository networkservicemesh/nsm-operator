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

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/apis/nsm/v1alpha1"
)

// NSMReconciler reconciles a NSM object
type NSMReconciler struct {
	client.Client
	Log    logr.Logger
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
	registryMemoryImage string = "ghcr.io/networkservicemesh/cmd-registry-memory:latest"
	registryK8sImage    string = "ghcr.io/networkservicemesh/ci/cmd-registry-k8s:latest"
	nsmgrImage          string = "ghcr.io/networkservicemesh/cmd-nsmgr"
)

// Reconcile for NSMs
func (r *NSMReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	_ = context.Background()
	reqLogger := r.Log.WithValues("nsm", req.NamespacedName)

	// Fetch the NSM instance
	nsm := &nsmv1alpha1.NSM{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, nsm)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// setting up deafult images for registry
	if nsm.Spec.Registry.Image == "" {
		switch nsm.Spec.Registry.Type {
		case "memory":
			nsm.Spec.Registry.Image = registryMemoryImage
		case "k8s":
			nsm.Spec.Registry.Image = registryK8sImage
		}
	}

	if nsm.Spec.NsmgrImage == "" {
		nsm.Spec.NsmgrImage = nsmgrImage
	}

	// Update the status field to creating
	if nsm.Status.Phase == nsmv1alpha1.NSMPhaseInitial {
		nsm.Status.Phase = nsmv1alpha1.NSMPhaseCreating
		if updateErr := r.Client.Status().Update(context.TODO(), nsm); updateErr != nil {
			reqLogger.Info("Failed to update status", "Error", err.Error())
		}
	}

	// Reconcile Deployment for registry-memory
	deploymentForNsmRegistry := &appsv1.Deployment{}
	objectMeta := setObjectMeta("nsm-registry", "nsm", map[string]string{"app": "nsm"})
	r.reconcileResource(r.deploymentForRegistryMemory, nsm, deploymentForNsmRegistry, objectMeta)

	// Reconcile Deployment for registry-service
	svcForRegistry := &corev1.Service{}
	objectMeta = setObjectMeta("nsm-registry-svc", "nsm", map[string]string{"app": "nsm"})
	r.reconcileResource(r.serviceForNsmRegistry, nsm, svcForRegistry, objectMeta)

	for _, fp := range nsm.Spec.Forwarders {
		// Reconcile Daemonset for forwarder
		dsForFP := &appsv1.DaemonSet{}
		objectMeta = setObjectMeta(fp.Name, "nsm", map[string]string{"app": "nsm"})
		r.reconcileResource(r.deamonSetForForwardingPlane, nsm, dsForFP, objectMeta)
	}

	// Reconcile Daemonset for nsmgr
	dsForNsmgr := &appsv1.DaemonSet{}
	objectMeta = setObjectMeta("nsmgr", "nsm", map[string]string{"app": "nsm"})
	r.reconcileResource(r.deamonSetForNSMGR, nsm, dsForNsmgr, objectMeta)

	// Update Status field after creating all resources
	if nsm.Status.Phase != nsmv1alpha1.NSMPhaseRunning {
		nsm.Status.Phase = nsmv1alpha1.NSMPhaseRunning
		if updateErr := r.Client.Status().Update(context.TODO(), nsm); updateErr != nil {
			reqLogger.Info("Failed to update status", "Error", err.Error())
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
