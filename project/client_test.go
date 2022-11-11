package project_test

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
	"github.com/nikoksr/doppler-go/project"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	client := project.Default()
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

func TestProject_Get(t *testing.T) {
	t.Parallel()

	// Table driven tests that mock the Doppler API using httptest.Server. Provide different project IDs and expected
	// responses.
	tests := []struct {
		name         string
		options      *doppler.ProjectGetOptions
		wantProject  *doppler.Project
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Get project with ID 1",
			options: &doppler.ProjectGetOptions{
				Name: "p1",
			},
			wantProject: &doppler.Project{
				ID:          pointer.To("1"),
				Slug:        pointer.To("p1"),
				Name:        pointer.To("Project 1"),
				Description: pointer.To("Project 1 description"),
				CreatedAt:   pointer.To("2020-01-01T00:00:00.000Z"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Get project with ID 2",
			options: &doppler.ProjectGetOptions{
				Name: "p2",
			},
			wantProject: &doppler.Project{
				ID:   pointer.To("2"),
				Slug: pointer.To("p2"),
				Name: pointer.To("Project 2"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Get unknown project",
			options: &doppler.ProjectGetOptions{
				Name: "unknown",
			},
			wantProject: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Messages:   []string{"Project not found"},
			},
			wantErr: true,
		},
		{
			name:         "Get project with options validation error",
			options:      &doppler.ProjectGetOptions{},
			wantProject:  nil,
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
				err := json.NewEncoder(w).Encode(&doppler.ProjectGetResponse{
					Project:     tt.wantProject,
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &project.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotProject, gotResponse, err := client.Get(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the project is expected.
			if diff := cmp.Diff(tt.wantProject, gotProject); diff != "" {
				t.Errorf("Unexpected project (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProject_List(t *testing.T) {
	t.Parallel()

	// Table driven tests that mock the Doppler API using httptest.Server. Provide different project IDs and expected
	// responses.
	tests := []struct {
		name         string
		options      *doppler.ProjectListOptions
		wantProjects []*doppler.Project
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "List projects",
			options: &doppler.ProjectListOptions{
				ListOptions: doppler.ListOptions{
					Page:    1,
					PerPage: 10,
				},
			},
			wantProjects: []*doppler.Project{
				{
					ID:          pointer.To("1"),
					Slug:        pointer.To("p1"),
					Name:        pointer.To("Project 1"),
					Description: pointer.To("Project 1 description"),
					CreatedAt:   pointer.To("2020-01-01T00:00:00.000Z"),
				},
				{
					ID:          pointer.To("2"),
					Slug:        pointer.To("p2"),
					Name:        pointer.To("Project 2"),
					Description: pointer.To("Project 2 description"),
					CreatedAt:   pointer.To("2020-01-01T00:00:00.000Z"),
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
			name: "List projects with page 2 and per page 1",
			options: &doppler.ProjectListOptions{
				ListOptions: doppler.ListOptions{
					Page:    2,
					PerPage: 1,
				},
			},
			wantProjects: []*doppler.Project{
				{
					ID:          pointer.To("2"),
					Slug:        pointer.To("p2"),
					Name:        pointer.To("Project 2"),
					Description: pointer.To("Project 2 description"),
					CreatedAt:   pointer.To("2020-01-01T00:00:00.000Z"),
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
			name: "List projects with invalid page",
			options: &doppler.ProjectListOptions{
				ListOptions: doppler.ListOptions{
					Page:    0,
					PerPage: 10,
				},
			},
			wantProjects: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"invalid page"},
				Page:       pointer.To(1),
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
				err := json.NewEncoder(w).Encode(&doppler.ProjectListResponse{
					APIResponse: tt.wantResponse,
					Projects:    tt.wantProjects,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &project.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Get method with the test ID.
			gotProjects, gotResponse, err := client.List(context.Background(), tt.options)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}
			// Check if the projects are expected.
			if diff := cmp.Diff(tt.wantProjects, gotProjects); diff != "" {
				t.Errorf("Unexpected projects (-want +got):\n%s", diff)
			}
			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProject_Create(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		name         string
		project      *doppler.ProjectCreateOptions
		wantProject  *doppler.Project
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Create project",
			project: &doppler.ProjectCreateOptions{
				Name:        "Project 1",
				Description: pointer.To("Project 1 description"),
			},
			wantProject: &doppler.Project{
				ID:          pointer.To("1"),
				Slug:        pointer.To("p1"),
				Name:        pointer.To("Project 1"),
				Description: pointer.To("Project 1 description"),
				CreatedAt:   pointer.To("2020-01-01T00:00:00.000Z"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "201 Created",
				StatusCode: http.StatusCreated,
			},
			wantErr: false,
		},
		{
			name: "Create project with invalid name",
			project: &doppler.ProjectCreateOptions{
				Name:        "unknown",
				Description: pointer.To("Project 1 description"),
			},
			wantProject: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"name is required"},
			},
			wantErr: true,
		},
		{
			name:         "Create project with options validation error",
			project:      &doppler.ProjectCreateOptions{},
			wantProject:  nil,
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
				err := json.NewEncoder(w).Encode(&doppler.ProjectCreateResponse{
					APIResponse: tt.wantResponse,
					Project:     tt.wantProject,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &project.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Create method with the test project.
			gotProject, gotResponse, err := client.Create(context.Background(), tt.project)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}

			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}

			// Check if the project is expected.
			if diff := cmp.Diff(tt.wantProject, gotProject); diff != "" {
				t.Errorf("Unexpected project (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProject_Update(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		name         string
		project      *doppler.ProjectUpdateOptions
		wantProject  *doppler.Project
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name: "Update project",
			project: &doppler.ProjectUpdateOptions{
				Name:           "p1",
				NewName:        "Project 1",
				NewDescription: pointer.To("Project 1 description"),
			},
			wantProject: &doppler.Project{
				ID:          pointer.To("1"),
				Slug:        pointer.To("p1"),
				Name:        pointer.To("Project 1"),
				Description: pointer.To("Project 1 description"),
				CreatedAt:   pointer.To("2020-01-01T00:00:00.000Z"),
			},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Update project with unknown name",
			project: &doppler.ProjectUpdateOptions{
				Name:           "unknown",
				NewName:        "Project 1",
				NewDescription: pointer.To("Project 1 description"),
			},
			wantProject: nil,
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"id is required"},
			},
			wantErr: true,
		},
		{
			name:         "Update project with options validation error",
			project:      &doppler.ProjectUpdateOptions{},
			wantProject:  nil,
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
				err := json.NewEncoder(w).Encode(&doppler.ProjectUpdateResponse{
					APIResponse: tt.wantResponse,
					Project:     tt.wantProject,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &project.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Update method with the test project.
			gotProject, gotResponse, err := client.Update(context.Background(), tt.project)
			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Unexpected error. Expected %t, got %t", tt.wantErr, err != nil)
				return
			}

			// Check if the API response is expected. Ignore the http.Header field, since it's variable.
			if diff := cmp.Diff(tt.wantResponse, gotResponse, cmpopts.IgnoreFields(doppler.APIResponse{}, "Header")); diff != "" {
				t.Errorf("Unexpected API response (-want +got):\n%s", diff)
			}

			// Check if the project is expected.
			if diff := cmp.Diff(tt.wantProject, gotProject); diff != "" {
				t.Errorf("Unexpected project (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProject_Delete(t *testing.T) {
	t.Parallel()

	// Define tests
	tests := []struct {
		name         string
		options      *doppler.ProjectDeleteOptions
		wantResponse doppler.APIResponse
		wantErr      bool
	}{
		{
			name:    "Delete project",
			options: &doppler.ProjectDeleteOptions{Name: "p1"},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(true),
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name:    "Delete project with unknown name",
			options: &doppler.ProjectDeleteOptions{Name: "unknown"},
			wantResponse: doppler.APIResponse{
				Success:    pointer.To(false),
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Messages:   []string{"id is required"},
			},
			wantErr: true,
		},
		{
			name:         "Delete project with options validation error",
			options:      &doppler.ProjectDeleteOptions{},
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
				err := json.NewEncoder(w).Encode(&doppler.ProjectDeleteResponse{
					APIResponse: tt.wantResponse,
				})
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Create a new Doppler client with the httptest.Server URL as base URL.
			client := &project.Client{
				Backend: doppler.GetBackendWithConfig(&doppler.BackendConfig{
					URL: pointer.To(ts.URL),
				}),
				Key: "test",
			}

			// Call the Delete method with the test project.
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
