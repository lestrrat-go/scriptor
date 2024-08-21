package scene

import "context"

type Action interface {
	Execute(context.Context) error
}

type ActionFunc func(context.Context) error

func (f ActionFunc) Execute(ctx context.Context) error {
	return f(ctx)
}

type repeater struct {
	action Action
	count  int
}

func Repeat(a Action, n int) Action {
	return &repeater{action: a, count: n}
}

func (r *repeater) Execute(ctx context.Context) error {
	for range r.count {
		if err := r.action.Execute(ctx); err != nil {
			return err
		}
	}
	return nil
}
