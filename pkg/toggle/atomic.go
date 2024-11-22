package toggle

import "context"


type AtomicToggler struct {
	CommonToggler
	Done chan bool
}

func (t *AtomicToggler) Setup(startFunc func(context.Context) error) {
	newStartFunc := func(ctx context.Context) error {
		t.Done = make(chan bool)
		if err := startFunc(ctx); err != nil {
			close(t.Done)
			t.Done = nil
			return err
		}
		return nil
	}
	stopFunc := func() error {
		if _, ok := <-t.Done; ok {
			close(t.Done)
		}
		t.Done = nil
		return nil
	}

	t.CommonToggler.Setup(newStartFunc, stopFunc)
}
