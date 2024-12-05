package boxes

import "project/pkg/utils"

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
		b.Cancel()
		return utils.ExecuteFuncIfNotNil(finalizeFunc)
	}
}
