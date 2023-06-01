package cronjob

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Schedule = cron.Schedule
type SpecSchedule = cron.SpecSchedule
type ConstantDelaySchedule = cron.ConstantDelaySchedule

func Every(duration time.Duration) ConstantDelaySchedule {
	return cron.Every(duration)
}
