// Package validate defines a Validate function that evaluates a NetworkValidatorSpec and returns a ValidationResponse.
package validate

import (
	"github.com/go-logr/logr"

	"github.com/validator-labs/validator-plugin-network/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-network/pkg/http"
	"github.com/validator-labs/validator-plugin-network/pkg/network"
	"github.com/validator-labs/validator/pkg/types"
)

// Validate validates the NetworkValidatorSpec and returns a ValidationResponse.
func Validate(spec v1alpha1.NetworkValidatorSpec, caPems [][]byte, auths map[string][]string, log logr.Logger) types.ValidationResponse {
	resp := types.ValidationResponse{
		ValidationRuleResults: make([]*types.ValidationRuleResult, 0, spec.ResultCount()),
		ValidationRuleErrors:  make([]error, 0, spec.ResultCount()),
	}

	svc := network.NewRuleService(log)

	// DNS rules
	for _, rule := range spec.DNSRules {
		vrr := svc.ReconcileDNSRule(rule)
		resp.AddResult(vrr, nil)
	}

	// ICMP rules
	for _, rule := range spec.ICMPRules {
		vrr := svc.ReconcileICMPRule(rule)
		resp.AddResult(vrr, nil)
	}

	// IP range rules
	for _, rule := range spec.IPRangeRules {
		vrr := svc.ReconcileIPRangeRule(rule)
		resp.AddResult(vrr, nil)
	}

	// MTU rules
	for _, rule := range spec.MTURules {
		vrr := svc.ReconcileMTURule(rule)
		resp.AddResult(vrr, nil)
	}

	// TCP connection rules
	for _, rule := range spec.TCPConnRules {
		tlsConfig := http.TLSConfig(caPems, rule.InsecureSkipTLSVerify, log)
		ruleSvc := network.NewRuleService(log, network.WithTLSConfig(tlsConfig))
		vrr := ruleSvc.ReconcileTCPConnRule(rule)
		resp.AddResult(vrr, nil)
	}

	// HTTP file rules
	for _, rule := range spec.HTTPFileRules {
		auth := auths[rule.Name()]
		transport := http.Transport(caPems, auth, rule.InsecureSkipTLSVerify, log)
		ruleSvc := network.NewRuleService(log, network.WithTransport(transport))
		vrr := ruleSvc.ReconcileHTTPFileRule(rule)
		resp.AddResult(vrr, nil)
	}

	return resp
}
