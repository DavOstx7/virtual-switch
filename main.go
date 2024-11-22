package main

import (
	"context"
	"project/pkg"
	"project/stubs"
	"time"
)

func main() {
	vSwitch := pkg.NewVirtualSwitch(
		[]*pkg.VirtualPort{
			pkg.NewVirtualPort(
				&pkg.VirtualPortConfig{
					PortName: "eth1",
					FrameSourceProvider: &stubs.FrameSourceProviderStub{
						MaxJitter:       5 * time.Second,
						SourceMACs:      []string{"AAA"},
						DestinationMACs: []string{"BBB", "CCC"},
					},
					FrameTransmitterProvider: &stubs.FrameTransmitterProviderStub{},
				},
			),
			pkg.NewVirtualPort(
				&pkg.VirtualPortConfig{
					PortName: "eth2",
					FrameSourceProvider: &stubs.FrameSourceProviderStub{
						MaxJitter:       10 * time.Second,
						SourceMACs:      []string{"BBB"},
						DestinationMACs: []string{"AAA", "CCC"},
					},
					FrameTransmitterProvider: &stubs.FrameTransmitterProviderStub{},
				},
			),
			pkg.NewVirtualPort(
				&pkg.VirtualPortConfig{
					PortName: "eth3",
					FrameSourceProvider: &stubs.FrameSourceProviderStub{
						MaxJitter:       20 * time.Second,
						SourceMACs:      []string{"CCC"},
						DestinationMACs: []string{"AAA", "BBB"},
					},
					FrameTransmitterProvider: &stubs.FrameTransmitterProviderStub{},
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
