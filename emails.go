package dugble

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type EmailAttachment struct {
	Content     string `json:"content"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type,omitempty"`
	ContentID   string `json:"content_id,omitempty"`
}

type EmailTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SendEmailParams struct {
	From        any               `json:"from,omitempty"`
	To          any               `json:"to"`
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	ReplyTo     any               `json:"reply_to,omitempty"`
	CC          any               `json:"cc,omitempty"`
	BCC         any               `json:"bcc,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
	Tags        []EmailTag        `json:"tags,omitempty"`
	ScheduledAt string            `json:"scheduled_at,omitempty"`
	Metadata    map[string]any    `json:"metadata,omitempty"`
}

type SendEmailResponse struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type Email struct {
	Object      string     `json:"object"`
	ID          string     `json:"id"`
	MessageID   *string    `json:"message_id"`
	To          []string   `json:"to"`
	From        string     `json:"from"`
	CreatedAt   string     `json:"created_at"`
	Subject     string     `json:"subject"`
	HTML        *string    `json:"html"`
	Text        *string    `json:"text"`
	BCC         []string   `json:"bcc"`
	CC          []string   `json:"cc"`
	ReplyTo     []string   `json:"reply_to"`
	LastEvent   string     `json:"last_event"`
	ScheduledAt *string    `json:"scheduled_at"`
	Tags        []EmailTag `json:"tags"`
}

type EmailSummary struct {
	ID          string  `json:"id"`
	ToEmail     string  `json:"to_email"`
	ToName      *string `json:"to_name"`
	Subject     string  `json:"subject"`
	Status      string  `json:"status"`
	Provider    *string `json:"provider"`
	QueuedAt    string  `json:"queued_at"`
	SubmittedAt *string `json:"submitted_at"`
	DeliveredAt *string `json:"delivered_at"`
	CreatedAt   string  `json:"created_at"`
}

type EmailEvent struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	OccurredAt string  `json:"occurred_at"`
	Provider   *string `json:"provider,omitempty"`
	Code       *string `json:"code,omitempty"`
	Message    *string `json:"message,omitempty"`
}

type EmailEventList struct {
	Object string       `json:"object"`
	Data   []EmailEvent `json:"data"`
}

type EmailAnalyticsRate struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type EmailAnalyticsPoint struct {
	Date      string `json:"date"`
	Total     int64  `json:"total"`
	Delivered int64  `json:"delivered"`
	Opened    int64  `json:"opened"`
	Clicked   int64  `json:"clicked"`
	Bounced   int64  `json:"bounced"`
}

type EmailAnalyticsWindow struct {
	Days   int                   `json:"days"`
	Rates  []EmailAnalyticsRate  `json:"rates"`
	Series []EmailAnalyticsPoint `json:"series"`
}

type EmailAnalytics struct {
	Object  string                 `json:"object"`
	Windows []EmailAnalyticsWindow `json:"windows"`
}

type ListEmailsParams struct {
	Limit  int
	Offset int
}

type ListEmailEventsParams struct {
	Limit  int
	Offset int
}

type UpdateEmailParams struct {
	ScheduledAt string `json:"scheduled_at"`
}

type EmailsService struct{ client *Client }

func (s *EmailsService) Send(ctx context.Context, params SendEmailParams, idempotencyKey string) (*SendEmailResponse, error) {
	headers := http.Header{}
	if idempotencyKey != "" {
		headers.Set("Idempotency-Key", idempotencyKey)
	}
	var result SendEmailResponse
	if err := s.client.request(ctx, http.MethodPost, "/emails", params, &result, headers); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *EmailsService) BatchSend(ctx context.Context, params []SendEmailParams, idempotencyKey string) ([]SendEmailResponse, error) {
	headers := http.Header{}
	if idempotencyKey != "" {
		headers.Set("Idempotency-Key", idempotencyKey)
	}
	var result []SendEmailResponse
	if err := s.client.request(ctx, http.MethodPost, "/emails/batch", params, &result, headers); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *EmailsService) Get(ctx context.Context, id string) (*Email, error) {
	var result Email
	if err := s.client.request(ctx, http.MethodGet, "/emails/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *EmailsService) Analytics(ctx context.Context) (*EmailAnalytics, error) {
	var result EmailAnalytics
	if err := s.client.request(ctx, http.MethodGet, "/emails/analytics", nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *EmailsService) Events(ctx context.Context, id string, params ListEmailEventsParams) (*EmailEventList, error) {
	query := url.Values{}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset != 0 {
		query.Set("offset", strconv.Itoa(params.Offset))
	}
	path := "/emails/" + url.PathEscape(id) + "/events"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var result EmailEventList
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *EmailsService) List(ctx context.Context, params ListEmailsParams) ([]EmailSummary, error) {
	query := url.Values{}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset != 0 {
		query.Set("offset", strconv.Itoa(params.Offset))
	}
	path := "/emails"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var result []EmailSummary
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *EmailsService) Update(ctx context.Context, id string, params UpdateEmailParams) (*SendEmailResponse, error) {
	var result SendEmailResponse
	if err := s.client.request(ctx, http.MethodPatch, "/emails/"+url.PathEscape(id), params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *EmailsService) Cancel(ctx context.Context, id string) (*SendEmailResponse, error) {
	var result SendEmailResponse
	if err := s.client.request(ctx, http.MethodPost, "/emails/"+url.PathEscape(id)+"/cancel", nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
