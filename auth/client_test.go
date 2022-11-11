package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nikoksr/doppler-go"
	"github.com/nikoksr/doppler-go/auth"
	"github.com/nikoksr/doppler-go/pointer"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := auth.Default()
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

func TestAuth_Revoke(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.AuthRevokeOptions
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Revoke one token",
			options: &doppler.AuthRevokeOptions{
				Tokens: []doppler.AuthToken{
					{Token: pointer.To("token1")},
				},
			},
			wantResponse: doppler.APIResponse{Status: "200 OK", StatusCode: http.StatusOK},
			wantErr:      false,
		},
		{
			name: "Revoke multiple tokens",
			options: &doppler.AuthRevokeOptions{
				Tokens: []doppler.AuthToken{
					{Token: pointer.To("token1")},
					{Token: pointer.To("token2")},
				},
			},
			wantResponse: doppler.APIResponse{Status: "200 OK", StatusCode: http.StatusOK},
			wantErr:      false,
		},
		{
			name: "Unknown token",
			options: &doppler.AuthRevokeOptions{
				Tokens: []doppler.AuthToken{
					{Token: pointer.To("unknown")},
				},
			},
			wantResponse: doppler.APIResponse{Status: "200 OK", StatusCode: http.StatusOK},
			wantErr:      false,
		},
		{
			name:         "Revoke failed with invalid options error",
			options:      &doppler.AuthRevokeOptions{},
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
				err := json.NewEncoder(w).Encode(&doppler.AuthRevokeResponse{
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &auth.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Delete method with the test auth.
			gotResponse, err := client.Revoke(context.Background(), tt.options)
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
