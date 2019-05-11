package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkConnectivityTest represents a network connectivity test
type NetworkConnectivityTest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NetworkConnectivityTestSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkConnectivityTestList is a list of network notifications
type NetworkConnectivityTestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []NetworkConnectivityTest `json:"items,omitempty"`
}

type NetworkConnectivityTestSpec struct {
	Tests []NetworkTest `json:"tests"`
}

type NetworkEndpointSource struct {
	Name string `json:"name,omitempty"`
	IP   string `json:"ip,omitempty"`
}

type NetworkEndpointDestination struct {
	Name string `json:"name,omitempty"`
	IP   string `json:"ip,omitempty"`
	Port string `json:"port,omitempty"`
}

type NetworkTest struct {
	Layer        string                       `json:"layer"`
	Sources      []NetworkEndpointSource      `json:"sources"`
	Destinations []NetworkEndpointDestination `json:"destinations"`
}
