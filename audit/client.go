package audit

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nikoksr/doppler-go"
)

// Client is the client used to invoke /v3/configs/config/secrets APIs.
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

func (c Client) workplaceGet(ctx context.Context, opts *doppler.AuditWorkplaceGetOptions) (*doppler.AuditWorkplace, doppler.APIResponse, error) {
	var resp doppler.AuditWorkplaceGetResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/workplace",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.AuditWorkplace, resp.APIResponse, err
}

// WorkplaceGet returns an audit log for a workplace.
func (c Client) WorkplaceGet(ctx context.Context, opts *doppler.AuditWorkplaceGetOptions) (*doppler.AuditWorkplace, doppler.APIResponse, error) {
	return c.workplaceGet(ctx, opts)
}

// WorkplaceGet returns an audit log for a workplace using the default client.
func WorkplaceGet(ctx context.Context, opts *doppler.AuditWorkplaceGetOptions) (*doppler.AuditWorkplace, doppler.APIResponse, error) {
	return Default().WorkplaceGet(ctx, opts)
}

func (c Client) workplaceUserGet(ctx context.Context, opts *doppler.AuditWorkplaceUserGetOptions) (*doppler.AuditWorkplaceUser, doppler.APIResponse, error) {
	if opts == nil {
		return nil, doppler.APIResponse{}, fmt.Errorf("options may not be nil")
	}

	path := fmt.Sprintf("/v3/workplace/users/%s", opts.UserID)

	var resp doppler.AuditWorkplaceUserGetResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    path,
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.AuditWorkplaceUser, resp.APIResponse, err
}

// WorkplaceUserGet returns an audit log for a workplace user.
func (c Client) WorkplaceUserGet(ctx context.Context, opts *doppler.AuditWorkplaceUserGetOptions) (*doppler.AuditWorkplaceUser, doppler.APIResponse, error) {
	return c.workplaceUserGet(ctx, opts)
}

// WorkplaceUserGet returns an audit log for a workplace user using the default client.
func WorkplaceUserGet(ctx context.Context, opts *doppler.AuditWorkplaceUserGetOptions) (*doppler.AuditWorkplaceUser, doppler.APIResponse, error) {
	return Default().WorkplaceUserGet(ctx, opts)
}

func (c Client) workplaceUserList(ctx context.Context, opts *doppler.AuditWorkplaceUserListOptions) ([]*doppler.AuditWorkplaceUser, doppler.APIResponse, error) {
	var resp doppler.AuditWorkplaceUserListResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/workplace/users",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.AuditWorkplaceUsers, resp.APIResponse, err
}

// WorkplaceUserList returns a list of audit logs for workplace users.
func (c Client) WorkplaceUserList(ctx context.Context, opts *doppler.AuditWorkplaceUserListOptions) ([]*doppler.AuditWorkplaceUser, doppler.APIResponse, error) {
	return c.workplaceUserList(ctx, opts)
}

// WorkplaceUserList returns a list of audit logs for workplace users using the default client.
func WorkplaceUserList(ctx context.Context, opts *doppler.AuditWorkplaceUserListOptions) ([]*doppler.AuditWorkplaceUser, doppler.APIResponse, error) {
	return Default().WorkplaceUserList(ctx, opts)
}
