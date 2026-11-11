package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/internal/errutils"
)

// DefaultBaseURL is the Gyro API base URL.
const DefaultBaseURL = "https://www.gyrohosting.com/api/v1/client/"

// Client is a Gyro DNS ACME challenge API client.
type Client struct {
	token      string
	BaseURL    *url.URL
	HTTPClient *http.Client
}

// NewClient creates a new Client.
func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("credentials missing")
	}

	baseURL, _ := url.Parse(DefaultBaseURL)

	return &Client{token: token, BaseURL: baseURL, HTTPClient: &http.Client{Timeout: 30 * time.Second}}, nil
}

// Challenge is the API answer for challenge operations.
type Challenge struct {
	Zone  string `json:"zone"`
	Name  string `json:"name"`
	Ready bool   `json:"ready"`
}

type apiError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// AddTXT adds value to the TXT record set at fqdn (other values stay).
func (c *Client) AddTXT(ctx context.Context, fqdn, value string) error {
	body, _ := json.Marshal(map[string]string{"name": strings.TrimSuffix(fqdn, "."), "value": value})

	_, err := c.do(ctx, http.MethodPost, nil, bytes.NewReader(body))

	return err
}

// RemoveTXT removes value from the TXT record set at fqdn.
func (c *Client) RemoveTXT(ctx context.Context, fqdn, value string) error {
	_, err := c.do(ctx, http.MethodDelete, url.Values{"name": {strings.TrimSuffix(fqdn, ".")}, "value": {value}}, nil)

	return err
}

// Ready reports whether both authoritative name servers publish value.
func (c *Client) Ready(ctx context.Context, fqdn, value string) (bool, error) {
	ch, err := c.do(ctx, http.MethodGet, url.Values{"name": {strings.TrimSuffix(fqdn, ".")}, "value": {value}}, nil)
	if err != nil {
		return false, err
	}

	return ch.Ready, nil
}

func (c *Client) do(ctx context.Context, method string, query url.Values, body io.Reader) (*Challenge, error) {
	endpoint := c.BaseURL.JoinPath("cloud", "dns", "acme-challenge")
	if query != nil {
		endpoint.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), body)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errutils.NewHTTPDoError(req, err)
	}

	defer func() { _ = resp.Body.Close() }()

	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		var e apiError
		if json.Unmarshal(raw, &e) == nil && e.Error.Code != "" {
			return nil, fmt.Errorf("%d: %s: %s", resp.StatusCode, e.Error.Code, e.Error.Message)
		}

		return nil, errutils.NewUnexpectedStatusCodeError(req, resp.StatusCode, raw)
	}

	var out struct {
		Data Challenge `json:"data"`
	}

	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, errutils.NewUnmarshalError(req, resp.StatusCode, raw, err)
	}

	return &out.Data, nil
}
