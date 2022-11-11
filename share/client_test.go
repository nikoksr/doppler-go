package share_test

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
	"github.com/nikoksr/doppler-go/share"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := share.Default()
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

func TestShare_PlainSecret(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.SharePlainOptions
		wantSecret   *doppler.SharePlain
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Share plain secret",
			options: &doppler.SharePlainOptions{
				Secret: "my-secret",
			},
			wantSecret: &doppler.SharePlain{
				URL:              pointer.To("https://example.com/share/secret/1234567890"),
				AuthenticatedURL: pointer.To("https://example.com/share/secret/1234567890?auth=1234567890"),
				Password:         pointer.To("1234567890"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Share plain secret with error",
			options: &doppler.SharePlainOptions{
				Secret: "my-secret",
			},
			wantSecret: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"some error"},
			},
			wantErr: true,
		},
		{
			name:         "Share plain secret with options validation error",
			options:      &doppler.SharePlainOptions{},
			wantSecret:   nil,
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
				// Write the expected response to the ResponseWriter.
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.wantResponse.StatusCode)
				err := json.NewEncoder(w).Encode(&doppler.SharePlainResponse{
					APIResponse: tt.wantResponse,
					Secret:      tt.wantSecret,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &share.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the share.PlainSecret method.
			gotSecret, gotResponse, err := client.PlainSecret(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the share is expected.
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

func TestShare_EncryptedSecret(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.ShareEncryptedOptions
		wantSecret   *doppler.ShareEncrypted
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Share encrypted secret",
			options: &doppler.ShareEncryptedOptions{
				Secret:     "my-secret",
				Password:   "my-password",
				KDF:        doppler.EncryptionKDF,
				SaltRounds: doppler.EncryptionSaltRounds,
			},
			wantSecret: &doppler.ShareEncrypted{
				URL: pointer.To("https://example.com/share/secret/1234567890"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Share encrypted secret with error",
			options: &doppler.ShareEncryptedOptions{
				Secret:     "my-secret",
				Password:   "my-password",
				KDF:        doppler.EncryptionKDF,
				SaltRounds: doppler.EncryptionSaltRounds,
			},
			wantSecret: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"some error"},
			},
			wantErr: true,
		},
		{
			name: "Share encrypted secret with invalid kdf",
			options: &doppler.ShareEncryptedOptions{
				Secret:     "my-secret",
				Password:   "my-password",
				KDF:        "invalid",
				SaltRounds: doppler.EncryptionSaltRounds,
			},
			wantSecret:   nil,
			wantResponse: doppler.APIResponse{},
			wantErr:      true,
		},
		{
			name: "Share encrypted secret with invalid salt rounds",
			options: &doppler.ShareEncryptedOptions{
				Secret:     "my-secret",
				Password:   "my-password",
				KDF:        doppler.EncryptionKDF,
				SaltRounds: -1,
			},
			wantSecret:   nil,
			wantResponse: doppler.APIResponse{},
			wantErr:      true,
		},
		{
			name:         "Share encrypted secret with options validation error",
			options:      &doppler.ShareEncryptedOptions{},
			wantSecret:   nil,
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
				// Write the expected response to the ResponseWriter.
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.wantResponse.StatusCode)
				err := json.NewEncoder(w).Encode(&doppler.ShareEncryptedResponse{
					APIResponse: tt.wantResponse,
					Secret:      tt.wantSecret,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &share.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the share.EncryptedSecret method.
			gotSecret, gotResponse, err := client.EncryptedSecret(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}

			// Check if the share is expected.
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
