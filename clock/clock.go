package clock

import (
	"context"
	"time"
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
	return context.WithValue(ctx, clockKey{}, v)
}

// FromContext returns the Clock from the context, if any.
// Make sure to check the return value for nil.
func FromContext(ctx context.Context) Clock {
	v := ctx.Value(clockKey{})
	if v == nil {
		return nil
	}
	if c, ok := v.(Clock); ok {
		return c
	}
	return nil
}
