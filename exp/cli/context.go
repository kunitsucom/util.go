package cliz

import "context"

type contextKeyCommand struct{}

func WithContext(ctx context.Context, cmd *Command) context.Context {
	return context.WithValue(ctx, contextKeyCommand{}, cmd)
}

func FromContext(ctx context.Context) (cmd *Command, ok bool) {
	c, ok := ctx.Value(contextKeyCommand{}).(*Command)
	return c, ok
}
