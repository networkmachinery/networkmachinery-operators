package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkNotification represents the notification resource
type NetworkNotification struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NetworkNotificationSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkNotificationList is a list of network notifications
type NetworkNotificationList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NetworkNotificationSpec `json:"spec,omitempty"`
}

type NetworkNotificationSpec struct {
	Events []NetworkEvent `json:"events"`
}

type NetworkEvent struct {
	Flow  Flow  `json:"flow"`
	Event Event `json:"event"`
}

type Event struct {
	ID         string  `json:"id"`
	Time       string  `json:"time"`
	Name       string  `json:"name"`
	Metric     string  `json:"metric"`
	Threshold  int32   `json:"threshold"`
	Value      float32 `json:"value"`
	Agent      float32 `json:"agent"`
	DataSource int32   `json:"datasource"`
}
