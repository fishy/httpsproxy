package httpsproxy

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"
)

// DefaultHTTPClient returns an http client that can be used in Mux function
// with:
//
// * certPool: the x509 cert pool to trust.
//
// * timeout: the http timeout.
//
// * checkRedirectFunc: the function to handle 3xx redirects, could be nil which
//                      means default behavior.
func DefaultHTTPClient(
	certPool *x509.CertPool,
	timeout time.Duration,
	checkRedirectFunc func(*http.Request, []*http.Request) error,
) *http.Client {
	return &http.Client{
		CheckRedirect: checkRedirectFunc,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
		Timeout: timeout,
	}
}

// NoRedirCheckRedirectFunc is a CheckRedirect function implemention can be used
// in http.Client. It does not follow any redirections.
func NoRedirCheckRedirectFunc(*http.Request, []*http.Request) error {
	// Don't follow any redirects.
	return http.ErrUseLastResponse
}
