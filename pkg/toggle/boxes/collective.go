package boxes

import (
	"context"
	"sync"
)

type CollectiveToggleBox struct {
	SafeToggleBox
	wg *sync.WaitGroup
}

func NewCollectiveToggleBox() *CollectiveToggleBox {
	b := new(CollectiveToggleBox)

	b.SetStopper(b.DefaultStopper())

	return b
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

func (b *CollectiveToggleBox) SetStarter(startFunc StartFunc) {
	if startFunc == nil {
		return
	}

	b.SafeToggleBox.SetStarter(func(ctx context.Context) error {
		b.wg = new(sync.WaitGroup)
		if err := startFunc(ctx); err != nil {
			b.wg = nil
			return err
		}
		return nil
	})
}

func (b *CollectiveToggleBox) SetStopper(stopFunc StopFunc) {
	if stopFunc == nil {
		return
	}

	b.SafeToggleBox.SetStopper(func() error {
		if err := stopFunc(); err != nil {
			return err
		}
		b.wg = nil
		return nil
	})
}

func (b *CollectiveToggleBox) DefaultStopper() StopFunc {
	return func() error {
		b.Cancel()
		b.Wait()
		return nil
	}
}

func (b *CollectiveToggleBox) NewStopperFromDefault(stopExtension func() error) StopFunc {
	return extendStopFunc(b.DefaultStopper(), stopExtension)
}
