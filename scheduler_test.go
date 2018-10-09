package cron

import (
	"time"

	"testing"
)

func TestIntervalSche(t *testing.T) {
	sche := FromInterval(6)
	for i := 0; i < 10; i++ {
		now := time.Now()
		next := sche.Next(time.Unix(0, 0))
		if next.Sub(now) != sche.(IntervalScheduler).interval {
			t.Error(now, next)
		}
		time.Sleep(time.Second * 1)
	}
}
