package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NSMSpec defines the desired state of NSM
// +k8s:openapi-gen=true
type NSMSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// nsmgr configs true or false
	Insecure bool `json:"insecure"`

	// Forwarding plane configs
	ForwardingPlaneName  string `json:"forwardingPlaneName"`
	ForwardingPlaneImage string `json:"forwardingPlaneImage"`

	// Version field for reference on Openshift UI
	Version string `json:"version"`

	// Enable Spire true or false - for future release
	// Spire bool `json:"spire"`

	// Enable Tracing true or false
	// JaegerTracing bool `json:"jaegerTracing"`

}

// NSMPhase is the type for the operator phases
type NSMPhase string

// Operator phases
const (
	NSMPhaseInitial     NSMPhase = ""
	NSMPhasePending     NSMPhase = "Pending"
	NSMPhaseCreating    NSMPhase = "Creating"
	NSMPhaseRunning     NSMPhase = "Running"
	NSMPhaseTerminating NSMPhase = "Terminating"
)

// NSMStatus defines the observed state of NSM
// +k8s:openapi-gen=true
type NSMStatus struct {
	// Operator phases during deployment
	Phase NSMPhase `json:"phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NSM is the Schema for the nsms API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nsms,scope=Namespaced
type NSM struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NSMSpec   `json:"spec,omitempty"`
	Status NSMStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NSMList contains a list of NSM
type NSMList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NSM `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NSM{}, &NSMList{})
}
