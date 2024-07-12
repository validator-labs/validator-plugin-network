// Package validators contains network plugin validators.
package validators

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/go-logr/logr"
	"github.com/validator-labs/validator-plugin-network/api/v1alpha1"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	vapitypes "github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

// doer is an interface that defines the Do method for an HTTP client.
type doer func(req *http.Request) (*http.Response, error)

// FakeHTTPClient is a fake HTTP client that implements the httpClient interface.
type fakeHTTPClient struct {
	doer doer
}

// Do is a method that satisfies the httpClient interface. Uses the doer function to return a
// response during testing.
func (f fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return f.doer(req)
}

// createDoer is a helper function that creates a doer function that returns either an error or a
// response with the specified status code.
func createDoer(statusCode int, err error) doer {
	if err != nil {
		return func(_ *http.Request) (*http.Response, error) {
			return nil, err
		}
	}
	return func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}, nil
	}
}

func TestPublicBlobRuleService_ReconcilePublicBlobRule(t *testing.T) {

	type testCase struct {
		name           string
		rule           v1alpha1.HTTPFileRule
		httpClientMock httpClient
		expectedError  error
		expectedResult vapitypes.ValidationRuleResult
	}

	testCases := []testCase{
		{
			name: "Pass (all files found - 1 file)",
			rule: v1alpha1.HTTPFileRule{
				RuleName: "rule-1",
				Paths:    []string{"https://file1"},
			},
			httpClientMock: fakeHTTPClient{
				doer: createDoer(http.StatusOK, nil),
			},
			expectedError: nil,
			expectedResult: vapitypes.ValidationRuleResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "http-file",
					ValidationRule: "rule-1",
					Message:        "All files are publicly accessible.",
					Details: []string{
						"Ensuring that files [https://file1] are publicly accessible.",
						"File 'https://file1' is publicly accessible.",
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
				Paths:    []string{"https://file1", "https://file2"},
			},
			httpClientMock: fakeHTTPClient{
				doer: createDoer(http.StatusOK, nil),
			},
			expectedError: nil,
			expectedResult: vapitypes.ValidationRuleResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "http-file",
					ValidationRule: "rule-1",
					Message:        "All files are publicly accessible.",
					Details: []string{
						"Ensuring that files [https://file1 https://file2] are publicly accessible.",
						"File 'https://file1' is publicly accessible.",
						"File 'https://file2' is publicly accessible.",
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
				Paths:    []string{"https://file1"},
			},
			httpClientMock: fakeHTTPClient{
				doer: createDoer(0, fmt.Errorf("Do error")),
			},
			expectedError: nil,
			expectedResult: vapitypes.ValidationRuleResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "http-file",
					ValidationRule: "rule-1",
					Message:        "One or more files not publicly accessible. See failures for details.",
					Details: []string{
						"Ensuring that files [https://file1] are publicly accessible.",
					},
					Failures: []string{
						"failed to check file 'https://file1': failed to send HTTP request: Do error",
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
				Paths:    []string{"https://file1"},
			},
			httpClientMock: fakeHTTPClient{
				doer: createDoer(http.StatusNotFound, nil),
			},
			expectedError: nil,
			expectedResult: vapitypes.ValidationRuleResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "http-file",
					ValidationRule: "rule-1",
					Message:        "One or more files not publicly accessible. See failures for details.",
					Details: []string{
						"Ensuring that files [https://file1] are publicly accessible.",
					},
					Failures: []string{
						"file 'https://file1' is not publicly accessible; '404' status code in response to HEAD request",
					},
					Status: corev1.ConditionFalse,
				},
				State: util.Ptr(vapi.ValidationFailed),
			},
		},
	}

	for _, tc := range testCases {
		svc := NewNetworkService(tc.httpClientMock, logr.Logger{})
		result := svc.ReconcileHTTPFileRule(tc.rule)
		util.CheckTestCase(t, result, tc.expectedResult, nil, nil)
	}
}
