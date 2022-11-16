package contextz

import (
	"context"
	"errors"
	"fmt"
	"os"
)

var ErrValueNotSet = errors.New("contextz: context value not set")

type key int

// nolint: gochecknoglobals
const (
	_ key = iota
	signalChannelKey
)

func WithSignalChannel(parent context.Context, signalChannel chan os.Signal) context.Context {
	return context.WithValue(parent, signalChannelKey, signalChannel)
}

func MustSignalChannel(ctx context.Context) (signalChannel chan os.Signal) {
	ch, ok := ctx.Value(signalChannelKey).(chan os.Signal)
	if !ok {
		panic(fmt.Errorf("MustSignalChannel: %w", ErrValueNotSet))
	}

	return ch
}
