package stubs

import (
	"context"
	"fmt"
	"project/network/frame"
)

type FrameTransmitter struct {
	PortName string
	closed   bool
}

func (t *FrameTransmitter) TransmitFrame(ctx context.Context, frame frame.Frame) error {
	fmt.Printf("transmitting frame %s on port '%s'\n", FrameToString(frame), t.PortName)
	return nil
}

func (t *FrameTransmitter) Close() {
	fmt.Printf("closing frame transmission on port '%s'\n", t.PortName)
	t.closed = true
}

type FrameTransmitterProvider struct {
	PortName string
}

func (ftf *FrameTransmitterProvider) NewFrameTransmitter(portName string) (frame.Transmitter, error) {
	return &FrameTransmitter{
		PortName: portName,
		closed:   false,
	}, nil
}
