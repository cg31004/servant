package servant

import (
	"sync"
)

type Job interface {
	Run()
}

type FuncJob func(ctx *CronContext)

func (f FuncJob) Run() {
	ctx := &CronContext{}
	f(ctx)
}

func NewCustomJobFunc(c *Cron, f func(ctx *CronContext), profile *Profile) *CustomJob {
	return NewCustomJob(c, FuncJob(f), profile)
}

func NewCustomJob(c *Cron, job Job, profile *Profile) *CustomJob {
	return &CustomJob{
		job:     job,
		profile: profile,
		wg:      c.wg,
		mx:      c.mx,
	}
}

type CustomJob struct {
	job     Job
	profile *Profile
	mx      *sync.Mutex
	wg      *sync.WaitGroup
}

func (cj *CustomJob) Run() {
	if !cj.canRun() {
		return
	}

	defer cj.finish()

	cj.wg.Add(1)
	cj.job.Run()
	cj.wg.Done()
}

func (cj *CustomJob) canRun() bool {
	cj.mx.Lock()
	defer cj.mx.Unlock()

	if !cj.profile.isRunning() {
		return false
	}

	// 不可重覆執行，且執行中
	if !cj.profile.isOverlapping() && cj.profile.isProcessing() {
		return false
	}

	cj.profile.processingAdd()

	return true
}

func (cj *CustomJob) finish() {
	cj.mx.Lock()
	cj.profile.processingDone()
	cj.mx.Unlock()
}
