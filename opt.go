package cronjob

type FuncCronOpt func(o *Profile)

func (f FuncCronOpt) Set(o *Profile) {
	f(o)
}

type CronOpt func()

// Overlapping 在同一時間可否重覆執行
func (o CronOpt) Overlapping(b bool) FuncCronOpt {
	return func(o *Profile) {
		o.Overlapping(b)
	}
}

// Running 運行排程
func (o CronOpt) Running(b bool) FuncCronOpt {
	return func(o *Profile) {
		if b {
			o.Start()
		} else {
			o.Stop()
		}
	}
}

// ------------------------------

func parseCronOpt(opts ...FuncCronOpt) *Profile {
	result := newProfile()

	for _, o := range opts {
		o.Set(result)
	}

	return result
}
