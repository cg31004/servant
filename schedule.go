package servant

import (
	"time"

	"github.com/robfig/cron"
)

type Schedule = cron.Schedule
type SpecSchedule = cron.SpecSchedule
type ConstantDelaySchedule = cron.ConstantDelaySchedule

func Every(duration time.Duration) ConstantDelaySchedule {
	return cron.Every(duration)
}
