package toggle

import (
	"context"
	"sync"
)

type CommonToggler struct {
	Mu        sync.Mutex
	startFunc func(context.Context) error
	stopFunc  func() error
	cancel    context.CancelFunc
}

func (t *CommonToggler) IsOn() bool {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t._isOn()
}

func (t *CommonToggler) On(ctx context.Context) error {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	if t._isOn() {
		return nil
	}

	newCtx, newCancel := context.WithCancel(ctx)
	if err := t.startFunc(newCtx); err != nil {
		newCancel()
		return err
	}

	t.cancel = newCancel
	return nil
}

func (t *CommonToggler) Off() error {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	if !t._isOn() {
		return nil
	}

	t.cancel()
	if err := t.stopFunc(); err != nil {
		return err
	}

	t.cancel = nil
	return nil
}

func (t *CommonToggler) Setup(startFunc func(context.Context) error, stopFunc func() error) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.startFunc = startFunc
	t.stopFunc = stopFunc
}

func (t *CommonToggler) _isOn() bool {
	return t.cancel != nil
}
