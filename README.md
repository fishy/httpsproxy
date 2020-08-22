[![PkgGoDev](https://pkg.go.dev/badge/github.com/fishy/httpsproxy)](https://pkg.go.dev/github.com/fishy/httpsproxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishy/httpsproxy)](https://goreportcard.com/report/github.com/fishy/httpsproxy)

# httpsproxy

`httpsproxy` is a [Go](https://golang.org)
library provides an http serve mux that can work as an HTTPS proxy for a site
with self-signed https certificate.

## Why?

The main user of this library is
[`blynk-proxy`](https://github.com/fishy/blynk-proxy),
please refer to its
[README](https://github.com/fishy/blynk-proxy/blob/master/README.md)
for more information.

This library is moved out of `blynk-proxy` project because I believe others
facing similar situation could benefit from it.

## Example

Please refer to
[pkg.go.dev example](https://pkg.go.dev/github.com/fishy/httpsproxy?tab=doc#example-package)
or
[blynk-proxy code](https://github.com/fishy/blynk-proxy/blob/master/main.go).

## License

BSD 3-Clause, refer to LICENSE file for more details.
