package proberesponder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	httpHeaderAccept           = "Accept"
	httpHeaderContentType      = "Content-Type"
	httpHeaderContentTypeJSON  = "application/json"
	httpHeaderContentTypeHTML  = "text/html"
	httpHeaderContentTypePlain = "text/plain"
)

func HTTPStartup(pres *ProbeResponder) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		if pres.NotStarted() {
			status = http.StatusServiceUnavailable
		}
		w.WriteHeader(status)
		respond(pres, w, r)
	}
}

func HTTPReady(pres *ProbeResponder) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		if pres.NotReady() {
			status = http.StatusServiceUnavailable
		}
		w.WriteHeader(status)
		respond(pres, w, r)
	}
}

func HTTPLive(pres *ProbeResponder) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		if pres.NotLive() {
			status = http.StatusServiceUnavailable
		}
		w.WriteHeader(status)
		respond(pres, w, r)
	}
}

func respond(pres *ProbeResponder, w http.ResponseWriter, r *http.Request) {
	contentType, bPayload := contentNeogiater(r, pres.HealthResponse())
	w.Header().Add(httpHeaderContentType, contentType)
	_, err := w.Write(bPayload)
	if err != nil {
		log.Println("failed to write response", err)
	}
}

func contentNeogiater(r *http.Request, payload map[string]string) (cType string, bPayload []byte) {
	ctype := r.Header.Get(httpHeaderAccept)
	if strings.Contains(ctype, httpHeaderContentTypeHTML) {
		cType = httpHeaderContentTypeHTML
		bPayload = responseAsHTML(payload)
	} else if strings.Contains(ctype, httpHeaderContentTypePlain) {
		cType = httpHeaderContentTypePlain
		bPayload = responseAsPlainText(payload)
	} else {
		// default is json
		cType = httpHeaderContentTypeJSON
		bPayload, _ = json.Marshal(payload)
	}

	return cType, bPayload
}

func responseAsHTML(payload map[string]string) []byte {
	buff := bytes.NewBufferString(`<table><tbody>`)
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
		buff.WriteString(fmt.Sprintf("%s: %s|", key, value))
	}
	return buff.Bytes()
}

func StartHTTPServer(pres *ProbeResponder, host string, port uint16) error {
	smux := http.NewServeMux()
	smux.Handle("/-/startup", HTTPStartup(pres))
	smux.Handle("/-/ready", HTTPReady(pres))
	smux.Handle("/-/live", HTTPLive(pres))
	return http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), smux)
}
