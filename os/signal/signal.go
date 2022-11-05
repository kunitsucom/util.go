package signalz

import (
	"os"
	"os/signal"
)

func Notify(c chan<- os.Signal, sig ...os.Signal) chan<- os.Signal {
	signal.Notify(c, sig...)

	return c
}
