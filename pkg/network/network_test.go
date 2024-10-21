package network

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"github.com/validator-labs/validator-plugin-network/api/v1alpha1"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/test"
	vapitypes "github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

func TestReconcileHTTPFileRule(t *testing.T) {
	testCases := []struct {
		name           string
		rule           v1alpha1.HTTPFileRule
		handler        http.HandlerFunc
		expectedResult vapitypes.ValidationRuleResult
	}{
		// Test cases use httptest, and we don't know the server URL (inc. its port) until the tests
		// run, so the test case Details and Failures are templates. When running the tests, it subs
		// in the real server URL of the test server before asserting.
		{
			name: "Pass (all files found - 1 file)",
			rule: v1alpha1.HTTPFileRule{
				RuleName: "rule-1",
				Paths:    []string{"/file1"},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectedResult: vapitypes.ValidationRuleResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "http-file",
					ValidationRule: "rule-1",
					Message:        "All files are accessible.",
					Details: []string{
						"Ensuring that files [{{serverUrl}}/file1] are accessible.",
						"File {{serverUrl}}/file1 is accessible.",
					},
					Failures: []string{},
					Status:   corev1.ConditionTrue,
				},
				State: util.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Pass (all files found - 2 files)",
			rule: v1alpha1.HTTPFileRule{
				RuleName: "rule-1",
				Paths:    []string{"/file1", "/file2"},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectedResult: vapitypes.ValidationRuleResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "http-file",
					ValidationRule: "rule-1",
					Message:        "All files are accessible.",
					Details: []string{
						"Ensuring that files [{{serverUrl}}/file1 {{serverUrl}}/file2] are accessible.",
						"File {{serverUrl}}/file1 is accessible.",
						"File {{serverUrl}}/file2 is accessible.",
					},
					Failures: []string{},
					Status:   corev1.ConditionTrue,
				},
				State: util.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Fail (error making HTTP request)",
			rule: v1alpha1.HTTPFileRule{
				RuleName: "rule-1",
				Paths:    []string{"/file1"},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Close the connection immediately to simulate an error
				hj, ok := w.(http.Hijacker)
				if !ok {
					http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
					return
				}
				conn, _, err := hj.Hijack()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				_ = conn.Close()
			},
			expectedResult: vapitypes.ValidationRuleResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "http-file",
					ValidationRule: "rule-1",
					Message:        "HTTPFile check(s) failed. See failures for details.",
					Details: []string{
						"Ensuring that files [{{serverUrl}}/file1] are accessible.",
					},
					Failures: []string{
						`failed to check file {{serverUrl}}/file1: failed to send HTTP request: Head "{{serverUrl}}/file1": EOF`,
					},
					Status: corev1.ConditionFalse,
				},
				State: util.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Fail (HTTP request succeeds but has non-200 OK response)",
			rule: v1alpha1.HTTPFileRule{
				RuleName: "rule-1",
				Paths:    []string{"/file1"},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectedResult: vapitypes.ValidationRuleResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "http-file",
					ValidationRule: "rule-1",
					Message:        "HTTPFile check(s) failed. See failures for details.",
					Details: []string{
						"Ensuring that files [{{serverUrl}}/file1] are accessible.",
					},
					Failures: []string{
						"file {{serverUrl}}/file1 is not accessible; '404' status code in response to HEAD request",
					},
					Status: corev1.ConditionFalse,
				},
				State: util.Ptr(vapi.ValidationFailed),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			// Update the rule's paths and the expected result to use the test server's URL
			for i := range tc.rule.Paths {
				tc.rule.Paths[i] = server.URL + tc.rule.Paths[i]
				for j := range tc.expectedResult.Condition.Details {
					tc.expectedResult.Condition.Details[j] =
						strings.ReplaceAll(tc.expectedResult.Condition.Details[j], "{{serverUrl}}", server.URL)
				}
				for k := range tc.expectedResult.Condition.Failures {
					tc.expectedResult.Condition.Failures[k] =
						strings.ReplaceAll(tc.expectedResult.Condition.Failures[k], "{{serverUrl}}", server.URL)
				}
			}

			// Create the network service with the default HTTP client
			svc := NewRuleService(logr.Logger{})

			// Test
			result := svc.ReconcileHTTPFileRule(tc.rule)
			test.CheckTestCase(t, result, tc.expectedResult, nil, nil)
		})
	}
}
