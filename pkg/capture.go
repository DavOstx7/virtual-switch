package pkg

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type FrameSource interface {
	Frames() <-chan Frame
	Close()
}

type FrameSourceProvider interface {
	FrameSource(portName string) (FrameSource, error)
}

type FrameCapture struct {
	portName       string
	inFrames       chan Frame
	sourceProvider FrameSourceProvider
	stopTimeout    time.Duration
	mu             sync.Mutex
	cancel         context.CancelFunc
	done           chan bool
}

func NewFrameCapture(portName string, sourceProvider FrameSourceProvider, stopTimeout time.Duration) *FrameCapture {
	return &FrameCapture{
		portName:       portName,
		inFrames:       make(chan Frame),
		sourceProvider: sourceProvider,
		stopTimeout:    stopTimeout,
	}
}

func (fc *FrameCapture) IsOn() bool {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	return fc._isOn()
}

func (fc *FrameCapture) On(ctx context.Context) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if fc._isOn() {
		return nil
	}

	frameSource, err := fc.sourceProvider.FrameSource(fc.portName)
	if err != nil {
		return err
	}

	newCtx, newCancel := context.WithCancel(ctx)
	fc.done = make(chan bool)

	go fc.startCapture(newCtx, frameSource)

	fc.cancel = newCancel
	return nil
}

func (fc *FrameCapture) Off() {
	fc.stopCapture(true)
}

func (fc *FrameCapture) Finish() {
	fc.stopCapture(false)
}

func (fc *FrameCapture) startCapture(ctx context.Context, frameSource FrameSource) {
	defer close(fc.done)

	frames := frameSource.Frames()
	defer frameSource.Close()

	for {
		select {
		case frame := <-frames:
			select {
			case fc.inFrames <- frame:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (fc *FrameCapture) stopCapture(shouldCancel bool) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if !fc._isOn() {
		return
	}

	if shouldCancel {
		fc.cancel()
	}

	fc.waitUntillStopped()
	fc.done = nil
	fc.cancel = nil
}

func (fc *FrameCapture) waitUntillStopped() {
	select {
	case <-fc.done:
	case <-time.After(fc.stopTimeout):
		msg := fmt.Sprintf(
			"frame-transmit: timed out while waiting %f seconds for capture to stop on port '%s'",
			fc.stopTimeout.Seconds(), fc.portName,
		)
		panic(msg)
	}
}

func (fc *FrameCapture) _isOn() bool {
	return fc.cancel != nil
}
