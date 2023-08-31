package validators

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	ktypes "k8s.io/apimachinery/pkg/types"

	"github.com/spectrocloud-labs/valid8or-plugin-network/api/v1alpha1"
	"github.com/spectrocloud-labs/valid8or-plugin-network/internal/constants"
	v8or "github.com/spectrocloud-labs/valid8or/api/v1alpha1"
	"github.com/spectrocloud-labs/valid8or/pkg/types"
	"github.com/spectrocloud-labs/valid8or/pkg/util/ptr"
)

const (
	ping     string = "ping"
	nc       string = "nc"
	nslookup string = "nslookup"

	icmpPacketSize int = 28
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
func (n *NetworkService) ReconcileDNSRule(nn ktypes.NamespacedName, rule v1alpha1.DNSRule) (*types.ValidationResult, error) {

	// Build the default ValidationResult for this DNS rule
	vr := buildValidationResult(rule, constants.ValidationTypeDNS)

	errMsg := "Failed to execute DNS check"
	args := []string{rule.Host}
	if rule.Server != "" {
		args = append(args, rule.Server)
	}
	if err := n.handleRuleExec(vr, rule, nslookup, errMsg, args...); err != nil {
		return vr, err
	}
	n.log.V(0).Info("DNS check passed", "stdout", vr.Condition.Message, "rule", rule.RuleName, "host", rule.Host)

	return vr, nil
}

// ReconcileICMPRule reconciles a ICMP rule from a NetworkValidator config
func (n *NetworkService) ReconcileICMPRule(nn ktypes.NamespacedName, rule v1alpha1.ICMPRule) (*types.ValidationResult, error) {

	// Build the default ValidationResult for this ICMP rule
	vr := buildValidationResult(rule, constants.ValidationTypeICMP)

	errMsg := "Failed to execute ICMP check"
	args := []string{"-c", "3", "-W", "3", rule.Host}
	if err := n.handleRuleExec(vr, rule, ping, errMsg, args...); err != nil {
		return vr, err
	}
	n.log.V(0).Info("ICMP check passed", "stdout", vr.Condition.Message, "rule", rule.RuleName, "host", rule.Host)

	return vr, nil
}

// ReconcileIPRangeRule reconciles an IP range rule from a NetworkValidator config
func (n *NetworkService) ReconcileIPRangeRule(nn ktypes.NamespacedName, rule v1alpha1.IPRangeRule) (*types.ValidationResult, error) {

	// Build the default ValidationResult for this IP range rule
	vr := buildValidationResult(rule, constants.ValidationTypeIPRange)

	errMsg := "Failed to execute IP range check"

	// parse the starting IP
	ip, err := toIpV4(rule.StartIP)
	if err != nil {
		vr.State = ptr.Ptr(v8or.ValidationFailed)
		vr.Condition.Failures = append(vr.Condition.Failures, fmt.Sprintf("invalid start IP: %s. error: %v", rule.StartIP, err))
		vr.Condition.Message = errMsg
		return vr, nil
	}

	// ping each IP in the range
	for i := 0; i < rule.Length; i++ {
		args := []string{"-c", "3", "-W", "3", ip.String()}
		stdout, stderr, _, err := execCmd(ping, args...)
		if err != nil || stderr != "" {
			if !strings.Contains(stdout, "3 packets transmitted, 0 received") {
				n.failResult(vr, err, ping, errMsg, stdout, stderr, rule.RuleName, args...)
				return vr, err
			}
			vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("%s %s failed", ping, args))
		} else {
			n.log.V(0).Info("IP allocated", "stdout", stdout, "rule", rule.RuleName)
			vr.State = ptr.Ptr(v8or.ValidationFailed)
			vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("%s %s succeeded", ping, args))
			vr.Condition.Failures = append(vr.Condition.Failures, fmt.Sprintf("stdout: %s", stdout))
			vr.Condition.Message = "One or more IPs in the provided range was allocated"
		}
		incIp(ip)
	}
	n.log.V(0).Info("IP range check passed", "rule", rule.RuleName)

	return vr, nil
}

// ReconcileMTURule reconciles an MTU rule from a NetworkValidator config
func (n *NetworkService) ReconcileMTURule(nn ktypes.NamespacedName, rule v1alpha1.MTURule) (*types.ValidationResult, error) {

	// Build the default ValidationResult for this MTU rule
	vr := buildValidationResult(rule, constants.ValidationTypeMTU)

	errMsg := "Failed to execute MTU check"
	size := strconv.FormatInt(int64(rule.MTU-icmpPacketSize), 10)
	args := []string{"-c", "3", "-W", "3", "-M", "do", "-s", size, rule.Host}

	stdout, stderr, _, err := execCmd(ping, args...)
	if err != nil || stderr != "" {
		n.failResult(vr, err, ping, errMsg, stdout, stderr, rule.RuleName, args...)
		return vr, err
	} else {
		vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("%s %s succeeded", ping, args))
		vr.Condition.Message = stdout
	}
	n.log.V(0).Info("MTU check passed", "rule", rule.RuleName)

	return vr, nil
}

// ReconcileTCPConnRule reconciles a TCP connection rule from a NetworkValidator config
func (n *NetworkService) ReconcileTCPConnRule(nn ktypes.NamespacedName, rule v1alpha1.TCPConnRule) (*types.ValidationResult, error) {

	// Build the default ValidationResult for this TCP connection rule
	vr := buildValidationResult(rule, constants.ValidationTypeTCPConn)

	errMsg := "Failed to execute TCP connection check"

	// construct args
	args := []string{"-w", "3"}
	if rule.ProxyProtocol != "" && rule.ProxyAddress == "" || rule.ProxyProtocol == "" && rule.ProxyAddress != "" {
		vr.State = ptr.Ptr(v8or.ValidationFailed)
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

	// execute netcat
	stdout, stderr, exitCode, err := execCmd(nc, args...)
	if err != nil || stderr != "" || exitCode != 0 {
		n.failResult(vr, err, ping, errMsg, stdout, stderr, rule.RuleName, args...)
		return vr, err
	} else {
		vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("%s %s succeeded", nc, args))
		vr.Condition.Message = stdout
	}
	n.log.V(0).Info("TCP connection check passed", "rule", rule.RuleName)

	return vr, nil
}

// buildValidationResult builds a default ValidationResult for a given validation type
func buildValidationResult(rule networkRule, validationType string) *types.ValidationResult {
	state := v8or.ValidationSucceeded
	latestCondition := v8or.DefaultValidationCondition()
	latestCondition.Details = make([]string, 0)
	latestCondition.Failures = make([]string, 0)
	latestCondition.Message = fmt.Sprintf("All %s checks passed", validationType)
	latestCondition.ValidationRule = rule.Name()
	latestCondition.ValidationType = validationType
	return &types.ValidationResult{Condition: &latestCondition, State: &state}
}

// handleRuleExec executes a rule's command and updates the rule's validation result accordingly
func (n *NetworkService) handleRuleExec(vr *types.ValidationResult, r networkRule, binary, errMsg string, args ...string) error {
	n.log.V(0).Info("Executing command: %s %s", binary, args)
	stdout, stderr, _, err := execCmd(binary, args...)
	if err != nil || stderr != "" {
		n.failResult(vr, err, binary, errMsg, stdout, stderr, r.Name(), args...)
		if err == nil {
			err = errors.New(errMsg)
		}
		return err
	}
	vr.Condition.Details = append(vr.Condition.Details, fmt.Sprintf("%s %s succeeded", binary, args))
	vr.Condition.Message = stdout
	return nil
}

// failResult updates a validation result with failure details
func (n *NetworkService) failResult(vr *types.ValidationResult, err error, binary, errMsg, stdout, stderr, ruleName string, args ...string) {
	n.log.V(0).Info(errMsg, "stdout", stdout, "stderr", stderr, "error", err.Error(), "rule", ruleName)
	failure := fmt.Sprintf("stdout: %s, stderr: %s, error: %v", stdout, stderr, err)
	vr.State = ptr.Ptr(v8or.ValidationFailed)
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
