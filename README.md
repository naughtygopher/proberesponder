<p align="center"><img src="https://github.com/user-attachments/assets/2bd99e22-d0fa-464f-8dca-3336ec7b6e0b" alt="proberesponder gopher" width="256px"/></p>

[![](https://github.com/naughtygopher/proberesponder/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/naughtygopher/proberesponder/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/naughtygopher/proberesponder.svg)](https://pkg.go.dev/github.com/naughtygopher/proberesponder)
[![Go Report Card](https://goreportcard.com/badge/github.com/naughtygopher/proberesponder?cache_invalidate=v0.3.0)](https://goreportcard.com/report/github.com/naughtygopher/proberesponder)
[![Coverage Status](https://coveralls.io/repos/github/naughtygopher/proberesponder/badge.svg?branch=main&cache_invalidate=v0.3.0)](https://coveralls.io/github/naughtygopher/proberesponder?branch=main)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/creativecreature/sturdyc/blob/master/LICENSE)

# Proberesponder

Probe-responder is a package to deal with handling [appropriate statuses for Kuberentes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/). Even though the statuses are specific to Kubernetes, it can be used in any context, and can later be updated to include more if required.

The sample below shows how to use proberesponder, all the statuses are "NOT OK" by default. This is intentional as the app is expected to explicitly update the respective status as OK which would be more accurate than being OK by default.

## Extras

By default a bare bones HTTP server can be setup to respond to probe request. The default HTTP handlers provided does content negotiation and provides appropriate response for JSON, HTML & plain text. For any unidentified content type, it will respond with JSON.

`AppendHealthResponse` is a helper function with which you can maintain statuses of a dependency or similar. All the custom statuses set using this and the native ones (startup, live, ready) can be fetched as a map[string]string using `HealthResponse`.

## Sample usage

```golang
package main
import (
    "github.com/naughtygopher/proberesponder"
)

func main() {
    pRes := proberesponder.New()
    // setup an HTTP server to handle probe requests
    go proberesponder.StartHTTPServer(pres, "localhost", 1234)

    // with set listener you can register a callback, for when any of the statuses
    // (startup, live, ready) is changed
    pRes.SetListener(func(status Statuskey, value bool) {
        fmt.Println(status, "changed to", value)
    })

    // Update the status of the app as Startup: OK
    pRes.SetNotStarted(false)

    // update the status of app as Live: OK
    pRes.SetNotLive(false)
    // update the status of app as Ready: OK
    pRes.SetNotReady(false)

    // set status of any service
    pRes.AppendHealthResponse("mydb", "OK")

    // retrieves all the statuses maintained by the proberesponder, it returns a map[string]string
    _ = pRes.HealthResponse()
}
```

Below is an example of probe responses for HTTP requests using curl.

```bash
$ curl -H 'Accept: text/plain' localhost:1234/-/startup
mydb: OK | probe->live: OK: 2025-01-09T17:45:24+01:00 | probe->ready: OK: 2025-01-09T17:45:24+01:00 | probe->startup: OK: 2025-01-09T17:45:24+01:00 |

$ curl -H 'Accept: text/plain' localhost:1234/-/ready
probe->ready: OK: 2025-01-09T17:45:24+01:00 | probe->startup: OK: 2025-01-09T17:45:24+01:00 | mydb: OK | probe->live: OK: 2025-01-09T17:45:24+01:00 |

$ curl -H 'Accept: text/plain' localhost:1234/-/live
probe->startup: OK: 2025-01-09T17:45:24+01:00 | mydb: OK | probe->live: OK: 2025-01-09T17:45:24+01:00 | probe->ready: OK: 2025-01-09T17:45:24+01:00 |
```

## The gopher

The gopher used here was created using [Gopherize.me](https://gopherize.me/). Just like the handyman gopher here, proberesponder helps setup, well, a probe responder!
