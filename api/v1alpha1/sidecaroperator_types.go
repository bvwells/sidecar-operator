package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SidecarOperatorSpec defines the desired state of SidecarOperator
type SidecarOperatorSpec struct {
	// Image is the image to inject as a sidecar.
	Image string `json:"image,omitempty"`
}

// SidecarOperatorStatus defines the observed state of SidecarOperator
type SidecarOperatorStatus struct {

	// Status is the status of the operator.
	Status string `json:"status"`
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
