package dugble

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type ContactPropertyType string

const (
	ContactPropertyString ContactPropertyType = "string"
	ContactPropertyNumber ContactPropertyType = "number"
)

type ContactProperty struct {
	Object        string              `json:"object"`
	ID            string              `json:"id"`
	Key           string              `json:"key"`
	Type          ContactPropertyType `json:"type"`
	FallbackValue any                 `json:"fallback_value,omitempty"`
	CreatedAt     string              `json:"created_at"`
}

type ContactPropertyList struct {
	Object  string            `json:"object"`
	HasMore bool              `json:"has_more"`
	Data    []ContactProperty `json:"data"`
}

type ContactPropertyMutationResponse struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type ContactPropertyDeleteResponse struct {
	ContactPropertyMutationResponse
	Deleted bool `json:"deleted"`
}

type CreateContactPropertyParams struct {
	Key           string              `json:"key"`
	Type          ContactPropertyType `json:"type"`
	FallbackValue any                 `json:"fallback_value,omitempty"`
}

type UpdateContactPropertyParams struct {
	FallbackValue any `json:"fallback_value"`
}

type ListContactPropertiesParams struct {
	Limit  int
	After  string
	Before string
}

type ContactPropertiesService struct{ client *Client }

func (s *ContactPropertiesService) Create(ctx context.Context, params CreateContactPropertyParams) (*ContactPropertyMutationResponse, error) {
	var result ContactPropertyMutationResponse
	if err := s.client.request(ctx, http.MethodPost, "/contact-properties", params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContactPropertiesService) List(ctx context.Context, params ListContactPropertiesParams) (*ContactPropertyList, error) {
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
	path := "/contact-properties"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var result ContactPropertyList
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContactPropertiesService) Get(ctx context.Context, id string) (*ContactProperty, error) {
	var result ContactProperty
	if err := s.client.request(ctx, http.MethodGet, "/contact-properties/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContactPropertiesService) Update(ctx context.Context, id string, params UpdateContactPropertyParams) (*ContactPropertyMutationResponse, error) {
	var result ContactPropertyMutationResponse
	if err := s.client.request(ctx, http.MethodPatch, "/contact-properties/"+url.PathEscape(id), params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContactPropertiesService) Delete(ctx context.Context, id string) (*ContactPropertyDeleteResponse, error) {
	var result ContactPropertyDeleteResponse
	if err := s.client.request(ctx, http.MethodDelete, "/contact-properties/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
