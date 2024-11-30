package boxes

import (
	"context"
	"sync"
	"project/pkg/utils"
)

type SafeToggleBox struct {
	ToggleBox
	mu sync.Mutex
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

/* ------------------------------------------------------------------------------ */

type AssistedSafeToggleBox struct {
	SafeToggleBox
}

func NewAssistedSafeToggleBox() *AssistedSafeToggleBox {
	return new(AssistedSafeToggleBox)
}

func (b *AssistedSafeToggleBox) Setup(startFunc StartFunc, finalizeFunc func() error) {
	stopFunc := b.wrapFinalizeFuncToStopFunc(finalizeFunc)
	b.SafeToggleBox.Setup(startFunc, stopFunc)
}

func (b *AssistedSafeToggleBox) BasicSetup(startFunc StartFunc) {
	b.Setup(startFunc, nil)
}

func (b *AssistedSafeToggleBox) wrapFinalizeFuncToStopFunc(finalizeFunc func() error) StopFunc {
	return func() error {
		b.cancel()
		return utils.ExecuteFunc(finalizeFunc)
	}
}
