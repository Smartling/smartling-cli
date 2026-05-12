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
// Per-stage timeouts cover connect/TLS/response-header so a hung remote can
// never block forever. The overall Client.Timeout is left unset on purpose:
// file upload and download can legitimately take minutes, and request lifetime
// is controlled by the caller's context.Context.
func NewHTTPClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
	}
	return &http.Client{
		Transport: transport,
	}
}
