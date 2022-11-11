package workplace

import (
	"context"
	"net/http"

	"github.com/nikoksr/doppler-go"
)

// Client is the client used to invoke /v3/workplace APIs.
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

func (c Client) get(ctx context.Context) (*doppler.Workplace, doppler.APIResponse, error) {
	var resp doppler.WorkplaceGetResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method: http.MethodGet,
		Path:   "/v3/workplace",
		Key:    c.Key,
	}, &resp)

	return resp.Workplace, resp.APIResponse, err
}

// Get returns an account and its respective info.
func (c Client) Get(ctx context.Context) (*doppler.Workplace, doppler.APIResponse, error) {
	return c.get(ctx)
}

// Get returns an account and its respective info using the default client.
func Get(ctx context.Context) (*doppler.Workplace, doppler.APIResponse, error) {
	return Default().Get(ctx)
}

func (c Client) update(ctx context.Context, opts *doppler.WorkplaceUpdateOptions) (*doppler.Workplace, doppler.APIResponse, error) {
	var resp doppler.WorkplaceUpdateResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/workplace",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Workplace, resp.APIResponse, err
}

// Update updates an existing workplace.
func (c Client) Update(ctx context.Context, workplace *doppler.WorkplaceUpdateOptions) (*doppler.Workplace, doppler.APIResponse, error) {
	return c.update(ctx, workplace)
}

// Update updates an existing workplace using the default client.
func Update(ctx context.Context, workplace *doppler.WorkplaceUpdateOptions) (*doppler.Workplace, doppler.APIResponse, error) {
	return Default().Update(ctx, workplace)
}
