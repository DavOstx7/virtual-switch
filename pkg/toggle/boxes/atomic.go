package boxes

import "context"

type SafeAtomicToggleBox struct {
	SafeToggleBox
	done chan bool
}

func NewSafeAtomicToggleBox() *SafeAtomicToggleBox {
	return new(SafeAtomicToggleBox)
}

func (b *SafeAtomicToggleBox) BasicSetup(startFunc StartFunc) {
	wrappedStartFunc := b.wrapStartFunc(startFunc)
	defaultFinalizeFunc := b.defaultFinalizeFunc()
	b.SafeToggleBox.Setup(wrappedStartFunc, defaultFinalizeFunc)
}

func (b *SafeAtomicToggleBox) Setup(startFunc StartFunc, finalizeFunc FinalizeFunc) {
	wrappedStartFunc := b.wrapStartFunc(startFunc)
	wrappedFinalizeFunc := b.wrapFinalizeFunc(finalizeFunc)
	b.SafeToggleBox.Setup(wrappedStartFunc, wrappedFinalizeFunc)
}

func (b *SafeAtomicToggleBox) MarkDone() {
	close(b.done)
}

func (b *SafeAtomicToggleBox) wrapStartFunc(startFunc StartFunc) StartFunc {
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

func (b *SafeAtomicToggleBox) wrapFinalizeFunc(finalizeFunc FinalizeFunc) FinalizeFunc {
	return func() error {
		<-b.done
		if err := finalizeFunc(); err != nil {
			return err
		}
		b.done = nil
		return nil
	}
}

func (b *SafeAtomicToggleBox) defaultFinalizeFunc() FinalizeFunc {
	return func() error {
		<-b.done
		b.done = nil
		return nil
	}
}
