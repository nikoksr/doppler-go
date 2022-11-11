package dynamicsecret_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nikoksr/doppler-go"
	dynamicsecret "github.com/nikoksr/doppler-go/dynamic_secret"
	"github.com/nikoksr/doppler-go/pointer"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := dynamicsecret.Default()
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

func TestDynamicSecret_IssueLease(t *testing.T) {
	t.Parallel()

	// Table driven tests that mock the Doppler API using httptest.Server. Provide different secret IDs and expected
	// responses.
	tests := []struct {
		name         string
		options      *doppler.DynamicSecretIssueLeaseOptions
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Lease issue successful",
			options: &doppler.DynamicSecretIssueLeaseOptions{
				Project:    "test",
				Config:     "test",
				Name:       "test",
				TTLSeconds: 3600,
			},
			wantResponse: doppler.APIResponse{
				Status:     "200 OK",
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "Lease issue failed - secret not found",
			options: &doppler.DynamicSecretIssueLeaseOptions{
				Project:    "test",
				Config:     "test",
				Name:       "unknown",
				TTLSeconds: 3600,
			},
			wantResponse: doppler.APIResponse{
				Status:     "400 Bad Request",
				StatusCode: 400,
				Success:    pointer.To(false),
				Messages:   []string{"Secret not found"},
			},
			wantErr: true,
		},
		{
			name:         "Lease issue failed - invalid options error",
			options:      &doppler.DynamicSecretIssueLeaseOptions{},
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
				err := json.NewEncoder(w).Encode(&doppler.DynamicSecretIssueLeaseResponse{
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &dynamicsecret.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method.
			gotResponse, err := client.IssueLease(context.Background(), tt.options)
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

func TestDynamicSecret_RevokeLease(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.DynamicSecretRevokeLeaseOptions
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Lease revoke successful",
			options: &doppler.DynamicSecretRevokeLeaseOptions{
				Project: "test",
				Config:  "test",
				Name:    "test",
				Slug:    "test",
			},
			wantResponse: doppler.APIResponse{
				Status:     "200 OK",
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "Revoke lease failed - secret not found",
			options: &doppler.DynamicSecretRevokeLeaseOptions{
				Project: "test",
				Config:  "test",
				Name:    "test",
				Slug:    "unknown",
			},
			wantResponse: doppler.APIResponse{
				Status:     "400 Bad Request",
				StatusCode: 400,
				Success:    pointer.To(false),
				Messages:   []string{"Secret not found"},
			},
			wantErr: true,
		},
		{
			name:         "Revoke lease failed - invalid options error",
			options:      &doppler.DynamicSecretRevokeLeaseOptions{},
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
				err := json.NewEncoder(w).Encode(&doppler.DynamicSecretRevokeLeaseResponse{
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &dynamicsecret.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method.
			gotResponse, err := client.RevokeLease(context.Background(), tt.options)
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
