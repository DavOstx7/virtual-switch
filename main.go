package main

import (
	"context"
	"project/network"
	"project/stubs"
	"time"
)

func main() {
	vSwitch := network.NewSwitch(
		[]*network.Port{
			network.NewPort(
				&network.PortConfig{
					PortName: "eth1",
					FrameSourceProvider: &stubs.FrameSourceProvider{
						MaxJitter:       5 * time.Second,
						SourceMACs:      []string{"AAA"},
						DestinationMACs: []string{"BBB", "CCC"},
					},
					FrameTransmitterProvider: &stubs.FrameTransmitterProvider{},
				},
			),
			network.NewPort(
				&network.PortConfig{
					PortName: "eth2",
					FrameSourceProvider: &stubs.FrameSourceProvider{
						MaxJitter:       10 * time.Second,
						SourceMACs:      []string{"BBB"},
						DestinationMACs: []string{"AAA", "CCC"},
					},
					FrameTransmitterProvider: &stubs.FrameTransmitterProvider{},
				},
			),
			network.NewPort(
				&network.PortConfig{
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
