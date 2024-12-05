package main

import (
	"context"
	"project/pkg/net/device"
	"project/stubs"
	"time"
)

const SwitchOnDuration = 10 * time.Second

func main() {
	vSwitch := device.NewVirtualSwitch(
		[]*device.VirtualPort{
			device.NewVirtualPort(
				&device.VirtualPortConfig{
					PortName: "eth1",
					FrameSourceFactory: &stubs.FrameSourceFactoryStub{
						MaxJitter:       5 * time.Second,
						SourceMACs:      []string{"AAA"},
						DestinationMACs: []string{"BBB", "CCC"},
					},
					FrameWriterFactory: &stubs.FrameWriterFactoryStub{},
				},
			),
			device.NewVirtualPort(
				&device.VirtualPortConfig{
					PortName: "eth2",
					FrameSourceFactory: &stubs.FrameSourceFactoryStub{
						MaxJitter:       10 * time.Second,
						SourceMACs:      []string{"BBB"},
						DestinationMACs: []string{"AAA", "CCC"},
					},
					FrameWriterFactory: &stubs.FrameWriterFactoryStub{},
				},
			),
			device.NewVirtualPort(
				&device.VirtualPortConfig{
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
	time.Sleep(SwitchOnDuration)
	vSwitch.Off()
}
