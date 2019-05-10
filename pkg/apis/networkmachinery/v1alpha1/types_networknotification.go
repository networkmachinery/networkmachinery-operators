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
	Event NetworkEvent `json:"networkEvent"`
}

type NetworkEvent struct {
	Flow  Flow  `json:"flow,omitempty"`
	Event Event `json:"event"`
}

type Event struct {
	EventID     int32   `json:"eventID"`
	Threshold   int32   `json:"threshold"`
	TimeStamp   int64   `json:"timestamp"`
	Value       float32 `json:"value"`
	Name        string  `json:"name"`
	Metric      string  `json:"metric"`
	ThresholdID string  `json:"thresholdID"`
	Agent       string  `json:"agent"`
	DataSource  string  `json:"dataSource"`
}
