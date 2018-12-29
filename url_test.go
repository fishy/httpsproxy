package httpsproxy_test

import (
	"net/url"
	"testing"

	"github.com/fishy/httpsproxy"
)

func TestRewriteURL(t *testing.T) {
	targetHost := "self-signed.badssl.com"
	selfURL, err := url.Parse("https://my-proxy.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Run(
		"wrong host",
		func(t *testing.T) {
			orig := "https://expired.badssl.com"
			expect := orig
			actual := httpsproxy.RewriteURL(nil, orig, targetHost, selfURL)
			if actual != expect {
				t.Errorf("RewriteURL(%q) expected %q, got %q", orig, expect, actual)
			}
		},
	)

	t.Run(
		"http",
		func(t *testing.T) {
			orig := "http://self-signed.badssl.com"
			expect := "https://my-proxy.com"
			actual := httpsproxy.RewriteURL(nil, orig, targetHost, selfURL)
			if actual != expect {
				t.Errorf("RewriteURL(%q) expected %q, got %q", orig, expect, actual)
			}
		},
	)

	t.Run(
		"https",
		func(t *testing.T) {
			orig := "https://self-signed.badssl.com/foo/bar?foo=bar&baz=qux#asdf"
			expect := "https://my-proxy.com/foo/bar?foo=bar&baz=qux#asdf"
			actual := httpsproxy.RewriteURL(nil, orig, targetHost, selfURL)
			if actual != expect {
				t.Errorf("RewriteURL(%q) expected %q, got %q", orig, expect, actual)
			}
		},
	)

	t.Run(
		"invalid url",
		func(t *testing.T) {
			orig := "al1i7y4hnelf  1lanlsu"
			expect := orig
			actual := httpsproxy.RewriteURL(nil, orig, targetHost, selfURL)
			if actual != expect {
				t.Errorf("RewriteURL(%q) expected %q, got %q", orig, expect, actual)
			}
		},
	)

	t.Run(
		"no selfURL",
		func(t *testing.T) {
			orig := "https://self-signed.badssl.com/foo/bar?foo=bar&baz=qux#asdf"
			expect := orig
			actual := httpsproxy.RewriteURL(nil, orig, targetHost, nil)
			if actual != expect {
				t.Errorf("RewriteURL(%q) expected %q, got %q", orig, expect, actual)
			}
		},
	)
}
