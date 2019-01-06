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
	rm      chan string
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
		rm:   make(chan string),
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
	for _, job := range c.h.array {
		job.pre = time.Now()
		job.next = job.Scheduler.Next(job.pre)
	}
	c.h.BuildHeap()

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
				if cj == nil || !cj.next.Before(now) {
					break
				} else {
					go c.runJobWithRecover(cj)
					cj.pre = cj.next
					cj.next = cj.Next(cj.pre)
					c.h.RestoreDown(0)
				}
			}

		case job := <-c.add:
			timer.Stop()
			idx := c.h.Walk(func(cj *CronJob) bool {
				if cj.Id == job.Id {
					return true
				}
				return false
			})
			if idx >= 0 {
				c.h.Remove(idx)
			}
			job.pre = time.Now()
			job.next = job.Next(job.pre)
			c.h.Push(job)

		case id := <-c.rm:
			timer.Stop()
			idx := c.h.Walk(func(cj *CronJob) bool {
				if cj.Id == id {
					return true
				}
				return false
			})
			if idx >= 0 {
				c.h.Remove(idx)
			}

		case <-c.stop:
			timer.Stop()
			return
		}
	}
}

func (c *Cron) runJobWithRecover(cj *CronJob) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("panic when executing cron job,stack info:\n%s", bytes2String(buf))
		}
	}()
	cj.Run()
}

func (c *Cron) AddJob(job *CronJob) {
	if !c.running {
		for i, j := range c.h.array {
			if j.Id == job.Id {
				c.h.array[i] = job
				return
			}
		}
		c.h.array = append(c.h.array, job)
		return
	}
	c.add <- job
}

func (c *Cron) RemoveJob(id string) {
	if !c.running {
		for i, j := range c.h.array {
			if j.Id == id {
				copy(c.h.array[i:], c.h.array[i+1:])
				c.h.array = c.h.array[:len(c.h.array)-1]
				return
			}
		}
		return
	}
	c.rm <- id
}
