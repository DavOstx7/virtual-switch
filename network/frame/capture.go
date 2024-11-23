package frame

import (
	"context"
	"project/toggle"
)

type Source interface {
	Frames() <-chan Frame
	Close()
}

type SourceProvider interface {
	NewFrameSource(portName string) (Source, error)
}

type Capture struct {
	toggle.AtomicToggler
	portName       string
	inFrames       chan Frame
	sourceProvider SourceProvider
}

func NewCapture(portName string, sourceProvider SourceProvider) *Capture {
	c := &Capture{
		portName:       portName,
		inFrames:       make(chan Frame),
		sourceProvider: sourceProvider,
	}

	c.Setup(c.startFrameCapture)

	return c
}

func (c *Capture) InFrames() <-chan Frame {
	return c.inFrames
}

func (c *Capture) startFrameCapture(ctx context.Context) error {
	source, err := c.sourceProvider.NewFrameSource(c.portName)
	if err != nil {
		return err
	}

	go c.captureFrames(ctx, source)
	return nil
}

func (c *Capture) captureFrames(ctx context.Context, source Source) {
	defer close(c.Done)

	frames := source.Frames()
	defer source.Close()

	for {
		select {
		case frame := <-frames:
			select {
			case c.inFrames <- frame:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
