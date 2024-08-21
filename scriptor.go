package scriptor

import (
	"context"
	"log/slog"
	"os"

	"github.com/lestrrat-go/scriptor/log"
	"github.com/lestrrat-go/scriptor/stash"
)

func DefaultContext(ctx context.Context) context.Context {
	newctx := stash.InjectContext(
		log.InjectContext(ctx, slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))),
		stash.New(),
	)
	return newctx
}
