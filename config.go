package taskgroup

type config struct {
	foldFunc FoldFunc
}

func newConfig(opts []Option) *config {
	conf := &config{
		foldFunc: AppendError,
	}
	for _, opt := range opts {
		opt.Apply(conf)
	}
	return conf
}

type Option interface {
	Apply(*config)
}

type FoldFunc func(acc, err error) error

func (f FoldFunc) Apply(conf *config) {
	conf.foldFunc = f
}
