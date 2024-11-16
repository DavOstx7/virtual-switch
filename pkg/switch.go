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
		wg:        new(sync.WaitGroup),
	}
}

func (vs *VirtualSwitch) IsOn() bool {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return vs._isOn()
}

func (vs *VirtualSwitch) On(ctx context.Context) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if vs._isOn() {
		return nil
	}

	newCtx, newCancel := context.WithCancel(ctx)
	vs.enableAllPorts(newCtx)
	vs.cancel = newCancel
	return nil
}

func (vs *VirtualSwitch) Off() {
	vs.stopForwarding(true)
}

func (vs *VirtualSwitch) Finalize() {
	vs.stopForwarding(false)
}

func (vs *VirtualSwitch) Wait() {
	vs.wg.Wait()
}

func (vs *VirtualSwitch) enableAllPorts(ctx context.Context) error {
	// could potentially be returning an error upon one of the ports returning an error (would require cleanup)

	for _, vPort := range vs.vPorts {
		if err := vPort.On(ctx); err != nil {
			fmt.Println(err)
			continue
		}

		vs.wg.Add(1)
		go vs.startForwarding(ctx, vPort)
	}

	return nil
}

func (vs *VirtualSwitch) startForwarding(ctx context.Context, vPort *VirtualPort) {
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

func (vs *VirtualSwitch) stopForwarding(shouldCancel bool) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if !vs._isOn() {
		return
	}

	if shouldCancel {
		vs.cancel()
	}

	vs.wg.Wait()
	vs.cancel = nil
}

func (vs *VirtualSwitch) _isOn() bool {
	return vs.cancel != nil
}
