package toggle

import (
	"context"
	"fmt"
)

type AtomicToggler struct {
	CommonToggler
	done chan bool
}

func (t *AtomicToggler) Setup(startFunc func(context.Context) error) {
	newStartFunc := func(ctx context.Context) error {
		t.done = make(chan bool)
		go func() {
			if err := startFunc(ctx); err != nil {
				fmt.Println(err)
			}
			close(t.done)
		} ()
		return nil
	}
	stopFunc := func() error {
		<-t.done
		t.done = nil
		return nil
	}

	t.CommonToggler.Setup(newStartFunc, stopFunc)
}
