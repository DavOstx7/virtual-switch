package net

import (
	"context"
	"fmt"
	"project/pkg/toggle"
	"project/pkg/toggle/boxes"
	"sync"
)

type VirtualSwitch struct {
	*toggle.TogglerAPI
	vPorts     []*VirtualPort
	sMACTable  *MACTable
	sToggleBox *boxes.SafeAggregateToggleBox
	mu         sync.Mutex
}

func NewVirtualSwitch(vPorts []*VirtualPort, outputTableChanges bool) *VirtualSwitch {
	sToggleBox := boxes.NewSafeAggregateToggleBox()

	vs := &VirtualSwitch{
		TogglerAPI: toggle.NewTogglerAPI(sToggleBox),
		vPorts:     vPorts,
		sMACTable:  NewMACTable(outputTableChanges),
		sToggleBox: sToggleBox,
	}

	sToggleBox.Setup(vs.startPortForwarding, vs.finalizePortForwarding)

	return vs
}

func (vs *VirtualSwitch) startPortForwarding(ctx context.Context) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	for _, vPort := range vs.vPorts {
		if err := vPort.On(ctx); err != nil {
			fmt.Println(err)
			continue
		}

		vs.sToggleBox.MarkOneActive()
		go vs.forwardFrames(ctx, vPort)
	}

	return nil
}

func (vs *VirtualSwitch) finalizePortForwarding() error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	for _, vPort := range vs.vPorts {
		if err := vPort.Off(); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (vs *VirtualSwitch) forwardFrames(ctx context.Context, vPort *VirtualPort) {
	defer vs.sToggleBox.MarkOneDone()

	frames := vPort.InFrames()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			select {
			case frame := <-frames:
				vs.forwardFrame(ctx, vPort, frame)
			case <-ctx.Done():
				return
			}
		}
	}
}

func (vs *VirtualSwitch) forwardFrame(ctx context.Context, vSourcePort *VirtualPort, frame Frame) {
	vs.sMACTable.LearnMAC(frame.SourceMAC(), vSourcePort)

	vDestinationPort := vs.sMACTable.LookupPort(frame.DestinationMAC())
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
		if vPort.Name() == vSourcePort.Name() {
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
