package pkg

import (
	"context"
	"project/pkg/toggle"
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
	toggle.AtomicToggler
	portName       string
	inFrames       chan Frame
	sourceProvider FrameSourceProvider
}

func NewFrameCapture(portName string, sourceProvider FrameSourceProvider, stopTimeout time.Duration) *FrameCapture {
	fc := &FrameCapture{
		portName:       portName,
		inFrames:       make(chan Frame),
		sourceProvider: sourceProvider,
	}

	fc.Setup(fc.startFrameCapture)

	return fc
}

func (fc *FrameCapture) startFrameCapture(ctx context.Context) error {
	frameSource, err := fc.sourceProvider.FrameSource(fc.portName)
	if err != nil {
		return err
	}

	go fc.captureFrames(ctx, frameSource)
	return nil
}

func (fc *FrameCapture) captureFrames(ctx context.Context, frameSource FrameSource) {
	defer close(fc.Done)

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
