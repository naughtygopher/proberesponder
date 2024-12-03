package proberesponder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthOK(t *testing.T) {
	type args[T ~string] struct {
		s T
	}
	tests := []struct {
		name string
		args args[string]
		want bool
	}{
		{
			name: "Health is OK",
			args: args[string]{
				s: HealthOK.String(),
			},
			want: true,
		},
		{
			name: "Health is not OK",
			args: args[string]{
				s: HealthNotOK.String(),
			},
			want: false,
		},
		{
			name: "any other string",
			args: args[string]{
				s: "hello",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHealthOK(tt.args.s); got != tt.want {
				t.Errorf("StatusOK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus(t *testing.T) {
	if StatusStartup.String() != "startup" {
		t.Errorf("expected 'startup', got: '%s'", StatusStartup)
	}

	if StatusReady.String() != "ready" {
		t.Errorf("expected 'ready', got: '%s'", StatusReady)
	}

	if StatusLive.String() != "live" {
		t.Errorf("expected 'live', got: '%s'", StatusLive)
	}
}

func TestProbeResponder_AppendHealthResponse(tt *testing.T) {
	tt.Run("adding new key", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		const key = "key_1"
		const val = "value_1"

		pRes.AppendHealthResponse(key, val)
		asserter.Equal(val, pRes.msgPayload[key])
	})

	tt.Run("adding existing key", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		const (
			key  = "key_1"
			val  = "value_1"
			val2 = "value_2"
		)

		pRes.AppendHealthResponse(key, val)
		asserter.Equal(val, pRes.msgPayload[key])

		pRes.AppendHealthResponse(key, val2)
		asserter.Equal(val2, pRes.msgPayload[key])
	})

	tt.Run("partially initialized responder", func(t *testing.T) {
		asserter := assert.New(t)
		defer func() {
			rec := recover()
			asserter.Contains(rec, "nil pointer dereference")
		}()
		pRes := &ProbeResponder{}
		const (
			key = "key_1"
			val = "value_1"
		)
		pRes.AppendHealthResponse(key, val)
	})

	tt.Run("uninitialized responder", func(t *testing.T) {
		asserter := assert.New(t)
		defer func() {
			rec := recover()
			asserter.Nil(rec)
		}()
		var pRes *ProbeResponder
		const (
			key = "key_1"
			val = "value_1"
		)
		pRes.AppendHealthResponse(key, val)
	})
}

func TestProbeResponder_Statuses(tt *testing.T) {
	tt.Run("test defaults", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()

		asserter.True(pRes.NotLive())
		asserter.True(pRes.NotReady())
		asserter.True(pRes.NotStarted())
	})

	tt.Run("test set statuses", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		pRes.SetNotLive(false)
		asserter.False(pRes.NotLive())
		pRes.SetNotReady(false)
		asserter.False(pRes.NotReady())
		pRes.SetNotStarted(false)
		asserter.False(pRes.NotStarted())
		for key, value := range pRes.HealthResponse() {
			asserter.Contains(key, "probe->")
			asserter.NotContains(value, "NOT OK: ")
			asserter.Contains(value, "OK: ")
		}

		pRes.SetNotLive(true)
		asserter.True(pRes.NotLive())
		pRes.SetNotReady(true)
		asserter.True(pRes.NotReady())
		pRes.SetNotStarted(true)
		asserter.True(pRes.NotStarted())
		for key, value := range pRes.HealthResponse() {
			asserter.Contains(key, "probe->")
			asserter.Contains(value, "NOT OK: ")
		}
	})

	tt.Run("test get and set statuses: uninitialized", func(t *testing.T) {
		asserter := assert.New(t)
		var pRes *ProbeResponder
		defer func() {
			rec := recover()
			asserter.Nil(rec)
		}()
		pRes.SetNotLive(false)
		pRes.SetNotReady(false)
		pRes.SetNotStarted(false)

		asserter.False(pRes.NotLive())
		asserter.False(pRes.NotReady())
		asserter.False(pRes.NotStarted())
	})

	tt.Run("change listener", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		var (
			lastChangedStatus      Statuskey
			lastChangedStatusValue bool
		)

		pRes.SetListener(func(status Statuskey, value bool) {
			lastChangedStatus = status
			lastChangedStatusValue = value
		})

		pRes.SetNotLive(false)
		asserter.Equal(lastChangedStatus, StatusLive)
		asserter.False(lastChangedStatusValue)
		pRes.SetNotLive(true)
		asserter.Equal(lastChangedStatus, StatusLive)
		asserter.True(lastChangedStatusValue)

		pRes.SetNotReady(false)
		asserter.Equal(lastChangedStatus, StatusReady)
		asserter.False(lastChangedStatusValue)
		pRes.SetNotReady(true)
		asserter.Equal(lastChangedStatus, StatusReady)
		asserter.True(lastChangedStatusValue)

		pRes.SetNotStarted(false)
		asserter.Equal(lastChangedStatus, StatusStartup)
		asserter.False(lastChangedStatusValue)
		pRes.SetNotStarted(true)
		asserter.Equal(lastChangedStatus, StatusStartup)
		asserter.True(lastChangedStatusValue)
	})
}

func Test_HealthResponses(tt *testing.T) {
	tt.Run("default and custom", func(t *testing.T) {
		asserter := assert.New(t)
		pRes := New()
		hRes := pRes.HealthResponse()
		asserter.Len(hRes, 3)
		keys := []string{
			fmt.Sprintf("probe->%s", StatusLive),
			fmt.Sprintf("probe->%s", StatusReady),
			fmt.Sprintf("probe->%s", StatusStartup),
		}
		for _, key := range keys {
			asserter.Contains(hRes[key], "NOT OK:")
		}

		pRes.SetNotStarted(false)
		pRes.SetNotReady(false)
		pRes.SetNotLive(false)
		hRes = pRes.HealthResponse()
		for _, key := range keys {
			asserter.NotContains(hRes[key], "NOT OK:")
			asserter.Contains(hRes[key], "OK:")
		}
	})

	tt.Run("uninitialized", func(t *testing.T) {
		asserter := assert.New(t)
		defer func() {
			rec := recover()
			asserter.Nil(rec)
		}()
		var pRes *ProbeResponder
		pRes.SetListener(nil)
		pRes.HealthResponse()
		pRes.AppendHealthResponse("key", "value")
	})
}
