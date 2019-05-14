package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NetcatResultState string

const (
	Refused   NetcatResultState = "Refused"
	Succeeded NetcatResultState = "Succeeded"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetcatStatus contains information related netcat command results.
type NetcatStatus struct {
	metav1.TypeMeta `json:",inline"`

	NetcatIPEndpoints      []NetcatIPEndpoint      `json:"ipEndpoints,omitempty"`
	NetcatPodEndpoints     []NetcatPodEndpoint     `json:"podEndpoints,omitempty"`
	NetcatServiceEndpoints []NetcatServiceEndpoint `json:"serviceEndpoints,omitempty"`
}

type NetcatIPEndpoint struct {
	IP           string       `json:"ip"`
	Port         string       `json:"port"`
	NetcatResult NetcatResult `json:"netcatResult"`
}

type NetcatPodEndpoint struct {
	PodParams    Params       `json:"podParams"`
	NetcatResult NetcatResult `json:"netcatResult"`
}

type NetcatServiceEndpoint struct {
	ServiceParams  Params             `json:"serviceParams"`
	ServiceResults []NetcatIPEndpoint `json:"serviceResults"`
}

type NetcatResult struct {
	State NetcatResultState `json:"state"`
}
