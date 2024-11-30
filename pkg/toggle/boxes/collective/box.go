package collective

import (
	"context"
	"project/pkg/toggle/boxes"
	"sync"
)

type ToggleBox struct {
	boxes.SafeToggleBox
	wg *sync.WaitGroup
}

func (tb *ToggleBox) MarkActive(amount int) {
	tb.wg.Add(amount)
}

func (tb *ToggleBox) MarkDone(amount int) {
	tb.wg.Add(-amount)
}

func (tb *ToggleBox) MarkOneActive() {
	tb.wg.Add(1)
}

func (b *ToggleBox) MarkOneDone() {
	b.wg.Done()
}

func (tb *ToggleBox) Wait() {
	tb.wg.Wait()
}

func (tb *ToggleBox) Setup(startFunc boxes.StartFunc, stopFunc boxes.StopFunc) {
	wrappedStartFunc := tb.wrapStartFunc(startFunc)
	wrappedStopFunc := tb.wrapStopFunc(stopFunc)
	tb.SafeToggleBox.Setup(wrappedStartFunc, wrappedStopFunc)
}

func (tb *ToggleBox) wrapStartFunc(startFunc boxes.StartFunc) boxes.StartFunc {
	return func(ctx context.Context) error {
		tb.wg = new(sync.WaitGroup)
		if err := startFunc(ctx); err != nil {
			tb.wg = nil
			return err
		}
		return nil
	}
}

func (b *ToggleBox) wrapStopFunc(stopFunc boxes.StopFunc) boxes.StopFunc {
	return func() error {
		if err := stopFunc(); err != nil {
			return err
		}
		b.wg = nil
		return nil
	}
}
