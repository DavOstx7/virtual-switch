package atomic

import (
	"project/pkg/toggle/boxes"
	"project/pkg/utils"
)

type AssistedToggleBox struct {
	ToggleBox
}

func NewAssistedToggleBox() *AssistedToggleBox {
	return new(AssistedToggleBox)
}

func (b *AssistedToggleBox) Setup(startFunc boxes.StartFunc, finalizeFunc func() error) {
	stopFunc := b.wrapFinalizeFuncToStopFunc(finalizeFunc)
	b.ToggleBox.Setup(startFunc, stopFunc)
}

func (b *AssistedToggleBox) BasicSetup(startFunc boxes.StartFunc) {
	b.Setup(startFunc, nil)
}

func (b *AssistedToggleBox) wrapFinalizeFuncToStopFunc(finalizeFunc func() error) boxes.StopFunc {
	return func() error {
		b.Cancel()
		b.Wait()
		return utils.ExecuteFunc(finalizeFunc)
	}
}
