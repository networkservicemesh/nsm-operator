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

	"fmt"
	"time"

	"github.com/go-logr/logr"
	admissionregv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/apis/nsm/v1alpha1"
)

// caCert variable holds the TLS Certificates to the mutatingWebhookConfiguration
var caCert []byte

// NSMReconciler reconciles a NSM object
type NSMReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=nsm.networkservicemesh.io,resources=nsms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nsm.networkservicemesh.io,resources=nsms/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=daemonsets;deployments;replicasets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resourceNames=nsm-operator,resources=deployments/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets;services;services/finalizers;configmaps;events;persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;create
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get

// +kubebuilder:rbac:groups=nsm.networkservicemesh.io,resources=nsms,verbs=get;list;watch;create;update;patch;delete,namespace=nsm-system
// +kubebuilder:rbac:groups=nsm.networkservicemesh.io,resources=nsms/status,verbs=get;update;patch,namespace=nsm-system
// +kubebuilder:rbac:groups=apps,resources=daemonsets;deployments;replicasets,verbs=get;list;watch;create;update;patch;delete,namespace=nsm-system
// +kubebuilder:rbac:groups=apps,resourceNames=nsm-operator,resources=deployments/finalizers,verbs=update,namespace=nsm-system
// +kubebuilder:rbac:groups=core,resources=secrets;services;services/finalizers;configmaps;events;persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete,namespace=nsm-system
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;create,namespace=nsm-system
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get,namespace=nsm-system
// +kubebuilder:rbac:groups=admissionregistration,resources=mutatingwebhookconfigurations;mutatingwebhookconfigurations/finalizers,verbs=get;list;watch;create;update;patch;delete

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

	// Update the status field to creating
	if nsm.Status.Phase == nsmv1alpha1.NSMPhaseInitial {
		nsm.Status.Phase = nsmv1alpha1.NSMPhaseCreating
		if updateErr := r.Client.Status().Update(context.TODO(), nsm); updateErr != nil {
			reqLogger.Info("Failed to update status", "Error", err.Error())
		}
	}
	// if it is OpenShift, create the kubeadm configMap with the network prefixes in use
	if r.isPlatformOpenShift() {
		cm := &corev1.ConfigMap{}
		err = r.Client.Get(context.TODO(), types.NamespacedName{Name: "kubeadm-config", Namespace: "kube-system"}, cm)
		fmt.Print(err)
		if err != nil && errors.IsNotFound(err) {
			cm = r.getNetworkConfigMap()
			err = r.Client.Create(context.TODO(), cm)
			if err != nil {
				reqLogger.Error(err, "Failed to create kubeadm ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
				return reconcile.Result{}, nil
			}
			// ConfigMap created successfully - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Reconcile the Admission Webhook Secret Containing the CABundle data if platform is not OpenShift
	// OpenShift uses the service-ca operator to get the secret
	if !r.isPlatformOpenShift() {
		secret := &corev1.Secret{}
		err = r.Client.Get(context.TODO(), types.NamespacedName{Name: webhookSecretName, Namespace: nsm.Namespace}, secret)
		if err != nil && errors.IsNotFound(err) {
			// Define a new Secret
			secret := r.secretForWebhook(nsm)
			reqLogger.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
			err = r.Client.Create(context.TODO(), secret)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
				return reconcile.Result{}, nil
			}
			// Secret created successfully - return and requeue
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to get Secret")
			return reconcile.Result{}, err
		}
	}
	// Reconcile the Admission Webhook Service
	service := &corev1.Service{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: webhookServiceName, Namespace: nsm.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		service := r.serviceForWebhook(nsm)
		reqLogger.Info("Creating a new service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.Client.Create(context.TODO(), service)
		time.Sleep(500 * time.Millisecond)
		if err != nil {
			reqLogger.Error(err, "Failed to create new service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return reconcile.Result{}, err
		}
		// service created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get service")
		return reconcile.Result{}, err
	}

	// Reconcile the Admission Webhook Deployment
	deploy := &appsv1.Deployment{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: webhookName, Namespace: nsm.Namespace}, deploy)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		deploy := r.deploymentForWebhook(nsm)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
		err = r.Client.Create(context.TODO(), deploy)
		time.Sleep(500 * time.Millisecond)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// Reconcile the Mutating Webhook Configuration Object
	mutatingConfig := &admissionregv1.MutatingWebhookConfiguration{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: webhookMutatingConfigName}, mutatingConfig)
	if err != nil && errors.IsNotFound(err) {
		// Define a new mutatingConfig
		mutatingConfig := r.mutatingConfigForWebhook(nsm)
		reqLogger.Info("Creating a new mutatingConfig", "MutatingConfig.Namespace", mutatingConfig.Namespace, "MutatingConfig.Name", mutatingConfig.Name)
		err = r.Client.Create(context.TODO(), mutatingConfig)

		if err != nil {
			reqLogger.Error(err, "Failed to create new mutatingConfig", "MutatingConfig.Namespace", mutatingConfig.Namespace, "MutatingConfig.Name", mutatingConfig.Name)
			return reconcile.Result{}, err
		}
		// mutatingConfig created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get mutatingConfig")
		return reconcile.Result{}, err
	}

	// Reconcile the Network Service Manager
	daemonsetForNSMGR := &appsv1.DaemonSet{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: "nsmgr", Namespace: nsm.Namespace}, daemonsetForNSMGR)
	if err != nil && errors.IsNotFound(err) {
		// Define a new daemonsetForNSMGR
		daemonsetForNSMGR := r.deamonSetForNSMGR(nsm)
		reqLogger.Info("Creating a new daemonsetForNSMGR", "Daemonset.Namespace", daemonsetForNSMGR.Namespace, "Daemonset.Name", daemonsetForNSMGR.Name)
		err = r.Client.Create(context.TODO(), daemonsetForNSMGR)
		if err != nil {
			reqLogger.Error(err, "Failed to create new daemonsetForNSMGR", "Daemonset.Namespace", daemonsetForNSMGR.Namespace, "Daemonset.Name", daemonsetForNSMGR.Name)
			return reconcile.Result{}, err
		}
		// daemonsetForNSMGR created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get daemonsetForNSMGR")
		return reconcile.Result{}, err
	}

	// Reconcile the Forwarding Plane DaemonSet
	daemonsetForFP := &appsv1.DaemonSet{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: "nsm-" + nsm.Spec.ForwardingPlaneName + "-forwarder", Namespace: nsm.Namespace}, daemonsetForFP)
	if err != nil && errors.IsNotFound(err) {
		// Define a new daemonsetForFP
		daemonsetForFP := r.deamonSetForForwardingPlane(nsm)
		reqLogger.Info("Creating a new daemonsetForFP", "Daemonset.Namespace", daemonsetForFP.Namespace, "Daemonset.Name", daemonsetForFP.Name)
		err = r.Client.Create(context.TODO(), daemonsetForFP)
		if err != nil {
			reqLogger.Error(err, "Failed to create new daemonsetForFP", "Daemonset.Namespace", daemonsetForFP.Namespace, "Daemonset.Name", daemonsetForFP.Name)
			return reconcile.Result{}, err
		}
		// daemonsetForFP created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get daemonsetForFP")
		return reconcile.Result{}, err
	}

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
		Owns(&corev1.Secret{}).
		Owns(&admissionregv1.MutatingWebhookConfiguration{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
