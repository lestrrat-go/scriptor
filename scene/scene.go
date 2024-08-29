package scene

import (
	"context"
	"errors"
)

var ErrEndOfScene = errors.New("end of scene")

type Scene struct {
	actions []Action
}

type sceneKey struct{}

func InjectContext(ctx context.Context, s *Scene) context.Context {
	return context.WithValue(ctx, sceneKey{}, s)
}

func FromContext(ctx context.Context) *Scene {
	v := ctx.Value(sceneKey{})
	if v == nil {
		return nil
	}
	s, ok := v.(*Scene)
	if !ok {
		return nil
	}
	return s
}

func New() *Scene {
	return &Scene{}
}

func (s *Scene) Add(a Action) *Scene {
	s.actions = append(s.actions, a)
	return s
}

func (s *Scene) Execute(ctx context.Context) error {
	ctx = InjectContext(ctx, s)

	for _, a := range s.actions {
		if err := a.Execute(ctx); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
	return nil
}
