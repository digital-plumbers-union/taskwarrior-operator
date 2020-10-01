package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TaskwarriorSpec defines the desired state of Taskwarrior
type TaskwarriorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Size int32 `json:"size"`
}

// TaskwarriorStatus defines the observed state of Taskwarrior
type TaskwarriorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Taskwarrior is the Schema for the taskwarriors API
type Taskwarrior struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskwarriorSpec   `json:"spec,omitempty"`
	Status TaskwarriorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TaskwarriorList contains a list of Taskwarrior
type TaskwarriorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Taskwarrior `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Taskwarrior{}, &TaskwarriorList{})
}
