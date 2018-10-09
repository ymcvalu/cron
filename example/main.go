package main

import (
	"cron"
	"fmt"
	"time"
)

func main() {
	c := cron.NewCron()
	c.Start()
	for i := 0; i < 5; i++ {
		i := i
		c.AddJob(&cron.CronJob{
			Id:        fmt.Sprintf("job-%d", i),
			Scheduler: cron.FromInterval(int64(i + 1)),
			Job: cron.FunJob(func() {
				fmt.Printf("job-%d run at %s\n", i, time.Now())
				if i == 4 {
					panic("test panic")
				}
			}),
		})
	}

	time.Sleep(time.Second * 10)
	c.RemoveJob(&cron.CronJob{Id: "job-4"})
	time.Sleep(time.Second * 10)
	c.Stop()
}
