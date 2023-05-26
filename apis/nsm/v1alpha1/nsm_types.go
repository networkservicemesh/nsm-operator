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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Forwarder struct {
	// Forwarder type
	// +kubebuilder:validation:Enum=vpp;ovs;sriov
	Type ForwarderType `json:"type"`
	// Forwarder descriptive name
	// (if empty then "forwarder-<type>" is used)
	Name string `json:"name,omitempty"`
	// Forwarder image string
	// (must be a complete image path with tag)
	Image string `json:"image,omitempty"`
	// EnvVars for Forwarder configuration
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`
}

// ForwarderType is the type of the forwarder
type ForwarderType string

// Forwarder types
const (
	ForwarderOvs   ForwarderType = "ovs"
	ForwarderSriov ForwarderType = "sriov"
	ForwarderVpp   ForwarderType = "vpp"
)

type Registry struct {
	// Number of replicas for the NSM Registry
	ReplicaCount int32 `json:"replicaCount,omitempty"`
	// Registry type
	// +kubebuilder:validation:Enum=k8s;memory
	Type string `json:"type"`
	// Registry Image with tag
	Image string `json:"image,omitempty"`
	// EnvVars for Registry configuration
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`
}

// Webhook
type Webhook struct {
	// admission-webhook-k8s image string
	// (must be a complete image path with tag)
	Image string `json:"image,omitempty"`
	// EnvVars for Webhook configuration
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`
}

type Nsmgr struct {
	// NSMGR image string
	// (must be a complete image path with tag)
	Image string `json:"image,omitempty"`
	// EnvVars for Nsmgr configuration
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`
}

type ExclPref struct {
	// exclude-prefixes-k8s image string
	// (must be a complete image path with tag)
	Image string `json:"exclPrefImage,omitempty"`
	// EnvVars for ExclPrefImage configuration
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`
}

// NSMSpec defines the desired state of NSM
type NSMSpec struct {
	// Tag represents the desired Network Service Mesh version
	Version string `json:"version"`
	// Pull policy for NSM images, defaults to IfNotPresent
	NsmPullPolicy corev1.PullPolicy `json:"nsmPullPolicy,omitempty"`
	// Log level of the NSM components, defaults to "INFO"
	NsmLogLevel string `json:"nsmLogLevel,omitempty"`
	// SPIRE agent socket for NSM components, must be set
	// according to the socket_path parameter of spire-agent
	SpireAgentSocket string `json:"spireAgentSocket,omitempty"`
	// Webhook for NSM
	Webhook Webhook `json:"webhook,omitempty"`
	// Registry for NSM
	Registry Registry `json:"registry"`
	// Network Service Manager
	Nsmgr Nsmgr `json:"nsmgr,omitempty"`
	// Exclude-prefixes-k8s
	ExclPref ExclPref `json:"exclPref,omitempty"`
	// List of forwarders to be used with NSM
	Forwarders []Forwarder `json:"forwarders"`
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
	// Operator phases during deployment
	Phase NSMPhase `json:"phase"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nsms

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
