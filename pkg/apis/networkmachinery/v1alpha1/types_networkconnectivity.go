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

	Spec   NetworkConnectivityTestSpec   `json:"spec,omitempty"`
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
	Source       NetworkSourceEndpoint        `json:"source"`
	Destinations []NetworkDestinationEndpoint `json:"destinations"`
	Frequency    string                       `json:"frequency,omitempty"`
}

type NetworkConnectivityTestStatus struct {
	TestStatus *runtime.RawExtension `json:"testStatus,omitempty"`
}
