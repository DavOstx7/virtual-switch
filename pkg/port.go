package pkg

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const DefaultStopTimeout time.Duration = 1 * time.Second

// Assumes Ethernet
type Frame interface {
	Bytes() []byte
	SourceMAC() string
	DestinationMAC() string
}

type VirtualPort struct {
	FrameCapture  *FrameCapture
	FrameTransmit *FrameTransmit
	portName      string
	mu            sync.Mutex
	cancel        context.CancelFunc
}

type VirtualPortConfig struct {
	PortName                 string
	FrameSourceProvider      FrameSourceProvider
	FrameTransmitterProvider FrameTransmitterProvider
}

func NewVirtualPort(config *VirtualPortConfig) *VirtualPort {
	return &VirtualPort{
		FrameCapture:  NewFrameCapture(config.PortName, config.FrameSourceProvider, DefaultStopTimeout),
		FrameTransmit: NewFrameTransmit(config.PortName, config.FrameTransmitterProvider, DefaultStopTimeout),
		portName:      config.PortName,
	}
}

func (vp *VirtualPort) PortName() string {
	return vp.FrameCapture.portName
}

func (vp *VirtualPort) InFrames() <-chan Frame {
	return vp.FrameCapture.inFrames
}

func (vp *VirtualPort) OutFrames() chan<- Frame {
	return vp.FrameTransmit.outFrames
}

func (vp *VirtualPort) IsOn() bool {
	vp.mu.Lock()
	defer vp.mu.Unlock()
	return vp._isOn()
}

func (vp *VirtualPort) On(ctx context.Context) error {
	vp.mu.Lock()
	defer vp.mu.Unlock()

	if vp._isOn() {
		return nil
	}

	newCtx, newCancel := context.WithCancel(ctx)
	vp.startTasks(newCtx)
	vp.cancel = newCancel
	return nil
}

func (vp *VirtualPort) Off() {
	vp.stopTasks(true)
}

func (vp *VirtualPort) Finalize() {
	vp.stopTasks(false)
}

func (vp *VirtualPort) startTasks(ctx context.Context) {
	// could potentially be sending an error to an error channel upon one of the tasks returning an error (would require cleanup)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		if err := vp.FrameCapture.On(ctx); err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()
	go func() {
		if err := vp.FrameTransmit.On(ctx); err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()

	wg.Wait()
	vp.cancel = nil
}

func (vp *VirtualPort) stopTasks(shouldCancel bool) {
	vp.mu.Lock()
	defer vp.mu.Unlock()

	if !vp._isOn() {
		return
	}

	if shouldCancel {
		vp.cancel()
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		vp.FrameCapture.Finalize()
		wg.Done()
	}()
	go func() {
		vp.FrameTransmit.Finalize()
		wg.Done()
	}()

	wg.Wait()
	vp.cancel = nil
}

func (vp *VirtualPort) _isOn() bool {
	return vp.cancel != nil
}
