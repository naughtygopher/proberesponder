package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/naughtygopher/proberesponder"
)

const (
	httpHeaderAccept           = "Accept"
	httpHeaderContentType      = "Content-Type"
	httpHeaderContentTypeJSON  = "application/json"
	httpHeaderContentTypeXML   = "application/xml"
	httpHeaderContentTypeHTML  = "text/html"
	httpHeaderContentTypePlain = "text/plain"

	HTTPPathStartup = "/-/startup"
	HTTPPathReady   = "/-/ready"
	HTTPPathLive    = "/-/live"
)

var (
	acceptedContentTypes = strings.Join([]string{
		httpHeaderContentTypeHTML,
		httpHeaderContentTypePlain,
		httpHeaderContentTypeJSON,
	}, ",")
)

type httpMethod string
type Handlers struct {
	httpMethod
	string
	http.HandlerFunc
}

func HTTPStartup(pres *proberesponder.ProbeResponder) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		if pres.NotStarted() {
			status = http.StatusServiceUnavailable
		}
		respond(pres, w, r, status)
	}
}

func HTTPReady(pres *proberesponder.ProbeResponder) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		if pres.NotReady() {
			status = http.StatusServiceUnavailable
		}
		respond(pres, w, r, status)
	}
}

func HTTPLive(pres *proberesponder.ProbeResponder) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		if pres.NotLive() {
			status = http.StatusServiceUnavailable
		}
		respond(pres, w, r, status)
	}
}

func respond(
	pres *proberesponder.ProbeResponder,
	w http.ResponseWriter,
	r *http.Request,
	status int,
) {
	contentType, bPayload := contentNeogiater(r, pres.HealthResponse())
	w.Header().Add(httpHeaderAccept, acceptedContentTypes)
	w.Header().Add(httpHeaderContentType, contentType)
	w.WriteHeader(status)
	_, err := w.Write(bPayload)
	if err != nil {
		log.Println("failed to write response", err)
	}
}

func contentNeogiater(r *http.Request, payload map[string]string) (cType string, bPayload []byte) {
	ctypes := strings.Split(r.Header.Get(httpHeaderAccept), ",")
	maxQfactor := 0.0

	for _, ct := range ctypes {
		qFactor := 0.0
		for _, part := range strings.Split(ct, ";") {
			part = strings.TrimSpace(part)
			if strings.Contains(part, "q=") || strings.Contains(part, "Q=") {
				qFactor, _ = strconv.ParseFloat(strings.Split(part, "=")[1], 32)
				if qFactor > 1.00 || qFactor < 0 {
					qFactor = 0
				}
			}
		}

		if cType == "" || qFactor > maxQfactor {
			maxQfactor = qFactor
			cType = ct
		}
	}

	if strings.Contains(cType, httpHeaderContentTypeHTML) {
		cType = httpHeaderContentTypeHTML
		bPayload = responseAsHTML(payload)
	} else if strings.Contains(cType, httpHeaderContentTypePlain) {
		cType = httpHeaderContentTypePlain
		bPayload = responseAsPlainText(payload)
	} else if strings.Contains(cType, httpHeaderContentTypeXML) {
		cType = httpHeaderContentTypeXML
		bPayload = responseAsXML(payload)
	} else {
		cType = httpHeaderContentTypeJSON
		bPayload, _ = json.Marshal(payload)
	}

	return cType, bPayload
}

func responseAsHTML(payload map[string]string) []byte {
	buff := bytes.NewBufferString(
		`<table><tbody>`,
	)
	for key, value := range payload {
		buff.WriteString(`<tr>` +
			`<th>` + key + `</th>` +
			`<td>` + value + `</td>` +
			`</tr>`)
	}
	buff.WriteString(`</tbody></table>`)
	return buff.Bytes()
}

func responseAsPlainText(payload map[string]string) []byte {
	buff := bytes.NewBuffer([]byte{})
	for key, value := range payload {
		buff.WriteString(fmt.Sprintf("%s: %s | ", key, value))
	}
	return buff.Bytes()
}

func responseAsXML(payload map[string]string) []byte {
	buff := bytes.NewBufferString(
		`<statuses>`,
	)
	for key, value := range payload {
		buff.WriteString(`<status name="` + key + `" value="` + value + `"></status>`)
	}
	buff.WriteString(`</statuses>`)
	return buff.Bytes()
}

// Server is a basic/standard Golang HTTP server with the 3 default handlers for probes
func Server(pres *proberesponder.ProbeResponder, host string, port uint16, handlers ...Handlers) *http.Server {
	smux := http.NewServeMux()
	if len(handlers) == 0 {
		handlers = []Handlers{
			{http.MethodGet, HTTPPathStartup, HTTPStartup(pres)},
			{http.MethodGet, HTTPPathReady, HTTPReady(pres)},
			{http.MethodGet, HTTPPathLive, HTTPLive(pres)},
		}
	} else {
		handlers = append(handlers, []Handlers{
			{http.MethodGet, HTTPPathStartup, HTTPStartup(pres)},
			{http.MethodGet, HTTPPathReady, HTTPReady(pres)},
			{http.MethodGet, HTTPPathLive, HTTPLive(pres)},
		}...)
	}

	for _, h := range handlers {
		smux.Handle(h.string, h.HandlerFunc)
	}

	return &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		Handler:           smux,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       time.Second,
		WriteTimeout:      time.Second * 5,
		IdleTimeout:       time.Minute,
	}
}

// StartHTTPServer directly initializes and starts a basic HTTP probe responder
func StartHTTPServer(pres *proberesponder.ProbeResponder, host string, port uint16) error {
	return Server(pres, host, port).ListenAndServe()
}
