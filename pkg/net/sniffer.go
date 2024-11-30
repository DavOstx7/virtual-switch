package net

import (
	"context"
	"project/pkg/toggle"
	"project/pkg/toggle/boxes"
)

type FrameSource interface {
	Frames() <-chan Frame
	Close()
}

type FrameSourceFactory interface {
	NewFrameSource(ctx context.Context, portName string) (FrameSource, error)
}

type FrameSniffer struct {
	*toggle.TogglerAPI
	portName      string
	inFrames      chan Frame
	sourceFactory FrameSourceFactory
	sToggleBox    *boxes.AssistedAtomicToggleBox
}

func NewFrameSniffer(portName string, sourceFactory FrameSourceFactory) *FrameSniffer {
	sToggleBox := boxes.NewAssistedAtomicToggleBox()

	fs := &FrameSniffer{
		TogglerAPI:    toggle.NewTogglerAPI(sToggleBox),
		portName:      portName,
		inFrames:      make(chan Frame),
		sourceFactory: sourceFactory,
		sToggleBox:    sToggleBox,
	}

	sToggleBox.BasicSetup(fs.startSniffingFrames)

	return fs
}

func (fs *FrameSniffer) InFrames() <-chan Frame {
	return fs.inFrames
}

func (fs *FrameSniffer) startSniffingFrames(ctx context.Context) error {
	source, err := fs.sourceFactory.NewFrameSource(ctx, fs.portName)
	if err != nil {
		return err
	}

	go fs.sniffFrames(ctx, source)
	return nil
}

func (fs *FrameSniffer) sniffFrames(ctx context.Context, source FrameSource) {
	defer fs.sToggleBox.MarkDone()

	frames := source.Frames()
	defer source.Close()

	for {
		select {
		case frame := <-frames:
			select {
			case fs.inFrames <- frame:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
