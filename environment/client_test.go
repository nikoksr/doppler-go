package environment_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nikoksr/doppler-go"
	"github.com/nikoksr/doppler-go/environment"
	"github.com/nikoksr/doppler-go/pointer"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := environment.Default()
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

func TestEnvironment_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		options         *doppler.EnvironmentGetOptions
		wantEnvironment *doppler.Environment
		wantResponse    doppler.APIResponse
		wantErr         bool
	}{
		{
			name: "Get environment with ID 1",
			options: &doppler.EnvironmentGetOptions{
				Project: "p1",
				Slug:    "e1",
			},
			wantEnvironment: &doppler.Environment{
				ID:             pointer.To("1"),
				Name:           pointer.To("p1"),
				Project:        pointer.To("p1"),
				InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
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
			name: "Get environment with unknown slug",
			options: &doppler.EnvironmentGetOptions{
				Project: "p1",
				Slug:    "unknown",
			},
			wantEnvironment: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"Invalid environment slug"},
			},
			wantErr: true,
		},
		{
			name:            "Get environment with invalid options error",
			options:         &doppler.EnvironmentGetOptions{},
			wantEnvironment: nil,
			wantResponse:    doppler.APIResponse{},
			wantErr:         true,
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
				err := json.NewEncoder(w).Encode(&doppler.EnvironmentGetResponse{
					Environment: tt.wantEnvironment,
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &environment.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotEnvironment, gotResponse, err := client.Get(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the environment is expected.
			if diff := cmp.Diff(tt.wantEnvironment, gotEnvironment); diff != "" {
				t.Errorf("Unexpected environment (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEnvironment_List(t *testing.T) {
	t.Parallel()

	// Table driven tests that mock the Doppler API using httptest.Server. Provide different environment IDs and expected
	// responses.
	tests := []struct {
		name             string
		options          *doppler.EnvironmentListOptions
		wantEnvironments []*doppler.Environment
		wantResponse     doppler.APIResponse
		wantErr          bool
	}{
		{
			name: "List environments for project p1",
			options: &doppler.EnvironmentListOptions{
				Project: "p1",
			},
			wantEnvironments: []*doppler.Environment{
				{
					ID:             pointer.To("1"),
					Name:           pointer.To("e1"),
					Project:        pointer.To("p1"),
					InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
					CreatedAt:      pointer.To("2021-01-01T00:00:00.000Z"),
				},
				{
					ID:             pointer.To("2"),
					Name:           pointer.To("e2"),
					Project:        pointer.To("p1"),
					InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
					CreatedAt:      pointer.To("2021-01-01T00:00:00.000Z"),
				},
				{
					ID:             pointer.To("3"),
					Name:           pointer.To("e3"),
					Project:        pointer.To("p1"),
					InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
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
			name: "List environments with unknown project",
			options: &doppler.EnvironmentListOptions{
				Project: "unknown",
			},
			wantEnvironments: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"Project not found"},
			},
			wantErr: true,
		},
		{
			name:             "List environments with invalid options error",
			options:          &doppler.EnvironmentListOptions{},
			wantEnvironments: nil,
			wantResponse:     doppler.APIResponse{},
			wantErr:          true,
		},
		{
			name:             "List environments with invalid options error",
			options:          &doppler.EnvironmentListOptions{},
			wantEnvironments: nil,
			wantResponse:     doppler.APIResponse{},
			wantErr:          true,
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
				err := json.NewEncoder(w).Encode(&doppler.EnvironmentListResponse{
					APIResponse:  tt.wantResponse,
					Environments: tt.wantEnvironments,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &environment.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotEnvironments, gotResponse, err := client.List(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the environments are expected.
			if diff := cmp.Diff(tt.wantEnvironments, gotEnvironments); diff != "" {
				t.Errorf("Unexpected environments (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEnvironment_Create(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		name            string
		options         *doppler.EnvironmentCreateOptions
		wantEnvironment *doppler.Environment
		wantResponse    doppler.APIResponse
		wantErr         bool
	}{
		{
			name: "Create environment",
			options: &doppler.EnvironmentCreateOptions{
				Project: "p1",
				Name:    "Environment-1",
				Slug:    "e1",
			},
			wantEnvironment: &doppler.Environment{
				ID:             pointer.To("Environtment-1"),
				Slug:           pointer.To("e1"),
				Name:           pointer.To("Environment-1"),
				Project:        pointer.To("p1"),
				InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
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
			name: "Create environment with unknown project",
			options: &doppler.EnvironmentCreateOptions{
				Project: "unknown",
				Name:    "test",
				Slug:    "e1",
			},
			wantEnvironment: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"Project not found"},
			},
			wantErr: true,
		},
		{
			name:            "Create environment with invalid options error",
			options:         &doppler.EnvironmentCreateOptions{},
			wantEnvironment: nil,
			wantResponse:    doppler.APIResponse{},
			wantErr:         true,
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
				err := json.NewEncoder(w).Encode(&doppler.EnvironmentCreateResponse{
					APIResponse: tt.wantResponse,
					Environment: tt.wantEnvironment,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &environment.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Create method with the test environment.
			gotEnvironment, gotResponse, err := client.Create(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}

			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}

			// Check if the environment is expected.
			if diff := cmp.Diff(tt.wantEnvironment, gotEnvironment); diff != "" {
				t.Errorf("Unexpected environment (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEnvironment_Rename(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		name            string
		options         *doppler.EnvironmentRenameOptions
		wantEnvironment *doppler.Environment
		wantResponse    doppler.APIResponse
		wantErr         bool
	}{
		{
			name: "Rename environment",
			options: &doppler.EnvironmentRenameOptions{
				Project: "p1",
				Slug:    "e1",
				NewName: pointer.To("Environment-1"),
				NewSlug: pointer.To("en1"),
			},
			wantEnvironment: &doppler.Environment{
				ID:             pointer.To("Environtment-1"),
				Slug:           pointer.To("en1"),
				Name:           pointer.To("Environment-1"),
				Project:        pointer.To("p1"),
				InitialFetchAt: pointer.To("2021-01-01T00:00:00.000Z"),
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
			name: "Rename environment with unknown project",
			options: &doppler.EnvironmentRenameOptions{
				Project: "unknown",
				Slug:    "e1",
			},
			wantEnvironment: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"name or slug is required"},
			},
			wantErr: true,
		},
		{
			name:            "Rename environment with invalid options error",
			options:         &doppler.EnvironmentRenameOptions{},
			wantEnvironment: nil,
			wantResponse:    doppler.APIResponse{},
			wantErr:         true,
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
				err := json.NewEncoder(w).Encode(&doppler.EnvironmentRenameResponse{
					APIResponse: tt.wantResponse,
					Environment: tt.wantEnvironment,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &environment.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Rename method with the test environment.
			gotEnvironment, gotResponse, err := client.Rename(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}

			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}

			// Check if the environment is expected.
			if diff := cmp.Diff(tt.wantEnvironment, gotEnvironment); diff != "" {
				t.Errorf("Unexpected environment (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEnvironment_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.EnvironmentDeleteOptions
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Delete environment",
			options: &doppler.EnvironmentDeleteOptions{
				Project: "p1",
				Slug:    "e1",
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Delete environment with unknown environment",
			options: &doppler.EnvironmentDeleteOptions{
				Project: "p1",
				Slug:    "unknown",
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Messages:   []string{"environment not found"},
			},
			wantErr: true,
		},
		{
			name:         "Delete environment with invalid options error",
			options:      &doppler.EnvironmentDeleteOptions{},
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
				err := json.NewEncoder(w).Encode(&doppler.EnvironmentDeleteResponse{
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &environment.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Delete method with the test environment.
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
