package dugble

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

const DefaultBaseURL = "https://api.dugble.com"

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	apiKey  string
	baseURL string
	http    HTTPDoer

	Topics     *TopicsService
	Segments   *SegmentsService
	Domains    *DomainsService
	Emails     *EmailsService
	SMS        *SmsService
	Templates  *TemplatesService
	Broadcasts *BroadcastsService
}

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

func WithHTTPClient(client HTTPDoer) Option {
	return func(c *Client) {
		if client != nil {
			c.http = client
		}
	}
}

func New(apiKey string, options ...Option) *Client {
	client := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, option := range options {
		option(client)
	}
	client.Topics = &TopicsService{client: client}
	client.Segments = &SegmentsService{client: client}
	client.Domains = &DomainsService{client: client}
	client.Emails = &EmailsService{client: client}
	client.SMS = &SmsService{client: client}
	versions := &TemplateVersionsService{client: client}
	client.Templates = &TemplatesService{client: client, Versions: versions}
	client.Broadcasts = &BroadcastsService{client: client}
	return client
}

type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	RequestID  string `json:"-"`
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Code == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (c *Client) request(ctx context.Context, method, path string, body any, out any, headers http.Header) error {
	endpoint, err := url.Parse(c.baseURL + path)
	if err != nil {
		return err
	}

	var reader io.Reader
	if body != nil {
		payload, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			return marshalErr
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	var raw struct {
		Success bool            `json:"success"`
		Data    json.RawMessage `json:"data"`
		Error   *APIError       `json:"error"`
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&raw); err != nil {
		if errors.Is(err, io.EOF) && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 || !raw.Success {
		apiErr := raw.Error
		if apiErr == nil {
			apiErr = &APIError{Message: http.StatusText(resp.StatusCode)}
		}
		apiErr.StatusCode = resp.StatusCode
		apiErr.RequestID = resp.Header.Get("x-request-id")
		return apiErr
	}

	if out == nil || len(raw.Data) == 0 || string(raw.Data) == "null" {
		return nil
	}
	return json.Unmarshal(raw.Data, out)
}

func remarshal(input any, output any) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}
