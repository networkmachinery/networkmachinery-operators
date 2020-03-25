package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type QperfResultState string

const (
	FailedQperfCmd  QperfResultState = "Failed"
	SuccessQperfCmd QperfResultState = "Success"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QperfStatus contains information related qperf command results.
type QperfStatus struct {
	metav1.TypeMeta    `json:",inline"`
	QperfPodEndpoints  []QperfPodEndpoint  `json:"podEndpoints,omitempty"`
	QperfNodeEndpoints []QperfNodeEndpoint `json:"nodeEndpoints,omitempty"`
}

type QperfPodEndpoint struct {
	PodParams   Params      `json:"podParams"`
	QperfResult QperfResult `json:"qperfResult"`
}
type QperfNodeEndpoint struct {
	NodeParams  NodeParams  `json:"nodeParams"`
	NodeResults QperfResult `json:"nodeResult"`
}

type QperfResult struct {
	State    QperfResultState `json:"state"`
	UDPDelay string           `json:"udpDelay,omitempty"`
	TCPDelay string           `json:"tcpDelay,omitempty"`
}
