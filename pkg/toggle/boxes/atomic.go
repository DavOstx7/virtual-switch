package boxes

import "context"

type AtomicToggleBox struct {
	SafeToggleBox
	done chan bool
}

func NewAtomicToggleBox() *AtomicToggleBox {
	b := new(AtomicToggleBox)

	b.SetStopper(b.DefaultStopper())

	return b
}

func (b *AtomicToggleBox) MarkDone() {
	close(b.done)
}

func (b *AtomicToggleBox) Wait() {
	<-b.done
}

func (b *AtomicToggleBox) SetStarter(startFunc StartFunc) {
	if startFunc == nil {
		return
	}

	b.SafeToggleBox.SetStarter(func(ctx context.Context) error {
		b.done = make(chan bool)
		if err := startFunc(ctx); err != nil {
			close(b.done)
			b.done = nil
			return err
		}
		return nil
	})
}

func (b *AtomicToggleBox) SetStopper(stopFunc StopFunc) {
	if stopFunc == nil {
		return
	}

	b.SafeToggleBox.SetStopper(func() error {
		if err := stopFunc(); err != nil {
			return err
		}
		b.done = nil
		return nil
	})
}

func (b *AtomicToggleBox) DefaultStopper() StopFunc {
	return func() error {
		b.Cancel()
		b.Wait()
		return nil
	}
}

func (b *AtomicToggleBox) NewStopperFromDefault(stopExtension func() error) StopFunc {
	return extendStopFunc(b.DefaultStopper(), stopExtension)
}
