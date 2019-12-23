package main

import (
	"time"
)

type TimeStamp struct {
	ReqTime  time.Time     `json:"request"`
	RespTime time.Time     `json:"response"`
	Elapsed  time.Duration `json:"elapsed"`
}

func NewTimestamp() (t TimeStamp) {
	return TimeStamp{
		ReqTime: time.Now(),
	}
}

func Timestamp(t time.Time) TimeStamp {
	return TimeStamp{
		ReqTime: t,
	}
}

// SetResponseTime sets according to the argumented passed to us, as a side
// effect the Elapsed property will also be set.
func (ts *TimeStamp) SetResponseTime(now time.Time) time.Duration {
	ts.RespTime = time.Now()
	ts.Elapsed = ts.RespTime.Sub(ts.ReqTime)
	return ts.Elapsed
}
