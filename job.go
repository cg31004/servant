package servant

import (
	"context"
	"sync"
)

type Job interface {
	Run()
}

func NewCustomJob(c *Cron, job Job, profile *Profile) *CustomJob {
	return &CustomJob{
		job:      job,
		profile:  profile,
		wg:       c.wg,
		mx:       c.mx,
		handlers: handlers,
		//ctx: &Context{
		//	Context:      context.Background(),
		//	handlerChain: handlers,
		//},
	}
}

type CustomJob struct {
	job          Job
	profile      *Profile
	mx           *sync.Mutex
	wg           *sync.WaitGroup
	handlerChain []handlerChain
}

func (cj *CustomJob) Run() {
	if !cj.canRun() {
		return
	}
	defer cj.finish()

	ctx := &Context{Context: context.Background(), handlerChain: cj.handlers}
	cj.wg.Add(1)
	ctx.Next()
	//cj.job.Run()
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
