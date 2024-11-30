package boxes

import (
	"context"
	"project/pkg/utils"
	"sync"
)

type CollectiveToggleBox struct {
	SafeToggleBox
	wg *sync.WaitGroup
}

func (b *CollectiveToggleBox) MarkActive(amount int) {
	b.wg.Add(amount)
}

func (b *CollectiveToggleBox) MarkDone(amount int) {
	b.wg.Add(-amount)
}

func (b *CollectiveToggleBox) MarkOneActive() {
	b.wg.Add(1)
}

func (b *CollectiveToggleBox) MarkOneDone() {
	b.wg.Done()
}

func (b *CollectiveToggleBox) Wait() {
	b.wg.Wait()
}

func (b *CollectiveToggleBox) Setup(startFunc StartFunc, stopFunc StopFunc) {
	wrappedStartFunc := b.wrapStartFunc(startFunc)
	wrappedStopFunc := b.wrapStopFunc(stopFunc)
	b.SafeToggleBox.Setup(wrappedStartFunc, wrappedStopFunc)
}

func (b *CollectiveToggleBox) wrapStartFunc(startFunc StartFunc) StartFunc {
	return func(ctx context.Context) error {
		b.wg = new(sync.WaitGroup)
		if err := startFunc(ctx); err != nil {
			b.wg = nil
			return err
		}
		return nil
	}
}

func (b *CollectiveToggleBox) wrapStopFunc(stopFunc StopFunc) StopFunc {
	return func() error {
		if err := stopFunc(); err != nil {
			return err
		}
		b.wg = nil
		return nil
	}
}

/* ------------------------------------------------------------------------------ */

type AssistedCollectiveToggleBox struct {
	CollectiveToggleBox
}

func NewAssistedCollectiveToggleBox() *AssistedCollectiveToggleBox {
	return new(AssistedCollectiveToggleBox)
}

func (b *AssistedCollectiveToggleBox) Setup(startFunc StartFunc, finalizeFunc func() error) {
	stopFunc := b.wrapFinalizeFuncToStopFunc(finalizeFunc)
	b.CollectiveToggleBox.Setup(startFunc, stopFunc)
}

func (b *AssistedCollectiveToggleBox) BasicSetup(startFunc StartFunc) {
	b.Setup(startFunc, nil)
}

func (b *AssistedCollectiveToggleBox) wrapFinalizeFuncToStopFunc(finalizeFunc func() error) StopFunc {
	return func() error {
		b.cancel()
		b.wg.Wait()
		return utils.ExecuteFunc(finalizeFunc)
	}
}
