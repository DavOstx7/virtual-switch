package main

import (
	"context"
	"project/pkg"
	"project/stubs"
)

func main() {
	vSwitch := pkg.NewVirtualSwitch(
		[]*pkg.VirtualPort{
			pkg.NewVirtualPort(
				&pkg.VirtualPortConfig{
					PortName: "eth1",
					FrameSourceProvider: &stubs.FrameSourceProviderStub{
						MaxJitterSeconds: 5,
						SourceMAC:        "AAA",
						DestinationMACs:  []string{"BBB", "CCC"},
					},
					FrameTransmitterProvider: &stubs.FrameTransmitterProviderStub{},
				},
			),
			pkg.NewVirtualPort(
				&pkg.VirtualPortConfig{
					PortName: "eth2",
					FrameSourceProvider: &stubs.FrameSourceProviderStub{
						MaxJitterSeconds: 10,
						SourceMAC:        "BBB",
						DestinationMACs:  []string{"AAA", "CCC"},
					},
					FrameTransmitterProvider: &stubs.FrameTransmitterProviderStub{},
				},
			),
			pkg.NewVirtualPort(
				&pkg.VirtualPortConfig{
					PortName: "eth3",
					FrameSourceProvider: &stubs.FrameSourceProviderStub{
						MaxJitterSeconds: 20,
						SourceMAC:        "CCC",
						DestinationMACs:  []string{"AAA", "BBB"},
					},
					FrameTransmitterProvider: &stubs.FrameTransmitterProviderStub{},
				},
			),
		},
	)

	ctx := context.Background()
	vSwitch.On(ctx)
	vSwitch.Wait()
}
