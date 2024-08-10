package cliz

import (
	"context"
	"fmt"
)

type contextKeyCommand struct{}

func WithContext(ctx context.Context, cmd *Command) context.Context {
	return context.WithValue(ctx, contextKeyCommand{}, cmd)
}

func FromContext(ctx context.Context) (*Command, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	c, ok := ctx.Value(contextKeyCommand{}).(*Command)
	if !ok {
		return nil, ErrCommandNotSetInContext
	}

	return c, nil
}

func MustFromContext(ctx context.Context) *Command {
	c, err := FromContext(ctx)
	if err != nil {
		err = fmt.Errorf("FromContext: %w", err)
		panic(err)
	}

	return c
}
