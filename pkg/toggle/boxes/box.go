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

func (b *ToggleBox) SetStarter(startFunc StartFunc) {
	b.startFunc = startFunc
}

func (b *ToggleBox) SetStopper(stopFunc StopFunc) {
	b.stopFunc = stopFunc
}

func (b *ToggleBox) Cancel() {
	b.cancel()
}

func (b *ToggleBox) IsOn() bool {
	return b.cancel != nil
}

func (b *ToggleBox) On(ctx context.Context) error {
	if b.IsOn() {
		return nil
	}

	newCtx, newCancel := context.WithCancel(ctx)
	if b.startFunc != nil {
		if err := b.startFunc(newCtx); err != nil {
			newCancel()
			return err
		}
	}

	b.cancel = newCancel
	return nil
}

func (b *ToggleBox) Off() error {
	if !b.IsOn() {
		return nil
	}

	if b.stopFunc != nil {
		if err := b.stopFunc(); err != nil {
			return err
		}
	}

	b.cancel = nil
	return nil
}
