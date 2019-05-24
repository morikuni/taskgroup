package taskgroup

type CallOption func(o *callOption)

type callOption struct {
	middlewares []middleware
}

func (co *callOption) apply(os []CallOption) {
	for _, o := range os {
		o(co)
	}
}

func Wait(ss ...Signal) CallOption {
	return func(o *callOption) {
		o.middlewares = append(o.middlewares, wait(ss))
	}
}
