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
	var dst *http.Request
	if stash.Fetch[*http.Request](ctx, prevRequestKey{}, &dst) {
		return dst
	}
	return nil
}

func PrevResponse(ctx context.Context) *http.Response {
	var dst *http.Response
	if stash.Fetch[*http.Response](ctx, prevResponseKey{}, &dst) {
		return dst
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

	logger.DebugContext(ctx, "actions.httpAction", slog.Any("request", req))
	logger.DebugContext(ctx, "actions.httpAction", slog.Any("response", res))
	stash.Set(ctx, prevRequestKey{}, req)
	stash.Set(ctx, prevResponseKey{}, res)
	return nil
}
