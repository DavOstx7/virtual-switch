package atomic

import (
	"context"
	"project/pkg/toggle/boxes"
)

type ToggleBox struct {
	boxes.SafeToggleBox
	done chan bool
}

func (b *ToggleBox) MarkDone() {
	close(b.done)
}

func (b *ToggleBox) Wait() {
	<-b.done
}

func (b *ToggleBox) Setup(startFunc boxes.StartFunc, stopFunc boxes.StopFunc) {
	wrappedStartFunc := b.wrapStartFunc(startFunc)
	wrappedStopFunc := b.wrapStopFunc(stopFunc)
	b.SafeToggleBox.Setup(wrappedStartFunc, wrappedStopFunc)
}

func (b *ToggleBox) wrapStartFunc(startFunc boxes.StartFunc) boxes.StartFunc {
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

func (b *ToggleBox) wrapStopFunc(stopFunc boxes.StopFunc) boxes.StopFunc {
	return func() error {
		if err := stopFunc(); err != nil {
			return err
		}
		b.done = nil
		return nil
	}
}
