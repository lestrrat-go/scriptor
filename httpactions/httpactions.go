package httpactions

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lestrrat-go/scriptor/log"
	"github.com/lestrrat-go/scriptor/scene"
	"github.com/lestrrat-go/scriptor/stash"
)

type prevRequestKey struct{}
type prevResponseKey struct{}

func PrevRequest(ctx context.Context) *http.Request {
	st := stash.FromContext(ctx)
	if st == nil {
		return nil
	}

	v, ok := st.Get(prevRequestKey{})
	if !ok {
		return nil
	}

	if req, ok := v.(*http.Request); ok {
		return req
	}
	return nil
}

func PrevResponse(ctx context.Context) *http.Response {
	st := stash.FromContext(ctx)
	if st == nil {
		return nil
	}

	v, ok := st.Get(prevResponseKey{})
	if !ok {
		return nil
	}

	if res, ok := v.(*http.Response); ok {
		return res
	}
	return nil
}

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{client: client}
}

func (h *Client) GetAction(u string) scene.Action {
	r, err := http.NewRequest(http.MethodGet, u, nil)
	return &httpAction{client: h.client, req: r, err: err, name: "GetAction"}
}

type httpAction struct {
	client *http.Client
	req    *http.Request
	err    error
	name   string
}

func (h *httpAction) Execute(ctx context.Context) error {
	logger := log.FromContext(ctx)
	if err := h.err; err != nil {
		logger.LogAttrs(ctx, slog.LevelInfo, "actions.httpAction", slog.String("error", err.Error()))
		return fmt.Errorf("actions.httpAction: error during setup: %w", err)
	}

	req := h.req.WithContext(ctx)
	//nolint:bodyclose
	res, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("actions.httpAction: error during request: %w", err)
	}

	st := stash.FromContext(ctx)

	logger.DebugContext(ctx, "actions.httpAction", slog.Any("request", req))
	logger.DebugContext(ctx, "actions.httpAction", slog.Any("response", res))
	st.Set(prevRequestKey{}, req)
	st.Set(prevResponseKey{}, res)
	return nil
}