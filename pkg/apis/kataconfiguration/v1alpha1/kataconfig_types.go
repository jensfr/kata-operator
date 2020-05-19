package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KataConfigSpec defines the desired state of KataConfig
type KataConfigSpec struct {
	// RuntimeClass blah blah
	RuntimeClass string `json:"runtimeClass"`

	// KataImage blah blah blah blah blah blah
	KataImage string `json:"kataImage"`
}

// KataConfigStatus defines the observed state of KataConfig
type KataConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	TotalNodesCount int `json:"totalNodesCount"`

	CompletedNodesCount int `json:"completedNodesCount"`

	InProgressNodesCount int `json:"inProgressNodesCount"`

	FailedNodes []FailedNode `json:"failedNodes"`
}

// FailedNode abc
type FailedNode struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KataConfig is the Schema for the kataconfigs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kataconfigs,scope=Namespaced
type KataConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KataConfigSpec   `json:"spec,omitempty"`
	Status KataConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KataConfigList contains a list of KataConfig
type KataConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KataConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KataConfig{}, &KataConfigList{})
}
