package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkTrafficShaper represents a network connectivity test
type NetworkTrafficShaper struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkTrafficShaperSpec   `json:"spec,omitempty"`
	Status NetworkTrafficShaperStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkTrafficShaperList is a list of network traffic shapers
type NetworkTrafficShaperList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []NetworkTrafficShaper `json:"items,omitempty"`
}

type NetworkTrafficShaperSpec struct {
	Targets []ShaperTarget `json:"targets"`
}

type ShaperTarget struct {
	Namespace      string                `json:"namespace"`
	Kind           EndpointKind          `json:"kind"`
	ShaperConfig   ShaperConfiguration   `json:"configuration"`
	Name           string                `json:"name,omitempty"`
	Container      string                `json:"container,omitempty"`
	SourceSelector *metav1.LabelSelector `json:"targetSelector,omitempty"`
}

type ShaperType string

const (
	Loss  ShaperType = "loss"
	Delay ShaperType = "delay"
)

type ShaperConfiguration struct {
	Type   ShaperType `json:"type"`
	Device string     `json:"device"`
	Value  string     `json:"value"`
}

type NetworkTrafficShaperStatus struct {
	Status `json:",inline"`
}
