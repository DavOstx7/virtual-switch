package toggle

import (
	"context"
	"fmt"
	"sync"
)

type AggregateToggler struct {
	CommonToggler
	wg *sync.WaitGroup
}

func (t *AggregateToggler) Setup(startFuncs []func(context.Context) error) {
	newStartFunc := func(ctx context.Context) error {
		t.wg = new(sync.WaitGroup)
		for _, startFunc := range startFuncs {
			startFunc := startFunc

			t.wg.Add(1)
			go func() {
				if err := startFunc(ctx); err != nil {
					fmt.Println(err)
				}
				t.wg.Done()
			}()
		}
		return nil
	}
	stopFunc := func() error {
		t.wg.Wait()
		t.wg = nil
		return nil
	}

	t.CommonToggler.Setup(newStartFunc, stopFunc)
}
