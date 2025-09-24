package mineru

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OCRRequest struct {
	FileURI string `json:"file_uri"`
	Lang    string `json:"lang,omitempty"`
}

type OCRResponse struct {
	JobID     string `json:"job_id,omitempty"`
	Status    string `json:"status,omitempty"`
	ResultURI string `json:"result_uri,omitempty"`
}

type Client interface {
	CreateJob(ctx context.Context, req OCRRequest) (jobID string, err error)
	GetJob(ctx context.Context, jobID string) (OCRResponse, error)
}

type Option func(*HTTPClient)

type HTTPClient struct {
	baseURL    *url.URL
	httpClient *http.Client
	createPath string
	statusPath string
}

func WithPaths(createPath, statusPath string) Option {
	return func(c *HTTPClient) {
		if strings.TrimSpace(createPath) != "" {
			c.createPath = strings.TrimSpace(createPath)
		}
		if strings.TrimSpace(statusPath) != "" {
			c.statusPath = strings.TrimSpace(statusPath)
		}
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *HTTPClient) {
		if client != nil {
			c.httpClient = client
		}
	}
}

func NewHTTPClient(baseURL string, timeout time.Duration, opts ...Option) (*HTTPClient, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, errors.New("mineru base url is required")
	}
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse mineru base url: %w", err)
	}
	if parsed.Scheme == "" {
		return nil, errors.New("mineru base url must include scheme")
	}

	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	client := &HTTPClient{
		baseURL: parsed,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		createPath: "/api/v1/ocr/jobs",
		statusPath: "/api/v1/ocr/jobs/%s",
	}
	for _, opt := range opts {
		opt(client)
	}
	if !strings.Contains(client.statusPath, "%s") {
		return nil, errors.New("mineru status path must contain %s placeholder for job id")
	}
	return client, nil
}

func (c *HTTPClient) CreateJob(ctx context.Context, req OCRRequest) (string, error) {
	if c == nil || c.httpClient == nil {
		return "", errors.New("mineru client is not initialized")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if strings.TrimSpace(req.FileURI) == "" {
		return "", errors.New("file uri is required")
	}

	endpoint, err := c.buildURL(c.createPath)
	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal mineru request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create mineru request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("call mineru create job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("mineru create job failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var out OCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode mineru response: %w", err)
	}
	if strings.TrimSpace(out.JobID) == "" {
		return "", errors.New("mineru response missing job_id")
	}

	return out.JobID, nil
}

func (c *HTTPClient) GetJob(ctx context.Context, jobID string) (OCRResponse, error) {
	if c == nil || c.httpClient == nil {
		return OCRResponse{}, errors.New("mineru client is not initialized")
	}
	if err := ctx.Err(); err != nil {
		return OCRResponse{}, err
	}
	if strings.TrimSpace(jobID) == "" {
		return OCRResponse{}, errors.New("job id is required")
	}

	endpoint, err := c.buildURL(c.statusPath, jobID)
	if err != nil {
		return OCRResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return OCRResponse{}, fmt.Errorf("create mineru get request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return OCRResponse{}, fmt.Errorf("call mineru get job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return OCRResponse{}, fmt.Errorf("mineru get job failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var out OCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return OCRResponse{}, fmt.Errorf("decode mineru get response: %w", err)
	}
	if strings.TrimSpace(out.JobID) == "" {
		out.JobID = strings.TrimSpace(jobID)
	}
	return out, nil
}

func (c *HTTPClient) buildURL(pattern string, args ...any) (string, error) {
	formatted := pattern
	if len(args) > 0 {
		formatted = fmt.Sprintf(pattern, args...)
	}
	formatted = strings.TrimSpace(formatted)
	if formatted == "" {
		return c.baseURL.String(), nil
	}
	if strings.HasPrefix(formatted, "http://") || strings.HasPrefix(formatted, "https://") {
		return formatted, nil
	}
	rel, err := url.Parse(formatted)
	if err != nil {
		return "", fmt.Errorf("parse mineru url path: %w", err)
	}
	return c.baseURL.ResolveReference(rel).String(), nil
}
