package main

import (
	"log"
	"runtime"
	"sync/atomic"
	"time"
)

type RuntimeStats struct {
	Started   time.Time
	status1xx uint64
	status2xx uint64
	status3xx uint64
	status4xx uint64
	status5xx uint64
}

func (rs *RuntimeStats) StartDate() string {
	return rs.Started.Format("2006-01-02 15:04")
}

func (rs *RuntimeStats) Version() string {
	return runtime.Version()
}

func (rs *RuntimeStats) NumGoroutine() int {
	return runtime.NumGoroutine()
}

func (rs *RuntimeStats) inc(addr *uint64) {
	atomic.AddUint64(addr, 1)
}

func (rs *RuntimeStats) IncStatus(code int) {
	codeClass := code / 100
	switch codeClass {
	case 1:
		rs.Inc1xx()
		return
	case 2:
		rs.Inc2xx()
		return
	case 3:
		rs.Inc3xx()
		return
	case 4:
		rs.Inc4xx()
		return
	case 5:
		rs.Inc5xx()
		return
	}

	log.Fatalf("Unexpected response code %v.", code)
}

func (rs *RuntimeStats) Inc1xx() { rs.inc(&rs.status1xx) }
func (rs *RuntimeStats) Inc2xx() { rs.inc(&rs.status2xx) }
func (rs *RuntimeStats) Inc3xx() { rs.inc(&rs.status3xx) }
func (rs *RuntimeStats) Inc4xx() { rs.inc(&rs.status4xx) }
func (rs *RuntimeStats) Inc5xx() { rs.inc(&rs.status5xx) }

func (rs *RuntimeStats) Status1xx() uint64 { return atomic.LoadUint64(&rs.status1xx) }
func (rs *RuntimeStats) Status2xx() uint64 { return atomic.LoadUint64(&rs.status2xx) }
func (rs *RuntimeStats) Status3xx() uint64 { return atomic.LoadUint64(&rs.status3xx) }
func (rs *RuntimeStats) Status4xx() uint64 { return atomic.LoadUint64(&rs.status4xx) }
func (rs *RuntimeStats) Status5xx() uint64 { return atomic.LoadUint64(&rs.status5xx) }

func NewStats() *RuntimeStats {
	return &RuntimeStats{
		Started: time.Now(),
	}
}
