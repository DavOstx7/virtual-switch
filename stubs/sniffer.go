package stubs

import (
	"context"
	"fmt"
	"math/rand"
	"project/pkg/net"
	"sync"
	"time"
)

const DefaultFrameSourceStopTimeout time.Duration = 5 * time.Second

type FrameSourceStub struct {
	SourceMACs      []string
	destinationMACs []string
	maxJitter       time.Duration
	portName        string
	frames          chan net.Frame
	done            chan bool
	stopTimeout     time.Duration
	once            sync.Once
}

/*
type FrameSource interface {
	Frames() <-chan Frame
	Close()
}
*/

func (fs *FrameSourceStub) Frames() <-chan net.Frame {
	return fs.frames
}

func (fs *FrameSourceStub) Close() {
	fs.once.Do(func() {
		select {
		case <-fs.done:
			fmt.Printf("closed frame source on port '%s'\n", fs.portName)
		case <-time.After(fs.stopTimeout):
			msg := fmt.Sprintf(
				"timed out after %f seconds while trying to close frame source on port '%s'",
				fs.stopTimeout.Seconds(), fs.portName,
			)
			panic(msg)
		}
	})
}

func (fs *FrameSourceStub) generateRandomFrames(ctx context.Context) {
	defer fs.close()

	for {
		if ok := fs.waitWithJitter(ctx); !ok {
			return
		}

		frame := &FrameStub{
			bytes:          nil,
			sourceMAC:      ChooseRandomly(fs.SourceMACs),
			destinationMAC: ChooseRandomly(fs.destinationMACs),
		}

		select {
		case fs.frames <- frame:
			fmt.Printf("received frame %s from port '%s'\n", FrameToString(frame), fs.portName)
		case <-ctx.Done():
			return
		}
	}
}

func (fs *FrameSourceStub) waitWithJitter(ctx context.Context) (ok bool) {
	waitTime := time.Duration(rand.Intn(int(fs.maxJitter)))

	timer := time.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func (fs *FrameSourceStub) close() {
	close(fs.frames)
	close(fs.done)
}

type FrameSourceFactoryStub struct {
	SourceMACs      []string
	DestinationMACs []string
	MaxJitter       time.Duration
}

/*
type FrameSourceFactory interface {
	NewFrameSource(ctx context.Context, portName string) (FrameSource, error)
}
*/

func (fsf *FrameSourceFactoryStub) NewFrameSource(ctx context.Context, portName string) (net.FrameSource, error) {
	if fsf.MaxJitter <= 0 {
		return nil, fmt.Errorf("max jitter duration must be positive")
	}
	if len(fsf.SourceMACs) == 0 || len(fsf.DestinationMACs) == 0 {
		return nil, fmt.Errorf("source & destination MAC addresses cannot be empty")
	}

	fs := &FrameSourceStub{
		portName:        portName,
		SourceMACs:      fsf.SourceMACs,
		destinationMACs: fsf.DestinationMACs,
		maxJitter:       fsf.MaxJitter,
		frames:          make(chan net.Frame),
		done:            make(chan bool),
		stopTimeout:     DefaultFrameSourceStopTimeout,
	}

	go fs.generateRandomFrames(ctx)

	return fs, nil
}
