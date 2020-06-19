package taskgroup

type config struct {
	runner      Runner
	interceptor Interceptor
	fold        func(acc, err error) error
}

func evaluateOptions(opts []Option) *config {
	c := &config{
		runner:      DefaultRunner,
		interceptor: nil,
		fold: func(acc, err error) error {
			if acc == nil {
				return err
			}
			return acc
		},
	}

	for _, o := range opts {
		o.apply(c)
	}

	return c
}

type Option interface {
	apply(*config)
}
