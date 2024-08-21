package ctxutil

import "context"

func InjectContext[T any](ctx context.Context, key any, value T) context.Context {
	return context.WithValue(ctx, key, value)
}

func FromContext[T any](ctx context.Context, key any, dst *T) bool {
	v := ctx.Value(key)
	if v == nil {
		return false
	}
	if t, ok := v.(T); ok {
		*dst = t
		return true
	}
	return false
}
