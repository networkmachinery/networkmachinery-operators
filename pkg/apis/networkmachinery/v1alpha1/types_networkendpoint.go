package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EndpointKind string

const (
	IP       EndpointKind = "ip"
	Pod      EndpointKind = "pod"
	Service  EndpointKind = "service"
	Selector EndpointKind = "selector"
	Node     EndpointKind = "node"
)

type NetworkSourceEndpoint struct {
	Kind           EndpointKind          `json:"kind"`
	Name           string                `json:"name,omitempty"`
	Namespace      string                `json:"namespace"`
	Container      string                `json:"container"`
	SourceSelector *metav1.LabelSelector `json:"sourceSelector,omitempty"`
}

type NetworkDestinationEndpoint struct {
	Kind      EndpointKind `json:"kind"`
	Name      string       `json:"name,omitempty"`
	Namespace string       `json:"namespace,omitempty"`
	IP        string       `json:"ip,omitempty"`
	Port      string       `json:"port,omitempty"`
	Node      string       `json:"node,omitempty"`
	Container string       `json:"container,omitempty"`
}
