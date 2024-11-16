package stubs

import (
	"fmt"
	"math/rand"
	"project/pkg"
	"time"
)

func SleepRandomly(maxDuration time.Duration) {
	randomDuration := time.Duration(rand.Intn(int(maxDuration)))

	time.Sleep(randomDuration)
}

func ChooseRandomly[T any](s []T) T {
	randomIndex := rand.Intn(len(s))
	return s[randomIndex]
}

func FrameToString(f pkg.Frame) string {
	return fmt.Sprintf("[Source: '%s', Destination: '%s']", f.SourceMAC(), f.DestinationMAC())
}
