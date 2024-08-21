package stash

import (
	"context"
	"errors"

	"github.com/lestrrat-go/scriptor/ctxutil"
)

type Stash interface {
	Set(any, any) Stash
	Get(any) (any, bool)
}

type stashKey struct{}

func FromContext(ctx context.Context) Stash {
	var dst Stash
	if ctxutil.FromContext[Stash](ctx, stashKey{}, &dst) {
		return dst
	}
	return nil
}

func InjectContext(ctx context.Context, s Stash) context.Context {
	return ctxutil.InjectContext[Stash](ctx, stashKey{}, s)
}

type stash struct {
	data map[any]any
}

func New() Stash {
	return &stash{
		data: make(map[any]any),
	}
}

func (s *stash) Set(k, v any) Stash {
	s.data[k] = v
	return s
}

func (s *stash) Get(k any) (any, bool) {
	v, ok := s.data[k]
	return v, ok
}

func Set(ctx context.Context, key any, value any) error {
	st := FromContext(ctx)
	if st == nil {
		return errors.New("stash.Set: no stash found in context")
	}
	st.Set(key, value)
	return nil
}

func Fetch[T any](ctx context.Context, key any, dst *T) bool {
	st := FromContext(ctx)
	if st == nil {
		return false
	}

	v, ok := st.Get(key)
	if !ok {
		return false
	}

	if val, ok := v.(T); ok {
		*dst = val
		return true
	}
	return false
}
