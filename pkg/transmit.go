package pkg

import (
	"context"
	"fmt"
)

type FrameTransmitter interface {
	TransmitFrame(frame Frame) error
	Close()
}

type FrameTransmitterProvider interface {
	FrameTransmitter(portName string) (FrameTransmitter, error)
}

type FrameTransmit struct {
	portName            string
	outFrames           chan Frame
	transmitterProvider FrameTransmitterProvider
	cancel              context.CancelFunc
	done                chan bool
}

func NewFrameTransmit(portName string, transmitterProvider FrameTransmitterProvider) *FrameTransmit {
	return &FrameTransmit{
		portName:            portName,
		outFrames:           make(chan Frame),
		transmitterProvider: transmitterProvider,
	}
}

func (ft *FrameTransmit) IsOn() bool {
	return ft.cancel != nil
}

func (ft *FrameTransmit) On(ctx context.Context) error {
	if ft.IsOn() {
		return nil
	}

	frameTransmitter, err := ft.transmitterProvider.FrameTransmitter(ft.portName)
	if err != nil {
		return err
	}

	newCtx, cancel := context.WithCancel(ctx)
	ft.done = make(chan bool)

	go ft.startTransmission(newCtx, frameTransmitter)

	ft.cancel = cancel
	return nil
}

func (ft *FrameTransmit) Off() {
	if !ft.IsOn() {
		return
	}

	ft.cancel()
	<-ft.done
	ft.Reset()
}

func (ft *FrameTransmit) Finalize() {
	if !ft.IsOn() {
		return
	}

	<-ft.done
	ft.Reset()
}

func (ft *FrameTransmit) Reset() {
	ft.done = nil
	ft.cancel = nil
}

func (ft *FrameTransmit) startTransmission(ctx context.Context, frameTransmitter FrameTransmitter) {
	defer close(ft.done)
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
