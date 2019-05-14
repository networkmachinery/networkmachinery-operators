package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PingResultState string

const (
	FailedPing  PingResultState = "Failed"
	SuccessPing PingResultState = "Success"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PingStatus contains information related ping command results.
type PingStatus struct {
	metav1.TypeMeta `json:",inline"`

	PingIPEndpoints     []PingIPEndpoint      `json:"ipEndpoints,omitempty"`
	PingPodEndpoints    []PingPodEndpoint     `json:"podEndpoints,omitempty"`
	PingServiceEndpoint []PingServiceEndpoint `json:"serviceEndpoints,omitempty"`
}

type PingIPEndpoint struct {
	IP         string     `json:"ip"`
	PingResult PingResult `json:"pingResult"`
}

type PingPodEndpoint struct {
	PodParams  Params     `json:"podParams"`
	PingResult PingResult `json:"pingResult"`
}

type PingServiceEndpoint struct {
	ServiceParams  Params           `json:"serviceParams"`
	ServiceResults []PingIPEndpoint `json:"serviceResults"`
}

type PingResult struct {
	State   PingResultState `json:"state"`
	Min     string          `json:"min,omitempty"`
	Max     string          `json:"max,omitempty"`
	Average string          `json:"average,omitempty"`
}
