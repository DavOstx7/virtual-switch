package main

import (
	"context"
	"project/network/virtual"
	"project/stubs"
	"time"
)

func main() {
	vSwitch := virtual.NewSwitch(
		[]*virtual.Port{
			virtual.NewPort(
				&virtual.PortConfig{
					PortName: "eth1",
					FrameSourceProvider: &stubs.FrameSourceProvider{
						MaxJitter:       5 * time.Second,
						SourceMACs:      []string{"AAA"},
						DestinationMACs: []string{"BBB", "CCC"},
					},
					FrameTransmitterProvider: &stubs.FrameTransmitterProvider{},
				},
			),
			virtual.NewPort(
				&virtual.PortConfig{
					PortName: "eth2",
					FrameSourceProvider: &stubs.FrameSourceProvider{
						MaxJitter:       10 * time.Second,
						SourceMACs:      []string{"BBB"},
						DestinationMACs: []string{"AAA", "CCC"},
					},
					FrameTransmitterProvider: &stubs.FrameTransmitterProvider{},
				},
			),
			virtual.NewPort(
				&virtual.PortConfig{
					PortName: "eth3",
					FrameSourceProvider: &stubs.FrameSourceProvider{
						MaxJitter:       20 * time.Second,
						SourceMACs:      []string{"CCC"},
						DestinationMACs: []string{"AAA", "BBB"},
					},
					FrameTransmitterProvider: &stubs.FrameTransmitterProvider{},
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
