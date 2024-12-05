package net

import (
	"context"
	"fmt"
	"project/pkg/toggle"
	"project/pkg/toggle/boxes"
)

type VirtualPort struct {
	*toggle.TogglerAPI
	FrameSniffer     *FrameSniffer
	FrameTransmitter *FrameTransmitter
	portName         string
	sToggleBox       *boxes.SafeToggleBox
}

type VirtualPortConfig struct {
	PortName           string
	FrameSourceFactory FrameSourceFactory
	FrameWriterFactory FrameWriterFactory
}

func NewVirtualPort(config *VirtualPortConfig) *VirtualPort {
	sToggleBox := boxes.NewSafeToggleBox()

	vp := &VirtualPort{
		TogglerAPI:       toggle.NewTogglerAPI(sToggleBox),
		FrameSniffer:     NewFrameSniffer(config.PortName, config.FrameSourceFactory),
		FrameTransmitter: NewFrameTransmitter(config.PortName, config.FrameWriterFactory),
		portName:         config.PortName,
		sToggleBox:       sToggleBox,
	}

	sToggleBox.SetStarter(vp.startProcessingFrames)
	sToggleBox.SetStopper(sToggleBox.NewStopperFromDefault(vp.stopProcessingFrames))

	return vp
}

func (vp *VirtualPort) Name() string {
	return vp.portName
}

func (vp *VirtualPort) InFrames() <-chan Frame {
	return vp.FrameSniffer.inFrames
}

func (vp *VirtualPort) OutFrames() chan<- Frame {
	return vp.FrameTransmitter.outFrames
}

func (vp *VirtualPort) startProcessingFrames(ctx context.Context) error {
	for err := range toggle.On(ctx, vp.FrameSniffer, vp.FrameTransmitter) {
		fmt.Println(err)
	}

	return nil
}

func (vp *VirtualPort) stopProcessingFrames() error {
	for err := range toggle.Off(vp.FrameSniffer, vp.FrameTransmitter) {
		fmt.Println(err)
	}

	return nil
}
