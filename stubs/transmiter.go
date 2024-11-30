package stubs

import (
	"context"
	"fmt"
	"project/pkg/net"
)

type FrameWriter struct {
	portName string
	closed   bool
}

/*
type FrameWriter interface {
	WriteFrame(ctx context.Context, frame Frame) error
	Close()
}
*/

func (fw *FrameWriter) WriteFrame(ctx context.Context, frame net.Frame) error {
	fmt.Printf("sent frame %s to port '%s'\n", FrameToString(frame), fw.portName)
	return nil
}

func (fw *FrameWriter) Close() {
	fw.closed = true
	fmt.Printf("closed frame writer on port '%s'\n", fw.portName)
}

type FrameWriterFactoryStub struct {
	PortName string
}

/*
type FrameWriterFactory interface {
	NewFrameWriter(ctx context.Context, portName string) (FrameWriter, error)
}
*/

func (fwf *FrameWriterFactoryStub) NewFrameWriter(ctx context.Context, portName string) (net.FrameWriter, error) {
	return &FrameWriter{
		portName: portName,
		closed:   false,
	}, nil
}
