// Package http contains helpers for configuring and creating HTTP clients.
package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
)

// TransportWithCA creates a new http.Transport with the provided CA certificates and
// insecureSkipVerify. If loading the system cert pool fails, it must be provided with at least one
// alternative cert to use. If no alternative certs are provided, the system cert pool must load.
func TransportWithCA(caPems [][]byte, insecureSkipVerify bool) (*http.Transport, error) {
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		if len(caPems) == 0 {
			return nil, fmt.Errorf("failed to get system cert pool and no alternative cert provided: %w", err)
		}
		// No error occurred and at least one alternative cert was provided, so start a new pool.
		caCertPool = x509.NewCertPool()
	}

	for _, caPem := range caPems {
		caCertPool.AppendCertsFromPEM(caPem)
	}

	return &http.Transport{
		// #nosec G402 -- InsecureSkipVerify defaults to false but users can set the property in the
		// spec to true to disable server certificate verification.
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
			RootCAs:            caCertPool,
			MinVersion:         tls.VersionTLS12,
		},
	}, nil
}
