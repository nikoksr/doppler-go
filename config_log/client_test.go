package configlog_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nikoksr/doppler-go"
	configlog "github.com/nikoksr/doppler-go/config_log"
	"github.com/nikoksr/doppler-go/pointer"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := configlog.Default()
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

func TestConfigLog_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		options       *doppler.ConfigLogGetOptions
		wantConfigLog *doppler.ConfigLog
		wantResponse  doppler.APIResponse
		wantErr       bool
	}{
		{
			name: "Get config log",
			options: &doppler.ConfigLogGetOptions{
				Project: "test",
				Config:  "test",
				ID:      "1",
			},
			wantConfigLog: &doppler.ConfigLog{
				ID:   pointer.To("1"),
				Text: pointer.To("Config log text"),
				HTML: pointer.To("Config log HTML"),
				User: &doppler.User{
					Email: pointer.To("test@example.com"),
					Name:  pointer.To("Test User"),
				},
				Project:     pointer.To("prj1"),
				Environment: pointer.To("env1"),
				Config:      pointer.To("cfg1"),
				CreatedAt:   pointer.To("2021-01-01T00:00:00Z"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Get config log with unknown project",
			options: &doppler.ConfigLogGetOptions{
				Project: "unknown",
				Config:  "test",
				ID:      "1",
			},
			wantConfigLog: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Messages:   []string{"Project not found"},
			},
			wantErr: true,
		},
		{
			name:          "Get config log with invalid options error",
			options:       &doppler.ConfigLogGetOptions{},
			wantConfigLog: nil,
			wantResponse:  doppler.APIResponse{},
			wantErr:       true,
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
				err := json.NewEncoder(w).Encode(&doppler.ConfigLogGetResponse{
					ConfigLog:   tt.wantConfigLog,
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &configlog.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotConfigLog, gotResponse, err := client.Get(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the configlog is expected.
			if diff := cmp.Diff(tt.wantConfigLog, gotConfigLog); diff != "" {
				t.Errorf("Unexpected configlog (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConfigLog_List(t *testing.T) {
	t.Parallel()

	// Table driven tests that mock the Doppler API using httptest.Server. Provide different config log IDs and expected
	// responses.
	tests := []struct {
		name           string
		options        *doppler.ConfigLogListOptions
		wantConfigLogs []*doppler.ConfigLog
		wantResponse   doppler.APIResponse
		wantErr        bool
	}{
		{
			name: "List config logs",
			options: &doppler.ConfigLogListOptions{
				Project:     "test",
				Config:      "test",
				ListOptions: doppler.ListOptions{Page: 1, PerPage: 2},
			},
			wantConfigLogs: []*doppler.ConfigLog{
				{
					ID:   pointer.To("1"),
					Text: pointer.To("Config log text"),
					HTML: pointer.To("Config log HTML"),
					User: &doppler.User{
						Email: pointer.To("test@example.com"),
						Name:  pointer.To("Test User"),
					},
					Project:     pointer.To("prj1"),
					Environment: pointer.To("env1"),
					Config:      pointer.To("cfg1"),
					CreatedAt:   pointer.To("2021-01-01T00:00:00Z"),
				},
				{
					ID:   pointer.To("2"),
					Text: pointer.To("Config log text"),
					HTML: pointer.To("Config log HTML"),
					User: &doppler.User{
						Email: pointer.To("test2@example.com"),
						Name:  pointer.To("Test User 2"),
					},
					Project:     pointer.To("prj2"),
					Environment: pointer.To("env2"),
					Config:      pointer.To("cfg2"),
					CreatedAt:   pointer.To("2021-01-01T00:00:00Z"),
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
			name: "List config logs with page 2",
			options: &doppler.ConfigLogListOptions{
				Project: "test",
				Config:  "test",
				ListOptions: doppler.ListOptions{
					Page:    2,
					PerPage: 1,
				},
			},
			wantConfigLogs: []*doppler.ConfigLog{
				{
					ID:   pointer.To("3"),
					Text: pointer.To("Config log text"),
					HTML: pointer.To("Config log HTML"),
					User: &doppler.User{
						Email: pointer.To("test3@example.com"),
						Name:  pointer.To("Test User 3"),
					},
					Project:     pointer.To("prj3"),
					Environment: pointer.To("env3"),
					Config:      pointer.To("cfg3"),
					CreatedAt:   pointer.To("2021-01-01T00:00:00Z"),
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
			name: "List config logs with no results",
			options: &doppler.ConfigLogListOptions{
				Project: "test",
				Config:  "test",
				ListOptions: doppler.ListOptions{
					Page:    1,
					PerPage: 1,
				},
			},
			wantConfigLogs: []*doppler.ConfigLog{},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Page:       pointer.To(1),
			},
			wantErr: false,
		},
		{
			name: "List config logs with invalid page",
			options: &doppler.ConfigLogListOptions{
				Project: "test",
				Config:  "test",
				ListOptions: doppler.ListOptions{
					Page:    -1,
					PerPage: 1,
				},
			},
			wantConfigLogs: []*doppler.ConfigLog{},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Page:       pointer.To(0),
				Messages:   []string{"Invalid page"},
			},
			wantErr: true,
		},
		{
			name:           "List config logs with invalid options error",
			options:        &doppler.ConfigLogListOptions{},
			wantConfigLogs: nil,
			wantResponse:   doppler.APIResponse{},
			wantErr:        true,
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
				err := json.NewEncoder(w).Encode(&doppler.ConfigLogListResponse{
					APIResponse: tt.wantResponse,
					ConfigLogs:  tt.wantConfigLogs,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &configlog.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotConfigLogs, gotResponse, err := client.List(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the configlogs are expected.
			if diff := cmp.Diff(tt.wantConfigLogs, gotConfigLogs); diff != "" {
				t.Errorf("Unexpected configlogs (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConfigLog_Rollback(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		options       *doppler.ConfigLogRollbackOptions
		wantConfigLog *doppler.ConfigLog
		wantResponse  doppler.APIResponse
		wantErr       bool
	}{
		{
			name: "Rollback config log",
			options: &doppler.ConfigLogRollbackOptions{
				Project: "test",
				Config:  "test",
				ID:      "test",
			},
			wantConfigLog: &doppler.ConfigLog{
				ID:   pointer.To("1"),
				Text: pointer.To("Config log text"),
				HTML: pointer.To("Config log HTML"),
				Diff: []doppler.ConfigLogDiff{
					{Name: pointer.To("test"), Added: pointer.To("test")},
				},
				Rollback:    pointer.To(false),
				User:        &doppler.User{Email: pointer.To("test@example.com")},
				Project:     pointer.To("test"),
				Environment: pointer.To("test"),
				Config:      pointer.To("test"),
				CreatedAt:   pointer.To(time.Unix(0, 0).String()),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Rollback config log with unknown project error",
			options: &doppler.ConfigLogRollbackOptions{
				Project: "unknown",
				Config:  "test",
				ID:      "test",
			},
			wantConfigLog: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Messages:   []string{"Project not found"},
			},
			wantErr: true,
		},
		{
			name: "Rollback config log with invalid options error",
			options: &doppler.ConfigLogRollbackOptions{
				Project: "test",
				Config:  "test",
			},
			wantConfigLog: nil,
			wantResponse:  doppler.APIResponse{},
			wantErr:       true,
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
				err := json.NewEncoder(w).Encode(&doppler.ConfigLogRollbackResponse{
					APIResponse: tt.wantResponse,
					ConfigLog:   tt.wantConfigLog,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &configlog.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotConfigLog, gotResponse, err := client.Rollback(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the configlog is expected.
			if diff := cmp.Diff(tt.wantConfigLog, gotConfigLog); diff != "" {
				t.Errorf("Unexpected configlog (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}
