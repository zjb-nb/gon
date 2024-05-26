package gonweb

import (
	"os"
	"syscall"
)

var ShutDownSignals []os.Signal = []os.Signal{
	syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGILL,
	syscall.SIGTERM,
}
