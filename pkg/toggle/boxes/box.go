package boxes

import "context"

type StartFunc func(context.Context) error
type FinalizeFunc func() error

type ToggleBox struct {
	startFunc    StartFunc
	finalizeFunc FinalizeFunc
	cancel       context.CancelFunc
}

func NewToggleBox() *ToggleBox {
	return new(ToggleBox)
}

func (b *ToggleBox) Setup(startFunc StartFunc, finalizeFunc FinalizeFunc) {
	b.startFunc = startFunc
	b.finalizeFunc = finalizeFunc
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

	b.cancel()
	if err := b.finalizeFunc(); err != nil {
		return err
	}

	b.cancel = nil
	return nil
}
