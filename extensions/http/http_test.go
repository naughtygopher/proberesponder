package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/naughtygopher/proberesponder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPStartup(tt *testing.T) {
	requirer := require.New(tt)
	const (
		httpHost = "localhost"
		httpPort = 1234
	)
	addr := fmt.Sprintf("http://%s:%d%s", httpHost, httpPort, HTTPPathStartup)
	tt.Run("NotStarted: true", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotStarted(true)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, addr, nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypePlain)
		requirer.NoError(err)
		Server(pRes, httpHost, httpPort).Handler.ServeHTTP(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->startup: NOT OK:")
		asserter.Equal(http.StatusServiceUnavailable, w.Result().StatusCode)
	})
	tt.Run("NotStarted: false", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotStarted(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, addr, nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypePlain)
		requirer.NoError(err)
		Server(pRes, httpHost, httpPort).Handler.ServeHTTP(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->startup: OK:")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: plain text", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotStarted(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypePlain)
		requirer.NoError(err)
		handler := HTTPStartup(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->startup:")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: HTML", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotStarted(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypeHTML)
		requirer.NoError(err)
		handler := HTTPStartup(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "<table><tbody>")
		asserter.Contains(text, "probe->startup")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: JSON", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotStarted(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypeJSON)
		requirer.NoError(err)
		handler := HTTPStartup(pRes)
		handler(w, r)

		payload := map[string]string{}
		jbytes := w.Body.Bytes()
		asserter.NoError(json.Unmarshal(jbytes, &payload))
		asserter.Len(payload, 3)
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
		asserter.Contains(payload["probe->startup"], "OK:")
	})
}

func TestHTTPReady(tt *testing.T) {
	requirer := require.New(tt)
	const (
		httpHost = "localhost"
		httpPort = 1234
	)
	addr := fmt.Sprintf("http://%s:%d%s", httpHost, httpPort, HTTPPathReady)
	tt.Run("NotReady: true", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotReady(true)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, addr, nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypePlain)
		requirer.NoError(err)
		Server(pRes, httpHost, httpPort).Handler.ServeHTTP(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->ready: NOT OK:")
		asserter.Equal(http.StatusServiceUnavailable, w.Result().StatusCode)
	})
	tt.Run("NotReady: false", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotReady(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, addr, nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypePlain)
		requirer.NoError(err)
		Server(pRes, httpHost, httpPort).Handler.ServeHTTP(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->ready: OK:")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: plain text", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotReady(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypePlain)
		requirer.NoError(err)
		handler := HTTPReady(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->ready:")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: HTML", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotReady(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypeHTML)
		requirer.NoError(err)
		handler := HTTPReady(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "<table><tbody>")
		asserter.Contains(text, "probe->ready")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: JSON", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotReady(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypeJSON)
		requirer.NoError(err)
		handler := HTTPReady(pRes)
		handler(w, r)

		payload := map[string]string{}
		jbytes := w.Body.Bytes()
		asserter.NoError(json.Unmarshal(jbytes, &payload))
		asserter.Len(payload, 3)
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
		asserter.Contains(payload["probe->ready"], "OK:")
	})
}

func TestHTTPLive(tt *testing.T) {
	requirer := require.New(tt)
	const (
		httpHost = "localhost"
		httpPort = 1234
	)
	addr := fmt.Sprintf("http://%s:%d%s", httpHost, httpPort, HTTPPathLive)
	tt.Run("NotLive: true", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotLive(true)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, addr, nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypePlain)
		requirer.NoError(err)
		Server(pRes, httpHost, httpPort).Handler.ServeHTTP(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->live: NOT OK:")
		asserter.Equal(http.StatusServiceUnavailable, w.Result().StatusCode)
	})
	tt.Run("NotLive: false", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotLive(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, addr, nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypePlain)
		requirer.NoError(err)
		Server(pRes, httpHost, httpPort).Handler.ServeHTTP(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->live: OK:")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: plain text", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotLive(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypePlain)
		requirer.NoError(err)
		handler := HTTPLive(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->live:")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: HTML", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotLive(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypeHTML)
		requirer.NoError(err)
		handler := HTTPLive(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "<table><tbody>")
		asserter.Contains(text, "probe->live")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: XML", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotLive(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypeXML)
		requirer.NoError(err)
		handler := HTTPLive(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "<statuses>")
		asserter.Contains(text, "probe->live")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})
	tt.Run("Content negotiation: JSON", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := proberesponder.New()
		pRes.SetNotLive(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderAccept, httpHeaderContentTypeJSON)
		requirer.NoError(err)
		handler := HTTPLive(pRes)
		handler(w, r)

		payload := map[string]string{}
		jbytes := w.Body.Bytes()
		asserter.NoError(json.Unmarshal(jbytes, &payload))
		asserter.Len(payload, 3)
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
		asserter.Contains(payload["probe->live"], "OK:")
	})
}

func Test_contentNeogiater(t *testing.T) {
	type args struct {
		r       *http.Request
		payload map[string]string
	}
	tests := []struct {
		name      string
		args      args
		wantCType string
	}{
		{
			name: "empty accept header",
			args: args{
				r: httpReq(""),
			},
			wantCType: httpHeaderContentTypeJSON,
		},
		{
			name: "accept plain text",
			args: args{
				r: httpReq(httpHeaderContentTypePlain),
			},
			wantCType: httpHeaderContentTypePlain,
		},
		{
			name: "accept HTML",
			args: args{
				r: httpReq(httpHeaderContentTypeHTML),
			},
			wantCType: httpHeaderContentTypeHTML,
		},
		{
			name: "accept XML",
			args: args{
				r: httpReq(httpHeaderContentTypeXML),
			},
			wantCType: httpHeaderContentTypeXML,
		},
		{
			name: "accept JSON",
			args: args{
				r: httpReq(httpHeaderContentTypeJSON),
			},
			wantCType: httpHeaderContentTypeJSON,
		},
		{
			name: "With Quality factor",
			args: args{
				r: httpReq("application/json;q=0.25,application/xml;q=0.5"),
			},
			wantCType: httpHeaderContentTypeXML,
		},
		{
			name: "With Quality factor less than 0",
			args: args{
				r: httpReq("application/json;q=0.25,application/xml;q=-0.5"),
			},
			wantCType: httpHeaderContentTypeJSON,
		},
		{
			name: "With Quality factor greater than 1",
			args: args{
				r: httpReq("application/json;q=0.25,application/xml;q=1.5"),
			},
			wantCType: httpHeaderContentTypeJSON,
		},
		{
			name: "With high Quality factor first",
			args: args{
				r: httpReq("application/json;q=0.55,application/xml;q=0.455"),
			},
			wantCType: httpHeaderContentTypeJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCType, _ := contentNeogiater(tt.args.r, tt.args.payload)
			if gotCType != tt.wantCType {
				t.Errorf("contentNeogiater() gotCType = %v, want %v", gotCType, tt.wantCType)
			}
		})
	}
}

func TestCustomHandler(tt *testing.T) {
	tt.Run("custom handler success", func(t *testing.T) {
		expectedResponse := "success"
		srv := Server(proberesponder.New(), "", 1234, Handler{http.MethodGet, "/mypath", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(expectedResponse))
		}})
		req, _ := http.NewRequest(http.MethodGet, "/mypath", nil)
		w := httptest.NewRecorder()
		srv.Handler.ServeHTTP(w, req)
		assert.Equal(tt, expectedResponse, w.Body.String())
	})
	tt.Run("unmatching method", func(t *testing.T) {
		srv := Server(proberesponder.New(), "", 1234, Handler{http.MethodPost, "/mypath", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("success"))
		}})
		req, _ := http.NewRequest(http.MethodGet, "/mypath", nil)
		w := httptest.NewRecorder()
		srv.Handler.ServeHTTP(w, req)
		assert.Equal(t, "", w.Body.String())
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}

func httpReq(acceptType string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:1234", nil)
	req.Header.Add(httpHeaderAccept, acceptType)
	return req
}
