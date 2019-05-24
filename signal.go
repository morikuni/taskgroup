package taskgroup

type Signal interface {
	Ready() <-chan struct{}
	done() <-chan struct{}
}

type signal struct {
	readyC <-chan struct{}
	doneC  <-chan struct{}
}

func (s *signal) Ready() <-chan struct{} {
	return s.readyC
}

func (s *signal) done() <-chan struct{} {
	return s.doneC
}
