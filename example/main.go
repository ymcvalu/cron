package main

import (
	"cron"
	"fmt"
	"time"
)

func main() {
	c := cron.NewCron()
	c.AddJob(&cron.CronJob{
		Id:        "id-xxx",
		Scheduler: cron.Every(time.Second * 3),
		Job: cron.FunJob(func() {
			fmt.Println("cron job running")
		}),
	})
	c.Run()
}
