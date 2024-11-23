package virtual

import (
	"context"
	"fmt"
	"project/network/frame"
	"project/toggle"
)

type Port struct {
	toggle.CommonToggler
	FrameCapture  *frame.Capture
	FrameTransmit *frame.Transmit
	name          string
}

type PortConfig struct {
	PortName                 string
	FrameSourceProvider      frame.SourceProvider
	FrameTransmitterProvider frame.TransmitterProvider
}

func NewPort(config *PortConfig) *Port {
	p := &Port{
		FrameCapture:  frame.NewCapture(config.PortName, config.FrameSourceProvider),
		FrameTransmit: frame.NewTransmit(config.PortName, config.FrameTransmitterProvider),
		name:          config.PortName,
	}

	p.Setup(p.startProcessingFrames, p.stopProcessingFrames)

	return p
}

func (p *Port) Name() string {
	return p.name
}

func (p *Port) InFrames() <-chan frame.Frame {
	return p.FrameCapture.InFrames()
}

func (p *Port) OutFrames() chan<- frame.Frame {
	return p.FrameTransmit.OutFrames()
}

func (p *Port) startProcessingFrames(ctx context.Context) error {
	for err := range toggle.On(ctx, p.FrameCapture, p.FrameTransmit) {
		fmt.Println(err)
	}
	return nil
}

func (p *Port) stopProcessingFrames() error {
	for err := range toggle.Off(p.FrameCapture, p.FrameTransmit) {
		fmt.Println(err)
	}
	return nil
}
