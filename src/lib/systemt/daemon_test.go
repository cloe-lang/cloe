package systemt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDaemonize(t *testing.T) {
	Daemonize(func() {})
	assert.Equal(t, 1, len(ds))
}

func TestRunDaemons(t *testing.T) {
	c := make(chan bool, 64)

	go RunDaemons()

	for i := 0; i < 10; i++ {
		Daemonize(func() {
			if i%2 == 0 {
				c <- true
			} else {
				c <- false
			}
		})
	}

	assert.Equal(t, 10, len(c))
}
