package dugble

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type SuppressionOrigin string

const (
	SuppressionOriginBounce    SuppressionOrigin = "bounce"
	SuppressionOriginComplaint SuppressionOrigin = "complaint"
	SuppressionOriginManual    SuppressionOrigin = "manual"
)

type Suppression struct {
	Object    string            `json:"object"`
	ID        string            `json:"id"`
	Email     string            `json:"email"`
	Origin    SuppressionOrigin `json:"origin"`
	SourceID  *string           `json:"source_id"`
	CreatedAt string            `json:"created_at"`
}

type SuppressionMutationResponse struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type SuppressionDeleteResponse struct {
	SuppressionMutationResponse
	Deleted bool `json:"deleted"`
}

type SuppressionList struct {
	Object  string        `json:"object"`
	HasMore bool          `json:"has_more"`
	Data    []Suppression `json:"data"`
}

type CreateSuppressionParams struct {
	Email string `json:"email"`
}

type ListSuppressionsParams struct {
	Limit  int
	After  string
	Before string
	Origin SuppressionOrigin
}

type BatchAddSuppressionsParams struct {
	Emails []string `json:"emails"`
}

type BatchRemoveSuppressionsParams struct {
	Emails []string `json:"emails,omitempty"`
	IDs    []string `json:"ids,omitempty"`
}

type BatchAddSuppressionsResponse struct {
	Data []SuppressionMutationResponse `json:"data"`
}

type BatchRemoveSuppressionsResponse struct {
	Data []SuppressionDeleteResponse `json:"data"`
}

type SuppressionsService struct{ client *Client }

func (s *SuppressionsService) Create(ctx context.Context, params CreateSuppressionParams) (*SuppressionMutationResponse, error) {
	var result SuppressionMutationResponse
	if err := s.client.request(ctx, http.MethodPost, "/suppressions", params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SuppressionsService) List(ctx context.Context, params ListSuppressionsParams) (*SuppressionList, error) {
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
	if params.Origin != "" {
		query.Set("origin", string(params.Origin))
	}
	path := "/suppressions"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var result SuppressionList
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SuppressionsService) Get(ctx context.Context, identifier string) (*Suppression, error) {
	var result Suppression
	if err := s.client.request(ctx, http.MethodGet, "/suppressions/"+url.PathEscape(identifier), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SuppressionsService) Delete(ctx context.Context, identifier string) (*SuppressionDeleteResponse, error) {
	var result SuppressionDeleteResponse
	if err := s.client.request(ctx, http.MethodDelete, "/suppressions/"+url.PathEscape(identifier), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SuppressionsService) BatchAdd(ctx context.Context, params BatchAddSuppressionsParams) (*BatchAddSuppressionsResponse, error) {
	var result BatchAddSuppressionsResponse
	if err := s.client.request(ctx, http.MethodPost, "/suppressions/batch/add", params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SuppressionsService) BatchRemove(ctx context.Context, params BatchRemoveSuppressionsParams) (*BatchRemoveSuppressionsResponse, error) {
	var result BatchRemoveSuppressionsResponse
	if err := s.client.request(ctx, http.MethodPost, "/suppressions/batch/remove", params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
