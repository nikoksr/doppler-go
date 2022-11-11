package activitylog_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nikoksr/doppler-go"
	activitylog "github.com/nikoksr/doppler-go/activity_log"
	"github.com/nikoksr/doppler-go/pointer"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := activitylog.Default()
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

func TestActivityLog_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		options         *doppler.ActivityLogGetOptions
		wantActivityLog *doppler.ActivityLog
		wantResponse    doppler.APIResponse
		wantErr         bool
	}{
		{
			name:    "Get activity log with ID 1",
			options: &doppler.ActivityLogGetOptions{ID: "1"},
			wantActivityLog: &doppler.ActivityLog{
				ID:   pointer.To("1"),
				Text: pointer.To("Activity log text"),
				HTML: pointer.To("Activity log HTML"),
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
			name:            "Get activity log with unknown ID",
			options:         &doppler.ActivityLogGetOptions{ID: "unknown"},
			wantActivityLog: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Messages:   []string{"Activity log not found"},
			},
			wantErr: true,
		},
		{
			name:            "Get activity log with invalid options error",
			options:         &doppler.ActivityLogGetOptions{},
			wantActivityLog: nil,
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
				err := json.NewEncoder(w).Encode(&doppler.ActivityLogGetResponse{
					ActivityLog: tt.wantActivityLog,
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &activitylog.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotActivityLog, gotResponse, err := client.Get(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the activitylog is expected.
			if diff := cmp.Diff(tt.wantActivityLog, gotActivityLog); diff != "" {
				t.Errorf("Unexpected activitylog (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestActivityLog_List(t *testing.T) {
	t.Parallel()

	// Table driven tests that mock the Doppler API using httptest.Server. Provide different activity log IDs and expected
	// responses.
	tests := []struct {
		name             string
		options          *doppler.ActivityLogListOptions
		wantActivityLogs []*doppler.ActivityLog
		wantResponse     doppler.APIResponse
		wantErr          bool
	}{
		{
			name: "List activity logs",
			options: &doppler.ActivityLogListOptions{
				ListOptions: doppler.ListOptions{
					Page:    1,
					PerPage: 2,
				},
			},
			wantActivityLogs: []*doppler.ActivityLog{
				{
					ID:   pointer.To("1"),
					Text: pointer.To("Activity log text"),
					HTML: pointer.To("Activity log HTML"),
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
					Text: pointer.To("Activity log text"),
					HTML: pointer.To("Activity log HTML"),
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
			name: "List activity logs with page 2",
			options: &doppler.ActivityLogListOptions{
				ListOptions: doppler.ListOptions{
					Page:    2,
					PerPage: 1,
				},
			},
			wantActivityLogs: []*doppler.ActivityLog{
				{
					ID:   pointer.To("3"),
					Text: pointer.To("Activity log text"),
					HTML: pointer.To("Activity log HTML"),
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
			name: "List activity logs with no results",
			options: &doppler.ActivityLogListOptions{
				ListOptions: doppler.ListOptions{
					Page:    1,
					PerPage: 1,
				},
			},
			wantActivityLogs: []*doppler.ActivityLog{},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Page:       pointer.To(1),
			},
			wantErr: false,
		},
		{
			name: "List activity logs with invalid page",
			options: &doppler.ActivityLogListOptions{
				ListOptions: doppler.ListOptions{
					Page:    0,
					PerPage: 1,
				},
			},
			wantActivityLogs: []*doppler.ActivityLog{},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Page:       pointer.To(0),
				Messages:   []string{"Invalid page"},
			},
			wantErr: true,
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
				err := json.NewEncoder(w).Encode(&doppler.ActivityLogListResponse{
					APIResponse:  tt.wantResponse,
					ActivityLogs: tt.wantActivityLogs,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &activitylog.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotActivityLogs, gotResponse, err := client.List(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the activitylogs are expected.
			if diff := cmp.Diff(tt.wantActivityLogs, gotActivityLogs); diff != "" {
				t.Errorf("Unexpected activitylogs (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}
