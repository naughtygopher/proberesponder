package proberesponder

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPStartup(tt *testing.T) {
	requirer := require.New(tt)
	tt.Run("NotStarted: true", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		pRes.SetNotStarted(true)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypePlain)
		requirer.NoError(err)
		handler := HTTPStartup(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->startup: NOT OK:")
		asserter.Equal(http.StatusServiceUnavailable, w.Result().StatusCode)
	})
	tt.Run("NotStarted: false", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		pRes.SetNotStarted(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypePlain)
		requirer.NoError(err)
		handler := HTTPStartup(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->startup: OK:")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: plain text", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		pRes.SetNotStarted(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypePlain)
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
		pRes := New()
		pRes.SetNotStarted(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypeHTML)
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
		pRes := New()
		pRes.SetNotStarted(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypeJSON)
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
	tt.Run("NotReady: true", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		pRes.SetNotReady(true)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypePlain)
		requirer.NoError(err)
		handler := HTTPReady(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->ready: NOT OK:")
		asserter.Equal(http.StatusServiceUnavailable, w.Result().StatusCode)
	})
	tt.Run("NotReady: false", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		pRes.SetNotReady(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypePlain)
		requirer.NoError(err)
		handler := HTTPReady(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->ready: OK:")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: plain text", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		pRes.SetNotReady(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypePlain)
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
		pRes := New()
		pRes.SetNotReady(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypeHTML)
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
		pRes := New()
		pRes.SetNotReady(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypeJSON)
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
	tt.Run("NotLive: true", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		pRes.SetNotLive(true)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypePlain)
		requirer.NoError(err)
		handler := HTTPLive(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->live: NOT OK:")
		asserter.Equal(http.StatusServiceUnavailable, w.Result().StatusCode)
	})
	tt.Run("NotLive: false", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		pRes.SetNotLive(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypePlain)
		requirer.NoError(err)
		handler := HTTPLive(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "|")
		asserter.Contains(text, "probe->live: OK:")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: plain text", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		pRes.SetNotLive(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypePlain)
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
		pRes := New()
		pRes.SetNotLive(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypeHTML)
		requirer.NoError(err)
		handler := HTTPLive(pRes)
		handler(w, r)

		text := w.Body.String()
		asserter.Contains(text, "<table><tbody>")
		asserter.Contains(text, "probe->live")
		asserter.Equal(http.StatusOK, w.Result().StatusCode)
	})

	tt.Run("Content negotiation: JSON", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		pRes.SetNotLive(false)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		r.Header.Set(httpHeaderContentType, httpHeaderContentTypeJSON)
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
