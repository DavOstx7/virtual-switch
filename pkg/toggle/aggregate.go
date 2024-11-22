package toggle

import (
	"context"
	"sync"
)

type AggregateToggler struct {
	CommonToggler
	Wg *sync.WaitGroup
}

func (t *AggregateToggler) Setup(startFunc func(context.Context) error) {
	newStartFunc := func(ctx context.Context) error {
		t.Wg = new(sync.WaitGroup)
		if err := startFunc(ctx); err != nil {
			t.Wg = nil
			return err
		}
		return nil
	}
	stopFunc := func() error {
		t.Wg.Wait()
		t.Wg = nil
		return nil
	}

	t.CommonToggler.Setup(newStartFunc, stopFunc)
}
