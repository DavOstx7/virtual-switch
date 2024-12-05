package network

import (
	"context"
	"project/network/frame"
	"project/network/mac"
	"project/toggle"
	"sync"
)

type Switch struct {
	toggle.AggregateToggler
	mu       sync.Mutex
	ports    []*Port
	macTable *mac.Table
}

func NewSwitch(ports []*Port, outputTableChanges bool) *Switch {
	s := &Switch{
		ports:    ports,
		macTable: mac.NewTable(outputTableChanges),
	}

	s.Setup(s.getStartPortForwardingFuncs())

	return s
}

func (s *Switch) getStartPortForwardingFuncs() []func(context.Context) error {
	startFuncs := make([]func(context.Context) error, len(s.ports))

	for i, port := range s.ports {
		startFuncs[i] = func(ctx context.Context) error {
			if err := port.On(ctx); err != nil {
				return err
			}

			s.forwardFrames(ctx, port)
			return nil
		}

	}

	return startFuncs
}

func (s *Switch) forwardFrames(ctx context.Context, port *Port) {
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

func (s *Switch) LookupPort(destinationMAC string) *Port {
	portName, exists := s.macTable.LookupPort(destinationMAC)
	if !exists {
		return nil
	}

	return s.GetPortByName(portName)

}

func (s *Switch) GetPortByName(portName string) *Port {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, port := range s.ports {
		if port.Name() == portName {
			return port
		}
	}

	return nil
}

func (s *Switch) broadcastFrame(ctx context.Context, sourcePort *Port, frame frame.Frame) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
