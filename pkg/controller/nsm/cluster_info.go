package nsm

import (
	"context"

	configv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type configMapData map[string]string

// Checks if cluster is OpenShift
func (r *ReconcileNSM) isPlatformOpenShift() bool {
	infra := &configv1.Infrastructure{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "cluster", Namespace: ""}, infra)
	if err != nil || errors.IsNotFound(err) {
		log.Info("Network object not found. This platform is probably not OpenShift")
		return false
	}
	// fmt.Println(infra.APIVersion)
	return true
}

// Get cluster subnets and recreate kubeadm config map for Openshift
// workaround NSMs way to capture network information
// TODO: This should be replaced by a pkg with the proper way to retrieve this information in the future
// Probably only using the Network CRD in the canonical OpenShfit API
func (r *ReconcileNSM) getNetworkConfigMap() *corev1.ConfigMap { //*configMapData {

	// get OpenShift Networks
	var clusterNetworkCIDRs []string
	network := &configv1.Network{}

	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "cluster"}, network)
	if err != nil {
		log.Info("Object not found or error")
		return nil
	}
	for _, clusterNetwork := range network.Spec.ClusterNetwork {
		clusterNetworkCIDRs = append(clusterNetworkCIDRs, clusterNetwork.CIDR)
	}

	// create kubeadm-config configMap with the network info
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubeadm-config",
			Namespace: "kube-system",
			Labels: map[string]string{
				"createdBy": "nsm-operator",
			},
		},
		Data: map[string]string{"clusterConfiguration": "networking:\n  podSubnet: " + clusterNetworkCIDRs[0] + "\n  serviceSubnet: " + network.Spec.ServiceNetwork[0] + "\n"},
	}

	// return kubeadm-config configMap
	return cm
}
