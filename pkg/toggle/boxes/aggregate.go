package boxes

import (
	"context"
	"sync"
)

type SafeAggregateToggleBox struct {
	SafeToggleBox
	wg *sync.WaitGroup
}

func NewSafeAggregateToggleBox() *SafeAggregateToggleBox {
	return new(SafeAggregateToggleBox)
}

func (b *SafeAggregateToggleBox) BasicSetup(startFunc StartFunc) {
	wrappedStartFunc := b.wrapStartFunc(startFunc)
	defaultFinalizeFunc := b.defaultFinalizeFunc()
	b.SafeToggleBox.Setup(wrappedStartFunc, defaultFinalizeFunc)
}

func (b *SafeAggregateToggleBox) Setup(startFunc StartFunc, finalizeFunc FinalizeFunc) {
	wrappedStartFunc := b.wrapStartFunc(startFunc)
	wrappedFinalizeFunc := b.wrapFinalizeFunc(finalizeFunc)
	b.SafeToggleBox.Setup(wrappedStartFunc, wrappedFinalizeFunc)
}

func (b *SafeAggregateToggleBox) MarkActive(amount int) {
	b.wg.Add(amount)
}

func (b *SafeAggregateToggleBox) MarkDone(amount int) {
	b.wg.Add(-amount)
}

func (b *SafeAggregateToggleBox) MarkOneActive() {
	b.wg.Add(1)
}

func (b *SafeAggregateToggleBox) MarkOneDone() {
	b.wg.Done()
}

func (b *SafeAggregateToggleBox) wrapStartFunc(startFunc StartFunc) StartFunc {
	return func(ctx context.Context) error {
		b.wg = new(sync.WaitGroup)
		if err := startFunc(ctx); err != nil {
			b.wg = nil
			return err
		}
		return nil
	}
}

func (b *SafeAggregateToggleBox) wrapFinalizeFunc(finalizeFunc FinalizeFunc) FinalizeFunc {
	return func() error {
		b.wg.Wait()
		if err := finalizeFunc(); err != nil {
			return err
		}
		b.wg = nil
		return nil
	}
}

func (b *SafeAggregateToggleBox) defaultFinalizeFunc() FinalizeFunc {
	return func() error {
		b.wg.Wait()
		b.wg = nil
		return nil
	}
}
