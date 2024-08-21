package log

import (
	"context"
	"log/slog"
)

type logKey struct{}

func FromContext(ctx context.Context) *slog.Logger {
	v := ctx.Value(logKey{})
	if v == nil {
		return nil
	}
	if l, ok := v.(*slog.Logger); ok {
		return l
	}
	return nil
}

func InjectContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, logKey{}, l)
}
