package servant

type FuncJob func(ctx *Context)
type handlerChain []FuncJob

func NewCustomJobFunc(c *Cron, f func(ctx *Context), profile *Profile) *CustomJob {
	job := &middlewareJob{handlers: append(c.handlers, f)}
	return NewCustomJob(c, Job(job), profile)
}

type middlewareJob struct {
	ctx      *Context
	handlers handlerChain
}

func (f *middlewareJob) Run() {
	f.ctx = &Context{}
	f.ctx.reset()
	f.ctx.handlerChain = f.handlers
	f.ctx.Next()
}
