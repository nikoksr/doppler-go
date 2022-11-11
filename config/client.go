package config

import (
	"context"
	"net/http"

	"github.com/nikoksr/doppler-go"
)

// Client is the client used to invoke /v3/configs APIs.
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

func (c Client) get(ctx context.Context, opts *doppler.ConfigGetOptions) (*doppler.Config, doppler.APIResponse, error) {
	var resp doppler.ConfigGetResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/configs/config",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Config, resp.APIResponse, err
}

// Get returns a config and its respective info.
func (c Client) Get(ctx context.Context, opts *doppler.ConfigGetOptions) (*doppler.Config, doppler.APIResponse, error) {
	return c.get(ctx, opts)
}

// Get returns a config and its respective info using the default client.
func Get(ctx context.Context, opts *doppler.ConfigGetOptions) (*doppler.Config, doppler.APIResponse, error) {
	return Default().Get(ctx, opts)
}

func (c Client) fetchList(ctx context.Context, opts *doppler.ConfigListOptions) ([]*doppler.Config, doppler.APIResponse, error) {
	// Make the request.
	var resp doppler.ConfigListResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/configs",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Configs, resp.APIResponse, err
}

// List returns a list of configs.
func (c Client) List(ctx context.Context, opts *doppler.ConfigListOptions) ([]*doppler.Config, doppler.APIResponse, error) {
	return c.fetchList(ctx, opts)
}

// List returns a list of configs using the default client.
func List(ctx context.Context, opts *doppler.ConfigListOptions) ([]*doppler.Config, doppler.APIResponse, error) {
	return Default().List(ctx, opts)
}

func (c Client) create(ctx context.Context, opts *doppler.ConfigCreateOptions) (*doppler.Config, doppler.APIResponse, error) {
	var resp doppler.ConfigCreateResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/configs",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Config, resp.APIResponse, err
}

// Create creates a new config.
func (c Client) Create(ctx context.Context, opts *doppler.ConfigCreateOptions) (*doppler.Config, doppler.APIResponse, error) {
	return c.create(ctx, opts)
}

// Create creates a new config using the default client.
func Create(ctx context.Context, opts *doppler.ConfigCreateOptions) (*doppler.Config, doppler.APIResponse, error) {
	return Default().Create(ctx, opts)
}

func (c Client) update(ctx context.Context, opts *doppler.ConfigUpdateOptions) (*doppler.Config, doppler.APIResponse, error) {
	var resp doppler.ConfigUpdateResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/configs/config",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Config, resp.APIResponse, err
}

// Update updates an existing config.
func (c Client) Update(ctx context.Context, config *doppler.ConfigUpdateOptions) (*doppler.Config, doppler.APIResponse, error) {
	return c.update(ctx, config)
}

// Update updates an existing config using the default client.
func Update(ctx context.Context, config *doppler.ConfigUpdateOptions) (*doppler.Config, doppler.APIResponse, error) {
	return Default().Update(ctx, config)
}

func (c Client) delete(ctx context.Context, opts *doppler.ConfigDeleteOptions) (doppler.APIResponse, error) {
	var resp doppler.ConfigDeleteResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodDelete,
		Path:    "/v3/configs/config",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.APIResponse, err
}

// Delete deletes an existing config.
func (c Client) Delete(ctx context.Context, opts *doppler.ConfigDeleteOptions) (doppler.APIResponse, error) {
	return c.delete(ctx, opts)
}

// Delete deletes an existing config using the default client.
func Delete(ctx context.Context, opts *doppler.ConfigDeleteOptions) (doppler.APIResponse, error) {
	return Default().Delete(ctx, opts)
}

func (c Client) lock(ctx context.Context, opts *doppler.ConfigLockOptions) (*doppler.Config, doppler.APIResponse, error) {
	var resp doppler.ConfigLockResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/configs/config/lock",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Config, resp.APIResponse, err
}

// Lock locks a config.
func (c Client) Lock(ctx context.Context, opts *doppler.ConfigLockOptions) (*doppler.Config, doppler.APIResponse, error) {
	return c.lock(ctx, opts)
}

// Lock locks a config using the default client.
func Lock(ctx context.Context, opts *doppler.ConfigLockOptions) (*doppler.Config, doppler.APIResponse, error) {
	return Default().Lock(ctx, opts)
}

func (c Client) unlock(ctx context.Context, opts *doppler.ConfigUnlockOptions) (*doppler.Config, doppler.APIResponse, error) {
	var resp doppler.ConfigUnlockResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/configs/config/unlock",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Config, resp.APIResponse, err
}

// Unlock unlocks a config.
func (c Client) Unlock(ctx context.Context, opts *doppler.ConfigUnlockOptions) (*doppler.Config, doppler.APIResponse, error) {
	return c.unlock(ctx, opts)
}

// Unlock unlocks a config using the default client.
func Unlock(ctx context.Context, opts *doppler.ConfigUnlockOptions) (*doppler.Config, doppler.APIResponse, error) {
	return Default().Unlock(ctx, opts)
}

func (c Client) clone(ctx context.Context, opts *doppler.ConfigCloneOptions) (*doppler.Config, doppler.APIResponse, error) {
	var resp doppler.ConfigCloneResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/configs/config/clone",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Config, resp.APIResponse, err
}

// Clone clones a config.
func (c Client) Clone(ctx context.Context, opts *doppler.ConfigCloneOptions) (*doppler.Config, doppler.APIResponse, error) {
	return c.clone(ctx, opts)
}

// Clone clones a config using the default client.
func Clone(ctx context.Context, opts *doppler.ConfigCloneOptions) (*doppler.Config, doppler.APIResponse, error) {
	return Default().Clone(ctx, opts)
}
