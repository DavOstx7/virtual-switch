package toggle

import (
	"context"
	"sync"
)

type Toggler interface {
	On(ctx context.Context) error
	Off() error
}

func On(ctx context.Context, togglers ...Toggler) <-chan error {
	errChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(len(togglers))

	for _, toggler := range togglers {
		toggler := toggler
		go func() {
			if err := toggler.On(ctx); err != nil {
				errChan <- err
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return errChan
}

func Off(togglers ...Toggler) <-chan error {
	errChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(len(togglers))

	for _, toggler := range togglers {
		toggler := toggler
		go func() {
			if err := toggler.Off(); err != nil {
				errChan <- err
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return errChan
}
