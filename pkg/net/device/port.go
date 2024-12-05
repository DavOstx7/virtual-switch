package device

import (
	"context"
	"fmt"
	"project/pkg/net/frame"
	"project/pkg/toggle"
	"project/pkg/toggle/boxes"
)

type VirtualPort struct {
	*toggle.TogglerAPI
	FrameSniffer     *frame.Sniffer
	FrameTransmitter *frame.Transmitter
	portName         string
	sToggleBox       *boxes.AssistedSafeToggleBox
}

type VirtualPortConfig struct {
	PortName           string
	FrameSourceFactory frame.SourceFactory
	FrameWriterFactory frame.WriterFactory
}

func NewVirtualPort(config *VirtualPortConfig) *VirtualPort {
	sToggleBox := boxes.NewAssistedSafeToggleBox()

	vp := &VirtualPort{
		TogglerAPI:       toggle.NewTogglerAPI(sToggleBox),
		FrameSniffer:     frame.NewSniffer(config.PortName, config.FrameSourceFactory),
		FrameTransmitter: frame.NewTransmitter(config.PortName, config.FrameWriterFactory),
		portName:         config.PortName,
		sToggleBox:       sToggleBox,
	}

	sToggleBox.Setup(vp.startProcessingFrames, vp.finalizeProcessingFrames)

	return vp
}

func (vp *VirtualPort) Name() string {
	return vp.portName
}

func (vp *VirtualPort) InFrames() <-chan frame.Frame {
	return vp.FrameSniffer.InFrames()
}

func (vp *VirtualPort) OutFrames() chan<- frame.Frame {
	return vp.FrameTransmitter.OutFrames()
}

func (vp *VirtualPort) startProcessingFrames(ctx context.Context) error {
	for err := range toggle.On(ctx, vp.FrameSniffer, vp.FrameTransmitter) {
		fmt.Println(err)
	}

	return nil
}

func (vp *VirtualPort) finalizeProcessingFrames() error {
	for err := range toggle.Off(vp.FrameSniffer, vp.FrameTransmitter) {
		fmt.Println(err)
	}

	return nil
}
