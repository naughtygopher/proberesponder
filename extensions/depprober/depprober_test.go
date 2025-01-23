// Package depprober, provides some utility functions to probe dependencies of an application.
// The statuses of probed dependencies are then made part of the proberesponder payload.
// Though the package requires Prober interface to be implemented
package depprober

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/naughtygopher/proberesponder"
	"github.com/stretchr/testify/assert"
)

type DummyPinger struct {
	serviceID      string
	affectedStatus []proberesponder.Statuskey
	err            error
}

func (dp *DummyPinger) Check(ctx context.Context) error {
	return dp.err
}

func (dp *DummyPinger) ServiceID() string {
	return dp.serviceID
}

func (dp *DummyPinger) AffectsStatuses() []proberesponder.Statuskey {
	return dp.affectedStatus
}

func newProbeRespWithAllOK() *proberesponder.ProbeResponder {
	pResp := proberesponder.New()
	// ensure proberesponder is all OK
	pResp.SetNotStarted(false)
	pResp.SetNotReady(false)
	pResp.SetNotLive(false)
	return pResp
}

func assertStatuses(asserter *assert.Assertions, pResp *proberesponder.ProbeResponder, probers ...Prober) {
	for _, dp := range probers {
		resp := pResp.HealthResponse()
		payload, hasKey := resp[dp.ServiceID()]
		asserter.True(hasKey)
		dpp := dp.(*DummyPinger)
		if dpp.err == nil {
			asserter.NotContains(payload, "NOT OK")
			continue
		}

		asserter.Contains(payload, "NOT OK")
		for _, afstatus := range dpp.AffectsStatuses() {
			switch afstatus {
			case proberesponder.StatusStartup:
				asserter.True(pResp.NotStarted(), "%s: %s", dpp.serviceID, payload)
			case proberesponder.StatusReady:
				asserter.True(pResp.NotReady(), "%s: %s", dpp.serviceID, payload)
			case proberesponder.StatusLive:
				asserter.True(pResp.NotLive(), "%s: %s", dpp.serviceID, payload)
			}
		}
	}
}

func TestStart(tt *testing.T) {
	const (
		delay               = time.Millisecond * 750
		waitBeforeAssertion = time.Second
	)

	tt.Run("no probers", func(t *testing.T) {
		asserter := assert.New(t)
		pResp := newProbeRespWithAllOK()

		Start(delay, pResp)
		time.Sleep(waitBeforeAssertion)

		asserter.False(pResp.NotStarted())
		asserter.False(pResp.NotReady())
		asserter.False(pResp.NotLive())
	})

	tt.Run("service affects no statuses", func(t *testing.T) {
		asserter := assert.New(t)
		pResp := newProbeRespWithAllOK()
		probers := []Prober{
			&DummyPinger{
				serviceID:      "service_affects_none",
				affectedStatus: nil,
				err:            nil,
			}}
		stopper := Start(delay, pResp, probers...)
		defer stopper.Stop()

		// wait for probe to complete at least 1 cycle
		time.Sleep(waitBeforeAssertion)
		assertStatuses(asserter, pResp, probers...)
	})
	tt.Run("service affects startup", func(t *testing.T) {
		asserter := assert.New(t)
		pResp := newProbeRespWithAllOK()
		probers := []Prober{
			&DummyPinger{
				serviceID:      "service_affects_start_no_error",
				affectedStatus: []proberesponder.Statuskey{proberesponder.StatusStartup},
				err:            nil,
			},
			&DummyPinger{
				serviceID:      "service_affects_start_has_error",
				affectedStatus: []proberesponder.Statuskey{proberesponder.StatusStartup},
				err:            errors.New("service down"),
			}}
		stopper := Start(delay, pResp, probers...)
		defer stopper.Stop()

		// wait for probe to complete at least 1 cycle
		time.Sleep(waitBeforeAssertion)
		assertStatuses(asserter, pResp, probers...)
	})

	tt.Run("service affects ready", func(t *testing.T) {
		asserter := assert.New(t)
		pResp := newProbeRespWithAllOK()
		probers := []Prober{
			&DummyPinger{
				serviceID:      "service_affects_ready_no_error",
				affectedStatus: []proberesponder.Statuskey{proberesponder.StatusReady},
				err:            nil,
			},
			&DummyPinger{
				serviceID:      "service_affects_ready_has_error",
				affectedStatus: []proberesponder.Statuskey{proberesponder.StatusReady},
				err:            errors.New("service down"),
			}}
		stopper := Start(delay, pResp, probers...)
		defer stopper.Stop()

		// wait for probe to complete at least 1 cycle
		time.Sleep(waitBeforeAssertion)
		assertStatuses(asserter, pResp, probers...)
	})

	tt.Run("service affects live", func(t *testing.T) {
		asserter := assert.New(t)
		pResp := newProbeRespWithAllOK()
		probers := []Prober{
			&DummyPinger{
				serviceID:      "service_affects_live_no_error",
				affectedStatus: []proberesponder.Statuskey{proberesponder.StatusLive},
				err:            nil,
			},
			&DummyPinger{
				serviceID:      "service_affects_live_has_error",
				affectedStatus: []proberesponder.Statuskey{proberesponder.StatusLive},
				err:            errors.New("service down"),
			}}
		stopper := Start(delay, pResp, probers...)
		defer stopper.Stop()

		// wait for probe to complete at least 1 cycle
		time.Sleep(waitBeforeAssertion)
		assertStatuses(asserter, pResp, probers...)
	})
}

func TestProber(tt *testing.T) {
	tt.Run("basic checks", func(t *testing.T) {
		asserter := assert.New(t)
		expectedServiceID := "service_1"
		expectedStatuses := []proberesponder.Statuskey{
			proberesponder.StatusLive,
		}
		pb := Probe{
			ID:               "service_1",
			AffectedStatuses: expectedStatuses,
			Checker: CheckerFunc(func(ctx context.Context) error {
				return nil
			}),
		}
		asserter.NoError(pb.Check(context.Background()))
		asserter.Equal(expectedServiceID, pb.ServiceID())
		asserter.Equal(expectedStatuses, pb.AffectsStatuses())
	})
	tt.Run("without checker", func(t *testing.T) {
		asserter := assert.New(t)
		expectedServiceID := "service_1"
		expectedStatuses := []proberesponder.Statuskey{
			proberesponder.StatusLive,
		}
		pb := Probe{
			ID:               "service_1",
			AffectedStatuses: expectedStatuses,
			Checker:          nil,
		}
		asserter.NoError(pb.Check(context.Background()))
		asserter.Equal(expectedServiceID, pb.ServiceID())
		asserter.Equal(expectedStatuses, pb.AffectsStatuses())
	})

	tt.Run("checker returns error", func(t *testing.T) {
		asserter := assert.New(t)
		expectedServiceID := "service_1"
		expectedStatuses := []proberesponder.Statuskey{
			proberesponder.StatusLive,
		}
		pb := Probe{
			ID:               "service_1",
			AffectedStatuses: expectedStatuses,
			Checker: CheckerFunc(func(ctx context.Context) error {
				return errors.New("has error")
			}),
		}
		asserter.Error(pb.Check(context.Background()))
		asserter.Equal(expectedServiceID, pb.ServiceID())
		asserter.Equal(expectedStatuses, pb.AffectsStatuses())
	})
}
