package net

import (
	"context"
	"fmt"
	"project/pkg/toggle"
	"project/pkg/toggle/boxes"
)

type FrameWriter interface {
	WriteFrame(ctx context.Context, frame Frame) error
	Close()
}

type FrameWriterFactory interface {
	NewFrameWriter(ctx context.Context, portName string) (FrameWriter, error)
}

type FrameTransmitter struct {
	*toggle.TogglerAPI
	portName      string
	outFrames     chan Frame
	writerFactory FrameWriterFactory
	sToggleBox    *boxes.AssistedAtomicToggleBox
}

func NewFrameTransmitter(portName string, writerFactory FrameWriterFactory) *FrameTransmitter {
	sToggleBox := boxes.NewAssistedAtomicToggleBox()

	ft := &FrameTransmitter{
		TogglerAPI:    toggle.NewTogglerAPI(sToggleBox),
		portName:      portName,
		outFrames:     make(chan Frame),
		writerFactory: writerFactory,
		sToggleBox:    sToggleBox,
	}

	sToggleBox.BasicSetup(ft.startTransmittingFrames)

	return ft
}

func (ft *FrameTransmitter) OutFrames() chan<- Frame {
	return ft.outFrames
}

func (ft *FrameTransmitter) startTransmittingFrames(ctx context.Context) error {
	writer, err := ft.writerFactory.NewFrameWriter(ctx, ft.portName)
	if err != nil {
		return err
	}

	go ft.transmitFrames(ctx, writer)
	return nil
}

func (ft *FrameTransmitter) transmitFrames(ctx context.Context, writer FrameWriter) {
	defer ft.sToggleBox.MarkDone()
	defer writer.Close()

	for {
		select {
		case frame := <-ft.outFrames:
			err := writer.WriteFrame(ctx, frame)
			if err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
