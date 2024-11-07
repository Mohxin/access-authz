package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type TestResponse struct {
	StatusCode int
	Body       any
}

type Requester struct {
	httpClient *http.Client
	port       string
}

func NewRequester(port string) Requester {
	return Requester{
		httpClient: new(http.Client),
		port:       port,
	}
}

func (r *Requester) DoRequest(path, method string, body, response any, headers map[string]string) (TestResponse, error) {
	// Encode body if it's provided
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return TestResponse{}, fmt.Errorf("error encoding body: %w", err)
		}
	}

	// Create request with URL and method
	req, err := http.NewRequest(method, r.CreateEndpointURL(path), buf)
	if err != nil {
		return TestResponse{}, fmt.Errorf("error making request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Execute the request
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return TestResponse{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Decode response if expected
	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return TestResponse{}, fmt.Errorf("error decoding response: %w", err)
		}
	}

	return TestResponse{
		StatusCode: resp.StatusCode,
		Body:       response,
	}, nil
}

func (r *Requester) CreateEndpointURL(path string) string {
	return fmt.Sprintf("http://localhost:%s/%s", r.port, path)
}

func (r *Requester) RetryPing(url string) (http.Response, error) {
	retryClient := retryablehttp.NewClient()

	req, _ := retryablehttp.NewRequest("GET", url, nil)

	resp, err := retryClient.Do(req)
	defer CloseIgnore(resp.Body)

	return *resp, err
}

// closes an io.Closer and ignores the returned error
func CloseIgnore(closer io.Closer) {
	_ = closer.Close()
}

func FakeId() int {
	lowRange := 10000000
	hiRange := 99999999
	return lowRange + rand.Intn(hiRange-lowRange)
}

func HostAndPathFromUrl(url string) (string, string) {
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	return req.URL.Host, req.URL.Path
}
