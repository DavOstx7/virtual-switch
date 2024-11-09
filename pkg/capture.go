package pkg

import (
	"context"
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
	cancel         context.CancelFunc
	done           chan bool
}

func NewFrameCapture(portName string, sourceProvider FrameSourceProvider) *FrameCapture {
	return &FrameCapture{
		portName:       portName,
		inFrames:       make(chan Frame),
		sourceProvider: sourceProvider,
		cancel:         nil,
		done:           nil,
	}
}

func (fc *FrameCapture) IsOn() bool {
	return fc.cancel != nil
}

func (fc *FrameCapture) On(ctx context.Context) error {
	if fc.IsOn() {
		return nil
	}

	frameSource, err := fc.sourceProvider.FrameSource(fc.portName)
	if err != nil {
		return err
	}

	newCtx, cancel := context.WithCancel(ctx)
	fc.done = make(chan bool)

	go fc.startCapture(newCtx, frameSource)

	fc.cancel = cancel
	return nil
}

func (fc *FrameCapture) Off() {
	if !fc.IsOn() {
		return
	}

	fc.cancel()
	<-fc.done
	fc.Reset()
}

func (fc *FrameCapture) Finalize() {
	if !fc.IsOn() {
		return
	}

	<-fc.done
	fc.Reset()
}

func (fc *FrameCapture) Reset() {
	fc.done = nil
	fc.cancel = nil
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
