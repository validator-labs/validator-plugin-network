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
	validatorv1alpha1 "github.com/spectrocloud-labs/validator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkValidatorSpec defines the desired state of NetworkValidator
type NetworkValidatorSpec struct {
	// DNSRules validate DNS name resolution of network hosts
	DNSRules []DNSRule `json:"dnsRules,omitempty" yaml:"dnsRules,omitempty"`
	// ICMPRules validate ICMP pings to network hosts
	ICMPRules []ICMPRule `json:"icmpRules,omitempty" yaml:"icmpRules,omitempty"`
	// IPRangeRules validate that all IPs in a given CIDR range are free (unallocated)
	IPRangeRules []IPRangeRule `json:"ipRangeRules,omitempty" yaml:"ipRangeRules,omitempty"`
	// MTURules validate that the default NIC has an MTU of at least X, where X is the provided MTU
	MTURules []MTURule `json:"mtuRules,omitempty" yaml:"mtuRules,omitempty"`
	// TCPConnRules validate arbitrary TCP connections, including proxied connections
	TCPConnRules []TCPConnRule `json:"tcpConnRules,omitempty" yaml:"tcpConnRules,omitempty"`
}

type DNSRule struct {
	RuleName string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Server   string `json:"server,omitempty" yaml:"server,omitempty"`
}

func (r DNSRule) Name() string {
	return r.RuleName
}

type ICMPRule struct {
	RuleName string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
}

func (r ICMPRule) Name() string {
	return r.RuleName
}

type IPRangeRule struct {
	RuleName string `json:"name" yaml:"name"`
	StartIP  string `json:"startIp" yaml:"startIp"`
	Length   int    `json:"length" yaml:"length"`
}

func (r IPRangeRule) Name() string {
	return r.RuleName
}

type MTURule struct {
	RuleName string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	MTU      int    `json:"mtu" yaml:"mtu"`
}

func (r MTURule) Name() string {
	return r.RuleName
}

type TCPConnRule struct {
	RuleName string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Ports    []int  `json:"ports" yaml:"ports"`
	// +kubebuilder:validation:Pattern=`^(4|5|connect)?$`
	ProxyProtocol string `json:"proxyProtocol,omitempty" yaml:"proxyProtocol,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty" yaml:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty" yaml:"proxyPort,omitempty"`
	// TODO: use socat for proxy validation using TLS CAFile & basic auth?
}

func (r TCPConnRule) Name() string {
	return r.RuleName
}

// NetworkValidatorStatus defines the observed state of NetworkValidator
type NetworkValidatorStatus struct {
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []validatorv1alpha1.ValidationCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
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
