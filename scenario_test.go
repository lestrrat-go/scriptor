package scriptor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lestrrat-go/scriptor"
	"github.com/lestrrat-go/scriptor/actions"
	"github.com/lestrrat-go/scriptor/clock"
	"github.com/lestrrat-go/scriptor/httpactions"
	"github.com/lestrrat-go/scriptor/scene"
	"github.com/stretchr/testify/require"
)

func TestStaticClock(t *testing.T) {
	var observed []time.Time
	action := scene.ActionFunc(func(ctx context.Context) error {
		cl := clock.FromContext(ctx)
		require.NotNil(t, cl, `clock.FromContext(ctx) should not return nil`)

		observed = append(observed, cl.Now())
		return nil
	})
	s := scene.New().
		Add(action).
		Add(actions.Delay(time.Millisecond * 100)).
		Add(action).
		Add(actions.Delay(time.Millisecond * 100)).
		Add(action)

	now := time.Now()
	s.Execute(clock.InjectContext(context.Background(), clock.Static(now)))

	require.Equal(t, 3, len(observed), `expected 3 actions to be executed`)
	for _, o := range observed {
		require.Equal(t, now, o, `expected all actions to be executed at the same time`)
	}
}

func TestRealClock(t *testing.T) {
	var observed []time.Time
	action := scene.ActionFunc(func(ctx context.Context) error {
		cl := clock.FromContext(ctx)
		require.NotNil(t, cl, `clock.FromContext(ctx) should not return nil`)

		observed = append(observed, cl.Now())
		return nil
	})
	s := scene.New().
		Add(action).
		Add(actions.Delay(time.Millisecond * 100)).
		Add(action).
		Add(actions.Delay(time.Millisecond * 100)).
		Add(action)

	s.Execute(clock.InjectContext(context.Background(), clock.RealClock()))

	require.Equal(t, 3, len(observed), `expected 3 actions to be executed`)
	for i := 1; i < len(observed); i++ {
		require.True(t, observed[i-1].Before(observed[i]), `expected actions to be executed in order`)
	}
}

func TestHTTP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	client := httpactions.NewClient(http.DefaultClient)

	s := scene.New().
		Add(client.GetAction(srv.URL)).
		Add(scene.ActionFunc(func(ctx context.Context) error {
			res := httpactions.PrevResponse(ctx)
			defer res.Body.Close()
			require.NotNil(t, res, `expected PrevResponse(ctx) to return a non-nil response`)
			require.Equal(t, http.StatusOK, res.StatusCode, `expected status code to be 200`)
			return nil
		}))

	s.Execute(scriptor.DefaultContext(context.Background()))
}
