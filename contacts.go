package dugble

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type SMSConsentStatus string
type SMSConsentSource string
type ContactTopicSubscription string

const (
	SMSConsentUnknown  SMSConsentStatus = "unknown"
	SMSConsentOptedIn  SMSConsentStatus = "opted_in"
	SMSConsentOptedOut SMSConsentStatus = "opted_out"

	SMSConsentSourceAPI    SMSConsentSource = "api"
	SMSConsentSourceImport SMSConsentSource = "import"
	SMSConsentSourceManual SMSConsentSource = "manual"

	ContactTopicOptIn  ContactTopicSubscription = "opt_in"
	ContactTopicOptOut ContactTopicSubscription = "opt_out"
)

type Contact struct {
	ID                  string            `json:"id"`
	TeamID              string            `json:"team_id"`
	Email               string            `json:"email"`
	Phone               *string           `json:"phone,omitempty"`
	NormalizedPhone     *string           `json:"normalized_phone,omitempty"`
	PhoneCountry        *string           `json:"phone_country,omitempty"`
	SMSConsentStatus    SMSConsentStatus  `json:"sms_consent_status"`
	SMSConsentUpdatedAt *string           `json:"sms_consent_updated_at,omitempty"`
	SMSConsentSource    *SMSConsentSource `json:"sms_consent_source,omitempty"`
	FirstName           *string           `json:"first_name,omitempty"`
	LastName            *string           `json:"last_name,omitempty"`
	Unsubscribed        bool              `json:"unsubscribed"`
	Properties          map[string]any    `json:"properties"`
	CreatedAt           string            `json:"created_at"`
	UpdatedAt           string            `json:"updated_at"`
}

type CreateContactParams struct {
	Email            string            `json:"email"`
	Phone            string            `json:"phone,omitempty"`
	SMSConsentStatus SMSConsentStatus  `json:"sms_consent_status,omitempty"`
	SMSConsentSource SMSConsentSource  `json:"sms_consent_source,omitempty"`
	FirstName        string            `json:"first_name,omitempty"`
	LastName         string            `json:"last_name,omitempty"`
	Unsubscribed     *bool             `json:"unsubscribed,omitempty"`
	Properties       map[string]any    `json:"properties,omitempty"`
}

type UpdateContactParams struct {
	Email            *string            `json:"email,omitempty"`
	Phone            *string            `json:"phone,omitempty"`
	SMSConsentStatus *SMSConsentStatus  `json:"sms_consent_status,omitempty"`
	SMSConsentSource *SMSConsentSource  `json:"sms_consent_source,omitempty"`
	FirstName        *string            `json:"first_name,omitempty"`
	LastName         *string            `json:"last_name,omitempty"`
	Unsubscribed     *bool              `json:"unsubscribed,omitempty"`
	Properties       *map[string]any    `json:"properties,omitempty"`
}

type ListContactsParams struct {
	Limit  int
	Offset int
}

type ContactTopic struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
	Description  *string                  `json:"description"`
	Subscription ContactTopicSubscription `json:"subscription"`
}

type ContactTopicList struct {
	Object  string         `json:"object"`
	HasMore bool           `json:"has_more"`
	Data    []ContactTopic `json:"data"`
}

type ListContactTopicsParams struct {
	Limit  int
	After  string
	Before string
}

type UpdateContactTopic struct {
	ID           string                   `json:"id"`
	Subscription ContactTopicSubscription `json:"subscription"`
}

type UpdateContactTopicsResponse struct {
	ID string `json:"id"`
}

type ContactSegmentMembership struct {
	ID         string `json:"id"`
	TeamID     string `json:"team_id"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
	AssignedAt string `json:"assigned_at"`
}

type ContactsService struct {
	client   *Client
	Topics   *ContactTopicsService
	Segments *ContactSegmentsService
}

type ContactTopicsService struct{ client *Client }
type ContactSegmentsService struct{ client *Client }

func (s *ContactsService) Create(ctx context.Context, params CreateContactParams) (*Contact, error) {
	var result Contact
	if err := s.client.request(ctx, http.MethodPost, "/contacts", params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContactsService) List(ctx context.Context, params ListContactsParams) ([]Contact, error) {
	query := url.Values{}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset != 0 {
		query.Set("offset", strconv.Itoa(params.Offset))
	}
	path := "/contacts"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var result []Contact
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ContactsService) Get(ctx context.Context, id string) (*Contact, error) {
	var result Contact
	if err := s.client.request(ctx, http.MethodGet, "/contacts/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContactsService) Update(ctx context.Context, id string, params UpdateContactParams) (*Contact, error) {
	var result Contact
	if err := s.client.request(ctx, http.MethodPatch, "/contacts/"+url.PathEscape(id), params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContactsService) Delete(ctx context.Context, id string) (*Contact, error) {
	var result Contact
	if err := s.client.request(ctx, http.MethodDelete, "/contacts/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContactTopicsService) List(ctx context.Context, contactID string, params ListContactTopicsParams) (*ContactTopicList, error) {
	query := url.Values{}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.After != "" {
		query.Set("after", params.After)
	}
	if params.Before != "" {
		query.Set("before", params.Before)
	}
	path := "/contacts/" + url.PathEscape(contactID) + "/topics"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var result ContactTopicList
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContactTopicsService) Update(ctx context.Context, contactID string, topics []UpdateContactTopic) (*UpdateContactTopicsResponse, error) {
	var result UpdateContactTopicsResponse
	if err := s.client.request(ctx, http.MethodPatch, "/contacts/"+url.PathEscape(contactID)+"/topics", topics, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContactSegmentsService) List(ctx context.Context, contactID string) ([]ContactSegmentMembership, error) {
	var result []ContactSegmentMembership
	if err := s.client.request(ctx, http.MethodGet, "/contacts/"+url.PathEscape(contactID)+"/segments", nil, &result, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ContactSegmentsService) Add(ctx context.Context, contactID, segmentID string) (*ContactSegmentMembership, error) {
	var result ContactSegmentMembership
	path := "/contacts/" + url.PathEscape(contactID) + "/segments/" + url.PathEscape(segmentID)
	if err := s.client.request(ctx, http.MethodPost, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContactSegmentsService) Remove(ctx context.Context, contactID, segmentID string) error {
	path := "/contacts/" + url.PathEscape(contactID) + "/segments/" + url.PathEscape(segmentID)
	return s.client.request(ctx, http.MethodDelete, path, nil, nil, nil)
}
