package client

import (
	"net"
	"net/http"
	"time"
)

// NewHTTPClient builds a *http.Client with production-grade timeouts suitable for
// CLI calls to the Smartling API. The returned client uses a fresh Transport so
// callers can mutate it (e.g. to set a proxy or TLS settings) without affecting
// other consumers.
//
// Per-stage timeouts cover connect/TLS/response-header. The overall
// Client.Timeout is left unset on purpose: file upload and download can
// legitimately take minutes.
func NewHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 5 * time.Minute,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
	}
	return &http.Client{
		Transport: transport,
	}
}
