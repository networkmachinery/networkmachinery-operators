package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Protocol string

const (
	TCP Protocol = "tcp_lat"
	UDP Protocol = "udp_lat"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkDelayTest represents a network connectivity test
type NetworkDelayTest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkDelayTestSpec   `json:"spec,omitempty"`
	Status NetworkDelayTestStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkDelayTestList is a list of network notifications
type NetworkDelayTestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []NetworkDelayTest `json:"items,omitempty"`
}

type NetworkDelayTestSpec struct {
	Protocol     string                       `json:"protocol"`
	Source       NetworkSourceEndpoint        `json:"source"`
	Destinations []NetworkDestinationEndpoint `json:"destinations"`
	Frequency    string                       `json:"frequency,omitempty"`
}

type NetworkDelayTestStatus struct {
	TestStatus *runtime.RawExtension `json:"testStatus,omitempty"`
}
