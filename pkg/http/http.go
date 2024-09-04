// Package http contains helpers for configuring and creating HTTP clients.
package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
)

// basicAuthTransport is an http.RoundTripper that adds HTTP Basic Authentication.
type basicAuthTransport struct {
	Username  string
	Password  string
	Transport http.RoundTripper
}

// RoundTrip executes a single HTTP transaction.
func (t *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add basic auth header
	auth := base64.StdEncoding.EncodeToString([]byte(t.Username + ":" + t.Password))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))

	// Delegate to the original transport
	return t.Transport.RoundTrip(req)
}

// Transport creates a new http.RoundTripper with insecureSkipVerify and optionally, additional CA certificates
// and/or HTTP basic authentication configured.
func Transport(caPems [][]byte, auth []string, insecureSkipVerify bool, log logr.Logger) http.RoundTripper {
	var transport http.RoundTripper

	transport = &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: TLSConfig(caPems, insecureSkipVerify, log),
	}
	if len(auth) > 0 {
		transport = &basicAuthTransport{
			Username:  auth[0],
			Password:  auth[1],
			Transport: transport,
		}
	}

	return transport
}

// TLSConfig creates a new tls.Config. If the system cert pool cannot be loaded, an empty pool is used.
func TLSConfig(caPems [][]byte, insecureSkipVerify bool, log logr.Logger) *tls.Config {
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
	return tlsConfig
}
