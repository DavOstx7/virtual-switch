package boxes

import (
	"context"
	"sync"
)

type StartFunc func(context.Context) error
type StopFunc func() error

type ToggleBox struct {
	startFunc StartFunc
	stopFunc  StopFunc
	cancel    context.CancelFunc
}

func NewToggleBox() *ToggleBox {
	return new(ToggleBox)
}

func (tb *ToggleBox) Setup(startFunc StartFunc, stopFunc StopFunc) {
	tb.startFunc = startFunc
	tb.stopFunc = stopFunc
}

func (tb *ToggleBox) IsOn() bool {
	return tb.cancel != nil
}

func (tb *ToggleBox) On(ctx context.Context) error {
	if tb.IsOn() {
		return nil
	}

	newCtx, newCancel := context.WithCancel(ctx)
	if err := tb.startFunc(newCtx); err != nil {
		newCancel()
		return err
	}

	tb.cancel = newCancel
	return nil
}

func (tb *ToggleBox) Off() error {
	if !tb.IsOn() {
		return nil
	}

	if err := tb.stopFunc(); err != nil {
		return err
	}

	tb.cancel = nil
	return nil
}

func (tb *ToggleBox) Cancel() {
	tb.cancel()
}

/* ------------------------------------------------------------------------------ */

type SafeToggleBox struct {
	ToggleBox
	mu sync.Mutex
}

func (tb *SafeToggleBox) IsOn() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.ToggleBox.IsOn()
}

func (tb *SafeToggleBox) On(ctx context.Context) error {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.ToggleBox.On(ctx)
}

func (tb *SafeToggleBox) Off() error {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.ToggleBox.Off()
}
