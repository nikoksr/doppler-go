package servicetoken_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nikoksr/doppler-go"
	"github.com/nikoksr/doppler-go/pointer"
	servicetoken "github.com/nikoksr/doppler-go/service_token"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := servicetoken.Default()
	if client == nil {
		t.Fatal("Expected client to be set")
	}
	if client.Backend == nil {
		t.Fatal("Expected client backend to be set")
	}
	if client.Key != doppler.Key {
		t.Fatalf("Expected client key to be %q, got %q", doppler.Key, client.Key)
	}
}

func TestServiceToken_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		options           *doppler.ServiceTokenListOptions
		wantServiceTokens []*doppler.ServiceToken
		wantResponse      doppler.APIResponse
		wantErr           bool
	}{
		{
			name: "List all service tokens",
			options: &doppler.ServiceTokenListOptions{
				Project: "test",
				Config:  "test",
			},
			wantServiceTokens: []*doppler.ServiceToken{
				{
					Name:        pointer.To("test"),
					Slug:        pointer.To("test"),
					Key:         pointer.To("test"),
					Project:     pointer.To("test"),
					Environment: pointer.To("test"),
					Config:      pointer.To("test"),
					Access:      pointer.To("test"),
					ExpiresAt:   pointer.To("test"),
					CreatedAt:   pointer.To("test"),
				},
				{
					Name:        pointer.To("test2"),
					Slug:        pointer.To("test2"),
					Key:         pointer.To("test2"),
					Project:     pointer.To("test2"),
					Environment: pointer.To("test2"),
					Config:      pointer.To("test2"),
					Access:      pointer.To("test2"),
					ExpiresAt:   pointer.To("test2"),
					CreatedAt:   pointer.To("test2"),
				},
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "List service tokens specific config",
			options: &doppler.ServiceTokenListOptions{
				Project: "test",
				Config:  "unknown",
			},
			wantServiceTokens: []*doppler.ServiceToken{
				{
					Name:        pointer.To("test"),
					Slug:        pointer.To("test"),
					Key:         pointer.To("test"),
					Project:     pointer.To("test"),
					Environment: pointer.To("test"),
					Config:      pointer.To("test"),
					Access:      pointer.To("test"),
					ExpiresAt:   pointer.To("test"),
					CreatedAt:   pointer.To("test"),
				},
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "List service tokens error",
			options: &doppler.ServiceTokenListOptions{
				Project: "test",
				Config:  "test",
			},
			wantServiceTokens: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"Config not found"},
			},
			wantErr: true,
		},
		{
			name:              "List service tokens options validation error",
			options:           &doppler.ServiceTokenListOptions{},
			wantServiceTokens: nil,
			wantResponse:      doppler.APIResponse{},
			wantErr:           true,
		},
	}

	// Run tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a new httptest.Server that will be used to mock the Doppler API.
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Set the expected response headers.
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.wantResponse.StatusCode)
				// Write the expected response.
				err := json.NewEncoder(w).Encode(&doppler.ServiceTokenListResponse{
					APIResponse: tt.wantResponse,
					Tokens:      tt.wantServiceTokens,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &servicetoken.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotServiceTokens, gotResponse, err := client.List(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the service tokens are expected.
			if diff := cmp.Diff(tt.wantServiceTokens, gotServiceTokens); diff != "" {
				t.Errorf("Unexpected service tokens (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestServiceToken_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		options          *doppler.ServiceTokenCreateOptions
		wantServiceToken *doppler.ServiceToken
		wantResponse     doppler.APIResponse
		wantErr          bool
	}{
		{
			name: "Create service token",
			options: &doppler.ServiceTokenCreateOptions{
				Project:   "test",
				Config:    "test",
				Name:      "test",
				Access:    pointer.To("read"),
				ExpiresAt: pointer.To("2021-01-01"),
			},
			wantServiceToken: &doppler.ServiceToken{
				Name:        pointer.To("test"),
				Slug:        pointer.To("test"),
				Key:         pointer.To("test"),
				Project:     pointer.To("test"),
				Environment: pointer.To("test"),
				Config:      pointer.To("test"),
				Access:      pointer.To("test"),
				ExpiresAt:   pointer.To("test"),
				CreatedAt:   pointer.To("test"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Create service token error",
			options: &doppler.ServiceTokenCreateOptions{
				Project:   "test",
				Config:    "unknown",
				Name:      "test",
				Access:    pointer.To("read"),
				ExpiresAt: pointer.To("2021-01-01"),
			},
			wantServiceToken: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"Config not found"},
			},
			wantErr: true,
		},
		{
			name: "Create service token options validation error",
			options: &doppler.ServiceTokenCreateOptions{
				Project: "test",
			},
			wantServiceToken: nil,
			wantResponse:     doppler.APIResponse{},
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a new httptest.Server that will be used to mock the Doppler API.
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Set the expected response headers.
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.wantResponse.StatusCode)
				// Write the expected response.
				err := json.NewEncoder(w).Encode(&doppler.ServiceTokenCreateResponse{
					APIResponse: tt.wantResponse,
					Token:       tt.wantServiceToken,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &servicetoken.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotServiceToken, gotResponse, err := client.Create(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the service token is expected.
			if diff := cmp.Diff(tt.wantServiceToken, gotServiceToken); diff != "" {
				t.Errorf("Unexpected service token (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestServiceToken_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.ServiceTokenDeleteOptions
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Delete service token",
			options: &doppler.ServiceTokenDeleteOptions{
				Project: "test",
				Config:  "test",
				Slug:    "test",
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Delete service token error",
			options: &doppler.ServiceTokenDeleteOptions{
				Project: "test",
				Config:  "unknown",
				Slug:    "test",
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"Config not found"},
			},
			wantErr: true,
		},
		{
			name: "Delete service token options validation error",
			options: &doppler.ServiceTokenDeleteOptions{
				Project: "test",
			},
			wantResponse: doppler.APIResponse{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a new httptest.Server that will be used to mock the Doppler API.
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Set the expected response headers.
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.wantResponse.StatusCode)
				// Write the expected response.
				err := json.NewEncoder(w).Encode(&doppler.ServiceTokenDeleteResponse{
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &servicetoken.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotResponse, err := client.Delete(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}
