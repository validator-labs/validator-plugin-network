// Package http contains helpers for configuring and creating HTTP clients.
package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
)

// Transport creates a new http.Transport with insecureSkipVerify and optionally, additional CA certificates.
func Transport(caPems [][]byte, insecureSkipVerify bool, log logr.Logger) (*http.Transport, error) {
	tlsConfig, err := TLSConfig(caPems, insecureSkipVerify, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS config: %w", err)
	}
	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: tlsConfig,
	}
	return transport, nil
}

// TLSConfig creates a new tls.Config. If the system cert pool cannot be loaded, an empty pool is used.
func TLSConfig(caPems [][]byte, insecureSkipVerify bool, log logr.Logger) (*tls.Config, error) {
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.V(0).Info("failed to load system cert pool, using empty pool", "error", err)
		caCertPool = x509.NewCertPool()
	}
	for _, caPem := range caPems {
		caCertPool.AppendCertsFromPEM(caPem)
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify, // #nosec G402
		RootCAs:            caCertPool,
		MinVersion:         tls.VersionTLS12,
	}
	return tlsConfig, nil
}
