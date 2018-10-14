package httpsproxy_test

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/fishy/httpsproxy"
)

func Example() {
	target := "https://self-signed.badssl.com/"
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	// Get it by the following command:
	// openssl s_client -showcerts -connect self-signed.badssl.com:443 </dev/null
	cert := `-----BEGIN CERTIFICATE-----
MIIE8DCCAtigAwIBAgIJAM28Wkrsl2exMA0GCSqGSIb3DQEBCwUAMH8xCzAJBgNV
BAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQHDA1TYW4gRnJhbmNp
c2NvMQ8wDQYDVQQKDAZCYWRTU0wxMjAwBgNVBAMMKUJhZFNTTCBJbnRlcm1lZGlh
dGUgQ2VydGlmaWNhdGUgQXV0aG9yaXR5MB4XDTE2MDgwODIxMTcwNVoXDTE4MDgw
ODIxMTcwNVowgagxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRYw
FAYDVQQHDA1TYW4gRnJhbmNpc2NvMTYwNAYDVQQKDC1CYWRTU0wgRmFsbGJhY2su
IFVua25vd24gc3ViZG9tYWluIG9yIG5vIFNOSS4xNDAyBgNVBAMMK2JhZHNzbC1m
YWxsYmFjay11bmtub3duLXN1YmRvbWFpbi1vci1uby1zbmkwggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQDCBOz4jO4EwrPYUNVwWMyTGOtcqGhJsCK1+ZWe
sSssdj5swEtgTEzqsrTAD4C2sPlyyYYC+VxBXRMrf3HES7zplC5QN6ZnHGGM9kFC
xUbTFocnn3TrCp0RUiYhc2yETHlV5NFr6AY9SBVSrbMo26r/bv9glUp3aznxJNEx
tt1NwMT8U7ltQq21fP6u9RXSM0jnInHHwhR6bCjqN0rf6my1crR+WqIW3GmxV0Tb
ChKr3sMPR3RcQSLhmvkbk+atIgYpLrG6SRwMJ56j+4v3QHIArJII2YxXhFOBBcvm
/mtUmEAnhccQu3Nw72kYQQdFVXz5ZD89LMOpfOuTGkyG0cqFAgMBAAGjRTBDMAkG
A1UdEwQCMAAwNgYDVR0RBC8wLYIrYmFkc3NsLWZhbGxiYWNrLXVua25vd24tc3Vi
ZG9tYWluLW9yLW5vLXNuaTANBgkqhkiG9w0BAQsFAAOCAgEAsuFs0K86D2IB20nB
QNb+4vs2Z6kECmVUuD0vEUBR/dovFE4PfzTr6uUwRoRdjToewx9VCwvTL7toq3dd
oOwHakRjoxvq+lKvPq+0FMTlKYRjOL6Cq3wZNcsyiTYr7odyKbZs383rEBbcNu0N
c666/ozs4y4W7ufeMFrKak9UenrrPlUe0nrEHV3IMSF32iV85nXm95f7aLFvM6Lm
EzAGgWopuRqD+J0QEt3WNODWqBSZ9EYyx9l2l+KI1QcMalG20QXuxDNHmTEzMaCj
4Zl8k0szexR8rbcQEgJ9J+izxsecLRVp70siGEYDkhq0DgIDOjmmu8ath4yznX6A
pYEGtYTDUxIvsWxwkraBBJAfVxkp2OSg7DiZEVlMM8QxbSeLCz+63kE/d5iJfqde
cGqX7rKEsVW4VLfHPF8sfCyXVi5sWrXrDvJm3zx2b3XToU7EbNONO1C85NsUOWy4
JccoiguV8V6C723IgzkSgJMlpblJ6FVxC6ZX5XJ0ZsMI9TIjibM2L1Z9DkWRCT6D
QjuKbYUeURhScofQBiIx73V7VXnFoc1qHAUd/pGhfkCUnUcuBV1SzCEhjiwjnVKx
HJKvc9OYjJD0ZuvZw9gBrY7qKyBX8g+sglEGFNhruH8/OhqrV8pBXX/EWY0fUZTh
iywmc6GTT7X94Ze2F7iB45jh7WQ=
-----END CERTIFICATE-----`

	certPool, failed, err := httpsproxy.NewCertPool(cert)
	if err != nil {
		log.Printf("WARNING: Cannot get system cert pool: %v", err)
	}
	if len(failed) > 0 {
		log.Printf("WARNING: Failed to add cert(s) to pool: %v", failed)
	}

	selfURL, err := url.Parse("https://my-proxy.com")
	if err != nil {
		log.Printf("WARNING: Cannot get parse self URL: %v", err)
	}

	mux := httpsproxy.Mux(
		httpsproxy.DefaultHTTPClient(
			certPool,
			30*time.Second,
			httpsproxy.NoRedirCheckRedirectFunc,
		),
		targetURL,
		selfURL,
		nil, // logger
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	// Start serving with:
	// log.Fatal(http.ListenAndServe(":"+port, mux))
	_ = mux
}
