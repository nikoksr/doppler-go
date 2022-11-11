package audit_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nikoksr/doppler-go"
	"github.com/nikoksr/doppler-go/audit"
	"github.com/nikoksr/doppler-go/pointer"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := audit.Default()
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

func TestAudit_WorkplaceGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.AuditWorkplaceGetOptions
		wantAudit    *doppler.AuditWorkplace
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name:    "Get audit",
			options: &doppler.AuditWorkplaceGetOptions{},
			wantAudit: &doppler.AuditWorkplace{
				ID:           pointer.To("123"),
				Name:         pointer.To("test"),
				BillingEmail: pointer.To("test@example.com"),
				SAMLEnabled:  pointer.To(true),
				SCIMEnabled:  pointer.To(true),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name:      "Get audit with error",
			options:   &doppler.AuditWorkplaceGetOptions{},
			wantAudit: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"error"},
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
				// Write the expected response to the ResponseWriter.
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.wantResponse.StatusCode)
				err := json.NewEncoder(w).Encode(&doppler.AuditWorkplaceGetResponse{
					AuditWorkplace: tt.wantAudit,
					APIResponse:    tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &audit.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotAudit, gotResponse, err := client.WorkplaceGet(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the audit is expected.
			if diff := cmp.Diff(tt.wantAudit, gotAudit); diff != "" {
				t.Errorf("Unexpected audit (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAudit_WorkplaceUserGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.AuditWorkplaceUserGetOptions
		wantAudit    *doppler.AuditWorkplaceUser
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Get audit",
			options: &doppler.AuditWorkplaceUserGetOptions{
				UserID: "123",
			},
			wantAudit: &doppler.AuditWorkplaceUser{
				ID:     pointer.To("123"),
				Access: pointer.To("read"),
				User: &doppler.User{
					Email:    pointer.To("test@example.com"),
					Name:     pointer.To("test"),
					UserName: pointer.To("test"),
				},
				CreatedAt: pointer.To("2020-01-01T00:00:00Z"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Get audit unknown user",
			options: &doppler.AuditWorkplaceUserGetOptions{
				UserID: "unknown",
			},
			wantAudit: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"user not found"},
			},
			wantErr: true,
		},
		{
			name:         "Get audit with invalid options error",
			options:      &doppler.AuditWorkplaceUserGetOptions{},
			wantAudit:    nil,
			wantResponse: doppler.APIResponse{},
			wantErr:      true,
		},
		{
			name:         "Get audit with nil options error", // Special case
			options:      nil,
			wantAudit:    nil,
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
				err := json.NewEncoder(w).Encode(&doppler.AuditWorkplaceUserGetResponse{
					AuditWorkplaceUser: tt.wantAudit,
					APIResponse:        tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &audit.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotAudit, gotResponse, err := client.WorkplaceUserGet(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the audit is expected.
			if diff := cmp.Diff(tt.wantAudit, gotAudit); diff != "" {
				t.Errorf("Unexpected audit (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAudit_WorkplaceUserList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      *doppler.AuditWorkplaceUserListOptions
		wantAudits   []*doppler.AuditWorkplaceUser
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name:    "List audits",
			options: &doppler.AuditWorkplaceUserListOptions{},
			wantAudits: []*doppler.AuditWorkplaceUser{
				{
					ID:     pointer.To("123"),
					Access: pointer.To("read"),
					User: &doppler.User{
						Email:    pointer.To("test@example.com"),
						Name:     pointer.To("test"),
						UserName: pointer.To("test"),
					},
					CreatedAt: pointer.To("2020-01-01T00:00:00Z"),
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
			name:       "List audits unknown access",
			options:    &doppler.AuditWorkplaceUserListOptions{},
			wantAudits: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"invalid access"},
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
				// Write the expected response to the ResponseWriter.
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.wantResponse.StatusCode)
				err := json.NewEncoder(w).Encode(&doppler.AuditWorkplaceUserListResponse{
					AuditWorkplaceUsers: tt.wantAudits,
					APIResponse:         tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &audit.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the List method with the test ID.
			gotAudits, gotResponse, err := client.WorkplaceUserList(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the audits are expected.
			if diff := cmp.Diff(tt.wantAudits, gotAudits); diff != "" {
				t.Errorf("Unexpected audits (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}
