[![](https://github.com/naughtygopher/proberesponder/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/naughtygopher/proberesponder/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/naughtygopher/proberesponder.svg)](https://pkg.go.dev/github.com/naughtygopher/proberesponder)
[![Go Report Card](https://goreportcard.com/badge/github.com/naughtygopher/proberesponder?cache_invalidate=v0.3.0)](https://goreportcard.com/report/github.com/naughtygopher/proberesponder)
[![Coverage Status](https://coveralls.io/repos/github/naughtygopher/proberesponder/badge.svg?branch=main&cache_invalidate=v0.3.0)](https://coveralls.io/github/naughtygopher/proberesponder?branch=main)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/creativecreature/sturdyc/blob/master/LICENSE)

# Proberesponder

Probe-responder is a package to deal with handling [appropriate statuses for Kuberentes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/). Even though the statuses are specific to Kubernetes, it can be used in any context, and can later be updated to include more if required.

# Sample usage

```golang
package main
import (
    "github.com/naughtygopher/proberesponder"
)

func main() {
    pRes := proberesponder.New()
    // do something to startup your app, once finished
    pRes.SetNotStarted(false)

    // if you're running an API server (HTTP, gRPC etc.), start the servers and then.
    pRes.SetNotLive(false)
	pRes.SetNotReady(false)
}
```
