package dugble

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type SmsStatus string

const (
	SmsStatusQueued      SmsStatus = "queued"
	SmsStatusProcessing  SmsStatus = "processing"
	SmsStatusSubmitted   SmsStatus = "submitted"
	SmsStatusSent        SmsStatus = "sent"
	SmsStatusDelivered   SmsStatus = "delivered"
	SmsStatusUndelivered SmsStatus = "undelivered"
	SmsStatusRejected    SmsStatus = "rejected"
	SmsStatusFailed      SmsStatus = "failed"
	SmsStatusExpired     SmsStatus = "expired"
	SmsStatusUnknown     SmsStatus = "unknown"
	SmsStatusCanceled    SmsStatus = "canceled"
)

type SendSmsParams struct {
	To          string         `json:"to"`
	From        string         `json:"from"`
	Body        string         `json:"body"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	ScheduledAt string         `json:"scheduled_at,omitempty"`
}

type SendSmsResponse struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type SmsFailure struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SmsDestination struct {
	Country string `json:"country"`
}

type SmsMessage struct {
	Object      string         `json:"object"`
	ID          string         `json:"id"`
	MessageID   *string        `json:"message_id"`
	To          string         `json:"to"`
	From        string         `json:"from"`
	Body        string         `json:"body"`
	LastEvent   string         `json:"last_event"`
	Destination SmsDestination `json:"destination"`
	Segments    int            `json:"segments"`
	Metadata    map[string]any `json:"metadata"`
	ScheduledAt *string        `json:"scheduled_at"`
	Failure     *SmsFailure    `json:"failure,omitempty"`
	SubmittedAt *string        `json:"submitted_at,omitempty"`
	DeliveredAt *string        `json:"delivered_at,omitempty"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

type SmsEvent struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	OccurredAt string  `json:"occurred_at"`
	Provider   *string `json:"provider,omitempty"`
	Code       *string `json:"code,omitempty"`
	Message    *string `json:"message,omitempty"`
}

type SmsEventList struct {
	Object string     `json:"object"`
	Data   []SmsEvent `json:"data"`
}

type SmsAnalyticsRate struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type SmsAnalyticsPoint struct {
	Date      string `json:"date"`
	Total     int64  `json:"total"`
	Delivered int64  `json:"delivered"`
	Failed    int64  `json:"failed"`
}

type SmsAnalyticsWindow struct {
	Days   int                 `json:"days"`
	Rates  []SmsAnalyticsRate  `json:"rates"`
	Series []SmsAnalyticsPoint `json:"series"`
}

type SmsCountryAnalytics struct {
	Country   string `json:"country"`
	Total     int64  `json:"total"`
	Delivered int64  `json:"delivered"`
	Failed    int64  `json:"failed"`
}

type SmsAnalytics struct {
	Object            string                `json:"object"`
	Windows           []SmsAnalyticsWindow  `json:"windows"`
	DeliveryByCountry []SmsCountryAnalytics `json:"delivery_by_country"`
}

type ListSmsParams struct {
	Limit     int
	Offset    int
	Status    SmsStatus
	Sender    string
	StartDate string
	EndDate   string
	Search    string
}

type ListSmsEventsParams struct {
	Limit int
}

type UpdateSmsParams struct {
	ScheduledAt string `json:"scheduled_at"`
}

type SmsService struct{ client *Client }

func (s *SmsService) Send(ctx context.Context, params SendSmsParams, idempotencyKey string) (*SendSmsResponse, error) {
	headers := http.Header{}
	if idempotencyKey != "" {
		headers.Set("Idempotency-Key", idempotencyKey)
	}
	var result SendSmsResponse
	if err := s.client.request(ctx, http.MethodPost, "/sms", params, &result, headers); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SmsService) BatchSend(ctx context.Context, params []SendSmsParams, idempotencyKey string) ([]SendSmsResponse, error) {
	headers := http.Header{}
	if idempotencyKey != "" {
		headers.Set("Idempotency-Key", idempotencyKey)
	}
	var result []SendSmsResponse
	if err := s.client.request(ctx, http.MethodPost, "/sms/batch", params, &result, headers); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SmsService) Get(ctx context.Context, id string) (*SmsMessage, error) {
	var result SmsMessage
	if err := s.client.request(ctx, http.MethodGet, "/sms/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SmsService) Analytics(ctx context.Context) (*SmsAnalytics, error) {
	var result SmsAnalytics
	if err := s.client.request(ctx, http.MethodGet, "/sms/analytics", nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SmsService) List(ctx context.Context, params ListSmsParams) ([]SmsMessage, error) {
	query := url.Values{}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset != 0 {
		query.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Status != "" {
		query.Set("status", string(params.Status))
	}
	if params.Sender != "" {
		query.Set("sender", params.Sender)
	}
	if params.StartDate != "" {
		query.Set("start_date", params.StartDate)
	}
	if params.EndDate != "" {
		query.Set("end_date", params.EndDate)
	}
	if params.Search != "" {
		query.Set("search", params.Search)
	}
	path := "/sms"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var result []SmsMessage
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SmsService) Update(ctx context.Context, id string, params UpdateSmsParams) (*SendSmsResponse, error) {
	var result SendSmsResponse
	if err := s.client.request(ctx, http.MethodPatch, "/sms/"+url.PathEscape(id), params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SmsService) Cancel(ctx context.Context, id string) (*SendSmsResponse, error) {
	var result SendSmsResponse
	if err := s.client.request(ctx, http.MethodPost, "/sms/"+url.PathEscape(id)+"/cancel", nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SmsService) Events(ctx context.Context, id string, params ListSmsEventsParams) (*SmsEventList, error) {
	query := url.Values{}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	path := "/sms/" + url.PathEscape(id) + "/events"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var result SmsEventList
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SmsService) SyncStatus(ctx context.Context, id string) (*SmsMessage, error) {
	var result SmsMessage
	if err := s.client.request(ctx, http.MethodPost, "/sms/"+url.PathEscape(id)+"/sync-status", nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
