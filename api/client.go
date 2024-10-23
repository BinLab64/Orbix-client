package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// SideType define side type of order
type SideType string

// Global enums
const (
	SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"
)

// Default Constants
const (
	DefaultBaseURL   = "https://satangcorp.com"
	DefaultUserAgent = "Bin64/1.0"
	DefaultTimeOut   = 10 * time.Second
)

type ClientAuth struct {
	ApiKey    string
	ApiSecret string
}

type Client struct {
	ClientAuth
	HttpClient *http.Client
	BaseURL    string
	UserAgent  string
	Logger     *slog.Logger
}

type ClientOptions struct {
	ClientAuth
	BaseURL   string
	UserAgent string
	Logger    *slog.Logger
}

func newDefaultLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

func NewClient(opts ClientOptions) *Client {
	if opts.BaseURL == "" {
		opts.BaseURL = DefaultBaseURL
	}

	if opts.UserAgent == "" {
		opts.UserAgent = DefaultUserAgent
	}

	if opts.Logger == nil {
		opts.Logger = newDefaultLogger()
	}

	return &Client{
		ClientAuth: opts.ClientAuth,
		HttpClient: http.DefaultClient,
		BaseURL:    opts.BaseURL,
		UserAgent:  opts.UserAgent,
		Logger:     opts.Logger,
	}
}

// parseRequest prepares the request and constructs the full URL.
func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// Set request options from user
	for _, opt := range opts {
		opt(r)
	}

	// Ensure the query param and form are initialized
	r.ensureInitialized()

	c.buildHeader(r)
	c.buildFullURL(r)

	return nil
}

func (c *Client) buildHeader(r *request) {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	r.header = header
}

// Constructs the full URL for the request.
func (c *Client) buildFullURL(r *request) {
	var sb strings.Builder
	sb.WriteString(c.BaseURL)
	sb.WriteString(r.endpoint)

	// Encode the query parameters
	queryString := r.query.Encode()
	if queryString != "" {
		sb.WriteString("?")
		sb.WriteString(queryString)
	}
	r.fullURL = sb.String()

}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, r.method, r.fullURL, r.body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = r.header

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			err = fmt.Errorf("failed to close response body: %w", cerr)
		}
	}()

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	c.Logger.Debug(
		"Orbix API Call",
		slog.String("method", r.method),
		slog.String("url", r.fullURL),
	)

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("API request failed with status %d: %s", res.StatusCode, string(data))
	}
	return data, nil
}

func (c *Client) SetBaseURL(url string) *Client {
	c.BaseURL = url
	return c
}

// Get endpoint status. When status is not ok, it is highly recommended to wait until the status changes back to ok.
//
// GET /api/status
// func (c *Client) NewEndpointStatusService() *EndpointStatusService {
// 	return &EndpointStatusService{c: c}
// }

// Get exchange information
//
// /api/v3/exchangeInfo
func (c *Client) NewExchangeInfoService() *ExchangeInfoService {
	return &ExchangeInfoService{c: c}
}
