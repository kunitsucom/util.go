package cliz

import "context"

type contextKeyCommand struct{}

func WithContext(ctx context.Context, cmd *Command) context.Context {
	return context.WithValue(ctx, contextKeyCommand{}, cmd)
}

func FromContext(ctx context.Context) (*Command, error) {
	c, ok := ctx.Value(contextKeyCommand{}).(*Command)
	if !ok {
		return nil, ErrCommandNotSetInContext
	}

	return c, nil
}
