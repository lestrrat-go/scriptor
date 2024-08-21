package clock

import (
	"context"
	"time"

	"github.com/lestrrat-go/scriptor/ctxutil"
)

// Clock is an interface that provides the current time.
type Clock interface {
	Now() time.Time
}

type static time.Time

func Static(t time.Time) Clock {
	return static(t)
}

func (s static) Now() time.Time {
	return time.Time(s)
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

func RealClock() Clock {
	return realClock{}
}

type clockKey struct{}

func InjectContext(ctx context.Context, v Clock) context.Context {
	return ctxutil.InjectContext[Clock](ctx, clockKey{}, v)
}

// FromContext returns the Clock from the context, if any.
// Make sure to check the return value for nil.
func FromContext(ctx context.Context) Clock {
	var dst Clock
	if ctxutil.FromContext[Clock](ctx, clockKey{}, &dst) {
		return dst
	}
	return nil
}
