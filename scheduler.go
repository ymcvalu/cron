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

func Every(duration time.Duration) Scheduler {
	return FuncScheduler(func(pre time.Time) time.Time {
		if pre.IsZero() {
			pre = time.Now()
		}
		return pre.Add(duration)
	})
}
