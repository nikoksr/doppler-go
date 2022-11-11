package config_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nikoksr/doppler-go"
	"github.com/nikoksr/doppler-go/config"
	"github.com/nikoksr/doppler-go/pointer"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := config.Default()
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

func TestConfig_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.ConfigGetOptions
		wantConfig   *doppler.Config
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Get config",
			options: &doppler.ConfigGetOptions{
				Project: "p1",
				Config:  "c1",
			},
			wantConfig: &doppler.Config{
				Name:           pointer.To("c1"),
				Project:        pointer.To("p1"),
				Environment:    pointer.To("dev"),
				Root:           pointer.To(true),
				Locked:         pointer.To(false),
				InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
				LastFetchAt:    pointer.To("2021-01-01T00:00:00.000Z"),
				CreatedAt:      pointer.To("2021-01-01T00:00:00.000Z"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Get config with unknown name",
			options: &doppler.ConfigGetOptions{
				Project: "p1",
				Config:  "unknown",
			},
			wantConfig: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Messages:   []string{"Config not found"},
			},
			wantErr: true,
		},
		{
			name:         "Get config with invalid options error",
			options:      &doppler.ConfigGetOptions{},
			wantConfig:   nil,
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
				err := json.NewEncoder(w).Encode(&doppler.ConfigGetResponse{
					Config:      tt.wantConfig,
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &config.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotConfig, gotResponse, err := client.Get(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the config is expected.
			if diff := cmp.Diff(tt.wantConfig, gotConfig); diff != "" {
				t.Errorf("Unexpected config (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConfig_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.ConfigListOptions
		wantConfigs  []*doppler.Config
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "List configs",
			options: &doppler.ConfigListOptions{
				Project: "p1",
				ListOptions: doppler.ListOptions{
					Page:    1,
					PerPage: 10,
				},
			},
			wantConfigs: []*doppler.Config{
				{
					Name:           pointer.To("c1"),
					Project:        pointer.To("p1"),
					Environment:    pointer.To("dev"),
					Root:           pointer.To(true),
					Locked:         pointer.To(false),
					InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
					LastFetchAt:    pointer.To("2021-01-01T00:00:00.000Z"),
					CreatedAt:      pointer.To("2021-01-01T00:00:00.000Z"),
				},
				{
					Name:           pointer.To("c2"),
					Project:        pointer.To("p1"),
					Environment:    pointer.To("dev"),
					Root:           pointer.To(true),
					Locked:         pointer.To(false),
					InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
					LastFetchAt:    pointer.To("2021-01-01T00:00:00.000Z"),
					CreatedAt:      pointer.To("2021-01-01T00:00:00.000Z"),
				},
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
			name: "List configs with page 2 and per page 1",
			options: &doppler.ConfigListOptions{
				Project: "p1",
				ListOptions: doppler.ListOptions{
					Page:    2,
					PerPage: 1,
				},
			},
			wantConfigs: []*doppler.Config{
				{
					Name:           pointer.To("c2"),
					Project:        pointer.To("p1"),
					Environment:    pointer.To("dev"),
					Root:           pointer.To(true),
					Locked:         pointer.To(false),
					InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
					LastFetchAt:    pointer.To("2021-01-01T00:00:00.000Z"),
					CreatedAt:      pointer.To("2021-01-01T00:00:00.000Z"),
				},
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Page:       pointer.To(2),
			},
			wantErr: false,
		},
		{
			name: "List configs with invalid page",
			options: &doppler.ConfigListOptions{
				Project: "p1",
				ListOptions: doppler.ListOptions{
					Page:    0,
					PerPage: 10,
				},
			},
			wantConfigs: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"invalid page"},
				Page:       pointer.To(1),
			},
			wantErr: true,
		},
		{
			name:         "List configs with invalid options error",
			options:      &doppler.ConfigListOptions{},
			wantConfigs:  nil,
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
				err := json.NewEncoder(w).Encode(&doppler.ConfigListResponse{
					APIResponse: tt.wantResponse,
					Configs:     tt.wantConfigs,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &config.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotConfigs, gotResponse, err := client.List(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the configs are expected.
			if diff := cmp.Diff(tt.wantConfigs, gotConfigs); diff != "" {
				t.Errorf("Unexpected configs (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConfig_Create(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		name         string
		options      *doppler.ConfigCreateOptions
		wantConfig   *doppler.Config
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Create config",
			options: &doppler.ConfigCreateOptions{
				Project:     "p1",
				Environment: "dev",
				Name:        "c1",
			},
			wantConfig: &doppler.Config{
				Name:           pointer.To("c1"),
				Project:        pointer.To("p1"),
				Environment:    pointer.To("dev"),
				Root:           pointer.To(true),
				Locked:         pointer.To(false),
				InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
				LastFetchAt:    pointer.To("2021-01-01T00:00:00.000Z"),
				CreatedAt:      pointer.To("2021-01-01T00:00:00.000Z"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "201 Created",
				StatusCode: http.StatusCreated,
			},
			wantErr: false,
		},
		{
			name: "Create config with unknown name",
			options: &doppler.ConfigCreateOptions{
				Project:     "p1",
				Environment: "dev",
				Name:        "unknown",
			},
			wantConfig: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"name is required"},
			},
			wantErr: true,
		},
		{
			name:         "Create config with invalid options error",
			options:      &doppler.ConfigCreateOptions{},
			wantConfig:   nil,
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
				err := json.NewEncoder(w).Encode(&doppler.ConfigCreateResponse{
					APIResponse: tt.wantResponse,
					Config:      tt.wantConfig,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &config.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Create method with the test config.
			gotConfig, gotResponse, err := client.Create(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}

			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}

			// Check if the config is expected.
			if diff := cmp.Diff(tt.wantConfig, gotConfig); diff != "" {
				t.Errorf("Unexpected config (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConfig_Update(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		name         string
		options      *doppler.ConfigUpdateOptions
		wantConfig   *doppler.Config
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Update config",
			options: &doppler.ConfigUpdateOptions{
				Project: "p1",
				Config:  "c1",
				NewName: "c2",
			},
			wantConfig: &doppler.Config{
				Name:           pointer.To("c2"),
				Project:        pointer.To("p1"),
				Environment:    pointer.To("dev"),
				Root:           pointer.To(true),
				Locked:         pointer.To(false),
				InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
				LastFetchAt:    pointer.To("2021-01-01T00:00:00.000Z"),
				CreatedAt:      pointer.To("2021-01-01T00:00:00.000Z"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Update config with unknown config",
			options: &doppler.ConfigUpdateOptions{
				Project: "p1",
				Config:  "unknown",
				NewName: "c2",
			},
			wantConfig: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"config name is required"},
			},
			wantErr: true,
		},
		{
			name:         "Update config with invalid options error",
			options:      &doppler.ConfigUpdateOptions{},
			wantConfig:   nil,
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
				err := json.NewEncoder(w).Encode(&doppler.ConfigUpdateResponse{
					APIResponse: tt.wantResponse,
					Config:      tt.wantConfig,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &config.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Update method with the test config.
			gotConfig, gotResponse, err := client.Update(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}

			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}

			// Check if the config is expected.
			if diff := cmp.Diff(tt.wantConfig, gotConfig); diff != "" {
				t.Errorf("Unexpected config (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConfig_Delete(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		name         string
		options      *doppler.ConfigDeleteOptions
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name:    "Delete config",
			options: &doppler.ConfigDeleteOptions{Project: "p1", Config: "c1"},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name:    "Delete config with unknown config",
			options: &doppler.ConfigDeleteOptions{Project: "p1", Config: "unknown"},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"id is required"},
			},
			wantErr: true,
		},
		{
			name:         "Delete config with invalid options error",
			options:      &doppler.ConfigDeleteOptions{},
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
				err := json.NewEncoder(w).Encode(&doppler.ConfigDeleteResponse{
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &config.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Delete method with the test config.
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

func TestConfig_Lock(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		name         string
		options      *doppler.ConfigLockOptions
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Lock config",
			options: &doppler.ConfigLockOptions{
				Project: "p1",
				Config:  "c1",
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Lock config with unknown config",
			options: &doppler.ConfigLockOptions{
				Project: "p1",
				Config:  "unknown",
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"id is required"},
			},
			wantErr: true,
		},
		{
			name:         "Lock config with invalid options error",
			options:      &doppler.ConfigLockOptions{},
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
				err := json.NewEncoder(w).Encode(&doppler.ConfigLockResponse{
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &config.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Lock method with the test config.
			gotConfig, gotResponse, err := client.Lock(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}

			// Check if the config is nil.
			if gotConfig != nil {
				t.Errorf("Unexpected config. Expected %v, got %v", nil, gotConfig)
			}

			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConfig_Unlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.ConfigUnlockOptions
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Unlock config",
			options: &doppler.ConfigUnlockOptions{
				Project: "p1",
				Config:  "c1",
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Unlock config with unknown config",
			options: &doppler.ConfigUnlockOptions{
				Project: "p1",
				Config:  "unknown",
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"id is required"},
			},
			wantErr: true,
		},
		{
			name:         "Unlock config with invalid options error",
			options:      &doppler.ConfigUnlockOptions{},
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
				err := json.NewEncoder(w).Encode(&doppler.ConfigUnlockResponse{
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &config.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Unlock method with the test config.
			gotConfig, gotResponse, err := client.Unlock(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}

			// Check if the config is nil.
			if gotConfig != nil {
				t.Errorf("Unexpected config. Expected %v, got %v", nil, gotConfig)
			}

			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConfig_Clone(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		name         string
		options      *doppler.ConfigCloneOptions
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Clone config",
			options: &doppler.ConfigCloneOptions{
				Project:   "p1",
				Config:    "c1",
				NewConfig: "c2",
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Clone config with unknown config",
			options: &doppler.ConfigCloneOptions{
				Project:   "p1",
				Config:    "unknown",
				NewConfig: "c2",
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"id is required"},
			},
			wantErr: true,
		},
		{
			name:         "Clone config with invalid options error",
			options:      &doppler.ConfigCloneOptions{},
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
				err := json.NewEncoder(w).Encode(&doppler.ConfigCloneResponse{
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &config.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Clone method with the test config.
			gotConfig, gotResponse, err := client.Clone(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}

			// Check if the config is nil.
			if gotConfig != nil {
				t.Errorf("Unexpected config. Expected %v, got %v", nil, gotConfig)
			}

			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}
