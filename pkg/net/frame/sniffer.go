package frame

import (
	"context"
	"project/pkg/toggle"
	"project/pkg/toggle/boxes/atomic"
)

type Source interface {
	Frames() <-chan Frame
	Close()
}

type SourceFactory interface {
	NewFrameSource(ctx context.Context, portName string) (Source, error)
}

type Sniffer struct {
	*toggle.TogglerAPI
	portName      string
	inFrames      chan Frame
	sourceFactory SourceFactory
	sToggleBox    *atomic.AssistedToggleBox
}

func NewSniffer(portName string, sourceFactory SourceFactory) *Sniffer {
	sToggleBox := atomic.NewAssistedToggleBox()

	s := &Sniffer{
		TogglerAPI:    toggle.NewTogglerAPI(sToggleBox),
		portName:      portName,
		inFrames:      make(chan Frame),
		sourceFactory: sourceFactory,
		sToggleBox:    sToggleBox,
	}

	sToggleBox.BasicSetup(s.startSniffingFrames)

	return s
}

func (s *Sniffer) InFrames() <-chan Frame {
	return s.inFrames
}

func (s *Sniffer) startSniffingFrames(ctx context.Context) error {
	source, err := s.sourceFactory.NewFrameSource(ctx, s.portName)
	if err != nil {
		return err
	}

	go s.sniffFrames(ctx, source)
	return nil
}

func (s *Sniffer) sniffFrames(ctx context.Context, source Source) {
	defer s.sToggleBox.MarkDone()

	frames := source.Frames()
	defer source.Close()

	for {
		select {
		case frame := <-frames:
			select {
			case s.inFrames <- frame:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
