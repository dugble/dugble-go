package dugble

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TopicSubscription string

type TopicVisibility string

const (
	TopicSubscriptionOptIn  TopicSubscription = "opt_in"
	TopicSubscriptionOptOut TopicSubscription = "opt_out"

	TopicVisibilityPublic  TopicVisibility = "public"
	TopicVisibilityPrivate TopicVisibility = "private"
)

type Topic struct {
	Object              string            `json:"object"`
	ID                  string            `json:"id"`
	Name                string            `json:"name"`
	Description         *string           `json:"description"`
	DefaultSubscription TopicSubscription `json:"default_subscription"`
	Visibility          TopicVisibility   `json:"visibility"`
	CreatedAt           string            `json:"created_at"`
}

type TopicMutationResponse struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type TopicDeleteResponse struct {
	TopicMutationResponse
	Deleted bool `json:"deleted"`
}

type TopicList struct {
	Object  string  `json:"object"`
	HasMore bool    `json:"has_more"`
	Data    []Topic `json:"data"`
}

type CreateTopicParams struct {
	Name                string            `json:"name"`
	Description         *string           `json:"description,omitempty"`
	DefaultSubscription TopicSubscription `json:"default_subscription"`
	Visibility          TopicVisibility   `json:"visibility,omitempty"`
}

type UpdateTopicParams struct {
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Visibility  *TopicVisibility `json:"visibility,omitempty"`
}

type ListTopicsParams struct {
	Limit  int
	Offset int
}

type TopicsService struct{ client *Client }

func (s *TopicsService) Create(ctx context.Context, params CreateTopicParams) (*TopicMutationResponse, error) {
	var result TopicMutationResponse
	if err := s.client.request(ctx, http.MethodPost, "/topics", params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *TopicsService) List(ctx context.Context, params ListTopicsParams) (*TopicList, error) {
	query := url.Values{}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset != 0 {
		query.Set("offset", strconv.Itoa(params.Offset))
	}
	path := "/topics"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var result TopicList
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *TopicsService) Get(ctx context.Context, id string) (*Topic, error) {
	var result Topic
	if err := s.client.request(ctx, http.MethodGet, "/topics/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *TopicsService) Update(ctx context.Context, id string, params UpdateTopicParams) (*TopicMutationResponse, error) {
	var result TopicMutationResponse
	if err := s.client.request(ctx, http.MethodPatch, "/topics/"+url.PathEscape(id), params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *TopicsService) Delete(ctx context.Context, id string) (*TopicDeleteResponse, error) {
	var result TopicDeleteResponse
	if err := s.client.request(ctx, http.MethodDelete, "/topics/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	if result.Object == "" && result.ID == "" && !result.Deleted {
		return nil, fmt.Errorf("dugble: empty topic delete response")
	}
	return &result, nil
}
