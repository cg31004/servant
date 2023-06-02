package servant

import (
	"context"
	"sync"
)

type Job interface {
	Run()
}

type FuncJob func(ctx *Context)

func (f FuncJob) Run() {
	ctx := &Context{Context: context.Background()}
	f(ctx)
}

func NewCustomJobFunc(c *Cron, f func(ctx *Context), profile *Profile) *CustomJob {
	return NewCustomJob(c, FuncJob(f), f, profile)
}

func NewCustomJob(c *Cron, job Job, f func(ctx *Context), profile *Profile) *CustomJob {
	handlers := append(c.handlers, f)
	return &CustomJob{
		job:     job,
		profile: profile,
		wg:      c.wg,
		mx:      c.mx,
		ctx: &Context{
			Context:      context.Background(),
			handlerChain: handlers,
		},
	}
}

type CustomJob struct {
	job      Job
	profile  *Profile
	handlers HandlersChain
	mx       *sync.Mutex
	wg       *sync.WaitGroup
	ctx      *Context
}

func (cj *CustomJob) Run() {
	if !cj.canRun() {
		return
	}

	defer cj.finish()

	cj.wg.Add(1)
	cj.ctx.Next()
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
