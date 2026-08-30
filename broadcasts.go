package dugble

import (
	"context"
	"net/http"
	"net/url"
)

type BroadcastStatus string

type BroadcastRecipientStatus string

type BroadcastExclusionReason string

const (
	BroadcastStatusDraft     BroadcastStatus = "draft"
	BroadcastStatusScheduled BroadcastStatus = "scheduled"
	BroadcastStatusQueued    BroadcastStatus = "queued"
	BroadcastStatusSent      BroadcastStatus = "sent"
	BroadcastStatusFailed    BroadcastStatus = "failed"
	BroadcastStatusCanceled  BroadcastStatus = "canceled"
)

type Broadcast struct {
	ID               string          `json:"id"`
	TeamID           string          `json:"team_id"`
	Name             string          `json:"name"`
	Status           BroadcastStatus `json:"status"`
	SegmentID        string          `json:"segment_id"`
	TopicID          *string         `json:"topic_id,omitempty"`
	FromEmail        string          `json:"from_email"`
	FromName         *string         `json:"from_name,omitempty"`
	ReplyToEmail     *string         `json:"reply_to_email,omitempty"`
	Subject          string          `json:"subject"`
	PreviewText      *string         `json:"preview_text,omitempty"`
	HTML             string          `json:"html"`
	Text             *string         `json:"text,omitempty"`
	VariableBindings map[string]any  `json:"variable_bindings"`
	ScheduledAt      *string         `json:"scheduled_at,omitempty"`
	QueuedAt         *string         `json:"queued_at,omitempty"`
	SentAt           *string         `json:"sent_at,omitempty"`
	CanceledAt       *string         `json:"canceled_at,omitempty"`
	AudienceCount    int64           `json:"audience_count"`
	EligibleCount    int64           `json:"eligible_count"`
	SuppressedCount  int64           `json:"suppressed_count"`
	QueuedCount      int64           `json:"queued_count"`
	FailedCount      int64           `json:"failed_count"`
	Revision         int64           `json:"revision"`
	CreatedAt        string          `json:"created_at"`
	UpdatedAt        string          `json:"updated_at"`
}

type CreateBroadcastParams struct {
	Name             string         `json:"name,omitempty"`
	SegmentID        string         `json:"segment_id"`
	TopicID          string         `json:"topic_id,omitempty"`
	FromEmail        string         `json:"from_email,omitempty"`
	FromName         string         `json:"from_name,omitempty"`
	ReplyToEmail     string         `json:"reply_to_email,omitempty"`
	Subject          string         `json:"subject"`
	PreviewText      string         `json:"preview_text,omitempty"`
	HTML             string         `json:"html"`
	Text             string         `json:"text,omitempty"`
	VariableBindings map[string]any `json:"variable_bindings,omitempty"`
	Send             bool           `json:"send,omitempty"`
	ScheduledAt      string         `json:"scheduled_at,omitempty"`
}

type UpdateBroadcastParams struct {
	Revision         int64           `json:"revision"`
	Name             *string         `json:"name,omitempty"`
	SegmentID        *string         `json:"segment_id,omitempty"`
	TopicID          **string        `json:"topic_id,omitempty"`
	FromEmail        **string        `json:"from_email,omitempty"`
	FromName         **string        `json:"from_name,omitempty"`
	ReplyToEmail     **string        `json:"reply_to_email,omitempty"`
	Subject          *string         `json:"subject,omitempty"`
	PreviewText      **string        `json:"preview_text,omitempty"`
	HTML             *string         `json:"html,omitempty"`
	Text             **string        `json:"text,omitempty"`
	VariableBindings *map[string]any `json:"variable_bindings,omitempty"`
}

type ListBroadcastsParams struct{ Limit, Offset int }
type SendBroadcastParams struct {
	ScheduledAt string `json:"scheduled_at,omitempty"`
}
type DuplicateBroadcastParams struct {
	Name string `json:"name,omitempty"`
}
type PreviewBroadcastParams struct {
	Variables map[string]any `json:"variables,omitempty"`
}
type ListBroadcastRecipientsParams struct{ Limit, Offset int }

type BroadcastRecipient struct {
	ID              string                    `json:"id"`
	BroadcastID     string                    `json:"broadcast_id"`
	ContactID       *string                   `json:"contact_id,omitempty"`
	Email           string                    `json:"email"`
	FirstName       *string                   `json:"first_name,omitempty"`
	LastName        *string                   `json:"last_name,omitempty"`
	ContactSnapshot map[string]any            `json:"contact_snapshot"`
	Status          BroadcastRecipientStatus  `json:"status"`
	ExclusionReason *BroadcastExclusionReason `json:"exclusion_reason,omitempty"`
	EmailMessageID  *string                   `json:"email_message_id,omitempty"`
	CreatedAt       string                    `json:"created_at"`
	QueuedAt        *string                   `json:"queued_at,omitempty"`
}

type BroadcastExclusionSummary struct {
	Object      string           `json:"object"`
	BroadcastID string           `json:"broadcast_id"`
	Total       int64            `json:"total"`
	Reasons     map[string]int64 `json:"reasons"`
}

type BroadcastAnalytics struct {
	Object      string `json:"object"`
	BroadcastID string `json:"broadcast_id"`
	Audience    int64  `json:"audience"`
	Eligible    int64  `json:"eligible"`
	Excluded    int64  `json:"excluded"`
	Queued      int64  `json:"queued"`
	Delivered   int64  `json:"delivered"`
	Bounced     int64  `json:"bounced"`
	Complained  int64  `json:"complained"`
	Failed      int64  `json:"failed"`
	Opened      int64  `json:"opened"`
	Clicked     int64  `json:"clicked"`
}

type BroadcastPreview struct {
	FromEmail    string  `json:"from_email"`
	FromName     *string `json:"from_name,omitempty"`
	ReplyToEmail *string `json:"reply_to_email,omitempty"`
	Subject      string  `json:"subject"`
	PreviewText  *string `json:"preview_text,omitempty"`
	HTML         string  `json:"html"`
	Text         *string `json:"text,omitempty"`
}

type BroadcastsService struct{ client *Client }

func (s *BroadcastsService) Create(ctx context.Context, params CreateBroadcastParams) (*Broadcast, error) {
	var result Broadcast
	if err := s.client.request(ctx, http.MethodPost, "/broadcasts", params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *BroadcastsService) List(ctx context.Context, params ListBroadcastsParams) ([]Broadcast, error) {
	var result []Broadcast
	if err := s.client.request(ctx, http.MethodGet, "/broadcasts"+paginationQuery(params.Limit, params.Offset), nil, &result, nil); err != nil {
		return nil, err
	}
	return result, nil
}
func (s *BroadcastsService) Get(ctx context.Context, id string) (*Broadcast, error) {
	var result Broadcast
	if err := s.client.request(ctx, http.MethodGet, "/broadcasts/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *BroadcastsService) Update(ctx context.Context, id string, params UpdateBroadcastParams) (*Broadcast, error) {
	var result Broadcast
	if err := s.client.request(ctx, http.MethodPatch, "/broadcasts/"+url.PathEscape(id), params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *BroadcastsService) Delete(ctx context.Context, id string) (*Broadcast, error) {
	var result Broadcast
	if err := s.client.request(ctx, http.MethodDelete, "/broadcasts/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *BroadcastsService) Send(ctx context.Context, id string, params SendBroadcastParams) (*Broadcast, error) {
	var body any
	if params.ScheduledAt != "" {
		body = params
	}
	var result Broadcast
	if err := s.client.request(ctx, http.MethodPost, "/broadcasts/"+url.PathEscape(id)+"/send", body, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *BroadcastsService) Cancel(ctx context.Context, id string) (*Broadcast, error) {
	var result Broadcast
	if err := s.client.request(ctx, http.MethodPost, "/broadcasts/"+url.PathEscape(id)+"/cancel", nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *BroadcastsService) Duplicate(ctx context.Context, id string, params ...DuplicateBroadcastParams) (*Broadcast, error) {
	body := DuplicateBroadcastParams{}
	if len(params) > 0 {
		body = params[0]
	}
	var result Broadcast
	if err := s.client.request(ctx, http.MethodPost, "/broadcasts/"+url.PathEscape(id)+"/duplicate", body, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *BroadcastsService) Preview(ctx context.Context, id string, params PreviewBroadcastParams) (*BroadcastPreview, error) {
	var body any
	if params.Variables != nil {
		body = params
	}
	var result BroadcastPreview
	if err := s.client.request(ctx, http.MethodPost, "/broadcasts/"+url.PathEscape(id)+"/preview", body, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *BroadcastsService) Recipients(ctx context.Context, id string, params ListBroadcastRecipientsParams) ([]BroadcastRecipient, error) {
	var result []BroadcastRecipient
	path := "/broadcasts/" + url.PathEscape(id) + "/recipients" + paginationQuery(params.Limit, params.Offset)
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return result, nil
}
func (s *BroadcastsService) Exclusions(ctx context.Context, id string) (*BroadcastExclusionSummary, error) {
	var result BroadcastExclusionSummary
	if err := s.client.request(ctx, http.MethodGet, "/broadcasts/"+url.PathEscape(id)+"/exclusions", nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *BroadcastsService) Analytics(ctx context.Context, id string) (*BroadcastAnalytics, error) {
	var result BroadcastAnalytics
	if err := s.client.request(ctx, http.MethodGet, "/broadcasts/"+url.PathEscape(id)+"/analytics", nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
