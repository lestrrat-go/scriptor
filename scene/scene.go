package scene

import (
	"context"
	"errors"
)

var ErrEndOfScene = errors.New("end of scene")

type Scene struct {
	actions []Action
}

func New() *Scene {
	return &Scene{}
}

func (s *Scene) Add(a Action) *Scene {
	s.actions = append(s.actions, a)
	return s
}

func (s *Scene) Execute(ctx context.Context) error {
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
