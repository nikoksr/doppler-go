package project

import (
	"context"
	"net/http"

	"github.com/nikoksr/doppler-go"
)

// Client is the client used to invoke /v3/projects APIs.
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

func (c Client) get(ctx context.Context, opts *doppler.ProjectGetOptions) (*doppler.Project, doppler.APIResponse, error) {
	var resp doppler.ProjectGetResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/projects/project",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Project, resp.APIResponse, err
}

// Get returns a project and its respective info.
func (c Client) Get(ctx context.Context, opts *doppler.ProjectGetOptions) (*doppler.Project, doppler.APIResponse, error) {
	return c.get(ctx, opts)
}

// Get returns a project and its respective info using the default client.
func Get(ctx context.Context, opts *doppler.ProjectGetOptions) (*doppler.Project, doppler.APIResponse, error) {
	return Default().Get(ctx, opts)
}

func (c Client) fetchList(ctx context.Context, opts *doppler.ProjectListOptions) ([]*doppler.Project, doppler.APIResponse, error) {
	// Make the request.
	var resp doppler.ProjectListResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/projects",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Projects, resp.APIResponse, err
}

// List returns a list of projects.
func (c Client) List(ctx context.Context, opts *doppler.ProjectListOptions) ([]*doppler.Project, doppler.APIResponse, error) {
	return c.fetchList(ctx, opts)
}

// List returns a list of projects using the default client.
func List(ctx context.Context, opts *doppler.ProjectListOptions) ([]*doppler.Project, doppler.APIResponse, error) {
	return Default().List(ctx, opts)
}

func (c Client) create(ctx context.Context, opts *doppler.ProjectCreateOptions) (*doppler.Project, doppler.APIResponse, error) {
	var resp doppler.ProjectCreateResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/projects",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Project, resp.APIResponse, err
}

// Create creates a new project.
func (c Client) Create(ctx context.Context, opts *doppler.ProjectCreateOptions) (*doppler.Project, doppler.APIResponse, error) {
	return c.create(ctx, opts)
}

// Create creates a new project using the default client.
func Create(ctx context.Context, opts *doppler.ProjectCreateOptions) (*doppler.Project, doppler.APIResponse, error) {
	return Default().Create(ctx, opts)
}

func (c Client) update(ctx context.Context, opts *doppler.ProjectUpdateOptions) (*doppler.Project, doppler.APIResponse, error) {
	var resp doppler.ProjectUpdateResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/projects/project",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Project, resp.APIResponse, err
}

// Update updates an existing project.
func (c Client) Update(ctx context.Context, project *doppler.ProjectUpdateOptions) (*doppler.Project, doppler.APIResponse, error) {
	return c.update(ctx, project)
}

// Update updates an existing project using the default client.
func Update(ctx context.Context, project *doppler.ProjectUpdateOptions) (*doppler.Project, doppler.APIResponse, error) {
	return Default().Update(ctx, project)
}

func (c Client) delete(ctx context.Context, opts *doppler.ProjectDeleteOptions) (doppler.APIResponse, error) {
	var resp doppler.ProjectDeleteResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodDelete,
		Path:    "/v3/projects/project",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.APIResponse, err
}

// Delete deletes an existing project.
func (c Client) Delete(ctx context.Context, opts *doppler.ProjectDeleteOptions) (doppler.APIResponse, error) {
	return c.delete(ctx, opts)
}

// Delete deletes an existing project using the default client.
func Delete(ctx context.Context, opts *doppler.ProjectDeleteOptions) (doppler.APIResponse, error) {
	return Default().Delete(ctx, opts)
}
