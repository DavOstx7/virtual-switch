package main

import (
	"context"
	"project/pkg/net"
	"project/stubs"
	"time"
)

func main() {
	vSwitch := net.NewVirtualSwitch(
		[]*net.VirtualPort{
			net.NewVirtualPort(
				&net.VirtualPortConfig{
					PortName: "eth1",
					FrameSourceFactory: &stubs.FrameSourceFactoryStub{
						MaxJitter:       5 * time.Second,
						SourceMACs:      []string{"AAA"},
						DestinationMACs: []string{"BBB", "CCC"},
					},
					FrameWriterFactory: &stubs.FrameWriterFactoryStub{},
				},
			),
			net.NewVirtualPort(
				&net.VirtualPortConfig{
					PortName: "eth2",
					FrameSourceFactory: &stubs.FrameSourceFactoryStub{
						MaxJitter:       10 * time.Second,
						SourceMACs:      []string{"BBB"},
						DestinationMACs: []string{"AAA", "CCC"},
					},
					FrameWriterFactory: &stubs.FrameWriterFactoryStub{},
				},
			),
			net.NewVirtualPort(
				&net.VirtualPortConfig{
					PortName: "eth3",
					FrameSourceFactory: &stubs.FrameSourceFactoryStub{
						MaxJitter:       20 * time.Second,
						SourceMACs:      []string{"CCC"},
						DestinationMACs: []string{"AAA", "BBB"},
					},
					FrameWriterFactory: &stubs.FrameWriterFactoryStub{},
				},
			),
		},
		true,
	)

	ctx := context.Background()
	vSwitch.On(ctx)
	for {

	}
}
