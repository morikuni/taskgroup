package taskgroup

type GroupOption func(o *groupOption)

func onFirstError(f func()) GroupOption {
	return func(o *groupOption) {
		o.onFirstError = func(_ error) {
			f()
		}
	}
}

type groupOption struct {
	onFirstError func(error)
}

func (co *groupOption) apply(os []GroupOption) {
	for _, o := range os {
		o(co)
	}
}
