package stubs

import (
	"fmt"
	"math/rand"
	"project/network/frame"
	"sync"
	"time"
)

const DefaultFrameSourceStopTimeout time.Duration = 5

type FrameSource struct {
	PortName        string
	SourceMACs      []string
	DestinationMACs []string
	MaxJitter       time.Duration
	once            sync.Once
	frames          chan frame.Frame
	done            chan bool
	stopTimeout     time.Duration
}

func (fs *FrameSource) Frames() <-chan frame.Frame {
	return fs.frames
}

func (fs *FrameSource) Close() {
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

func (fs *FrameSource) generateRandomFrames() {
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

func (fs *FrameSource) waitWithJitter() (ok bool) {
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

func (s *FrameSource) close() {
	fmt.Printf("closing frame capture on port '%s'\n", s.PortName)
	close(s.frames)
	close(s.done)
}

type FrameSourceProvider struct {
	SourceMACs      []string
	DestinationMACs []string
	MaxJitter       time.Duration
}

func (sp *FrameSourceProvider) NewFrameSource(portName string) (frame.Source, error) {
	if sp.MaxJitter <= 0 {
		return nil, fmt.Errorf("frame-source-provider-stub: max jitter duration must be positive")
	}
	if len(sp.SourceMACs) == 0 || len(sp.DestinationMACs) == 0 {
		return nil, fmt.Errorf("frame-source-provider-stub: SourceMACs and DestinationMACs cannot be empty")
	}

	fs := &FrameSource{
		PortName:        portName,
		SourceMACs:      sp.SourceMACs,
		DestinationMACs: sp.DestinationMACs,
		MaxJitter:       sp.MaxJitter,
		frames:          make(chan frame.Frame),
		done:            make(chan bool),
		stopTimeout:     DefaultFrameSourceStopTimeout,
	}

	go fs.generateRandomFrames()

	return fs, nil
}
