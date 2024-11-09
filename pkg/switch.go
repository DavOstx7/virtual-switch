package pkg

import (
	"context"
	"fmt"
	"sync"
)

type VirtualSwitch struct {
	vPorts []*VirtualPort
	mu     sync.Mutex
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func NewVirtualSwitch(vPorts []*VirtualPort) *VirtualSwitch {
	return &VirtualSwitch{
		vPorts: vPorts,
		cancel: nil,
		wg:     new(sync.WaitGroup),
		mu:     sync.Mutex{},
	}
}

func (vp *VirtualSwitch) IsOn() bool {
	return vp.cancel != nil
}

func (vs *VirtualSwitch) On(ctx context.Context) error {
	newCtx, cancel := context.WithCancel(ctx)

	vs.mu.Lock()
	for _, vPort := range vs.vPorts {
		if err := vPort.On(newCtx); err != nil {
			fmt.Println(err)
		}

		vs.wg.Add(1)
		go vs.forwardFrames(newCtx, vPort)
	}
	vs.mu.Unlock()

	vs.cancel = cancel
	return nil
}

func (vs *VirtualSwitch) Off() {
	if !vs.IsOn() {
		return
	}

	vs.cancel()
	vs.wg.Wait()
}

func (vs *VirtualSwitch) Wait() {
	if !vs.IsOn() {
		return
	}

	vs.wg.Wait()
}

func (vs *VirtualSwitch) forwardFrames(ctx context.Context, vPort *VirtualPort) {
	defer vs.wg.Done()

	frames := vPort.InFrames()
	for {
		select {
		case frame := <-frames:
			vs.forwardFrame(ctx, vPort, frame)
		case <-ctx.Done():
			return
		}
	}
}

func (vs *VirtualSwitch) forwardFrame(ctx context.Context, vSourcePort *VirtualPort, frame Frame) {
	/*
		sMacTable.UpdateSourceEntry(frame.SourceMAC(), vSourcePort.PortName())

		vDestinationPort := sMACTable.GetDestinationPort(frame.DestinationMAC())
		if vDestinationPort != nil {
			vDestinationPort.OutFrames() <- frame
		} else {
		 	vs.broadcastFrame(ctx, vSourcePort, frame)
		}
	*/
	vs.broadcastFrame(ctx, vSourcePort, frame)
}

func (vs *VirtualSwitch) broadcastFrame(ctx context.Context, vSourcePort *VirtualPort, frame Frame) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	for _, vPort := range vs.vPorts {
		if vPort.PortName() == vSourcePort.PortName() {
			continue
		}

		select {
		case vPort.OutFrames() <- frame:
		case <-ctx.Done():
			return
		}
	}
}
