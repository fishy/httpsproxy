package httpsproxy_test

import (
	"crypto/x509"
	"fmt"
	"strings"
	"testing"

	"github.com/fishy/httpsproxy"
)

var validCert = `-----BEGIN CERTIFICATE-----
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

var invalidCert = "asdf"

func TestCertPool(t *testing.T) {
	basePool, baseErr := getBaselinePool(t)

	t.Run(
		"no new certs",
		func(t *testing.T) {
			newPool, failedCerts, newErr := httpsproxy.NewCertPool()
			if len(failedCerts) != 0 {
				t.Errorf("failedCerts expected to be empty, got: %v", failedCerts)
			}
			if newErr != baseErr {
				t.Errorf("Expected error %v, got %v", baseErr, newErr)
			}
			comparePools(t, basePool, newPool, 0)
		},
	)

	t.Run(
		"valid cert",
		func(t *testing.T) {
			newPool, failedCerts, newErr := httpsproxy.NewCertPool(validCert)
			if len(failedCerts) != 0 {
				t.Errorf("failedCerts expected to be empty, got: %v", failedCerts)
			}
			if newErr != baseErr {
				t.Errorf("Expected error %v, got %v", baseErr, newErr)
			}
			comparePools(t, basePool, newPool, 1)
		},
	)

	t.Run(
		"invalid cert",
		func(t *testing.T) {
			newPool, failedCerts, newErr := httpsproxy.NewCertPool(invalidCert)
			if len(failedCerts) != 1 || failedCerts[0] != invalidCert {
				t.Errorf("failedCerts expected 1, got: %v", failedCerts)
			}
			if newErr != baseErr {
				t.Errorf("Expected error %v, got %v", baseErr, newErr)
			}
			comparePools(t, basePool, newPool, 0)
		},
	)

	t.Run(
		"mixed certs",
		func(t *testing.T) {
			newPool, failedCerts, newErr := httpsproxy.NewCertPool(
				validCert,
				invalidCert,
			)
			if len(failedCerts) != 1 || failedCerts[0] != invalidCert {
				t.Errorf("failedCerts expected 1, got: %v", failedCerts)
			}
			if newErr != baseErr {
				t.Errorf("Expected error %v, got %v", baseErr, newErr)
			}
			comparePools(t, basePool, newPool, 1)
		},
	)
}

func getBaselinePool(t *testing.T) (*x509.CertPool, error) {
	t.Helper()

	pool, err := x509.SystemCertPool()
	if err != nil {
		pool = x509.NewCertPool()
	}
	return pool, err
}

func comparePools(t *testing.T, base, newPool *x509.CertPool, diff uint) {
	t.Helper()

	inBase, inNew := subjectDiff(base.Subjects(), newPool.Subjects())

	if len(inBase) > 0 {
		t.Errorf(
			"The following certs are not in new pool: %s",
			subjectsToString(inBase),
		)
	}

	if len(inNew) != int(diff) {
		t.Errorf(
			"Exepcted %d new certs, got %s",
			len(inNew),
			subjectsToString(inNew),
		)
	}

	t.Logf("New certs: %s", subjectsToString(inNew))
}

func subjectDiff(a, b [][]byte) (inA, inB [][]byte) {
	mapA := subjectsToMap(a)
	mapB := subjectsToMap(b)

	for _, sub := range a {
		if !mapB[string(sub)] {
			inA = append(inA, sub)
		}
	}

	for _, sub := range b {
		if !mapA[string(sub)] {
			inB = append(inB, sub)
		}
	}

	return
}

func subjectsToMap(subs [][]byte) map[string]bool {
	ret := make(map[string]bool)
	for _, sub := range subs {
		ret[string(sub)] = true
	}
	return ret
}

func subjectsToString(subs [][]byte) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%d:[", len(subs)))
	for i, sub := range subs {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%q", sub))
	}
	builder.WriteString("]")
	return builder.String()
}
