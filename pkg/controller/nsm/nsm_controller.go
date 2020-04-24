package nsm

import (
	"context"
	"time"

	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/pkg/apis/nsm/v1alpha1"
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_nsm")
var caCert []byte

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new NSM Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNSM{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("nsm-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource NSM
	err = c.Watch(&source.Kind{Type: &nsmv1alpha1.NSM{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	//Watch for secondary resource admission webhook deployment
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &nsmv1alpha1.NSM{},
	})
	if err != nil {
		return err
	}

	//Watch for secondary resource admission webhook secret
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &nsmv1alpha1.NSM{},
	})
	if err != nil {
		return err
	}

	//Watch for secondary resource admission webhook service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &nsmv1alpha1.NSM{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &admissionregv1beta1.MutatingWebhookConfiguration{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &nsmv1alpha1.NSM{},
	})
	if err != nil {
		return err
	}

	// Watch for secondary resources nsmgr and forwading plane deamonsets
	err = c.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &nsmv1alpha1.NSM{},
	})
	if err != nil {
		return err
	}
	return nil
}

// blank assignment to verify that ReconcileNSM implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNSM{}

// ReconcileNSM reconciles a NSM object
type ReconcileNSM struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NSM object and makes changes based on the state read
// and what is in the NSM.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNSM) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling NSM")

	// Fetch the NSM instance
	nsm := &nsmv1alpha1.NSM{}
	err := r.client.Get(context.TODO(), request.NamespacedName, nsm)

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
		if updateErr := r.client.Status().Update(context.TODO(), nsm); updateErr != nil {
			reqLogger.Info("Failed to update status", "Error", err.Error())
		}
	}

	// reconcile secrets for admission webhook
	secret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: webhookSecretName, Namespace: nsm.Namespace}, secret)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Secret
		secret := r.secretForWebhook(nsm)
		reqLogger.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		err = r.client.Create(context.TODO(), secret)
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

	// reconcile deployment for admission webhook
	deploy := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: webhookName, Namespace: nsm.Namespace}, deploy)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		deploy := r.deploymentForWebhook(nsm)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
		err = r.client.Create(context.TODO(), deploy)
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

	// reconcile service for admission webhook
	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: webhookServiceName, Namespace: nsm.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		service := r.serviceForWebhook(nsm)
		reqLogger.Info("Creating a new service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
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

	// reconcile mutatingConfig for admission webhook
	mutatingConfig := &admissionregv1beta1.MutatingWebhookConfiguration{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: webhookMutatingConfigName}, mutatingConfig)
	if err != nil && errors.IsNotFound(err) {
		// Define a new mutatingConfig
		mutatingConfig := r.mutatingConfigForWebhook(nsm)
		reqLogger.Info("Creating a new mutatingConfig", "MutatingConfig.Namespace", mutatingConfig.Namespace, "MutatingConfig.Name", mutatingConfig.Name)
		err = r.client.Create(context.TODO(), mutatingConfig)

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

	// reconcile daemonset for nsmgr
	daemonsetForNSMGR := &appsv1.DaemonSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "nsmgr", Namespace: nsm.Namespace}, daemonsetForNSMGR)
	if err != nil && errors.IsNotFound(err) {
		// Define a new daemonsetForNSMGR
		daemonsetForNSMGR := r.deamonSetForNSMGR(nsm)
		reqLogger.Info("Creating a new daemonsetForNSMGR", "Daemonset.Namespace", daemonsetForNSMGR.Namespace, "Daemonset.Name", daemonsetForNSMGR.Name)
		err = r.client.Create(context.TODO(), daemonsetForNSMGR)
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

	// reconcile daemonset for forwarding plane
	daemonsetForFP := &appsv1.DaemonSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "nsm-" + nsm.Spec.ForwardingPlaneName + "-forwarder", Namespace: nsm.Namespace}, daemonsetForFP)
	if err != nil && errors.IsNotFound(err) {
		// Define a new daemonsetForFP
		daemonsetForFP := r.deamonSetForForwardingPlane(nsm)
		reqLogger.Info("Creating a new daemonsetForFP", "Daemonset.Namespace", daemonsetForFP.Namespace, "Daemonset.Name", daemonsetForFP.Name)
		err = r.client.Create(context.TODO(), daemonsetForFP)
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
		if updateErr := r.client.Status().Update(context.TODO(), nsm); updateErr != nil {
			reqLogger.Info("Failed to update status", "Error", err.Error())
		}
	}

	return reconcile.Result{}, nil
}
