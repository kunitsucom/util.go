package signalz

import (
	"context"
	"os"
	"os/signal"
)

func Notify(c chan os.Signal, sig ...os.Signal) chan os.Signal {
	signal.Notify(c, sig...)

	return c
}

type ContextSignalHandler func(signal os.Signal, stop context.CancelFunc)

func DefaultContextSignalHandler(_ os.Signal, stop context.CancelFunc) {
	stop()
}

func NotifyContext(parent context.Context, handler ContextSignalHandler, signals ...os.Signal) (ctx context.Context, stop context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	if ctx.Err() == nil {
		go func() {
			for {
				select {
				case sig := <-ch:
					if handler == nil {
						handler = DefaultContextSignalHandler
					}
					handler(sig, cancel)
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	return ctx, cancel
}
