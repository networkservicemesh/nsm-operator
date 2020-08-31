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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NSMSpec defines the desired state of NSM
type NSMSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

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
type NSMStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Operator phases during deployment
	Phase NSMPhase `json:"phase"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// NSM is the Schema for the nsms API
type NSM struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NSMSpec   `json:"spec,omitempty"`
	Status NSMStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NSMList contains a list of NSM
type NSMList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NSM `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NSM{}, &NSMList{})
}
