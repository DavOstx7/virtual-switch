package pkg

import (
	"context"
	"fmt"
	"project/pkg/toggle"
)

type VirtualSwitch struct {
	toggle.AggregateToggler
	vPorts    []*VirtualPort
	sMACTable *MACTable
}

func NewVirtualSwitch(vPorts []*VirtualPort, outputTableChanges bool) *VirtualSwitch {
	vs := &VirtualSwitch{
		vPorts:    vPorts,
		sMACTable: NewMACTable(outputTableChanges),
	}

	vs.Setup(vs.startFrameForwarding)

	return vs
}

func (vs *VirtualSwitch) startFrameForwarding(ctx context.Context) error {
	// could potentially be returning an error upon one of the ports returning an error (would require cleanup)

	for _, vPort := range vs.vPorts {
		if err := vPort.On(ctx); err != nil {
			fmt.Println(err)
			continue
		}

		vs.Wg.Add(1)
		go vs.forwardFrames(ctx, vPort)
	}

	return nil
}

func (vs *VirtualSwitch) forwardFrames(ctx context.Context, vPort *VirtualPort) {
	defer vs.Wg.Done()

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
	vs.Mu.Lock()
	defer vs.Mu.Unlock()

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
