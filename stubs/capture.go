package stubs

import (
	"fmt"
	"math/rand"
	"project/pkg"
	"time"
)

const MinJitterSeconds = 1

/*
type FrameSource interface {
	Frames() <-chan Frame
	Close()
}
*/

type FrameSourceStub struct {
	MaxJitterSeconds int
	PortName         string
	SourceMACs       []string
	DestinationMACs  []string
	Closed           bool

	frames chan pkg.Frame
}

func (s *FrameSourceStub) Frames() <-chan pkg.Frame {
	s.frames = make(chan pkg.Frame)

	go func() {
		for {
			randomDuration := time.Duration(MinJitterSeconds+rand.Intn(s.MaxJitterSeconds)) * time.Second

			time.Sleep(randomDuration)

			frame := &FrameStub{
				bytes:          nil,
				sourceMAC:      ChooseRandomly(s.SourceMACs),
				destinationMAC: ChooseRandomly(s.DestinationMACs),
			}

			fmt.Printf("captured frame %s on port '%s'\n", FrameToString(frame), s.PortName)
			s.frames <- frame
		}
	}()

	return s.frames
}

func (s *FrameSourceStub) Close() {
	fmt.Printf("closing frame capture on port '%s'\n", s.PortName)
	close(s.frames)
	s.Closed = true
}

/*
type FrameSourceProvider interface {
	FrameSource(portName string) (FrameSource, error)
}
*/

type FrameSourceProviderStub struct {
	MaxJitterSeconds int
	SourceMACs       []string
	DestinationMACs  []string
}

func (sp *FrameSourceProviderStub) FrameSource(portName string) (pkg.FrameSource, error) {
	return &FrameSourceStub{
		MaxJitterSeconds: sp.MaxJitterSeconds,
		PortName:         portName,
		SourceMACs:       sp.SourceMACs,
		DestinationMACs:  sp.DestinationMACs,
	}, nil
}
