package doppler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"

	"github.com/nikoksr/doppler-go/logging"
	"github.com/nikoksr/doppler-go/pointer"
)

const (
	// SDKVersion is the version of the SDK
	SDKVersion = "0.2.0"

	// APIURL is the base URL for the API
	APIURL string = "https://api.doppler.com"

	// UnknownPlatform is the platform name for unknown platforms
	UnknownPlatform = "unknown platform"

	// defaultHTTPTimeout is the default HTTP timeout
	defaultHTTPTimeout = 60 * time.Second

	// headerXRequestID is the name of the header containing the request ID
	headerXRequestID = "X-Request-Id"

	// headerRateLimitLimit is the name of the header containing the rate limit
	headerRateLimitLimit = "X-RateLimit-Limit"

	// headerRateLimitRemaining is the name of the header containing the remaining rate limit
	headerRateLimitRemaining = "X-RateLimit-Remaining"

	// headerRateLimitReset is the name of the header containing the rate limit reset time
	headerRateLimitReset = "X-RateLimit-Reset"
)

// Key is the API key used to authenticate with the API
var Key string

// EnableValidation enables validation of the payload. This is enabled by default. If you want to disable validation,
// set this to false.
var EnableValidation = true

// Internal standard logger and validator.
var (
	stdValidator = validator.New()
)

// BackendConfig is the configuration for the backend.
type BackendConfig struct {
	// Client is the HTTP client to use for requests. If nil, a default client will be used.
	Client *http.Client

	// URL is the base URL for the API. If empty, the default URL will be used.
	URL *string

	// Logger is the logger to use for logging. If nil, a noop logger will be used.
	Logger logging.Logger
}

// Backend is the backend used by the SDK. It is used to make requests to the API.
type Backend interface {
	Call(ctx context.Context, req *Request, resp Response) error
	CallRaw(ctx context.Context, req *Request) (*http.Response, error)
}

// backendImplementation is the default backend implementation. It satisfies the Backend interface.
type backendImplementation struct {
	URL        string
	HTTPClient *http.Client
	Logger     logging.Logger
}

// Compile-time check to ensure that backendImplementation implements the Backend interface.
var _ Backend = (*backendImplementation)(nil)

// Request is the base request type all Backend calls. Internally, it is converted to an HTTP request and sent to the
// API.
type Request struct {
	// Header are the headers to send with the request.
	Header http.Header `json:"-"`

	// Method is the HTTP method to use. e.g. "GET"
	Method string `json:"method"`

	// Path is the path to the API endpoint. e.g. "/projects"
	Path string `json:"path"`

	// Key is the API key to use.
	Key string `json:"-"`

	// Payload is expected to be a struct holding all necessary query and body parameters. Every field in the struct
	// is expected to be clearly tagged with the corresponding API parameter name and either "url" or "json", not both.
	// The "url" tag is used for query parameters, the "json" tag is used for body parameters.
	//
	// Example:
	//  type ExamplePayload struct {
	//  	// Query parameter
	//  	ProjectID string `url:"project_id" json:"-"`
	//
	//  	// Body parameter
	//  	ProjectName string `url:"-" json:"project_name"`
	//  }
	//
	//  As you can see, we have to be explicit about which parameters are query and which are body parameters. One of
	//  the parameter types ("url" or "json") must be set to "-". If this is not the case, the request may potentially
	//  be malformed and fail. However, we do not check for this, since there may be cases where you want to send a
	//  parameter in both the query and the body.
	//
	//  The above example will result in a request with the following query parameters:
	//      ?project_id=123
	//
	//  And the following body:
	//      { "project_name": "my project" }
	//
	Payload any `json:"payload,omitempty"`
}

// Response is the base response type all Backend calls. It's meant to bind the response body and parts of the HTTP
// response, hence the required WithDetails method.
type Response interface {
	WithDetails(resp *http.Response)
	Error() error
}

// RateLimit is the ratelimit information returned by the API.
type RateLimit struct {
	// Limit is the maximum number of requests allowed per period.
	Limit int `json:"limit"`

	// Remaining is the number of requests remaining in the current period.
	Remaining int `json:"remaining"`

	// Reset is the time when the current period ends.
	Reset time.Time `json:"reset"`
}

// APIResponse is the base response type for all Doppler API responses.
type APIResponse struct {
	// Header is the HTTP response header of the response.
	Header http.Header `json:"header,omitempty"`

	// RequestID is the ID of the request. It can be used to identify the request in the logs; useful for debugging.
	RequestID string `json:"request_id,omitempty"`

	// RateLimit is the ratelimit information returned by the API.
	RateLimit *RateLimit `json:"ratelimit,omitempty"`

	// Status is the HTTP status code of the response, e.g. 200 OK
	Status string `json:"status,omitempty"`

	// StatusCode is the HTTP status code of the response, e.g. 200
	StatusCode int `json:"status_code,omitempty"`

	// Success is true if the request was successful. This gets set by Doppler.
	Success *bool `json:"success,omitempty"`

	// Messages is a list of potential messages from the Doppler API.
	Messages []string `json:"messages,omitempty"`

	// Page is the current page of results.
	Page *int `json:"page,omitempty"`
}

func extractRateLimitFromHeader(header http.Header) *RateLimit {
	var rateLimit RateLimit
	var err error

	// Try to convert all three values to their respective types. If any of them fail, we return nil. This is because
	// the API may not return all three values, in which case we don't want to return a partial RateLimit.

	// Limit
	limit := header.Get(headerRateLimitLimit)
	if rateLimit.Limit, err = strconv.Atoi(limit); err != nil {
		return nil
	}

	// Remaining
	remaining := strings.TrimSpace(header.Get(headerRateLimitRemaining))
	if rateLimit.Remaining, err = strconv.Atoi(remaining); err != nil {
		return nil
	}

	// Reset; its value is a Unix timestamp
	reset := strings.TrimSpace(header.Get(headerRateLimitReset))
	var resetUnix int64
	if resetUnix, err = strconv.ParseInt(reset, 10, 64); err != nil {
		return nil
	}

	// Convert the Unix timestamp to a time.Time
	rateLimit.Reset = time.Unix(resetUnix, 0)

	return &rateLimit
}

// WithDetails binds important details from the HTTP response to the APIResponse.
func (r *APIResponse) WithDetails(resp *http.Response) {
	if resp == nil {
		return
	}

	r.Status = resp.Status
	r.StatusCode = resp.StatusCode

	if resp.Header != nil {
		r.Header = resp.Header
		r.RequestID = resp.Header.Get(headerXRequestID)
		r.RateLimit = extractRateLimitFromHeader(resp.Header)
	}
}

// Error checks if the APIResponse contains any errors. If so, it returns an error containing all messages.
func (r *APIResponse) Error() error {
	if len(r.Messages) > 0 {
		return errors.New(strings.Join(r.Messages, ": "))
	}

	return nil
}

// defaultClient is the default HTTP client used by the SDK.
var defaultClient = &http.Client{
	Timeout: defaultHTTPTimeout,
}

// normalizeURL returns a sanitized URL for the API. If the URL is empty, it does nothing.
func normalizeURL(url string) string {
	url = strings.TrimSuffix(url, "/")

	// Trim current API major version
	url = strings.TrimSuffix(url, "/v3")

	// Necessary for share endpoints
	url = strings.TrimSuffix(url, "/v1")

	return url
}

// newBackendImplementation returns a new backend implementation.
func newBackendImplementation(config *BackendConfig) Backend {
	// HTTP client
	if config.Client == nil {
		config.Client = defaultClient
	}

	// Base API URL
	if config.URL == nil {
		config.URL = pointer.To(APIURL)
	}
	config.URL = pointer.To(normalizeURL(*config.URL))

	// Logger
	if config.Logger == nil {
		config.Logger = &logging.NopLogger{}
	}

	return &backendImplementation{
		HTTPClient: config.Client,
		URL:        *config.URL,
		Logger:     config.Logger,
	}
}

// GetBackendWithConfig returns a new backend with the given configuration. This is the preferred way to create a new
// backend.
func GetBackendWithConfig(config *BackendConfig) Backend {
	return newBackendImplementation(config)
}

// GetBackend returns a new backend with the default configuration.
func GetBackend() Backend {
	return GetBackendWithConfig(&BackendConfig{Client: defaultClient})
}

func (req *Request) getQueryParameters() (parameters, error) {
	return extractQueryParameters(req.Payload) // Checks for nil
}

func (req *Request) getBody() (*bytes.Buffer, error) {
	encodedBody := new(bytes.Buffer)
	if req.Payload == nil {
		return encodedBody, nil
	}

	err := json.NewEncoder(encodedBody).Encode(&req.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "encode request body")
	}

	return encodedBody, nil
}

// isPayloadValid uses the go-playground validator to validate the payload. It is expected that the payload is a struct
// and validation tags are used to define the validation rules. Using this function to make the code more readable.
// Validation is skipped if EnableValidation is false or the payload is nil.
func isPayloadValid(payload any) error {
	if !EnableValidation || payload == nil {
		return nil
	}

	return stdValidator.Struct(payload)
}

// prepareRequest creates a new HTTP request from the given Request. The returned request is ready to be sent to the
// API.
func (b *backendImplementation) prepareRequest(ctx context.Context, req *Request) (*http.Request, error) {
	// Validate the request's payload before we do anything else.
	if err := isPayloadValid(req.Payload); err != nil {
		return nil, errors.Wrap(err, "validate request payload")
	}

	// Normalize URL
	if !strings.HasPrefix(req.Path, "/") {
		req.Path = "/" + req.Path
	}
	req.Path = b.URL + req.Path

	// Create basic HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.Path, nil)
	if err != nil {
		return nil, err
	}

	// Parse query parameters from request payload
	params, err := req.getQueryParameters()
	if err != nil {
		return nil, errors.Wrap(err, "get query parameters from request payload")
	}
	if params != nil {
		httpReq.URL.RawQuery = params.Encode()
	}

	// If it's not a GET request, we need to parse the body from the request payload and set the content type.
	if req.Method != http.MethodGet {
		body, err := req.getBody()
		if err != nil {
			return nil, errors.Wrap(err, "get body from request payload")
		}

		httpReq.Body = nopReadCloser{body}
		httpReq.Header.Add("Content-Type", "application/json")
	}

	// Set headers
	httpReq.SetBasicAuth(req.Key, "")
	httpReq.Header.Add("Accept", "application/json")
	httpReq.Header.Add("User-Agent", encodedUserAgent)

	// Set custom headers; doing this last so that we can override the default headers
	for key, values := range req.Header {
		for _, value := range values {
			httpReq.Header.Add(key, value)
		}
	}

	return httpReq, nil
}

func (b *backendImplementation) call(ctx context.Context, req *Request) (*http.Response, error) {
	// Translate our internal request to an HTTP request
	httpReq, err := b.prepareRequest(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "prepare request")
	}

	b.Logger.Infow("Sending HTTP request", "method", httpReq.Method, "url", httpReq.URL.String())

	return b.HTTPClient.Do(httpReq)
}

// CallRaw sends the given request to the API and returns the raw HTTP response. This is useful if you want to handle
// the response yourself. Otherwise, you should use Call. The returned response is not closed, so you need to close it
// yourself.
func (b *backendImplementation) CallRaw(ctx context.Context, req *Request) (*http.Response, error) {
	return b.call(ctx, req)
}

// Call sends the given request to the API and returns the parsed response. It does the same as CallRaw, but it also
// parses the response body and closes the response. This is the preferred way to send requests to the API, unless you
// need to handle the response yourself. If the response contains any errors, it returns an error containing all
// messages. The target Response (gotResponse) may be nil, in which case we skip parsing the response body completely.
func (b *backendImplementation) Call(ctx context.Context, req *Request, resp Response) error {
	httpResp, err := b.call(ctx, req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	// If the response body or target response is nil, we can return early.
	if httpResp.ContentLength == 0 || httpResp == nil || resp == nil || reflect.ValueOf(resp).IsNil() {
		return nil
	}

	// Attach details to the response object
	resp.WithDetails(httpResp)

	// Handle binding the response body to the response object (if it's not nil) based on the content type.
	if strings.HasPrefix(httpResp.Header.Get("Content-Type"), "application/json") {
		err = json.NewDecoder(httpResp.Body).Decode(resp)
		if err != nil {
			return errors.Wrap(err, "decode response body")
		}
	} else {
		b.Logger.Warnw("Response body is not JSON", "content-type", httpResp.Header.Get("Content-Type"))
	}

	// Check for errors in the response. This checks for API specific errors hidden in the Messages field.
	if rerr := resp.Error(); rerr != nil {
		if err == nil {
			err = rerr
		} else {
			err = errors.Wrap(err, rerr.Error())
		}
	}

	// If we have an error, log some information about the request and response.
	if err != nil {
		reqJSON, _ := json.Marshal(req)
		respJSON, _ := json.Marshal(resp)
		b.Logger.Debugw("HTTP request failed", "request", string(reqJSON), "response", string(respJSON))
	}

	return err
}

// AppInfo contains information about the "app" which this integration belongs to.
type AppInfo struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Version string `json:"version"`
}

// formatUserAgent formats an AppInfo in a way that's suitable to be appended to a User-Agent string. Note that this
// format is shared between all libraries so if it's changed, it should be changed everywhere.
func (a *AppInfo) formatUserAgent() string {
	userAgent := a.Name
	if a.Version != "" {
		userAgent += "/" + a.Version
	}
	if a.URL != "" {
		userAgent += " (" + a.URL + ")"
	}

	return userAgent
}

var (
	appInfo          *AppInfo
	encodedUserAgent string
)

// SetAppInfo sets the information about the "app" which this integration belongs to.
func SetAppInfo(info *AppInfo) {
	if info != nil && info.Name == "" {
		panic("info.NewConfig must not be empty")
	}

	appInfo = info

	// We need to re-init since we have a new app info.
	initUserAgent()
}

// initUserAgent initializes the encodedDopplerUserAgent and encodedUserAgent variables.
func initUserAgent() {
	encodedUserAgent = "doppler-go/" + SDKVersion
	if appInfo != nil {
		encodedUserAgent += " " + appInfo.formatUserAgent()
	}
}

func init() {
	initUserAgent()
}

// getUname returns a string containing the uname information. This is used to add additional debugging information. It
// returns UnknownPlatform if the uname command is not available. This is not a problem since the uname command is only
// used for debugging. The uname command is not available on all platforms.
func getUname() string {
	path, err := exec.LookPath("uname")
	if err != nil {
		return UnknownPlatform
	}

	cmd := exec.Command(path, "-a")
	var out bytes.Buffer
	cmd.Stderr = nil // goes to os.DevNull
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return UnknownPlatform
	}

	return out.String()
}

// Compile time check to make sure the type implements the interface.
var _ io.ReadCloser = nopReadCloser{}

// nopReadCloser is an implementation of `io.ReadCloser` that wraps an `io.Reader`. This does not alter the underlying
// `io.Reader`'s behavior. It just adds a `Close` method that does nothing. This is needed to make `http.Request`'s
// `Body` method work.
type nopReadCloser struct {
	io.Reader
}

// Close does nothing. It's here to satisfy the `io.ReadCloser` interface.
func (nopReadCloser) Close() error { return nil }
