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
	return new(SafeToggleBox)
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
