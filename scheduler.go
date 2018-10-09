package cron

import (
	"time"
)

type Scheduler interface {
	Next(time.Time) time.Time
}

type FuncScheduler func(time.Time) time.Time

func (f FuncScheduler) Next(pre time.Time) time.Time {
	return f(pre)
}

type IntervalScheduler struct {
	interval time.Duration
}

func FromInterval(secs int64) Scheduler {
	return IntervalScheduler{time.Duration(secs * int64(time.Second))}
}

func (s IntervalScheduler) Next(pre time.Time) time.Time {
	if pre.IsZero() {
		pre = time.Now()
	}
	next := pre.Add(s.interval)
	if next.Before(time.Now()) {
		next = time.Now().Add(s.interval)
	}
	return next
}
