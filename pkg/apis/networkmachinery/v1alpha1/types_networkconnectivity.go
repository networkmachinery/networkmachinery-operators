package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkConnectivityTest represents a network connectivity test
type NetworkConnectivityTest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NetworkConnectivityTestSpec `json:"spec,omitempty"`
	Status NetworkConnectivityTestStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkConnectivityTestList is a list of network notifications
type NetworkConnectivityTestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []NetworkConnectivityTest `json:"items,omitempty"`
}

type NetworkConnectivityTestSpec struct {
	Layer        string                       `json:"layer"`
	Source         NetworkSourceEndpoint      `json:"source"`
	Destinations []NetworkDestinationEndpoint `json:"destinations"`
}

type NetworkConnectivityTestStatus struct {
	TestStatus *runtime.RawExtension `json:"testStatus,omitempty"`
}


type NetworkSourceEndpoint struct {
	Name string `json:"name,omitempty"`
	Namespace string `json:"namespace"`
	Container string `json:"container"`
	SourceSelector *metav1.LabelSelector `json:"sourceSelector,omitempty"`
}


type NetworkDestinationEndpoint struct {
	Kind string `json:"kind"`
	Name string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	IP string `json:"ip,omitempty"`
	Port string `json:"port,omitempty"`
}




