package pkg

import (
	"context"
	"fmt"
	"sync"
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
	portName            string
	outFrames           chan Frame
	transmitterProvider FrameTransmitterProvider
	stopTimeout         time.Duration
	mu                  sync.Mutex
	cancel              context.CancelFunc
	done                chan bool
}

func NewFrameTransmit(portName string, transmitterProvider FrameTransmitterProvider, stopTimeout time.Duration) *FrameTransmit {
	return &FrameTransmit{
		portName:            portName,
		outFrames:           make(chan Frame),
		transmitterProvider: transmitterProvider,
		stopTimeout:         stopTimeout,
	}
}

func (ft *FrameTransmit) IsOn() bool {
	ft.mu.Lock()
	defer ft.mu.Unlock()
	return ft._isOn()
}

func (ft *FrameTransmit) On(ctx context.Context) error {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	if ft._isOn() {
		return nil
	}

	frameTransmitter, err := ft.transmitterProvider.FrameTransmitter(ft.portName)
	if err != nil {
		return err
	}

	newCtx, newCancel := context.WithCancel(ctx)
	ft.done = make(chan bool)

	go ft.startTransmission(newCtx, frameTransmitter)

	ft.cancel = newCancel
	return nil
}

func (ft *FrameTransmit) Off() {
	ft.stopTransmission(true)
}

func (ft *FrameTransmit) Finish() {
	ft.stopTransmission(false)
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

func (ft *FrameTransmit) stopTransmission(shouldCancel bool) {
	ft.mu.Lock()
	defer ft.mu.Unlock()

	if !ft._isOn() {
		return
	}

	if shouldCancel {
		ft.cancel()
	}

	ft.waitUntillStopped()
	ft.done = nil
	ft.cancel = nil
}

func (ft *FrameTransmit) waitUntillStopped() {
	select {
	case <-ft.done:
	case <-time.After(ft.stopTimeout):
		msg := fmt.Sprintf(
			"frame-transmit: timed out while waiting %f seconds for transmission to stop on port '%s'",
			ft.stopTimeout.Seconds(), ft.portName,
		)
		panic(msg)

	}
}

func (ft *FrameTransmit) _isOn() bool {
	return ft.cancel != nil
}
