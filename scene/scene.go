package scene

import (
	"context"
	"errors"
	"fmt"
	"os"
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
	fmt.Fprintf(os.Stderr, "Add called with %T\n", a)
	s.actions = append(s.actions, a)
	return s
}

func (s *Scene) Execute(ctx context.Context) error {
	ctx = InjectContext(ctx, s)

	//nolint:intrange // this is on purpose
	for i := 0; i < len(s.actions); i++ {
		a := s.actions[i]
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
