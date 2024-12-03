package proberesponder

import (
	"fmt"
	"sync"
	"time"
)

type Statuskey string

func (sk Statuskey) String() string {
	return string(sk)
}

const (
	StatusStartup Statuskey = "startup"
	StatusReady   Statuskey = "ready"
	StatusLive    Statuskey = "live"
)

type healthstatus string

func (hs healthstatus) String() string {
	return string(hs)
}

func (hs healthstatus) Equal(s healthstatus) bool {
	return s == hs
}

func IsHealthOK[T ~string](s T) bool {
	return HealthOK.Equal(healthstatus(s))
}

const (
	HealthOK    healthstatus = "OK"
	HealthNotOK healthstatus = "NOT OK"
)

type StatusChangeListener func(status Statuskey, value bool)

// ProbeStatuses are maintained primarily for K8s probe responses. Though it can be used
// for any prober.
type ProbeResponder struct {
	notReady       bool
	notLive        bool
	notStarted     bool
	locker         *sync.Mutex
	msgPayload     map[string]string
	changeListener StatusChangeListener
}

func (pr *ProbeResponder) AppendHealthResponse(key, value string) {
	if pr == nil {
		return
	}
	pr.locker.Lock()
	defer pr.locker.Unlock()

	pr.appendHealthRespWithoutLock(key, value)
}

func (pr *ProbeResponder) appendHealthRespWithoutLock(key, value string) {
	pr.msgPayload[key] = value
}

func (pr *ProbeResponder) HealthResponse() map[string]string {
	if pr == nil {
		return nil
	}
	pr.locker.Lock()
	defer pr.locker.Unlock()

	copied := map[string]string{}
	for k, v := range pr.msgPayload {
		copied[k] = v
	}

	return copied
}

func (pr *ProbeResponder) onChange(status Statuskey, value bool) {
	hs := HealthOK
	if value {
		hs = HealthNotOK
	}

	pr.appendHealthRespWithoutLock(
		"probe->"+status.String(),
		fmt.Sprintf("%s: %s", hs, time.Now().Format(time.RFC3339)),
	)

	if pr.changeListener == nil {
		return
	}

	pr.changeListener(status, value)
}

func (pr *ProbeResponder) SetNotReady(b bool) {
	if pr == nil {
		return
	}

	pr.locker.Lock()
	defer pr.locker.Unlock()

	pr.notReady = b
	pr.onChange(StatusReady, b)
}

func (pr *ProbeResponder) SetNotLive(b bool) {
	if pr == nil {
		return
	}

	pr.locker.Lock()
	defer pr.locker.Unlock()

	pr.notLive = b
	pr.onChange(StatusLive, b)
}

func (pr *ProbeResponder) SetNotStarted(b bool) {
	if pr == nil {
		return
	}

	pr.locker.Lock()
	defer pr.locker.Unlock()

	pr.notStarted = b
	pr.onChange(StatusStartup, b)
}

// SetListener is used to set a callback function which will be invoked every time
// any of the statuses change (e.g. liveness)
func (pr *ProbeResponder) SetListener(l StatusChangeListener) {
	if pr == nil {
		return
	}

	pr.locker.Lock()
	defer pr.locker.Unlock()

	pr.changeListener = l
}

func (pr *ProbeResponder) NotReady() bool {
	return pr != nil && pr.notReady
}

func (pr *ProbeResponder) NotLive() bool {
	return pr != nil && pr.notLive
}

func (pr *ProbeResponder) NotStarted() bool {
	return pr != nil && pr.notStarted
}

func New() *ProbeResponder {
	pRes := &ProbeResponder{
		locker:     &sync.Mutex{},
		msgPayload: map[string]string{},
	}

	pRes.SetNotLive(true)
	pRes.SetNotReady(true)
	pRes.SetNotStarted(true)

	return pRes
}
