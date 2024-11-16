package stubs

import (
	"fmt"
	"math/rand"
	"project/pkg"
	"sync"
	"time"
)

const DefaultFrameSourceStopTimeout time.Duration = 5

type FrameSourceStub struct {
	PortName        string
	SourceMACs      []string
	DestinationMACs []string
	MaxJitter       time.Duration
	once            sync.Once
	frames          chan pkg.Frame
	done            chan bool
	stopTimeout     time.Duration
}

func NewFrameSourceStub(portName string, sourceMACs, destinationMACs []string, maxJitter time.Duration) *FrameSourceStub {
	fs := &FrameSourceStub{
		PortName:        portName,
		SourceMACs:      sourceMACs,
		DestinationMACs: destinationMACs,
		MaxJitter:       maxJitter,
		frames:          make(chan pkg.Frame),
		done:            make(chan bool),
		stopTimeout:     DefaultFrameSourceStopTimeout,
	}

	go fs.sendRandomFrames()

	return fs
}

func (fs *FrameSourceStub) Frames() <-chan pkg.Frame {
	return fs.frames
}

func (fs *FrameSourceStub) Close() {
	fs.once.Do(func() {
		select {
		case fs.done <- true:
		case <-time.After(fs.stopTimeout):
			msg := fmt.Sprintf(
				"frame-source-stab: timed out while waiting %f seconds to signal done on port '%s'",
				fs.stopTimeout.Seconds(), fs.PortName,
			)
			panic(msg)
		}
	})
}

func (fs *FrameSourceStub) sendRandomFrames() {
	defer fs.close()

	for {
		if ok := fs.waitWithJitter(); !ok {
			return
		}

		frame := &FrameStub{
			bytes:          nil,
			sourceMAC:      ChooseRandomly(fs.SourceMACs),
			destinationMAC: ChooseRandomly(fs.DestinationMACs),
		}

		select {
		case fs.frames <- frame:
			fmt.Printf("captured frame %s on port '%s'\n", FrameToString(frame), fs.PortName)
		case <-fs.done:
			return
		}
	}
}

func (fs *FrameSourceStub) waitWithJitter() (ok bool) {
	waitTime := time.Duration(rand.Intn(int(fs.MaxJitter)))

	timer := time.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		return true
	case <-fs.done:
		return false
	}
}

func (s *FrameSourceStub) close() {
	fmt.Printf("closing frame capture on port '%s'\n", s.PortName)
	close(s.frames)
	close(s.done)
}

type FrameSourceProviderStub struct {
	SourceMACs      []string
	DestinationMACs []string
	MaxJitter       time.Duration
}

func (sp *FrameSourceProviderStub) FrameSource(portName string) (pkg.FrameSource, error) {
	if sp.MaxJitter <= 0 {
		return nil, fmt.Errorf("frame-source-provider-stub: max jitter duration must be positive")
	}
	if len(sp.SourceMACs) == 0 || len(sp.DestinationMACs) == 0 {
		return nil, fmt.Errorf("frame-source-provider-stub: SourceMACs and DestinationMACs cannot be empty")
	}

	return NewFrameSourceStub(portName, sp.SourceMACs, sp.DestinationMACs, sp.MaxJitter), nil
}
