package dugble

import (
	"context"
	"net/http"
	"net/url"
)

type SenderIDStatus string

const (
	SenderIDStatusPending   SenderIDStatus = "pending"
	SenderIDStatusApproved  SenderIDStatus = "approved"
	SenderIDStatusRejected  SenderIDStatus = "rejected"
	SenderIDStatusSuspended SenderIDStatus = "suspended"
	SenderIDStatusInactive  SenderIDStatus = "inactive"
)

type SenderID struct {
	ID              string         `json:"id"`
	TeamID          string         `json:"team_id"`
	Name            string         `json:"name"`
	CountryCode     string         `json:"country_code"`
	Purpose         string         `json:"purpose"`
	Status          SenderIDStatus `json:"status"`
	RejectionReason *string        `json:"rejection_reason,omitempty"`
	ApprovedAt      *string        `json:"approved_at,omitempty"`
	RejectedAt      *string        `json:"rejected_at,omitempty"`
	SuspendedAt     *string        `json:"suspended_at,omitempty"`
	CreatedBy       *string        `json:"created_by,omitempty"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
}

type CreateSenderIDParams struct {
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	Purpose     string `json:"purpose"`
	Provider    string `json:"provider,omitempty"`
}

type SenderIDsService struct{ client *Client }

func (s *SenderIDsService) Create(ctx context.Context, params CreateSenderIDParams) (*SenderID, error) {
	var result SenderID
	if err := s.client.request(ctx, http.MethodPost, "/sender-ids", params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SenderIDsService) List(ctx context.Context) ([]SenderID, error) {
	var result []SenderID
	if err := s.client.request(ctx, http.MethodGet, "/sender-ids", nil, &result, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SenderIDsService) Get(ctx context.Context, id string) (*SenderID, error) {
	var result SenderID
	if err := s.client.request(ctx, http.MethodGet, "/sender-ids/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SenderIDsService) Delete(ctx context.Context, id string) (*SenderID, error) {
	var result SenderID
	if err := s.client.request(ctx, http.MethodDelete, "/sender-ids/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
