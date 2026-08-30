package dugble

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type Segment struct {
	ID        string `json:"id"`
	TeamID    string `json:"team_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type SegmentAudienceSize struct {
	SegmentID string `json:"segment_id"`
	Count     int64  `json:"count"`
}

type SegmentContact struct {
	ID           string  `json:"id"`
	TeamID       string  `json:"team_id"`
	Email        string  `json:"email"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Unsubscribed bool    `json:"unsubscribed"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type CreateSegmentParams struct {
	Name string `json:"name"`
}

type ListSegmentsParams struct {
	Limit  int
	Offset int
}

type ListSegmentContactsParams struct {
	Limit  int
	Offset int
}

type SegmentsService struct{ client *Client }

func (s *SegmentsService) Create(ctx context.Context, params CreateSegmentParams) (*Segment, error) {
	var result Segment
	if err := s.client.request(ctx, http.MethodPost, "/segments", params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SegmentsService) List(ctx context.Context, params ListSegmentsParams) ([]Segment, error) {
	path := "/segments" + paginationSuffix(params.Limit, params.Offset)
	var result []Segment
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SegmentsService) Get(ctx context.Context, id string) (*Segment, error) {
	var result Segment
	if err := s.client.request(ctx, http.MethodGet, "/segments/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SegmentsService) AudienceSize(ctx context.Context, id string) (*SegmentAudienceSize, error) {
	var result SegmentAudienceSize
	path := "/segments/" + url.PathEscape(id) + "/audience-size"
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SegmentsService) Contacts(ctx context.Context, id string, params ListSegmentContactsParams) ([]SegmentContact, error) {
	path := "/segments/" + url.PathEscape(id) + "/contacts" + paginationSuffix(params.Limit, params.Offset)
	var result []SegmentContact
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SegmentsService) Delete(ctx context.Context, id string) (*Segment, error) {
	var result Segment
	if err := s.client.request(ctx, http.MethodDelete, "/segments/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func paginationSuffix(limit, offset int) string {
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
