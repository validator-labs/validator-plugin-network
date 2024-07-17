/*
Copyright 2024.

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
)

// NetworkValidatorSpec defines the desired state of NetworkValidator
type NetworkValidatorSpec struct {
	// DNSRules validate DNS name resolution of network hosts
	// +kubebuilder:validation:MaxItems=5
	// +kubebuilder:validation:XValidation:message="DNSRules must have unique names",rule="self.all(e, size(self.filter(x, x.name == e.name)) == 1)"
	DNSRules []DNSRule `json:"dnsRules,omitempty" yaml:"dnsRules,omitempty"`
	// ICMPRules validate ICMP pings to network hosts
	// +kubebuilder:validation:MaxItems=5
	// +kubebuilder:validation:XValidation:message="ICMPRules must have unique names",rule="self.all(e, size(self.filter(x, x.name == e.name)) == 1)"
	ICMPRules []ICMPRule `json:"icmpRules,omitempty" yaml:"icmpRules,omitempty"`
	// IPRangeRules validate that all IPs in a given CIDR range are free (unallocated)
	// +kubebuilder:validation:MaxItems=5
	// +kubebuilder:validation:XValidation:message="IPRangeRules must have unique names",rule="self.all(e, size(self.filter(x, x.name == e.name)) == 1)"
	IPRangeRules []IPRangeRule `json:"ipRangeRules,omitempty" yaml:"ipRangeRules,omitempty"`
	// MTURules validate that the default NIC has an MTU of at least X, where X is the provided MTU
	// +kubebuilder:validation:MaxItems=5
	// +kubebuilder:validation:XValidation:message="MTURules must have unique names",rule="self.all(e, size(self.filter(x, x.name == e.name)) == 1)"
	MTURules []MTURule `json:"mtuRules,omitempty" yaml:"mtuRules,omitempty"`
	// TCPConnRules validate arbitrary TCP connections, including proxied connections
	// +kubebuilder:validation:MaxItems=5
	// +kubebuilder:validation:XValidation:message="TCPConnRules must have unique names",rule="self.all(e, size(self.filter(x, x.name == e.name)) == 1)"
	TCPConnRules []TCPConnRule `json:"tcpConnRules,omitempty" yaml:"tcpConnRules,omitempty"`
	// HTTPFileRules validate that files are publicly available via HTTP
	// +kubebuilder:validation:MaxItems=5
	// +kubebuilder:validation:XValidation:message="HTTPFileRules must have unique names",rule="self.all(e, size(self.filter(x, x.name == e.name)) == 1)"
	HTTPFileRules []HTTPFileRule `json:"httpFileRules,omitempty" yaml:"httpFileRules,omitempty"`
	// CACerts allow additional CA certificates to be used for TLS. Applies to TCPConnRules and HTTPFileRules.
	CACerts CACertificates `json:"caCerts,omitempty" yaml:"caCerts,omitempty"`
}

// CACertificates contains configuration for additional CA certificates to use for TLS. Can be certs
// provided inline or secret references. Secrets are assumed to be in the same namespace as the
// NetworkValidator.
type CACertificates struct {
	// Certs is a list of certificates to use.
	// +kubebuilder:validation:MaxItems=500
	Certs []Certificate `json:"certs,omitempty" yaml:"certs,omitempty"`
	// SecretRefs is a list of secret references to use.
	// +kubebuilder:validation:MaxItems=500
	SecretRefs []SecretReference `json:"secretRefs,omitempty" yaml:"secretRefs,omitempty"`
}

// Certificate is a certificate specified inline.
// +kubebuilder:validation:MinLength=1
type Certificate string

// SecretReference is a secret's name and the key to use to get the data.
type SecretReference struct {
	// Name is the name of the secret.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name" yaml:"name"`
	// Key is the key in the secret data.
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key" yaml:"key"`
}

// ResultCount returns the number of validation results expected for a NetworkValidatorSpec.
func (s NetworkValidatorSpec) ResultCount() int {
	return len(s.DNSRules) + len(s.ICMPRules) + len(s.IPRangeRules) + len(s.MTURules) + len(s.TCPConnRules) + len(s.HTTPFileRules)
}

// DNSRule defines a DNS validation rule.
type DNSRule struct {
	// RuleName is a unique identifier for the rule in the validator. Used to ensure conditions do not overwrite each other.
	// +kubebuilder:validation:MaxLength=500
	RuleName string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Server   string `json:"server,omitempty" yaml:"server,omitempty"`
}

// Name returns the name of a DNSRule.
func (r DNSRule) Name() string {
	return r.RuleName
}

// ICMPRule defines an ICMP validation rule.
type ICMPRule struct {
	// RuleName is a unique identifier for the rule in the validator. Used to ensure conditions do not overwrite each other.
	// +kubebuilder:validation:MaxLength=500
	RuleName string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
}

// Name returns the name of an ICMPRule.
func (r ICMPRule) Name() string {
	return r.RuleName
}

// IPRangeRule defines an IP range validation rule.
type IPRangeRule struct {
	// RuleName is a unique identifier for the rule in the validator. Used to ensure conditions do not overwrite each other.
	// +kubebuilder:validation:MaxLength=500
	RuleName string `json:"name" yaml:"name"`
	StartIP  string `json:"startIp" yaml:"startIp"`
	Length   int    `json:"length" yaml:"length"`
}

// Name returns the name of an IPRangeRule.
func (r IPRangeRule) Name() string {
	return r.RuleName
}

// MTURule defines an MTU validation rule.
type MTURule struct {
	// RuleName is a unique identifier for the rule in the validator. Used to ensure conditions do not overwrite each other.
	// +kubebuilder:validation:MaxLength=500
	RuleName string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	MTU      int    `json:"mtu" yaml:"mtu"`
	// Optionally specify the size in bytes of the packet headers for the MTU ping packet.
	// This varies by medium, e.g. Ethernet, WiFi, etc.) and defaults to 28 bytes
	// (20 bytes IP header + 8 bytes ICMP header)
	PacketHeadersSize int `json:"packetHeadersSize,omitempty" yaml:"packetHeadersSize,omitempty"`
}

// Name returns the name of an MTURule.
func (r MTURule) Name() string {
	return r.RuleName
}

// TCPConnRule defines a TCP connection validation rule.
type TCPConnRule struct {
	// RuleName is a unique identifier for the rule in the validator. Used to ensure conditions do not overwrite each other.
	// +kubebuilder:validation:MaxLength=500
	RuleName string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Ports    []int  `json:"ports" yaml:"ports"`
	// +kubebuilder:validation:Pattern=`^(4|5|connect)?$`
	ProxyProtocol string `json:"proxyProtocol,omitempty" yaml:"proxyProtocol,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty" yaml:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty" yaml:"proxyPort,omitempty"`
	// TODO: use socat for proxy validation using TLS CAFile & basic auth?
}

// Name returns the name of a TCPConnRule.
func (r TCPConnRule) Name() string {
	return r.RuleName
}

// HTTPFileRule defines an HTTP file rule.
type HTTPFileRule struct {
	// RuleName is a unique identifier for the rule in the validator. Used to ensure conditions do not overwrite each other.
	// +kubebuilder:validation:MaxLength=500
	RuleName string `json:"name" yaml:"name"`
	// Paths is a list of file paths to check. When performing HTTP requests, if any of the paths result in a non-200 OK response code, the rule fails validation.
	// +kubebuilder:validation:MaxItems=1000
	Paths []string `json:"paths" yaml:"paths"`
	// InsecureSkipVerify controls whether the HTTP client used validate the rule skips TLS
	// certificate verification. Defaults to false.
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty" yaml:"insecureSkipVerify,omitempty"`
}

// Name returns the name of a HTTPFileRule.
func (r HTTPFileRule) Name() string {
	return r.RuleName
}

// NetworkValidatorStatus defines the observed state of NetworkValidator
type NetworkValidatorStatus struct{}

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
