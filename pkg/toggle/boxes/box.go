package boxes

import "context"

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

func (b *ToggleBox) Setup(startFunc StartFunc, stopFunc StopFunc) {
	b.startFunc = startFunc
	b.stopFunc = stopFunc
}

func (b *ToggleBox) IsOn() bool {
	return b.cancel != nil
}

func (b *ToggleBox) On(ctx context.Context) error {
	if b.IsOn() {
		return nil
	}

	newCtx, newCancel := context.WithCancel(ctx)
	if err := b.startFunc(newCtx); err != nil {
		newCancel()
		return err
	}

	b.cancel = newCancel
	return nil
}

func (b *ToggleBox) Off() error {
	if !b.IsOn() {
		return nil
	}

	if err := b.stopFunc(); err != nil {
		return err
	}

	b.cancel = nil
	return nil
}

func (b *ToggleBox) Cancel() {
	b.cancel()
}
