/*
Copyright 2023.

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
	valid8orv1alpha1 "github.com/spectrocloud-labs/valid8or/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkValidatorSpec defines the desired state of NetworkValidator
type NetworkValidatorSpec struct {
	// DNSRules validate DNS name resolution of network hosts
	DNSRules []DNSRule `json:"dnsRules,omitempty"`
	// ICMPRules validate ICMP pings to network hosts
	ICMPRules []ICMPRule `json:"icmpRules,omitempty"`
	// IPRangeRules validate that all IPs in a given CIDR range are free (unallocated)
	IPRangeRules []IPRangeRule `json:"ipRangeRules,omitempty"`
	// MTURules validate that the default NIC has an MTU of at least X, where X is the provided MTU
	MTURules []MTURule `json:"mtuRules,omitempty"`
	// NetcatRules validate arbitrary TCP & UDP connections
	NetcatRules []NetcatRule `json:"netcatRules,omitempty"`
}

type DNSRule struct {
	Host   string `json:"host"`
	Server string `json:"server,omitempty"`
}

type ICMPRule struct {
	Host string `json:"host"`
}

type IPRangeRule struct {
	CIDR string `json:"cidr"`
}

type MTURule struct {
	Host string `json:"host"`
	MTU  int    `json:"mtu"`
}

type NetcatRule struct {
	Host      string `json:"host"`
	Ports     []int  `json:"ports"`
	EnableTLS bool   `json:"enableTls,omitempty"`
	CABase64  string `json:"caBase64,omitempty"`
	// +kubebuilder:validation:Pattern=`^(4|5|connect)?$`
	ProxyProtocol string `json:"proxyProtocol,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty"`
}

// NetworkValidatorStatus defines the observed state of NetworkValidator
type NetworkValidatorStatus struct {
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []valid8orv1alpha1.ValidationCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NetworkValidator is the Schema for the networkvalidators API
type NetworkValidator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkValidatorSpec   `json:"spec,omitempty"`
	Status NetworkValidatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NetworkValidatorList contains a list of NetworkValidator
type NetworkValidatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkValidator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NetworkValidator{}, &NetworkValidatorList{})
}
