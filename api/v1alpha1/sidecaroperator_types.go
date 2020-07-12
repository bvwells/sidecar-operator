package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SidecarOperatorSpec defines the desired state of SidecarOperator
type SidecarOperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Image is the image to inject as a sidecar.
	Image string `json:"image,omitempty"`
}

// SidecarOperatorStatus defines the observed state of SidecarOperator
type SidecarOperatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SidecarOperator is the Schema for the sidecaroperators API
type SidecarOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SidecarOperatorSpec   `json:"spec,omitempty"`
	Status SidecarOperatorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SidecarOperatorList contains a list of SidecarOperator
type SidecarOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SidecarOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SidecarOperator{}, &SidecarOperatorList{})
}
