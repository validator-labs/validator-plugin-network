// Package network contains network validation methods.
package network

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"

	"github.com/validator-labs/validator-plugin-network/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-network/pkg/constants"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
)

const (
	// See conclusion section of https://www.comparitech.com/net-admin/determine-mtu-size-using-ping/
	defaultPacketHeadersSize int = 28

	nslookup string = "nslookup"
	ping     string = "ping"
)

type networkRule interface {
	Name() string
}

// RuleService is a service for network validation.
type RuleService struct {
	log        logr.Logger
	httpClient *http.Client
	tlsConfig  *tls.Config
}

// Option allows customizing the RuleService.
type Option func(*RuleService)

// WithTLSConfig allows callers to provide a custom TLS config for the RuleService.
func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(n *RuleService) {
		n.tlsConfig = tlsConfig
	}
}

// WithTransport allows callers to provide a custom transport for the RuleService's HTTP client.
func WithTransport(transport http.RoundTripper) Option {
	return func(n *RuleService) {
		n.httpClient.Transport = transport
	}
}

// NewRuleService creates a new RuleService.
func NewRuleService(log logr.Logger, opts ...Option) *RuleService {
	n := &RuleService{
		httpClient: http.DefaultClient,
		log:        log,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

// ReconcileDNSRule reconciles a DNS rule from a NetworkValidator config.
func (n *RuleService) ReconcileDNSRule(rule v1alpha1.DNSRule) *types.ValidationRuleResult {

	// Build the default ValidationResult for this DNS rule
	vr := BuildValidationResult(rule, constants.ValidationTypeDNS)

	args := []string{rule.Host}
	if rule.Server != "" {
		args = append(args, rule.Server)
	}
	n.handleRuleExec(vr, rule, nslookup, "DNS check failed", args...)

	n.log.V(0).Info("DNS check complete", "message", vr.Condition.Message, "rule", rule.RuleName, "host", rule.Host)
	return vr
}

// ReconcileICMPRule reconciles a ICMP rule from a NetworkValidator config.
func (n *RuleService) ReconcileICMPRule(rule v1alpha1.ICMPRule) *types.ValidationRuleResult {

	// Build the default ValidationResult for this ICMP rule
	vr := BuildValidationResult(rule, constants.ValidationTypeICMP)

	args := []string{"-c", "3", "-W", "3", rule.Host}
	n.handleRuleExec(vr, rule, ping, "ICMP check failed", args...)

	n.log.V(0).Info("ICMP check complete", "message", vr.Condition.Message, "rule", rule.RuleName, "host", rule.Host)
	return vr
}

// ReconcileIPRangeRule reconciles an IP range rule from a NetworkValidator config.
func (n *RuleService) ReconcileIPRangeRule(rule v1alpha1.IPRangeRule) *types.ValidationRuleResult {

	// Build the default ValidationResult for this IP range rule
	vr := BuildValidationResult(rule, constants.ValidationTypeIPRange)
	vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf(
		"Ensuring that %s and %d subsequent IPs are all unallocated",
		rule.StartIP, rule.Length,
	))

	defer func() {
		n.log.V(0).Info("IP range check complete", "message", vr.Condition.Message, "rule", rule.RuleName)
	}()

	// parse the starting IP
	ip, err := toIPV4(rule.StartIP)
	if err != nil {
		vr.Condition.Failures = append(vr.Condition.Failures, fmt.Sprintf("invalid start IP: %s. error: %v", rule.StartIP, err))
		vr.Condition.Message = "IP range check failed. See failures for details."
		vr.Condition.Status = corev1.ConditionFalse
		vr.State = util.Ptr(vapi.ValidationFailed)
		return vr
	}

	// ping each IP in the range
	anySucceeded := false
	for i := 0; i < rule.Length; i++ {
		args := []string{"-c", "3", "-W", "3", ip.String()}
		stdout, stderr, _, err := execCmd(ping, args...)
		if err != nil || stderr != "" {
			vr.Condition.Details = append(
				vr.Condition.Details,
				fmt.Sprintf("%s %s failed, IP is available; err: %v, stderr: %s", ping, args, err, stderr),
			)
		} else {
			n.log.V(0).Info("IP allocated", "stdout", stdout, "rule", rule.RuleName)
			anySucceeded = true
			vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("%s %s succeeded", ping, args))
			vr.Condition.Failures = append(vr.Condition.Failures, fmt.Sprintf("stdout: %s", stdout))
		}
		incIP(ip)
	}
	if anySucceeded {
		vr.Condition.Message = "IP range check failed: one or more IPs in the provided range was allocated"
		vr.Condition.Status = corev1.ConditionFalse
		vr.State = util.Ptr(vapi.ValidationFailed)
		return vr
	}

	return vr
}

// ReconcileMTURule reconciles an MTU rule from a NetworkValidator config.
func (n *RuleService) ReconcileMTURule(rule v1alpha1.MTURule) *types.ValidationRuleResult {

	// Build the default ValidationResult for this MTU rule
	vr := BuildValidationResult(rule, constants.ValidationTypeMTU)

	packetHeadersSize := defaultPacketHeadersSize
	if rule.PacketHeadersSize != 0 {
		packetHeadersSize = rule.PacketHeadersSize
	}
	size := strconv.FormatInt(int64(rule.MTU-packetHeadersSize), 10)
	vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf(
		"MTU check ping packet size: %s = %d (MTU) - %d (packet headers bytes)",
		size, rule.MTU, packetHeadersSize,
	))

	args := []string{"-c", "3", "-W", "3", "-M", "do", "-s", size, rule.Host}
	n.handleRuleExec(vr, rule, ping, "MTU check failed", args...)

	n.log.V(0).Info("MTU check complete", "message", vr.Condition.Message, "rule", rule.RuleName, "host", rule.Host)
	return vr
}

// ReconcileTCPConnRule reconciles a TCP connection rule from a NetworkValidator config.
func (n *RuleService) ReconcileTCPConnRule(rule v1alpha1.TCPConnRule) *types.ValidationRuleResult {

	// Build the default ValidationResult for this TCP connection rule
	vr := BuildValidationResult(rule, constants.ValidationTypeTCPConn)
	vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf(
		"Ensuring that TCP connection(s) can be established to %s on port(s) %v",
		rule.Host, rule.Ports,
	))

	tlsDialer := &tls.Dialer{
		NetDialer: &net.Dialer{
			Timeout: time.Duration(rule.Timeout) * time.Second,
		},
		Config: n.tlsConfig,
	}

	for _, port := range rule.Ports {
		conn, err := tlsDialer.Dial("tcp", fmt.Sprintf("%s:%d", rule.Host, port))
		if err != nil {
			vr.Condition.Failures = append(vr.Condition.Failures, fmt.Sprintf("failed to connect to %s on port %d: %v", rule.Host, port, err))
			continue
		}
		if err := conn.Close(); err != nil {
			vr.Condition.Failures = append(vr.Condition.Failures, fmt.Sprintf("failed to close connection to %s on port %d: %v", rule.Host, port, err))
			continue
		}
		vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("Successfully connected to %s on port %d", rule.Host, port))
	}

	if len(vr.Condition.Failures) > 0 {
		vr.Condition.Message = "TCP connection check(s) failed. See failures for details."
		vr.Condition.Status = corev1.ConditionFalse
		vr.State = util.Ptr(vapi.ValidationFailed)
	}

	n.log.V(0).Info("TCP connection check complete", "message", vr.Condition.Message, "rule", rule.RuleName, "host", rule.Host)
	return vr
}

// ReconcileHTTPFileRule reconciles an HTTP file rule from a NetworkValidator config.
func (n *RuleService) ReconcileHTTPFileRule(rule v1alpha1.HTTPFileRule) *types.ValidationRuleResult {

	// Build the default ValidationResult for this HTTP file rule
	vr := BuildValidationResult(rule, constants.ValidationTypeHTTPFile)
	vr.Condition.Message = "All files are accessible."
	vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf(
		"Ensuring that files %v are accessible.",
		rule.Paths,
	))

	for _, path := range rule.Paths {
		checkFileErrMsg, err := n.checkFile(path)
		if err != nil {
			vr.Condition.Failures = append(vr.Condition.Failures, fmt.Sprintf("failed to check file %s: %s", path, err))
			continue
		}
		if checkFileErrMsg != "" {
			vr.Condition.Failures = append(vr.Condition.Failures, fmt.Sprintf("file %s is not accessible; %s", path, checkFileErrMsg))
			continue
		}
		vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("File %s is accessible.", path))
	}

	if len(vr.Condition.Failures) > 0 {
		vr.Condition.Message = "HTTPFile check(s) failed. See failures for details."
		vr.Condition.Status = corev1.ConditionFalse
		vr.State = util.Ptr(vapi.ValidationFailed)
	}

	return vr
}

func (n *RuleService) checkFile(path string) (string, error) {
	// Create a new HTTP HEAD request for the file.
	req, err := http.NewRequest("HEAD", path, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Send the request. If a 200 response comes back, we consider the file accessible.
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				n.log.Error(closeErr, "failed to close HTTP response body")
			}
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("'%d' status code in response to HEAD request", resp.StatusCode), nil
	}

	// File accessible.
	return "", err
}

// BuildValidationResult builds a default ValidationResult for a given validation type.
func BuildValidationResult(rule networkRule, validationType string) *types.ValidationRuleResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", validationType)
	latestCondition.ValidationRule = rule.Name()
	latestCondition.ValidationType = validationType
	return &types.ValidationRuleResult{Condition: &latestCondition, State: &state}
}

// handleRuleExec executes a rule's command and updates the rule's validation result accordingly.
func (n *RuleService) handleRuleExec(vr *types.ValidationRuleResult, r networkRule, binary, errMsg string, args ...string) {
	n.log.V(0).Info("Executing command: %s %s", binary, args)
	stdout, stderr, exitCode, err := execCmd(binary, args...)
	if err != nil || stderr != "" || exitCode != 0 {
		n.failResult(vr, err, binary, errMsg, stdout, stderr, r.Name(), args...)
		return
	}
	vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("%s %s succeeded", binary, args))
	vr.Condition.Message = stdout
}

// failResult updates a validation result with failure details.
func (n *RuleService) failResult(vr *types.ValidationRuleResult, err error, binary, errMsg, stdout, stderr, ruleName string, args ...string) {
	n.log.V(0).Info(errMsg, "stdout", stdout, "stderr", stderr, "error", err.Error(), "rule", ruleName)
	failure := fmt.Sprintf("stdout: %s, stderr: %s, error: %v", stdout, stderr, err)
	vr.State = util.Ptr(vapi.ValidationFailed)
	vr.Condition.Status = corev1.ConditionFalse
	vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("%s %s failed", binary, args))
	vr.Condition.Failures = append(vr.Condition.Failures, failure)
	vr.Condition.Message = errMsg
}

// execCmd executes an arbitrary binary with arbitrary arguments.
func execCmd(binary string, args ...string) (string, string, int, error) {
	cmd := exec.Command(binary, args...)
	stderr, stdout := bytes.Buffer{}, bytes.Buffer{}
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	exitCode := 0
	err := cmd.Run()
	if err != nil {
		exitCode = 1
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	return stdout.String(), stderr.String(), exitCode, err
}

// toIPV4 converts a string to a net.IP.
func toIPV4(s string) (ip net.IP, err error) {
	ip = net.ParseIP(s).To4()
	if ip == nil {
		err = fmt.Errorf("failed to parse %s as an IPv4 address", s)
	}
	return
}

// incIP increments a net.IP.
func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
