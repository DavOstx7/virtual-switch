package pkg

import (
	"context"
	"fmt"
	"project/pkg/toggle"
	"time"
)

type FrameTransmitter interface {
	TransmitFrame(frame Frame) error
	Close()
}

type FrameTransmitterProvider interface {
	FrameTransmitter(portName string) (FrameTransmitter, error)
}

type FrameTransmit struct {
	toggle.AtomicToggler
	portName            string
	outFrames           chan Frame
	transmitterProvider FrameTransmitterProvider
}

func NewFrameTransmit(portName string, transmitterProvider FrameTransmitterProvider, stopTimeout time.Duration) *FrameTransmit {
	fc := &FrameTransmit{
		portName:            portName,
		outFrames:           make(chan Frame),
		transmitterProvider: transmitterProvider,
	}

	fc.Setup(fc.startFrameTransmission)

	return fc
}

func (ft *FrameTransmit) startFrameTransmission(ctx context.Context) error {
	frameTransmitter, err := ft.transmitterProvider.FrameTransmitter(ft.portName)
	if err != nil {
		return err
	}

	go ft.transmitFrames(ctx, frameTransmitter)
	return nil
}

func (ft *FrameTransmit) transmitFrames(ctx context.Context, frameTransmitter FrameTransmitter) {
	defer close(ft.Done)
	defer frameTransmitter.Close()

	for {
		select {
		case frame := <-ft.outFrames:
			err := frameTransmitter.TransmitFrame(frame)
			if err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
