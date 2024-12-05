package frame

import (
	"context"
	"fmt"
	"project/pkg/toggle"
	"project/pkg/toggle/boxes/atomic"
)

type Writer interface {
	WriteFrame(ctx context.Context, frame Frame) error
	Close()
}

type WriterFactory interface {
	NewFrameWriter(ctx context.Context, portName string) (Writer, error)
}

type Transmitter struct {
	*toggle.TogglerAPI
	portName      string
	outFrames     chan Frame
	writerFactory WriterFactory
	sToggleBox    *atomic.AssistedToggleBox
}

func NewTransmitter(portName string, writerFactory WriterFactory) *Transmitter {
	sToggleBox := atomic.NewAssistedToggleBox()

	t := &Transmitter{
		TogglerAPI:    toggle.NewTogglerAPI(sToggleBox),
		portName:      portName,
		outFrames:     make(chan Frame),
		writerFactory: writerFactory,
		sToggleBox:    sToggleBox,
	}

	sToggleBox.BasicSetup(t.startTransmittingFrames)

	return t
}

func (t *Transmitter) OutFrames() chan<- Frame {
	return t.outFrames
}

func (t *Transmitter) startTransmittingFrames(ctx context.Context) error {
	writer, err := t.writerFactory.NewFrameWriter(ctx, t.portName)
	if err != nil {
		return err
	}

	go t.transmitFrames(ctx, writer)
	return nil
}

func (t *Transmitter) transmitFrames(ctx context.Context, writer Writer) {
	defer t.sToggleBox.MarkDone()
	defer writer.Close()

	for {
		select {
		case frame := <-t.outFrames:
			err := writer.WriteFrame(ctx, frame)
			if err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
