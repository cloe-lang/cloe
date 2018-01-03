package systemt

const daemonChannelCapacity = 1024

var ds = make(chan func(), daemonChannelCapacity)

// Daemonize daemonizes a function running it in a goroutine.
func Daemonize(f func()) {
	ds <- f
}

// RunDaemons runs daemons.
// The current implementation doesn't limit a number of spawn goroutines
// because it can lead to dead lock. However, it should be done to keep
// reference locality for cache in one way or another.
func RunDaemons() {
	for d := range ds {
		go func(d func()) {
			d()
		}(d)
	}
}
