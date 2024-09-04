package v1alpha1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHTTPFileAuthBytesDirect(t *testing.T) {
	tests := []struct {
		name     string
		spec     NetworkValidatorSpec
		expected map[string][]string
	}{
		{
			name: "No HTTPFileRules",
			spec: NetworkValidatorSpec{
				HTTPFileRules: []HTTPFileRule{},
			},
			expected: map[string][]string{},
		},
		{
			name: "HTTPFileRule without basic auth",
			spec: NetworkValidatorSpec{
				HTTPFileRules: []HTTPFileRule{
					{
						RuleName: "rule1",
					},
				},
			},
			expected: map[string][]string{},
		},
		{
			name: "Single HTTPFileRule with basic auth",
			spec: NetworkValidatorSpec{
				HTTPFileRules: []HTTPFileRule{
					{
						RuleName: "rule1",
						Auth: Auth{
							Basic: &BasicAuth{
								Username: "user1",
								Password: "pass1",
							},
						},
					},
				},
			},
			expected: map[string][]string{
				"rule1": {"user1", "pass1"},
			},
		},
		{
			name: "Multiple HTTPFileRules with and without basic auth",
			spec: NetworkValidatorSpec{
				HTTPFileRules: []HTTPFileRule{
					{
						RuleName: "rule1",
						Auth: Auth{
							Basic: &BasicAuth{
								Username: "user1",
								Password: "pass1",
							},
						},
					},
					{
						RuleName: "rule2",
						Auth: Auth{
							Basic: &BasicAuth{
								Username: "user2",
								Password: "pass2",
							},
						},
					},
					{
						RuleName: "rule3",
					},
				},
			},
			expected: map[string][]string{
				"rule1": {"user1", "pass1"},
				"rule2": {"user2", "pass2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.spec.HTTPFileAuthBytesDirect()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
