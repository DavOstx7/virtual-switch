package stubs

import (
	"fmt"
	"project/pkg"
)

/*
type FrameTransmitter interface {
	TransmitFrame(frame Frame) error
	Close()
}
*/

type FrameTransmitterStub struct {
	PortName string
	Closed   bool
}

func (t *FrameTransmitterStub) TransmitFrame(frame pkg.Frame) error {
	fmt.Printf("transmitting frame %s on port '%s'\n", FrameToString(frame), t.PortName)
	return nil
}

func (t *FrameTransmitterStub) Close() {
	fmt.Printf("closing frame transmission on port '%s'\n", t.PortName)
	t.Closed = true
}

/*
type FrameTransmitterProvider interface {
	FrameTransmitter(portName string) (FrameTransmitter, error)
}
*/

type FrameTransmitterProviderStub struct {
	PortName string
}

func (ftf *FrameTransmitterProviderStub) FrameTransmitter(portName string) (pkg.FrameTransmitter, error) {
	return &FrameTransmitterStub{
		PortName: portName,
	}, nil
}
