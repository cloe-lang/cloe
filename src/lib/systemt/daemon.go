package systemt

const maxConcurrency = 256
const daemonChannelCapacity = 1024

var sem = make(chan bool, maxConcurrency)
var ds = make(chan func(), daemonChannelCapacity)

// Daemonize daemonize a function running it in a goroutine.
func Daemonize(f func()) {
	ds <- f
}

// RunDaemons runs daemons.
func RunDaemons() {
	for d := range ds {
		sem <- true
		go func(d func()) {
			d()
			<-sem
		}(d)
	}
}
