package workplace_test

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
	"github.com/nikoksr/doppler-go/workplace"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := workplace.Default()
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

func TestWorkplace_Get(t *testing.T) {
	t.Parallel()

	// Table driven tests that mock the Doppler API using httptest.Server. Provide different workplace IDs and expected
	// responses.
	tests := []struct {
		name          string
		wantWorkplace *doppler.Workplace
		wantResponse  doppler.APIResponse
		wantErr       bool
	}{
		{
			name: "Get workplace",
			wantWorkplace: &doppler.Workplace{
				ID:           pointer.To("w1"),
				Name:         pointer.To("Workplace 1"),
				BillingEmail: pointer.To("test@example.com"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
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
				err := json.NewEncoder(w).Encode(&doppler.WorkplaceGetResponse{
					Workplace:   tt.wantWorkplace,
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &workplace.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the GetByID method with the test ID.
			gotWorkplace, gotResponse, err := client.Get(context.Background())
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the workplace is expected.
			if diff := cmp.Diff(tt.wantWorkplace, gotWorkplace); diff != "" {
				t.Errorf("Unexpected workplace (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClient_Update(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		name          string
		workplace     *doppler.WorkplaceUpdateOptions
		wantWorkplace *doppler.Workplace
		wantResponse  doppler.APIResponse
		wantErr       bool
	}{
		{
			name: "Update workplace",
			workplace: &doppler.WorkplaceUpdateOptions{
				NewName:         pointer.To("Workplace 1"),
				NewBillingEmail: pointer.To("example@test.com"),
			},
			wantWorkplace: &doppler.Workplace{
				ID:           pointer.To("w1"),
				Name:         pointer.To("Workplace 1"),
				BillingEmail: pointer.To("example@test.com"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Update workplace name",
			workplace: &doppler.WorkplaceUpdateOptions{
				NewName: pointer.To("Workplace 1 updated"),
			},
			wantWorkplace: &doppler.Workplace{
				ID:           pointer.To("w1"),
				Name:         pointer.To("Workplace 1 updated"),
				BillingEmail: pointer.To("example@test.com"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Update workplace with error",
			workplace: &doppler.WorkplaceUpdateOptions{
				NewName: pointer.To(""),
			},
			wantWorkplace: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"Name is empty"},
			},
			wantErr: true,
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
				err := json.NewEncoder(w).Encode(&doppler.WorkplaceUpdateResponse{
					APIResponse: tt.wantResponse,
					Workplace:   tt.wantWorkplace,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &workplace.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Update method with the test workplace.
			gotWorkplace, gotResponse, err := client.Update(context.Background(), tt.workplace)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}

			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}

			// Check if the workplace is expected.
			if diff := cmp.Diff(tt.wantWorkplace, gotWorkplace); diff != "" {
				t.Errorf("Unexpected workplace (-want +got):\n%s", diff)
			}
		})
	}
}
