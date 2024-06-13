package validators

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strconv"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	ktypes "k8s.io/apimachinery/pkg/types"

	"github.com/validator-labs/validator-plugin-network/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-network/internal/constants"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
)

const (
	// See conclusion section of https://www.comparitech.com/net-admin/determine-mtu-size-using-ping/
	defaultPacketHeadersSize int = 28

	nc       string = "nc"
	nslookup string = "nslookup"
	ping     string = "ping"
)

type networkRule interface {
	Name() string
}

type NetworkService struct {
	log logr.Logger
}

func NewNetworkService(log logr.Logger) *NetworkService {
	return &NetworkService{
		log: log,
	}
}

// ReconcileDNSRule reconciles a DNS rule from a NetworkValidator config
func (n *NetworkService) ReconcileDNSRule(nn ktypes.NamespacedName, rule v1alpha1.DNSRule) *types.ValidationRuleResult {

	// Build the default ValidationResult for this DNS rule
	vr := buildValidationResult(rule, constants.ValidationTypeDNS)

	args := []string{rule.Host}
	if rule.Server != "" {
		args = append(args, rule.Server)
	}
	n.handleRuleExec(vr, rule, nslookup, "DNS check failed", args...)

	n.log.V(0).Info("DNS check complete", "message", vr.Condition.Message, "rule", rule.RuleName, "host", rule.Host)
	return vr
}

// ReconcileICMPRule reconciles a ICMP rule from a NetworkValidator config
func (n *NetworkService) ReconcileICMPRule(nn ktypes.NamespacedName, rule v1alpha1.ICMPRule) *types.ValidationRuleResult {

	// Build the default ValidationResult for this ICMP rule
	vr := buildValidationResult(rule, constants.ValidationTypeICMP)

	args := []string{"-c", "3", "-W", "3", rule.Host}
	n.handleRuleExec(vr, rule, ping, "ICMP check failed", args...)

	n.log.V(0).Info("ICMP check complete", "message", vr.Condition.Message, "rule", rule.RuleName, "host", rule.Host)
	return vr
}

// ReconcileIPRangeRule reconciles an IP range rule from a NetworkValidator config
func (n *NetworkService) ReconcileIPRangeRule(nn ktypes.NamespacedName, rule v1alpha1.IPRangeRule) *types.ValidationRuleResult {

	// Build the default ValidationResult for this IP range rule
	vr := buildValidationResult(rule, constants.ValidationTypeIPRange)
	vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf(
		"Ensuring that %s and %d subsequent IPs are all unallocated",
		rule.StartIP, rule.Length,
	))

	errMsg := "IP range check failed"

	defer func() {
		n.log.V(0).Info("IP range check complete", "message", vr.Condition.Message, "rule", rule.RuleName)
	}()

	// parse the starting IP
	ip, err := toIpV4(rule.StartIP)
	if err != nil {
		vr.Condition.Failures = append(vr.Condition.Failures, fmt.Sprintf("invalid start IP: %s. error: %v", rule.StartIP, err))
		vr.Condition.Message = errMsg
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
		incIp(ip)
	}
	if anySucceeded {
		vr.Condition.Message = "IP range check failed: one or more IPs in the provided range was allocated"
		vr.Condition.Status = corev1.ConditionFalse
		vr.State = util.Ptr(vapi.ValidationFailed)
		return vr
	}

	return vr
}

// ReconcileMTURule reconciles an MTU rule from a NetworkValidator config
func (n *NetworkService) ReconcileMTURule(nn ktypes.NamespacedName, rule v1alpha1.MTURule) *types.ValidationRuleResult {

	// Build the default ValidationResult for this MTU rule
	vr := buildValidationResult(rule, constants.ValidationTypeMTU)

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

// ReconcileTCPConnRule reconciles a TCP connection rule from a NetworkValidator config
func (n *NetworkService) ReconcileTCPConnRule(nn ktypes.NamespacedName, rule v1alpha1.TCPConnRule) *types.ValidationRuleResult {

	// Build the default ValidationResult for this TCP connection rule
	vr := buildValidationResult(rule, constants.ValidationTypeTCPConn)
	vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf(
		"Ensuring that TCP connection(s) can be established to %s on port(s) %v",
		rule.Host, rule.Ports,
	))

	errMsg := "TCP connection check failed"

	// construct args
	args := []string{"-w", "3"}
	if rule.ProxyProtocol != "" && rule.ProxyAddress == "" || rule.ProxyProtocol == "" && rule.ProxyAddress != "" {
		vr.State = util.Ptr(vapi.ValidationFailed)
		vr.Condition.Failures = append(vr.Condition.Failures, "invalid rule: proxyProtocol & proxyAddress must both be defined, or both undefined")
		vr.Condition.Message = errMsg
	}
	if rule.ProxyProtocol != "" && rule.ProxyAddress != "" {
		args = append(args, "-X", rule.ProxyProtocol, "-x", rule.ProxyAddress)
	}
	args = append(args, rule.Host)
	for _, p := range rule.Ports {
		args = append(args, strconv.FormatInt(int64(p), 10))
	}
	n.handleRuleExec(vr, rule, nc, errMsg, args...)

	n.log.V(0).Info("TCP connection check complete", "message", vr.Condition.Message, "rule", rule.RuleName, "host", rule.Host)
	return vr
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule networkRule, validationType string) *types.ValidationRuleResult {
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", validationType)
	latestCondition.ValidationRule = rule.Name()
	latestCondition.ValidationType = validationType
	return &types.ValidationRuleResult{Condition: &latestCondition, State: &state}
}

// handleRuleExec executes a rule's command and updates the rule's validation result accordingly
func (n *NetworkService) handleRuleExec(vr *types.ValidationRuleResult, r networkRule, binary, errMsg string, args ...string) {
	n.log.V(0).Info("Executing command: %s %s", binary, args)
	stdout, stderr, exitCode, err := execCmd(binary, args...)
	if err != nil || stderr != "" || exitCode != 0 {
		n.failResult(vr, err, binary, errMsg, stdout, stderr, r.Name(), args...)
		return
	}
	vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("%s %s succeeded", binary, args))
	vr.Condition.Message = stdout
}

// failResult updates a validation result with failure details
func (n *NetworkService) failResult(vr *types.ValidationRuleResult, err error, binary, errMsg, stdout, stderr, ruleName string, args ...string) {
	n.log.V(0).Info(errMsg, "stdout", stdout, "stderr", stderr, "error", err.Error(), "rule", ruleName)
	failure := fmt.Sprintf("stdout: %s, stderr: %s, error: %v", stdout, stderr, err)
	vr.State = util.Ptr(vapi.ValidationFailed)
	vr.Condition.Status = corev1.ConditionFalse
	vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("%s %s failed", binary, args))
	vr.Condition.Failures = append(vr.Condition.Failures, failure)
	vr.Condition.Message = errMsg
}

// execCmd executes an arbitrary binary with arbitrary arguments
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

// toIpV4 converts a string to a net.IP
func toIpV4(s string) (ip net.IP, err error) {
	ip = net.ParseIP(s).To4()
	if ip == nil {
		err = fmt.Errorf("failed to parse %s as an IPv4 address", s)
	}
	return
}

// incIp increments a net.IP
func incIp(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
