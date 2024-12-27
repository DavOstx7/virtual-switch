package boxes

import (
	"context"
	"sync"
)

type SafeToggleBox struct {
	ToggleBox
	mu sync.Mutex
}

func NewSafeToggleBox() *SafeToggleBox {
	b := new(SafeToggleBox)

	b.SetStopper(b.DefaultStopper())

	return b
}

func (b *SafeToggleBox) IsOn() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.ToggleBox.IsOn()
}

func (b *SafeToggleBox) On(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.ToggleBox.On(ctx)
}

func (b *SafeToggleBox) Off() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.ToggleBox.Off()
}

func (b *SafeToggleBox) DefaultStopper() StopFunc {
	return func() error {
		b.Cancel()
		return nil
	}
}

func (b *SafeToggleBox) NewStopperFromDefault(stopExtension func() error) StopFunc {
	return extendStopFunc(b.DefaultStopper(), stopExtension)
}
