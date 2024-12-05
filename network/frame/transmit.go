package frame

import (
	"context"
	"fmt"
	"project/toggle"
)

type Transmitter interface {
	TransmitFrame(ctx context.Context, frame Frame) error
	Close()
}

type TransmitterProvider interface {
	NewFrameTransmitter(ctx context.Context, portName string) (Transmitter, error)
}

type Transmit struct {
	toggle.AtomicToggler
	portName            string
	outFrames           chan Frame
	transmitterProvider TransmitterProvider
}

func NewTransmit(portName string, transmitterProvider TransmitterProvider) *Transmit {
	t := &Transmit{
		portName:            portName,
		outFrames:           make(chan Frame),
		transmitterProvider: transmitterProvider,
	}

	t.Setup(t.startFrameTransmit)

	return t
}

func (t *Transmit) OutFrames() chan<- Frame {
	return t.outFrames
}

func (t *Transmit) startFrameTransmit(ctx context.Context) error {
	transmitter, err := t.transmitterProvider.NewFrameTransmitter(ctx, t.portName)
	if err != nil {
		return err
	}

	t.transmitFrames(ctx, transmitter)
	return nil
}

func (t *Transmit) transmitFrames(ctx context.Context, transmitter Transmitter) {
	defer transmitter.Close()

	for {
		select {
		case frame := <-t.outFrames:
			err := transmitter.TransmitFrame(ctx, frame)
			if err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
