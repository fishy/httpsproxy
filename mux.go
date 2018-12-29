package httpsproxy

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"net/url"
)

var requestHeadersToCopy = []string{
	"Content-Type",
	"User-Agent",
}

var client *http.Client

// Mux creates an http serve mux to do the proxy job.
//
// The returned mux contains a single handler for "/" to the handler generated
// by ProxyRootHandler to do the proxy.
// You can add your own handlers to handle cases like health check.
//
// Refer to the doc of ProxyRootHandler for the more detailed explanations of
// the args.
func Mux(
	client *http.Client,
	targetURL, selfURL *url.URL,
	logger *log.Logger,
) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", ProxyRootHandler(client, targetURL, selfURL, logger))
	return mux
}

// ProxyRootHandler generates the http handler function to be used to serve root
// ("/") in http mux.
//
// client is the http client to use. You can either use DefaultHTTPClient
// function to get a default implementation, or refer to its code to create your
// own. You migh also find github.com/fishy/badcerts library useful when
// creating your own client.
//
// targetURL is the target URL this mux proxies to. Only its scheme and host
// will be used.
//
// selfURL is for 3xx redirect rewrite. It could be nil, which means this mux
// won't rewrite any 3xx responses.
//
// logger is the logger to log errors. It could be nil, which means no errors
// will be logged.
func ProxyRootHandler(
	client *http.Client,
	targetURL, selfURL *url.URL,
	logger *log.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		newURL := &url.URL{
			Scheme: targetURL.Scheme,
			Host:   targetURL.Host,
			// In incoming r.URL only these 2 fields are set:
			Path:     r.URL.Path,
			RawQuery: r.URL.RawQuery,
		}
		req, err := http.NewRequest(r.Method, newURL.String(), r.Body)
		if CheckError(logger, w, err) {
			return
		}
		req.Header.Set("X-Forwarded-For", r.RemoteAddr)
		CopyRequestHeaders(r, req, requestHeadersToCopy)

		resp, err := client.Do(req)
		if CheckError(logger, w, err) {
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if CheckError(logger, w, err) {
			return
		}

		header := w.Header()
		for key, values := range resp.Header {
			canonicalKey := textproto.CanonicalMIMEHeaderKey(key)
			for _, value := range values {
				if canonicalKey == "Location" {
					value = RewriteURL(logger, value, targetURL.Host, selfURL)
				}
				header.Add(canonicalKey, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		if _, err := w.Write(body); err != nil {
			if logger != nil {
				logger.Print(err)
			}
		}
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

// CheckError checks error. If error is non-nil, it writes HTTP status code 502
// (bad gateway) and the error message to the response and returns true.
func CheckError(logger *log.Logger, w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	if logger != nil {
		logger.Print(err)
	}
	w.WriteHeader(http.StatusBadGateway)
	if _, err := w.Write([]byte(err.Error())); err != nil {
		if logger != nil {
			logger.Print(err)
		}
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return true
}

// CopyRequestHeaders copies specified headers from one http.Request to another.
func CopyRequestHeaders(from, to *http.Request, headers []string) {
	for _, header := range headers {
		value := from.Header.Get(header)
		if value != "" {
			to.Header.Set(header, value)
		}
	}
}

// RewriteURL rewrites all targetHost URLs to us (selfURL).
func RewriteURL(
	logger *log.Logger,
	origURL, targetHost string,
	selfURL *url.URL,
) string {
	if selfURL == nil {
		return origURL
	}

	u, err := url.Parse(origURL)
	if err != nil {
		if logger != nil {
			logger.Print(err)
		}
		return origURL
	}
	if u.Host == targetHost {
		u.Scheme = selfURL.Scheme
		u.Host = selfURL.Host
		return u.String()
	}
	return origURL
}
