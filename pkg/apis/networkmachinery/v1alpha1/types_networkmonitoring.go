/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// EventTypeReconciliation an event reason to describe network monitor reconciliation.
	EventTypeReconciliation string = "NetworkMonitorReconciliation"
	// EventTypeDeletion an event reason to describe network monitor deletion.
	EventTypeDeletion string = "NetworkMonitorDeletion"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkMonitor is the top-level type for flow monitoring
type NetworkMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkMonitorSpec   `json:"spec,omitempty"`
	Status NetworkMonitorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkMonitorList is a list of Network Monitors .
type NetworkMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items is the list of Cluster.
	Items []NetworkMonitor `json:"items"`
}

// NetworkMonitorSpec defines the spec for the network monitor resource
type NetworkMonitorSpec struct {
	MonitoringEndpoint MonitoringEndpoint `json:"monitoringEndpoint"`
	Flows              []Flow             `json:"flows"`
	Thresholds         []Threshold        `json:"thresholds"`
	EventsConfig       EventsConfig       `json:"eventsConfig"`
}

// NetworkMonitorSpec defines the spec for the network monitor resource
type NetworkMonitorStatus struct {
	Status `json:",inline"`
}

type MonitoringEndpoint struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

//TODO: Replace string for native non-string types

// Flow defines the monitoring flow to be installed on the monitoring system
type Flow struct {
	Name          string `json:"name"`
	Keys          string `json:"keys"`
	Value         string `json:"value"`
	Filter        string `json:"filter,omitempty"`
	ActiveTimeout string `json:"activeTimeout,omitempty"`
	Log           string `json:"log,omitempty"`
	FlowStart     string `json:"flowStart,omitempty"`
}

// Threshold is the threshold to define for the flows
type Threshold struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	FlowName string `json:"flowName"`
	Metric   string `json:"metric,omitempty"`
	ByFlow   string `json:"byFlow,omitempty"`
}

// EventsConfig contains configuration parameters for event queries
type EventsConfig struct {
	MaxEvents string `json:"maxEvents"`
	Timeout   string `json:"timeout"`
}
