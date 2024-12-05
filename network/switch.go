package network

import (
	"context"
	"fmt"
	"project/network/frame"
	"project/network/mac"
	"project/toggle"
)

type Switch struct {
	toggle.AggregateToggler
	ports   []*Port
	macTable *mac.Table
}

func NewSwitch(ports []*Port, outputTableChanges bool) *Switch {
	s := &Switch{
		ports:   ports,
		macTable: mac.NewTable(outputTableChanges),
	}

	s.Setup(s.startPortForwarding)

	return s
}

func (s *Switch) GetPortByName(portName string) *Port {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for _, port := range s.ports {
		if port.Name() == portName {
			return port
		}
	}

	return nil
}

func (s *Switch) LookupPort(destinationMAC string) *Port {
	portName, exists := s.macTable.LookupPort(destinationMAC)
	if !exists {
		return nil
	}

	return s.GetPortByName(portName)

}

func (s *Switch) startPortForwarding(ctx context.Context) error {
	// could potentially be returning an error upon one of the ports returning an error (would require cleanup)

	for _, port := range s.ports {
		if err := port.On(ctx); err != nil {
			fmt.Println(err)
			continue
		}

		s.Wg.Add(1)
		go s.forwardFrames(ctx, port)
	}

	return nil
}

func (s *Switch) forwardFrames(ctx context.Context, port *Port) {
	defer s.Wg.Done()

	frames := port.InFrames()
	for {
		select {
		case frame := <-frames:
			s.forwardFrame(ctx, port, frame)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Switch) forwardFrame(ctx context.Context, sourcePort *Port, frame frame.Frame) {
	s.macTable.LearnMAC(frame.SourceMAC(), sourcePort.Name())

	destinationPort := s.LookupPort(frame.DestinationMAC())
	if destinationPort == nil {
		s.broadcastFrame(ctx, sourcePort, frame)
	} else {
		s.sendFrame(ctx, destinationPort, frame)
	}
}

func (s *Switch) broadcastFrame(ctx context.Context, sourcePort *Port, frame frame.Frame) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for _, port := range s.ports {
		if port.Name() == sourcePort.Name() {
			continue
		}

		if ok := s.sendFrame(ctx, port, frame); !ok {
			return
		}
	}
}

func (s *Switch) sendFrame(ctx context.Context, port *Port, frame frame.Frame) (ok bool) {
	select {
	case port.OutFrames() <- frame:
		return true
	case <-ctx.Done():
		return false
	}
}
