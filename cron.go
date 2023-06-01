package cronjob

import (
	"sync"

	"github.com/robfig/cron/v3"
)

func New() *Cron {
	return &Cron{
		pkgCron: cron.New(cron.WithSeconds()),
		wg:      &sync.WaitGroup{},
		mx:      &sync.Mutex{},
	}
}

type Cron struct {
	pkgCron *cron.Cron
	Opt     CronOpt
	mx      *sync.Mutex
	wg      *sync.WaitGroup
	running bool
}

// AddFunc 將Func 加入排程器，並依據字串規則執行任務。
func (c *Cron) AddFunc(spec string, cmd func(), opt ...FuncCronOpt) (*Profile, error) {
	job := NewCustomJobFunc(c, cmd, parseCronOpt(opt...))
	return c.addJob(spec, job)
}

// AddJob 將Job類型物件加入排程器，並依據字串規則執行任務。
func (c *Cron) AddJob(spec string, cmd Job, opt ...FuncCronOpt) (*Profile, error) {
	job := NewCustomJob(c, cmd, parseCronOpt(opt...))
	return c.addJob(spec, job)
}

func (c *Cron) addJob(spec string, job *CustomJob) (*Profile, error) {
	entryID, err := c.pkgCron.AddJob(spec, cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)).Then(job))
	if err != nil {
		return nil, newCronError(ErrRegistered, ErrorWithErrorMessage(err))
	}

	job.profile.setEntryID(entryID)

	return job.profile, nil
}

// AddScheduleFunc 將Func 加入排程器，並依據Schedule物件規則執行任務。
func (c *Cron) AddScheduleFunc(schedule Schedule, cmd func(), opt ...FuncCronOpt) (*Profile, error) {
	job := NewCustomJobFunc(c, cmd, parseCronOpt(opt...))
	return c.addScheduleJob(schedule, job)
}

// AddScheduleJob 將Job類型物件加入排程器，並依據Schedule物件規則執行任務。
func (c *Cron) AddScheduleJob(schedule Schedule, cmd Job, opt ...FuncCronOpt) (*Profile, error) {
	job := NewCustomJob(c, cmd, parseCronOpt(opt...))
	return c.addScheduleJob(schedule, job)
}

func (c *Cron) addScheduleJob(schedule Schedule, job *CustomJob) (*Profile, error) {
	if schedule == nil {
		return nil, newCronError(ErrScheduleIsNil)
	}
	entryID := c.pkgCron.Schedule(schedule, cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)).Then(job))
	job.profile.setEntryID(entryID)
	return job.profile, nil
}

// Remove 將任務從排程器中刪除
func (c *Cron) Remove(profile *Profile) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if profile.GetJobStatus() == JobStatusRemove {
		return newCronError(ErrJobAlreadyRemove)
	}

	profile.remove()
	c.pkgCron.Remove(profile.GetEntryID())

	return nil
}

// Start 啟動全部排程
func (c *Cron) Start() {
	c.mx.Lock()
	if !c.running {
		c.pkgCron.Start()
		c.running = true
	}
	c.mx.Unlock()
}

// Stop 暫停全部排程，但是不能關閉已經在執行中的。
func (c *Cron) Stop() {
	c.mx.Lock()
	if c.running {
		c.pkgCron.Stop()
		c.running = false
	}
	c.mx.Unlock()

	c.wg.Wait()
}
