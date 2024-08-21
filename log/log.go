package log

import (
	"context"
	"log/slog"

	"github.com/lestrrat-go/scriptor/ctxutil"
)

type logKey struct{}

func InjectContext(ctx context.Context, l *slog.Logger) context.Context {
	return ctxutil.InjectContext[*slog.Logger](ctx, logKey{}, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	var dst *slog.Logger
	if ctxutil.FromContext[*slog.Logger](ctx, logKey{}, &dst) {
		return dst
	}
	return nil
}
