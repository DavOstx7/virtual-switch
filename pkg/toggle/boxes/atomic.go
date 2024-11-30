package boxes

import (
	"context"
	"project/pkg/utils"
)

type AtomicToggleBox struct {
	SafeToggleBox
	done chan bool
}

func (b *AtomicToggleBox) MarkDone() {
	close(b.done)
}

func (b *AtomicToggleBox) Wait() {
	<-b.done
}

func (b *AtomicToggleBox) Setup(startFunc StartFunc, stopFunc StopFunc) {
	wrappedStartFunc := b.wrapStartFunc(startFunc)
	wrappedStopFunc := b.wrapStopFunc(stopFunc)
	b.SafeToggleBox.Setup(wrappedStartFunc, wrappedStopFunc)
}

func (b *AtomicToggleBox) wrapStartFunc(startFunc StartFunc) StartFunc {
	return func(ctx context.Context) error {
		b.done = make(chan bool)
		if err := startFunc(ctx); err != nil {
			close(b.done)
			b.done = nil
			return err
		}
		return nil
	}
}

func (b *AtomicToggleBox) wrapStopFunc(stopFunc StopFunc) StopFunc {
	return func() error {
		if err := stopFunc(); err != nil {
			return err
		}
		b.done = nil
		return nil
	}
}

/* ------------------------------------------------------------------------------ */

type AssistedAtomicToggleBox struct {
	AtomicToggleBox
}

func NewAssistedAtomicToggleBox() *AssistedAtomicToggleBox {
	return new(AssistedAtomicToggleBox)
}

func (b *AssistedAtomicToggleBox) Setup(startFunc StartFunc, finalizeFunc func() error) {
	stopFunc := b.wrapFinalizeFuncToStopFunc(finalizeFunc)
	b.AtomicToggleBox.Setup(startFunc, stopFunc)
}

func (b *AssistedAtomicToggleBox) BasicSetup(startFunc StartFunc) {
	b.Setup(startFunc, nil)
}

func (b *AssistedAtomicToggleBox) wrapFinalizeFuncToStopFunc(finalizeFunc func() error) StopFunc {
	return func() error {
		b.cancel()
		<-b.done
		return utils.ExecuteFunc(finalizeFunc)
	}
}
