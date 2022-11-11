package doppler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nikoksr/doppler-go/logging"
	"github.com/nikoksr/doppler-go/pointer"
)

func TestAppInfo_formatUserAgent(t *testing.T) {
	t.Parallel()

	type fields struct {
		Name    string
		Version string
		URL     string
	}
	cases := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Test formatUserAgent",
			fields: fields{Name: "doppler", Version: "0.1.0", URL: "https://example.com"},
			want:   "doppler/0.1.0 (https://example.com)",
		},
		{
			name:   "Test formatUserAgent",
			fields: fields{Name: "doppler", Version: "0.1.0"},
			want:   "doppler/0.1.0",
		},
		{
			name:   "Test formatUserAgent",
			fields: fields{Name: "doppler"},
			want:   "doppler",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			a := &AppInfo{
				Name:    tc.fields.Name,
				Version: tc.fields.Version,
				URL:     tc.fields.URL,
			}
			if got := a.formatUserAgent(); got != tc.want {
				t.Errorf("AppInfo.formatUserAgent() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSetAppInfo(t *testing.T) {
	t.Parallel()

	type args struct {
		info *AppInfo
	}
	cases := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name:      "Test SetAppInfo",
			args:      args{info: &AppInfo{Name: "doppler", Version: "0.1.0"}},
			wantPanic: false,
		},
		{
			name:      "Test SetAppInfo",
			args:      args{info: &AppInfo{Name: "doppler", Version: ""}},
			wantPanic: false,
		},
		{
			name:      "Test SetAppInfo",
			args:      args{info: &AppInfo{Name: "", Version: "0.1.0"}},
			wantPanic: true,
		},
		{
			name:      "Test SetAppInfo",
			args:      args{info: &AppInfo{Name: "", Version: ""}},
			wantPanic: true,
		},
	}
	for _, tc := range cases { //nolint:paralleltest // Test is accessing global variable
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("SetAppInfo() did not panic")
					}
				}()
			}
			SetAppInfo(tc.args.info)
		})
	}
}

func Test_newBackendImplementation(t *testing.T) {
	t.Parallel()

	const localURL = "http://localhost:8080"

	cases := []struct {
		name   string
		config *BackendConfig
		want   Backend
	}{
		{
			name:   "Test newBackendImplementation 1",
			config: &BackendConfig{Client: &http.Client{}, URL: pointer.To(localURL), Logger: &logging.NopLogger{}},
			want:   &backendImplementation{HTTPClient: &http.Client{}, URL: localURL, Logger: &logging.NopLogger{}},
		},
		{
			name:   "Test newBackendImplementation 2",
			config: &BackendConfig{Client: defaultClient, URL: pointer.To(localURL)},
			want:   &backendImplementation{HTTPClient: defaultClient, URL: localURL, Logger: &logging.NopLogger{}},
		},
		{
			name:   "Test newBackendImplementation 3",
			config: &BackendConfig{Client: defaultClient, URL: nil, Logger: nil},
			want:   &backendImplementation{HTTPClient: defaultClient, URL: APIURL, Logger: &logging.NopLogger{}},
		},
		{
			name:   "Test newBackendImplementation 4",
			config: &BackendConfig{Client: nil, URL: nil, Logger: nil},
			want:   &backendImplementation{HTTPClient: defaultClient, URL: APIURL, Logger: &logging.NopLogger{}},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := newBackendImplementation(tc.config)
			if got == nil {
				t.Fatal("newBackendImplementation() returned nil")
			}

			diff := cmp.Diff(got, tc.want)
			if diff != "" {
				t.Errorf("newBackendImplementation failed:\n%s", diff)
			}
		})
	}
}

func TestGetBackend(t *testing.T) {
	t.Parallel()

	want := &backendImplementation{
		HTTPClient: defaultClient,
		URL:        APIURL,
		Logger:     &logging.NopLogger{},
	}

	got := GetBackend()

	if got == nil {
		t.Fatal("GetBackend() returned nil")
	}

	diff := cmp.Diff(want, got)
	if diff != "" {
		t.Errorf("GetBackend failed:\n%s", diff)
	}
}

func TestNopReadCloser_Close(t *testing.T) {
	t.Parallel()

	// This test basically does nothing, but it's here to cheat up the test coverage.

	err := nopReadCloser{}.Close()
	if err != nil {
		t.Errorf("nopReadCloser.Close() = %v, want nil", err)
	}
}

func Test_getUname(t *testing.T) {
	t.Parallel()

	if strings.EqualFold(getUname(), "") {
		t.Fatal("getUname() returned empty string")
	}
}

func Test_extractRateLimitFromHeader(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		header http.Header
		want   *RateLimit
	}{
		{
			name: "Extract valid rate limit",
			header: http.Header{
				textproto.CanonicalMIMEHeaderKey(headerRateLimitLimit):     []string{"100"},
				textproto.CanonicalMIMEHeaderKey(headerRateLimitRemaining): []string{"50"},
				textproto.CanonicalMIMEHeaderKey(headerRateLimitReset):     []string{"1234567890"},
			},
			want: &RateLimit{
				Limit:     100,
				Remaining: 50,
				Reset:     time.Unix(1234567890, 0),
			},
		},
		{
			name: "Extract invalid rate limit - missing header",
			header: http.Header{
				textproto.CanonicalMIMEHeaderKey(headerRateLimitRemaining): []string{"50"},
				textproto.CanonicalMIMEHeaderKey(headerRateLimitReset):     []string{"1234567890"},
			},
			want: nil,
		},
		{
			name: "Extract invalid rate limit - invalid header 1",
			header: http.Header{
				textproto.CanonicalMIMEHeaderKey(headerRateLimitLimit):     []string{"invalid"},
				textproto.CanonicalMIMEHeaderKey(headerRateLimitRemaining): []string{"50"},
				textproto.CanonicalMIMEHeaderKey(headerRateLimitReset):     []string{"1234567890"},
			},
			want: nil,
		},
		{
			name: "Extract invalid rate limit - invalid header 2",
			header: http.Header{
				textproto.CanonicalMIMEHeaderKey(headerRateLimitLimit):     []string{"100"},
				textproto.CanonicalMIMEHeaderKey(headerRateLimitRemaining): []string{"50"},
				textproto.CanonicalMIMEHeaderKey(headerRateLimitReset):     []string{"invalid"},
			},
			want: nil,
		},
		{
			name: "Extract invalid rate limit - invalid header 3",
			header: http.Header{
				textproto.CanonicalMIMEHeaderKey(headerRateLimitLimit):     []string{"100"},
				textproto.CanonicalMIMEHeaderKey(headerRateLimitRemaining): []string{"invalid"},
				textproto.CanonicalMIMEHeaderKey(headerRateLimitReset):     []string{"1234567890"},
			},
			want: nil,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := extractRateLimitFromHeader(tc.header)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("extractRateLimitFromHeader() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAPIResponse(t *testing.T) {
	t.Parallel()

	// Test basic APIResponse functionality. First, try assigning a nil http.Response.
	resp := &APIResponse{}
	var fakeHTTPResponse *http.Response

	resp.WithDetails(fakeHTTPResponse)

	// gotResponse should be unchanged.
	if diff := cmp.Diff(&APIResponse{}, resp); diff != "" {
		t.Errorf("APIResponse.WithDetails() mismatch (-want +got):\n%s", diff)
	}

	// Now, set up a fake http.Response and try again.
	fakeHTTPResponse = &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Header: http.Header{
			textproto.CanonicalMIMEHeaderKey(headerRateLimitLimit):     []string{"100"},
			textproto.CanonicalMIMEHeaderKey(headerRateLimitRemaining): []string{"50"},
			textproto.CanonicalMIMEHeaderKey(headerRateLimitReset):     []string{"1234567890"},
			textproto.CanonicalMIMEHeaderKey(headerXRequestID):         []string{"1234567890"},
		},
	}

	// Bind the http response
	resp.WithDetails(fakeHTTPResponse)

	// Validate that the fields of the APIResponse are set correctly.
	if resp.StatusCode != 200 {
		t.Fatal("APIResponse.StatusCode is not set correctly")
	}
	if resp.Status != "200 OK" {
		t.Fatal("APIResponse.Status is not set correctly")
	}
	if resp.RateLimit.Limit != 100 {
		t.Fatal("APIResponse.RateLimit.Limit is not set correctly")
	}
	if resp.RateLimit.Remaining != 50 {
		t.Fatal("APIResponse.RateLimit.Remaining is not set correctly")
	}
	if resp.RateLimit.Reset != time.Unix(1234567890, 0) {
		t.Fatal("APIResponse.RateLimit.Reset is not set correctly")
	}
	if resp.RequestID != "1234567890" {
		t.Fatal("APIResponse.RequestID is not set correctly")
	}
	if diff := cmp.Diff(resp.Header, fakeHTTPResponse.Header); diff != "" {
		t.Errorf("APIResponse.Header mismatch (-want +got):\n%s", diff)
	}

	// Now add fake response data from the Doppler API
	resp.Success = pointer.To(true)
	resp.Page = pointer.To(1)
	resp.Messages = []string{"message1", "message2"}

	// Validate that the fields of the APIResponse are set correctly.
	if !*resp.Success {
		t.Fatal("APIResponse.Success is not set correctly")
	}
	if *resp.Page != 1 {
		t.Fatal("APIResponse.Page is not set correctly")
	}
	if diff := cmp.Diff(resp.Messages, []string{"message1", "message2"}); diff != "" {
		t.Errorf("APIResponse.Messages mismatch (-want +got):\n%s", diff)
	}

	// The response contains messages from the Doppler API, so Error() should generate an error.
	if resp.Error() == nil {
		t.Fatal("APIResponse.Error() should return an error")
	}

	// Reset list of messages and validate that Error() returns nil.
	resp.Messages = []string{}
	if resp.Error() != nil {
		t.Fatal("APIResponse.Error() should return nil")
	}
}

func TestRequest(t *testing.T) {
	t.Parallel()

	// Stub a base request object.
	req := &Request{
		Header: http.Header{},
		Method: "GET",
		Path:   "/v3/workplace",
		Key:    "fake-key",
	}

	// In the following we'll test various payloads.
	// First, a nil payload. This should not cause any issues.
	req.Payload = nil

	// First, try extracting the query parameters from the request.
	params, err := req.getQueryParameters()
	if err != nil {
		t.Fatalf("Request.getQueryParameters() returned an error: %v", err)
	}
	if len(params) != 0 {
		t.Fatalf("Request.getQueryParameters() returned %d parameters, expected 0", len(params))
	}

	// Now, try extracting the request body.
	body, err := req.getBody()
	if err != nil {
		t.Fatalf("Request.getBody() returned an error: %v", err)
	}
	// Body should be an empty bytes.Buffer.
	if body == nil || body.Len() != 0 {
		t.Fatalf("Request.getBody() returned a non-empty body")
	}

	// Now, set up an actual payload and try again. Payloads are usually Option objects for API endpoints, so we'll
	// take SecretGetOptions as an example. The fields are all query parameters.
	getOptions := &SecretGetOptions{
		Project: "fake-project",
		Config:  "fake-config",
		Name:    "fake-name",
	}
	req.Payload = getOptions

	// Try extracting the query parameters from the request.
	params, err = req.getQueryParameters()
	if err != nil {
		t.Fatalf("Request.getQueryParameters() returned an error: %v", err)
	}
	if len(params) != 3 {
		t.Fatalf("Request.getQueryParameters() returned %d parameters, expected 3", len(params))
	}
	// Compare the extracted parameters with the expected values.
	if params.Get("project") != getOptions.Project {
		t.Fatalf("Request.getQueryParameters() returned an unexpected value for project: %s", params.Get("project"))
	}
	if params.Get("config") != getOptions.Config {
		t.Fatalf("Request.getQueryParameters() returned an unexpected value for config: %s", params.Get("config"))
	}
	if params.Get("name") != getOptions.Name {
		t.Fatalf("Request.getQueryParameters() returned an unexpected value for name: %s", params.Get("name"))
	}

	// Request body should be empty.
	body, err = req.getBody()
	if err != nil {
		t.Fatalf("Request.getBody() returned an error: %v", err)
	}
	// Body should be an empty map. Decode it and compare with the expected values.
	var rawBody map[string]string
	_ = json.NewDecoder(body).Decode(&rawBody)
	if len(rawBody) != 0 {
		t.Fatalf("Request.getBody() returned a non-empty body")
	}

	// Now, set up a payload that's an Options object with only body parameters.
	updateOptions := &SecretUpdateOptions{
		Project: "fake-project",
		Config:  "fake-config",
		NewSecrets: map[string]string{
			"fake-name": "fake-value",
		},
	}

	req.Payload = updateOptions

	// Try extracting the query parameters from the request. This should not return any parameters.
	params, err = req.getQueryParameters()
	if err != nil {
		t.Fatalf("Request.getQueryParameters() returned an error: %v", err)
	}
	if len(params) != 0 {
		t.Fatalf("Request.getQueryParameters() returned %d parameters, expected 0", len(params))
	}

	// Request body should not be empty.
	body, err = req.getBody()
	if err != nil {
		t.Fatalf("Request.getBody() returned an error: %v", err)
	}

	// Decode the body into a doppler.SecretUpdateOptions object.
	var decodedUpdateOptions SecretUpdateOptions
	if err = json.NewDecoder(body).Decode(&decodedUpdateOptions); err != nil {
		t.Fatalf("Request.getBody() returned an invalid body: %v", err)
	}

	// Compare the decoded object with the original.
	if diff := cmp.Diff(decodedUpdateOptions, *updateOptions); diff != "" {
		t.Errorf("Request.getBody() returned an unexpected body (-want +got):\n%s", diff)
	}

	// Lastly, test the JSON parsing itself. Make the internal JSON encoder fail by setting a non-serializable value.
	req.Payload = make(chan int)
	_, err = req.getBody()
	if err == nil {
		t.Fatal("Request.getBody() should return an error")
	}
}

func Test_normalizeURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		originalURL string
		wantURL     string
	}{
		{
			name:        "normalizeURL - no change",
			originalURL: "https://api.doppler.com",
			wantURL:     "https://api.doppler.com",
		},
		{
			name:        "normalizeURL - trailing slash",
			originalURL: "https://api.doppler.com/",
			wantURL:     "https://api.doppler.com",
		},
		{
			name:        "normalizeURL - trailing major version v1",
			originalURL: "https://api.doppler.com/v1",
			wantURL:     "https://api.doppler.com",
		},
		{
			name:        "normalizeURL - trailing slash and major version v1",
			originalURL: "https://api.doppler.com/v1/",
			wantURL:     "https://api.doppler.com",
		},
		{
			name:        "normalizeURL - trailing major version v3",
			originalURL: "https://api.doppler.com/v3",
			wantURL:     "https://api.doppler.com",
		},
		{
			name:        "normalizeURL - trailing slash and major version v3",
			originalURL: "https://api.doppler.com/v3/",
			wantURL:     "https://api.doppler.com",
		},
		{
			name:        "normalizeURL - empty string",
			originalURL: "",
			wantURL:     "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if gotURL := normalizeURL(tt.originalURL); gotURL != tt.wantURL {
				t.Errorf("normalizeURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func Test_isPayloadValid(t *testing.T) {
	t.Parallel()

	// First, test a nil payload.
	if err := isPayloadValid(nil); err != nil {
		t.Fatalf("isPayloadValid(nil) returned an error: %v", err)
	}

	// Test a valid payload. All three fields of SecretGetOptions are required. The validator should not return an error.
	getOptions := &SecretGetOptions{
		Project: "fake-project",
		Config:  "fake-config",
		Name:    "fake-name",
	}
	if err := isPayloadValid(getOptions); err != nil {
		t.Fatalf("isPayloadValid() returned an error: %v", err)
	}

	// Now, remove the project field. The validator should return an error.
	getOptions.Project = ""
	if err := isPayloadValid(getOptions); err == nil {
		t.Fatal("isPayloadValid() should return an error")
	}
}

func Test_Call(t *testing.T) {
	t.Parallel()

	// This functionality is broadly covered already by all the client packages, so we'll just test the basics here.

	tests := []struct {
		name    string
		req     *Request
		resp    *SecretGetResponse
		respOut *SecretGetResponse
		wantErr bool
	}{
		{
			name: "Call - success",
			req: &Request{
				Method: "GET",
				Path:   "/fake-path",
				Key:    "fake-key",
				Payload: &SecretGetOptions{
					Project: "fake-project",
					Config:  "fake-config",
					Name:    "fake-name",
				},
			},
			resp: &SecretGetResponse{
				Secret: &Secret{
					Name:  pointer.To("fake-name"),
					Value: &SecretValue{Raw: pointer.To("fake-value"), Computed: pointer.To("fake-value")},
				},
				APIResponse: APIResponse{
					StatusCode: 200,
					Status:     "OK",
					RequestID:  "fake-request-id",
					RateLimit:  &RateLimit{Limit: 100, Remaining: 99, Reset: time.Unix(0, 0)},
					Success:    pointer.To(true),
				},
			},
			respOut: &SecretGetResponse{},
			wantErr: false,
		},
		{
			name: "Call - Fake API error",
			req: &Request{
				Method: "GET",
				Path:   "/fake-path",
				Key:    "fake-key",
				Payload: &SecretGetOptions{
					Project: "fake-project",
					Config:  "fake-config",
					Name:    "fake-name",
				},
			},
			resp: &SecretGetResponse{
				APIResponse: APIResponse{
					StatusCode: 400,
					Status:     "Bad Request",
					RequestID:  "fake-request-id",
					RateLimit:  &RateLimit{Limit: 100, Remaining: 99, Reset: time.Unix(0, 0)},
					Success:    pointer.To(false),
					Messages:   []string{"fake-error"},
				},
			},
			respOut: &SecretGetResponse{},
			wantErr: true,
		},
		{
			name: "Call - invalid payload",
			req: &Request{
				Method: "GET",
				Path:   "/fake-path",
				Key:    "fake-key",
				Payload: &SecretGetOptions{
					Project: "fake-project",
					Config:  "fake-config",
				},
			},
			resp:    &SecretGetResponse{},
			respOut: &SecretGetResponse{},
			wantErr: true,
		},
		{
			name: "Call - success, but no out-response",
			req: &Request{
				Method: "GET",
				Path:   "/fake-path",
				Key:    "fake-key",
				Payload: &SecretGetOptions{
					Project: "fake-project",
					Config:  "fake-config",
					Name:    "fake-name",
				},
			},
			resp: &SecretGetResponse{
				Secret: &Secret{
					Name:  pointer.To("fake-name"),
					Value: &SecretValue{Raw: pointer.To("fake-value"), Computed: pointer.To("fake-value")},
				},
				APIResponse: APIResponse{
					StatusCode: 200,
					Status:     "OK",
					RequestID:  "fake-request-id",
					RateLimit:  &RateLimit{Limit: 100, Remaining: 99, Reset: time.Unix(0, 0)},
					Success:    pointer.To(true),
				},
			},
			respOut: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a fake HTTP server to handle the request.
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set(headerXRequestID, tt.resp.RequestID)

				if tt.resp.RateLimit != nil {
					limit := strconv.Itoa(tt.resp.RateLimit.Limit)
					remaining := strconv.Itoa(tt.resp.RateLimit.Remaining)
					reset := strconv.FormatInt(tt.resp.RateLimit.Reset.Unix(), 10)

					w.Header().Set(headerRateLimitLimit, limit)
					w.Header().Set(headerRateLimitRemaining, remaining)
					w.Header().Set(headerRateLimitReset, reset)
				}

				w.WriteHeader(tt.resp.StatusCode)
				err := json.NewEncoder(w).Encode(tt.resp)
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer server.Close()

			// Create a client with the fake server URL.
			client := GetBackendWithConfig(&BackendConfig{
				URL: pointer.To(server.URL),
			})

			// Call the API.
			err := client.Call(context.Background(), tt.req, tt.respOut)

			// Check if the error is expected.
			if (err != nil) != tt.wantErr {
				t.Errorf("Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// We're done if we didn't want the response to be parsed.
			if tt.respOut == nil {
				return
			}

			// Check if the response is expected. Ignore the Headers field.
			if diff := cmp.Diff(tt.respOut, tt.resp, cmpopts.IgnoreFields(APIResponse{}, "Header")); diff != "" {
				t.Fatalf("Call() response mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
