package secret_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nikoksr/doppler-go"
	"github.com/nikoksr/doppler-go/pointer"
	"github.com/nikoksr/doppler-go/secret"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := secret.Default()
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

func TestSecret_Get(t *testing.T) {
	t.Parallel()

	// Table driven tests that mock the Doppler API using httptest.Server. Provide different secret IDs and expected
	// responses.
	tests := []struct {
		name         string
		options      *doppler.SecretGetOptions
		wantSecret   *doppler.Secret
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Get secret",
			options: &doppler.SecretGetOptions{
				Project: "my-project",
				Config:  "my-config",
				Name:    "my-secret",
			},
			wantSecret: &doppler.Secret{
				Name: pointer.To("test"),
				Value: &doppler.SecretValue{
					Raw:      pointer.To("value1"),
					Computed: pointer.To("value1"),
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
			name: "Get secret with error",
			options: &doppler.SecretGetOptions{
				Project: "my-project",
				Config:  "my-config",
				Name:    "unknown-secret",
			},
			wantSecret: &doppler.Secret{
				Name: pointer.To("test"),
				Value: &doppler.SecretValue{
					Raw:      pointer.To("value1"),
					Computed: pointer.To("value1"),
				},
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"Secret not found"},
			},
			wantErr: true,
		},
		{
			name:         "Get secret with options validation error",
			options:      &doppler.SecretGetOptions{},
			wantSecret:   nil,
			wantResponse: doppler.APIResponse{},
			wantErr:      true,
		},
	}

	// Run tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a new httptest.Server that will be used to mock the Doppler API.
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Write the expected response to the ResponseWriter.
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.wantResponse.StatusCode)
				err := json.NewEncoder(w).Encode(&doppler.SecretGetResponse{
					Secret:      tt.wantSecret,
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &secret.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotSecret, gotResponse, err := client.Get(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the secret is expected.
			if diff := cmp.Diff(tt.wantSecret, gotSecret); diff != "" {
				t.Errorf("Unexpected secret (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSecret_List(t *testing.T) {
	t.Parallel()

	// Table driven tests that mock the Doppler API using httptest.Server. Provide different secret IDs and expected
	// responses.
	tests := []struct {
		name         string
		options      *doppler.SecretListOptions
		wantSecrets  map[string]*doppler.SecretValue
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "List all secrets",
			options: &doppler.SecretListOptions{
				Project:        "test",
				Config:         "test",
				IncludeDynamic: pointer.To(true),
			},
			wantSecrets: map[string]*doppler.SecretValue{
				"test":  {Raw: pointer.To("value1"), Computed: pointer.To("value1")},
				"test2": {Raw: pointer.To("value2"), Computed: pointer.To("value2")},
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Page:       pointer.To(1),
			},
			wantErr: false,
		},
		{
			name: "List secrets invalid ttl",
			options: &doppler.SecretListOptions{
				Project:           "test",
				Config:            "test",
				IncludeDynamic:    pointer.To(true),
				DynamicTTLSeconds: pointer.To[int32](-1),
			},
			wantSecrets: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"dynamic_ttl_seconds must be greater than or equal to 0"},
			},
			wantErr: true,
		},
		{
			name: "List secrets with options validation error",
			options: &doppler.SecretListOptions{
				Project: "test",
			},
			wantSecrets:  nil,
			wantResponse: doppler.APIResponse{},
			wantErr:      true,
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
				err := json.NewEncoder(w).Encode(&doppler.SecretListResponse{
					APIResponse: tt.wantResponse,
					Secrets:     tt.wantSecrets,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &secret.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotSecrets, gotResponse, err := client.List(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the secrets are expected.
			if diff := cmp.Diff(tt.wantSecrets, gotSecrets); diff != "" {
				t.Errorf("Unexpected secrets (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSecret_Update(t *testing.T) {
	t.Parallel()

	// Table driven tests that mock the Doppler API using httptest.Server. Provide different secret IDs and expected
	// responses.
	tests := []struct {
		name         string
		options      *doppler.SecretUpdateOptions
		wantSecrets  map[string]string
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Update secret",
			options: &doppler.SecretUpdateOptions{
				Project: "test",
				Config:  "test",
				NewSecrets: map[string]string{
					"test_name_1": "test_value_1",
					"test_name_2": "test_value_2",
				},
			},
			wantSecrets: map[string]string{
				"test_name_1": "test_value_1",
				"test_name_2": "test_value_2",
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Page:       pointer.To(1),
			},
			wantErr: false,
		},
		{
			name: "Update secret unknown secret",
			options: &doppler.SecretUpdateOptions{
				Project: "test",
				Config:  "test",
				NewSecrets: map[string]string{
					"test_name_1": "test_value_1",
					"unknown":     "test_value_2",
				},
			},
			wantSecrets: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"dynamic_ttl_seconds must be greater than or equal to 0"},
			},
			wantErr: true,
		},
		{
			name: "Update secret validation error",
			options: &doppler.SecretUpdateOptions{
				Project: "test",
			},
			wantSecrets:  nil,
			wantResponse: doppler.APIResponse{},
			wantErr:      true,
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
				err := json.NewEncoder(w).Encode(&doppler.SecretUpdateResponse{
					APIResponse: tt.wantResponse,
					Secrets:     tt.wantSecrets,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &secret.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotSecrets, gotResponse, err := client.Update(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the secret is expected.
			if diff := cmp.Diff(tt.wantSecrets, gotSecrets); diff != "" {
				t.Errorf("Unexpected secret (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSecret_Download(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.SecretDownloadOptions
		wantSecrets  string
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Download secrets in json format",
			options: &doppler.SecretDownloadOptions{
				Project: "test",
				Config:  "test",
				Format:  pointer.To("json"),
			},
			wantSecrets: `"secrets": {
        "test": {
          "raw": "value",
          "computed": "value",
        },
        "test2": {
          "raw": "value2",
          "computed": "value2",
        },
      }`,
			wantResponse: doppler.APIResponse{
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Download secrets in docker format",
			options: &doppler.SecretDownloadOptions{
				Project: "test",
				Config:  "test",
				Format:  pointer.To("docker"),
			},
			wantSecrets: `TEST1=value
TEST2=value2`,
			wantResponse: doppler.APIResponse{
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Download secrets empty",
			options: &doppler.SecretDownloadOptions{
				Project: "test",
				Config:  "test",
			},
			wantSecrets: "",
			wantResponse: doppler.APIResponse{
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Download secrets validation error",
			options: &doppler.SecretDownloadOptions{
				Project: "test",
			},
			wantSecrets:  "",
			wantResponse: doppler.APIResponse{},
			wantErr:      true,
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
				_, err := fmt.Fprint(w, tt.wantSecrets)
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &secret.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotSecrets, gotResponse, err := client.Download(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the secret is expected.
			if diff := cmp.Diff(tt.wantSecrets, gotSecrets); diff != "" {
				t.Errorf("Unexpected secret (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}
