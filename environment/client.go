package environment

import (
	"context"
	"net/http"

	"github.com/nikoksr/doppler-go"
)

// Client is the client used to invoke /v3/environments APIs.
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

func (c Client) get(ctx context.Context, opts *doppler.EnvironmentGetOptions) (*doppler.Environment, doppler.APIResponse, error) {
	var resp doppler.EnvironmentGetResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/environments/environment",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Environment, resp.APIResponse, err
}

// Get returns an environment and its respective info.
func (c Client) Get(ctx context.Context, opts *doppler.EnvironmentGetOptions) (*doppler.Environment, doppler.APIResponse, error) {
	return c.get(ctx, opts)
}

// Get returns an environment and its respective info using the default client.
func Get(ctx context.Context, opts *doppler.EnvironmentGetOptions) (*doppler.Environment, doppler.APIResponse, error) {
	return Default().Get(ctx, opts)
}

func (c Client) fetchList(ctx context.Context, opts *doppler.EnvironmentListOptions) ([]*doppler.Environment, doppler.APIResponse, error) {
	// Make the request.
	var resp doppler.EnvironmentListResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/environments",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Environments, resp.APIResponse, err
}

// List returns a list of environments.
func (c Client) List(ctx context.Context, opts *doppler.EnvironmentListOptions) ([]*doppler.Environment, doppler.APIResponse, error) {
	return c.fetchList(ctx, opts)
}

// List returns a list of environments using the default client.
func List(ctx context.Context, opts *doppler.EnvironmentListOptions) ([]*doppler.Environment, doppler.APIResponse, error) {
	return Default().List(ctx, opts)
}

func (c Client) create(ctx context.Context, opts *doppler.EnvironmentCreateOptions) (*doppler.Environment, doppler.APIResponse, error) {
	var resp doppler.EnvironmentCreateResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/environments",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Environment, resp.APIResponse, err
}

// Create creates a new environment.
func (c Client) Create(ctx context.Context, opts *doppler.EnvironmentCreateOptions) (*doppler.Environment, doppler.APIResponse, error) {
	return c.create(ctx, opts)
}

// Create creates a new environment using the default client.
func Create(ctx context.Context, opts *doppler.EnvironmentCreateOptions) (*doppler.Environment, doppler.APIResponse, error) {
	return Default().Create(ctx, opts)
}

func (c Client) rename(ctx context.Context, opts *doppler.EnvironmentRenameOptions) (*doppler.Environment, doppler.APIResponse, error) {
	var resp doppler.EnvironmentRenameResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPut,
		Path:    "/v3/environments/environment",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Environment, resp.APIResponse, err
}

// Rename renames an existing environment.
func (c Client) Rename(ctx context.Context, environment *doppler.EnvironmentRenameOptions) (*doppler.Environment, doppler.APIResponse, error) {
	return c.rename(ctx, environment)
}

// Rename renames an existing environment using the default client.
func Rename(ctx context.Context, environment *doppler.EnvironmentRenameOptions) (*doppler.Environment, doppler.APIResponse, error) {
	return Default().Rename(ctx, environment)
}

func (c Client) delete(ctx context.Context, opts *doppler.EnvironmentDeleteOptions) (doppler.APIResponse, error) {
	var resp doppler.EnvironmentDeleteResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodDelete,
		Path:    "/v3/environments/environment",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.APIResponse, err
}

// Delete deletes an existing environment.
func (c Client) Delete(ctx context.Context, opts *doppler.EnvironmentDeleteOptions) (doppler.APIResponse, error) {
	return c.delete(ctx, opts)
}

// Delete deletes an existing environment using the default client.
func Delete(ctx context.Context, opts *doppler.EnvironmentDeleteOptions) (doppler.APIResponse, error) {
	return Default().Delete(ctx, opts)
}
