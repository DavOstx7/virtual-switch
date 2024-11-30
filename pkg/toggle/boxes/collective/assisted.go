package collective

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

func (tb *AssistedToggleBox) Setup(startFunc boxes.StartFunc, finalizeFunc func() error) {
	stopFunc := tb.wrapFinalizeFuncToStopFunc(finalizeFunc)
	tb.ToggleBox.Setup(startFunc, stopFunc)
}

func (tb *AssistedToggleBox) BasicSetup(startFunc boxes.StartFunc) {
	tb.Setup(startFunc, nil)
}

func (tb *AssistedToggleBox) wrapFinalizeFuncToStopFunc(finalizeFunc func() error) boxes.StopFunc {
	return func() error {
		tb.Cancel()
		tb.Wait()
		return utils.ExecuteFunc(finalizeFunc)
	}
}
