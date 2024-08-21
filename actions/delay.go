package actions

import (
	"context"
	"time"

	"github.com/lestrrat-go/scriptor/scene"
)

type delay time.Duration

func (d delay) Execute(ctx context.Context) error {
	select {
	case <-time.After(time.Duration(d)):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func Delay(d time.Duration) scene.Action {
	return delay(d)
}
