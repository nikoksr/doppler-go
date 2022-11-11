package activitylog

import (
	"context"
	"net/http"

	"github.com/nikoksr/doppler-go"
)

// Client is the client used to invoke /v3/logs APIs.
type Client struct {
	Backend doppler.Backend
	Key     string
}

// Default returns a new client based on the SDK's default backend and API key.
func Default() *Client {
	return &Client{
		Backend: doppler.GetBackend(),
		Key:     doppler.Key,
	}
}

func (c Client) get(ctx context.Context, opts *doppler.ActivityLogGetOptions) (*doppler.ActivityLog, doppler.APIResponse, error) {
	var resp doppler.ActivityLogGetResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/logs/log",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.ActivityLog, resp.APIResponse, err
}

// Get returns a activity log and its respective info.
func (c Client) Get(ctx context.Context, opts *doppler.ActivityLogGetOptions) (*doppler.ActivityLog, doppler.APIResponse, error) {
	return c.get(ctx, opts)
}

// Get returns a activity log and its respective info using the default client.
func Get(ctx context.Context, opts *doppler.ActivityLogGetOptions) (*doppler.ActivityLog, doppler.APIResponse, error) {
	return Default().Get(ctx, opts)
}

func (c Client) fetchList(ctx context.Context, opts *doppler.ActivityLogListOptions) ([]*doppler.ActivityLog, doppler.APIResponse, error) {
	// Make the request.
	var resp doppler.ActivityLogListResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/logs",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.ActivityLogs, resp.APIResponse, err
}

// List returns a list of activity logs.
func (c Client) List(ctx context.Context, opts *doppler.ActivityLogListOptions) ([]*doppler.ActivityLog, doppler.APIResponse, error) {
	return c.fetchList(ctx, opts)
}

// List returns a list of activity logs using the default client.
func List(ctx context.Context, opts *doppler.ActivityLogListOptions) ([]*doppler.ActivityLog, doppler.APIResponse, error) {
	return Default().List(ctx, opts)
}
