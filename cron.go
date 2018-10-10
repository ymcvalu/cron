package cron

import (
	"log"
	"runtime"
	"time"
)

type Job interface {
	Run()
}

type FunJob func()

func (f FunJob) Run() {
	f()
}

type CronJob struct {
	Id string
	Job
	Scheduler
	next time.Time //下一次执行时间
	pre  time.Time //上一次执行时间
}

type Cron struct {
	h       *KaryHeap
	stop    chan struct{}
	add     chan *CronJob
	rm      chan *CronJob
	running bool
}

func NewCron() *Cron {
	cron := &Cron{
		h: NewKaryHeap(4, func(ci, cj *CronJob) bool {
			if ci.next.IsZero() {
				return false
			}
			if cj.next.IsZero() {
				return true
			}
			return ci.next.Before(cj.next)
		}),
		stop: make(chan struct{}),
		add:  make(chan *CronJob),
		rm:   make(chan *CronJob),
	}
	return cron
}

func (c *Cron) Start() {
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

func (c *Cron) Stop() {
	if !c.running {
		return
	}
	c.stop <- struct{}{}
}

func (c *Cron) Run() {
	if c.running {
		return
	}
	c.running = true
	c.run()
}

func (c *Cron) run() {
	for {
		var timer *time.Timer
		if c.h.Len() == 0 {
			timer = time.NewTimer(time.Hour * 1)
		} else {
			d := c.h.Peek(0).next.Sub(time.Now())
			timer = time.NewTimer(d)
		}
		select {
		case now := <-timer.C:
			for {
				cj := c.h.Peek(0)
				if cj == nil {
					break
				}

				if cj.next.Before(now) {
					go c.RunJobWithRecover(cj)
					cj.pre = cj.next
					cj.next = cj.Next(cj.pre)
					c.h.RestoreDown(0)
				} else {
					break
				}
			}
		case job := <-c.add:
			timer.Stop()
			job.next = job.Next(job.pre)
			c.h.Push(job)
		case job := <-c.rm:
			timer.Stop()
			idx := c.h.Walk(func(cj *CronJob) bool {

				if cj.Id == job.Id {
					return true
				}
				return false
			})
			c.h.Remove(idx)
		case <-c.stop:
			timer.Stop()
			return
		}
	}
}

func (c *Cron) RunJobWithRecover(cj *CronJob) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("panic when execute cron job,stack info:\n%s", bytes2String(buf))
		}
	}()
	cj.Run()
}

func (c *Cron) AddJob(job *CronJob) {
	go func() {
		c.add <- job
	}()
}

func (c *Cron) RemoveJob(job *CronJob) {
	go func() {
		c.rm <- job
	}()
}
