package pkg

import (
	"context"
	"fmt"
)

// Assumes Ethernet
type Frame interface {
	Bytes() []byte
	SourceMAC() string
	DestinationMAC() string
}

type VirtualPort struct {
	portName      string
	cancel        context.CancelFunc
	FrameCapture  *FrameCapture
	FrameTransmit *FrameTransmit
}

type VirtualPortConfig struct {
	PortName                 string
	FrameSourceProvider      FrameSourceProvider
	FrameTransmitterProvider FrameTransmitterProvider
}

func NewVirtualPort(config *VirtualPortConfig) *VirtualPort {
	return &VirtualPort{
		portName:      config.PortName,
		cancel:        nil,
		FrameCapture:  NewFrameCapture(config.PortName, config.FrameSourceProvider),
		FrameTransmit: NewFrameTransmit(config.PortName, config.FrameTransmitterProvider),
	}
}

func (vp *VirtualPort) PortName() string {
	return vp.FrameCapture.portName
}

func (vp *VirtualPort) InFrames() <-chan Frame {
	return vp.FrameCapture.inFrames
}

func (vp *VirtualPort) OutFrames() chan<- Frame {
	return vp.FrameTransmit.outFrames
}

func (vp *VirtualPort) IsOn() bool {
	return vp.cancel != nil
}

func (vp *VirtualPort) On(ctx context.Context) error {
	if vp.IsOn() {
		return nil
	}

	newCtx, cancel := context.WithCancel(ctx)

	if err := vp.FrameCapture.On(newCtx); err != nil {
		fmt.Println(err)
	}

	if err := vp.FrameTransmit.On(newCtx); err != nil {
		fmt.Println(err)
	}

	vp.cancel = cancel
	return nil
}

func (vp *VirtualPort) Off() {
	if !vp.IsOn() {
		return
	}

	vp.cancel()
	vp.FrameCapture.Finalize()
	vp.FrameTransmit.Finalize()
}
