package pkg

import (
	"context"
	"fmt"
	"sync"
)

type VirtualSwitch struct {
	vPorts    []*VirtualPort
	sMACTable *MACTable
	mu        sync.Mutex
	cancel    context.CancelFunc
	wg        *sync.WaitGroup
}

func NewVirtualSwitch(vPorts []*VirtualPort, outputTableChanges bool) *VirtualSwitch {
	return &VirtualSwitch{
		vPorts:    vPorts,
		sMACTable: NewMACTable(outputTableChanges),
		cancel:    nil,
		wg:        new(sync.WaitGroup),
		mu:        sync.Mutex{},
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
	vs.sMACTable.UpdateSourceEntry(frame.SourceMAC(), vSourcePort)

	vDestinationPort := vs.sMACTable.GetDestinationPort(frame.DestinationMAC())
	if vDestinationPort == nil {
		vs.broadcastFrame(ctx, vSourcePort, frame)
	} else {
		vs.sendFrame(ctx, vDestinationPort, frame)
	}
}

func (vs *VirtualSwitch) broadcastFrame(ctx context.Context, vSourcePort *VirtualPort, frame Frame) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	for _, vPort := range vs.vPorts {
		if vPort.PortName() == vSourcePort.PortName() {
			continue
		}

		if ok := vs.sendFrame(ctx, vPort, frame); !ok {
			return
		}
	}
}

func (vs *VirtualSwitch) sendFrame(ctx context.Context, vPort *VirtualPort, frame Frame) (ok bool) {
	select {
	case vPort.OutFrames() <- frame:
		return true
	case <-ctx.Done():
		return false
	}
}
