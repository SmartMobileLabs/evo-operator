/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"smartmobilelabs.com/evo/evo-operator/controllers/k8sdynamic"
)

type AppStatus string

type PrivateNetworkAccess struct {
	// the ID of this application
	ApnUUID string `json:"apnUUID"`
	// the network on which this application will run
	ApplicationNetwork string `json:"applicationNetwork"`
	// additional needed IP routes
	AdditionalRoutes []string `json:"additionalRoutes,omitempty"`
	// the IP address under which the EVO will be reachable. If set, it will ONLY be reachable under this address
	AppPodFixIp string `json:"appPodFixIp,omitempty"`
}

const (
	AppStatusNotSet     = "UNSET"
	AppStatusNotRunning = "NOT_RUNNING"
	AppStatusRunning    = "RUNNING"
	AppStatusFrozen     = "FROZEN"
)

type AppReporteData struct {
	// metrics information
	MetricsClusterIp string `json:"metricsClusterIp,omitempty"`
	//Ip addresses of the services that received IP address from the private network
	PrivateNetworkIpAddress map[string]string `json:"privateNetworkIpAddresses,omitempty"`
}

// SmlEvoStatus defines the observed state of SmlEvo
type SmlEvoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	PrevSpec         *SmlEvoSpec                     `json:"prevSpec,omitempty"`
	AppStatus        AppStatus                       `json:"appStatus,omitempty"`
	AppReportedData  AppReporteData                  `json:"appReportedData,omitempty"`
	AppliedResources []k8sdynamic.ResourceDescriptor `json:"appliedResources,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type EvoPorts struct {
	// the tcp/ip port for the application to listen on for https REST and GUI
	UiPort int `json:"uiPort,omitempty"`
	// the lowest UDP port of the range of UDP ports to listen to for ports that are to be seen in the outside of the application
	UdpPortRangeLow int `json:"udpPortLow,omitempty"`
	// the highest UDP port of the range of UDP ports to listen to for ports that are to be seen in the outside of the application
	UdpPortRangeHigh int `json:"udpPortHigh,omitempty"`
}

type SmlEvoSpec struct {
	// ports to be used for the application
	Ports EvoPorts `json:"ports"`
	// the domain name for the metrics to report to
	MetricsDomainName string `json:"metricsDomainName,omitempty"`
	// information about into which network the application is to be placed
	PrivateNetworkAccess *PrivateNetworkAccess `json:"privateNetworkAccess,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SmlEvo is the Schema for the smlevoes API
type SmlEvo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SmlEvoSpec   `json:"spec,omitempty"`
	Status SmlEvoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SmlEvoList contains a list of SmlEvo
type SmlEvoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SmlEvo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SmlEvo{}, &SmlEvoList{})
}
