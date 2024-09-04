package http

import (
	"net/http"
	"testing"

	"github.com/go-logr/logr"
)

type mockRoundTripper struct{}

func (m mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, nil
}

func TestRoundTrip(t *testing.T) {
	transport := &basicAuthTransport{
		Username:  "user",
		Password:  "pass",
		Transport: mockRoundTripper{},
	}
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	_, err := transport.RoundTrip(req)
	if err != nil {
		t.Error("expected nil, got", err)
	}
}

func TestTransport(t *testing.T) {
	tests := []struct {
		name                   string
		caPems                 [][]byte
		auth                   []string
		insecureSkipVerify     bool
		wantBasicAuthTransport bool
	}{
		{
			name:                   "Pass: no auth",
			caPems:                 nil,
			auth:                   nil,
			insecureSkipVerify:     false,
			wantBasicAuthTransport: false,
		},
		{
			name:                   "Pass: with auth",
			caPems:                 nil,
			auth:                   []string{"user", "pass"},
			insecureSkipVerify:     false,
			wantBasicAuthTransport: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Transport(tt.caPems, tt.auth, tt.insecureSkipVerify, logr.Discard())
			if tt.wantBasicAuthTransport {
				if _, ok := got.(*basicAuthTransport); !ok {
					t.Errorf("Transport() = %v, want *basicAuthTransport", got)
				}
			} else {
				if _, ok := got.(*http.Transport); !ok {
					t.Errorf("Transport() = %v, want *http.Transport", got)
				}
			}
		})
	}
}
