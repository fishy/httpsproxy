[![PkgGoDev](https://pkg.go.dev/badge/go.yhsif.com/httpsproxy)](https://pkg.go.dev/go.yhsif.com/httpsproxy)
[![Go Report Card](https://goreportcard.com/badge/go.yhsif.com/httpsproxy)](https://goreportcard.com/report/go.yhsif.com/httpsproxy)

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
[pkg.go.dev example](https://pkg.go.dev/go.yhsif.com/httpsproxy?tab=doc#example-package)
or
[blynk-proxy code](https://github.com/fishy/blynk-proxy/blob/master/main.go).

## License

[BSD 3-Clause](LICENSE).
