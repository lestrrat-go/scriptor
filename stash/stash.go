package stash

import "context"

type Stash interface {
	Set(any, any) Stash
	Get(any) (any, bool)
}

type stashKey struct{}

func FromContext(ctx context.Context) Stash {
	v := ctx.Value(stashKey{})
	if v == nil {
		return nil
	}
	if s, ok := v.(Stash); ok {
		return s
	}
	return nil
}

func InjectContext(ctx context.Context, s Stash) context.Context {
	return context.WithValue(ctx, stashKey{}, s)
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
