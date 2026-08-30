package dugble

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type TemplateStatus string
type TemplateCategory string
type TemplateVariableType string

const (
	TemplateStatusDraft     TemplateStatus = "draft"
	TemplateStatusPublished TemplateStatus = "published"

	TemplateCategoryOTP          TemplateCategory = "otp"
	TemplateCategoryWelcome      TemplateCategory = "welcome"
	TemplateCategoryReceipt      TemplateCategory = "receipt"
	TemplateCategoryAlert        TemplateCategory = "alert"
	TemplateCategoryNotification TemplateCategory = "notification"
	TemplateCategoryCustom       TemplateCategory = "custom"

	TemplateVariableString TemplateVariableType = "string"
	TemplateVariableNumber TemplateVariableType = "number"
)

type TemplateVariableParam struct {
	Key           string               `json:"key"`
	Type          TemplateVariableType `json:"type"`
	FallbackValue any                  `json:"fallback_value,omitempty"`
}

type TemplateVariable struct {
	ID            string               `json:"id"`
	Key           string               `json:"key"`
	Type          TemplateVariableType `json:"type"`
	FallbackValue any                  `json:"fallback_value"`
	CreatedAt     string               `json:"created_at"`
	UpdatedAt     string               `json:"updated_at"`
}

type Template struct {
	Object                 string             `json:"object"`
	ID                     string             `json:"id"`
	CurrentVersionID       string             `json:"current_version_id"`
	Alias                  *string            `json:"alias"`
	Name                   string             `json:"name"`
	Category               TemplateCategory   `json:"category"`
	CreatedAt              string             `json:"created_at"`
	UpdatedAt              string             `json:"updated_at"`
	Status                 TemplateStatus     `json:"status"`
	PublishedAt            *string            `json:"published_at"`
	From                   *string            `json:"from"`
	Subject                *string            `json:"subject"`
	ReplyTo                []string           `json:"reply_to"`
	HTML                   string             `json:"html"`
	Text                   *string            `json:"text"`
	Variables              []TemplateVariable `json:"variables"`
	HasUnpublishedVersions bool               `json:"has_unpublished_versions"`
}

type TemplateListItem struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Category    TemplateCategory `json:"category"`
	Status      TemplateStatus   `json:"status"`
	PublishedAt *string          `json:"published_at"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
	Alias       *string          `json:"alias"`
}

type TemplateList struct {
	Object  string             `json:"object"`
	Data    []TemplateListItem `json:"data"`
	HasMore bool               `json:"has_more"`
}

type TemplateMutationResponse struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type TemplateDeleteResponse struct {
	TemplateMutationResponse
	Deleted bool `json:"deleted"`
}

type TemplateVersionVariable struct {
	Key           string               `json:"key"`
	Type          TemplateVariableType `json:"type"`
	FallbackValue any                  `json:"fallback_value,omitempty"`
}

type TemplateVersion struct {
	ID               string                    `json:"id"`
	TeamID           string                    `json:"team_id"`
	TemplateID       string                    `json:"template_id"`
	VersionNumber    int32                     `json:"version_number"`
	FromEmail        *string                   `json:"from_email,omitempty"`
	FromName         *string                   `json:"from_name,omitempty"`
	ReplyToEmail     *string                   `json:"reply_to_email,omitempty"`
	Subject          string                    `json:"subject"`
	HTML             string                    `json:"html"`
	Text             *string                   `json:"text,omitempty"`
	Variables        []TemplateVersionVariable `json:"variables"`
	BasedOnVersionID *string                   `json:"based_on_version_id,omitempty"`
	ChangeNote       *string                   `json:"change_note,omitempty"`
	CreatedAt        string                    `json:"created_at"`
}

type TemplateRevertResponse struct {
	ID                    string           `json:"id"`
	TeamID                string           `json:"team_id"`
	Name                  string           `json:"name"`
	Alias                 *string          `json:"alias"`
	Category              TemplateCategory `json:"category"`
	CurrentVersionID      *string          `json:"current_version_id,omitempty"`
	PublishedVersionID    *string          `json:"published_version_id,omitempty"`
	PublishedAt           *string          `json:"published_at,omitempty"`
	HasUnpublishedChanges bool             `json:"has_unpublished_changes"`
	CreatedAt             string           `json:"created_at"`
	UpdatedAt             string           `json:"updated_at"`
}

type TemplatePreview struct {
	TemplateID string  `json:"template_id"`
	VersionID  string  `json:"version_id"`
	Subject    string  `json:"subject"`
	HTML       string  `json:"html"`
	Text       *string `json:"text,omitempty"`
	FromEmail  *string `json:"from_email,omitempty"`
	FromName   *string `json:"from_name,omitempty"`
	ReplyTo    *string `json:"reply_to,omitempty"`
}

type CreateTemplateParams struct {
	Name      string                  `json:"name"`
	HTML      string                  `json:"html"`
	Category  TemplateCategory        `json:"category"`
	Alias     string                  `json:"alias,omitempty"`
	From      string                  `json:"from,omitempty"`
	Subject   string                  `json:"subject,omitempty"`
	ReplyTo   any                     `json:"reply_to,omitempty"`
	Text      string                  `json:"text,omitempty"`
	Variables []TemplateVariableParam `json:"variables,omitempty"`
}

type UpdateTemplateParams struct {
	Name      *string                  `json:"name,omitempty"`
	HTML      *string                  `json:"html,omitempty"`
	Alias     *string                  `json:"alias,omitempty"`
	Category  *TemplateCategory        `json:"category,omitempty"`
	From      *string                  `json:"from,omitempty"`
	Subject   *string                  `json:"subject,omitempty"`
	ReplyTo   any                      `json:"reply_to,omitempty"`
	Text      *string                  `json:"text,omitempty"`
	Variables *[]TemplateVariableParam `json:"variables,omitempty"`
}

type ListTemplatesParams struct{ Limit, Offset int }
type ListTemplateVersionsParams struct{ Limit, Offset int }

type PreviewTemplateParams struct {
	VersionID string         `json:"version_id,omitempty"`
	Variables map[string]any `json:"variables,omitempty"`
}

type TestSendTemplateParams struct {
	To        string         `json:"to"`
	VersionID string         `json:"version_id,omitempty"`
	Variables map[string]any `json:"variables,omitempty"`
}

type TemplatesService struct {
	client   *Client
	Versions *TemplateVersionsService
}

type TemplateVersionsService struct{ client *Client }

func (s *TemplatesService) Create(ctx context.Context, params CreateTemplateParams) (*TemplateMutationResponse, error) {
	var result TemplateMutationResponse
	if err := s.client.request(ctx, http.MethodPost, "/templates", params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *TemplatesService) List(ctx context.Context, params ListTemplatesParams) (*TemplateList, error) {
	var result TemplateList
	if err := s.client.request(ctx, http.MethodGet, "/templates"+paginationQuery(params.Limit, params.Offset), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *TemplatesService) Get(ctx context.Context, template string) (*Template, error) {
	var result Template
	if err := s.client.request(ctx, http.MethodGet, "/templates/"+url.PathEscape(template), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *TemplatesService) Update(ctx context.Context, template string, params UpdateTemplateParams) (*TemplateMutationResponse, error) {
	var result TemplateMutationResponse
	if err := s.client.request(ctx, http.MethodPatch, "/templates/"+url.PathEscape(template), params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *TemplatesService) Delete(ctx context.Context, template string) (*TemplateDeleteResponse, error) {
	var result TemplateDeleteResponse
	if err := s.client.request(ctx, http.MethodDelete, "/templates/"+url.PathEscape(template), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *TemplatesService) Publish(ctx context.Context, template string) (*TemplateMutationResponse, error) {
	var result TemplateMutationResponse
	if err := s.client.request(ctx, http.MethodPost, "/templates/"+url.PathEscape(template)+"/publish", nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *TemplatesService) Duplicate(ctx context.Context, template string) (*TemplateMutationResponse, error) {
	var result TemplateMutationResponse
	if err := s.client.request(ctx, http.MethodPost, "/templates/"+url.PathEscape(template)+"/duplicate", nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *TemplatesService) Preview(ctx context.Context, template string, params PreviewTemplateParams) (*TemplatePreview, error) {
	var body any
	if params.VersionID != "" || params.Variables != nil {
		body = params
	}
	var result TemplatePreview
	if err := s.client.request(ctx, http.MethodPost, "/templates/"+url.PathEscape(template)+"/preview", body, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *TemplatesService) TestSend(ctx context.Context, template string, params TestSendTemplateParams) (*SendEmailResponse, error) {
	var result SendEmailResponse
	if err := s.client.request(ctx, http.MethodPost, "/templates/"+url.PathEscape(template)+"/test-send", params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *TemplateVersionsService) List(ctx context.Context, template string, params ListTemplateVersionsParams) ([]TemplateVersion, error) {
	var result []TemplateVersion
	path := "/templates/" + url.PathEscape(template) + "/versions" + paginationQuery(params.Limit, params.Offset)
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return result, nil
}
func (s *TemplateVersionsService) Get(ctx context.Context, template, versionID string) (*TemplateVersion, error) {
	var result TemplateVersion
	path := "/templates/" + url.PathEscape(template) + "/versions/" + url.PathEscape(versionID)
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *TemplateVersionsService) Revert(ctx context.Context, template, versionID string) (*TemplateRevertResponse, error) {
	var result TemplateRevertResponse
	path := "/templates/" + url.PathEscape(template) + "/versions/" + url.PathEscape(versionID) + "/revert"
	if err := s.client.request(ctx, http.MethodPost, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func paginationQuery(limit, offset int) string {
	query := url.Values{}
	if limit != 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if offset != 0 {
		query.Set("offset", strconv.Itoa(offset))
	}
	if encoded := query.Encode(); encoded != "" {
		return "?" + encoded
	}
	return ""
}
