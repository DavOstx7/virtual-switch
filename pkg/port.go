package pkg

import (
	"context"
	"fmt"
	"project/pkg/toggle"
	"time"
)

const DefaultStopTimeout time.Duration = 1 * time.Second

// Assumes Ethernet
type Frame interface {
	Bytes() []byte
	SourceMAC() string
	DestinationMAC() string
}

type VirtualPort struct {
	toggle.CommonToggler
	FrameCapture  *FrameCapture
	FrameTransmit *FrameTransmit
	portName      string
}

type VirtualPortConfig struct {
	PortName                 string
	FrameSourceProvider      FrameSourceProvider
	FrameTransmitterProvider FrameTransmitterProvider
}

func NewVirtualPort(config *VirtualPortConfig) *VirtualPort {
	vp := &VirtualPort{
		FrameCapture:  NewFrameCapture(config.PortName, config.FrameSourceProvider, DefaultStopTimeout),
		FrameTransmit: NewFrameTransmit(config.PortName, config.FrameTransmitterProvider, DefaultStopTimeout),
		portName:      config.PortName,
	}

	vp.Setup(vp.startFrameOperations, vp.stopFrameOperations)

	return vp
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

func (vp *VirtualPort) startFrameOperations(ctx context.Context) error {
	for err := range toggle.On(ctx, vp.FrameCapture, vp.FrameTransmit) {
		fmt.Println(err)
	}
	return nil
}

func (vp *VirtualPort) stopFrameOperations() error {
	for err := range toggle.Off(vp.FrameCapture, vp.FrameTransmit) {
		fmt.Println(err)
	}
	return nil
}
