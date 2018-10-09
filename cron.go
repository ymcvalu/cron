package cron

import (
	"cron/utils"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
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
	running int32
}

func NewCron() *Cron {
	cron := &Cron{
		h: NewKaryHeap(func(i, j interface{}) bool {
			ci := i.(*CronJob)
			cj := j.(*CronJob)
			if ci.next.IsZero() {
				return false
			}
			if cj.next.IsZero() {
				return true
			}
			return ci.next.Before(cj.next)
		}, Kary(4), Locker(&sync.Mutex{})),
		stop:    make(chan struct{}),
		add:     make(chan *CronJob),
		rm:      make(chan *CronJob),
		running: 0,
	}
	return cron
}

func (c *Cron) Start() {
	if atomic.LoadInt32(&c.running) > 0 {
		return
	}
	atomic.StoreInt32(&c.running, 1)
	go c.run()
}

func (c *Cron) Stop() {
	c.stop <- struct{}{}
}

func (c *Cron) Run() {
	if atomic.LoadInt32(&c.running) > 0 {
		return
	}
	c.run()
}
func (c *Cron) run() {
	for {
		var timer *time.Timer
		if c.h.Len() == 0 {
			timer = time.NewTimer(time.Hour * 1)
		} else {
			d := c.h.Peek(0).(*CronJob).next.Sub(time.Now())
			timer = time.NewTimer(d)
		}
		select {
		case now := <-timer.C:
			for {
				j := c.h.Peek(0)
				if j == nil {
					break
				}
				cj := j.(*CronJob)
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
			c.h.WalkRm(func(v interface{}) (bool, bool) {
				cj := v.(*CronJob)
				if cj.Id == job.Id {
					return true, true
				} else {
					return false, false
				}
			})
		case <-c.stop:
			timer.Stop()
			atomic.StoreInt32(&c.running, 0)
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
			log.Printf("panic when execute cron job,stack info:\n%s", utils.Bytes2String(buf))
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
