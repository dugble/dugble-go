package dugble

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type DomainStatus string
type DomainHealthStatus string
type DomainRegion string
type DomainTLSMode string
type DomainRecordStatus string

const (
	DomainStatusNotStarted        DomainStatus = "not_started"
	DomainStatusPending           DomainStatus = "pending"
	DomainStatusVerified          DomainStatus = "verified"
	DomainStatusPartiallyVerified DomainStatus = "partially_verified"
	DomainStatusPartiallyFailed   DomainStatus = "partially_failed"
	DomainStatusFailed            DomainStatus = "failed"
	DomainStatusTemporaryFailure  DomainStatus = "temporary_failure"
	DomainStatusDisabled          DomainStatus = "disabled"

	DomainHealthUnknown  DomainHealthStatus = "unknown"
	DomainHealthHealthy  DomainHealthStatus = "healthy"
	DomainHealthDegraded DomainHealthStatus = "degraded"

	DomainRegionUSEast1 DomainRegion = "us-east-1"
	DomainRegionEUNorth1 DomainRegion = "eu-north-1"

	DomainTLSOpportunistic DomainTLSMode = "opportunistic"
	DomainTLSEnforced      DomainTLSMode = "enforced"

	DomainRecordPending  DomainRecordStatus = "pending"
	DomainRecordVerified DomainRecordStatus = "verified"
	DomainRecordFailed   DomainRecordStatus = "failed"
)

type DomainVerificationRecord struct {
	Record   string             `json:"record"`
	Name     string             `json:"name"`
	Value    string             `json:"value"`
	Type     string             `json:"type"`
	Status   DomainRecordStatus `json:"status"`
	TTL      string             `json:"ttl"`
	Priority *int               `json:"priority,omitempty"`
}

type Domain struct {
	ID                        string                     `json:"id"`
	TeamID                    string                     `json:"team_id"`
	Name                      string                     `json:"name"`
	Provider                  *string                    `json:"provider,omitempty"`
	ProviderAccount           *string                    `json:"provider_account,omitempty"`
	Region                    DomainRegion               `json:"region"`
	ProviderExternalID        *string                    `json:"provider_external_id,omitempty"`
	Status                    DomainStatus               `json:"status"`
	ProviderStatus            *string                    `json:"provider_status,omitempty"`
	Records                   []DomainVerificationRecord `json:"records"`
	TLS                       DomainTLSMode              `json:"tls"`
	FailureReason             *string                    `json:"failure_reason,omitempty"`
	HealthStatus              DomainHealthStatus         `json:"health_status"`
	ConsecutiveHealthFailures int                        `json:"consecutive_health_failures"`
	LastCheckedAt             *string                    `json:"last_checked_at,omitempty"`
	LastHealthCheckedAt       *string                    `json:"last_health_checked_at,omitempty"`
	LastHealthFailureAt       *string                    `json:"last_health_failure_at,omitempty"`
	VerifiedAt                *string                    `json:"verified_at,omitempty"`
	DisabledAt                *string                    `json:"disabled_at,omitempty"`
	CreatedBy                 *string                    `json:"created_by,omitempty"`
	CreatedAt                 string                     `json:"created_at"`
	UpdatedAt                 string                     `json:"updated_at"`
}

type CreateDomainParams struct {
	Name   string        `json:"name"`
	Region DomainRegion  `json:"region"`
	TLS    *DomainTLSMode `json:"tls,omitempty"`
}

type UpdateDomainParams struct {
	TLS *DomainTLSMode `json:"tls,omitempty"`
}

type ListDomainsParams struct {
	Limit  int
	Offset int
}

type DomainProvisioningResponse struct {
	Status            string `json:"status"`
	Message           string `json:"message"`
	RetryAfterSeconds int    `json:"retry_after_seconds"`
}

// CreateDomainResponse covers both an immediately returned domain and the
// provider provisioning response. Check Provisioning != nil first.
type CreateDomainResponse struct {
	Domain       *Domain
	Provisioning *DomainProvisioningResponse
}

type DomainsService struct{ client *Client }

func (s *DomainsService) Create(ctx context.Context, params CreateDomainParams) (*CreateDomainResponse, error) {
	var raw map[string]any
	if err := s.client.request(ctx, http.MethodPost, "/domains", params, &raw, nil); err != nil {
		return nil, err
	}
	if status, _ := raw["status"].(string); status == "provisioning" {
		var result DomainProvisioningResponse
		if err := remarshal(raw, &result); err != nil {
			return nil, err
		}
		return &CreateDomainResponse{Provisioning: &result}, nil
	}
	var result Domain
	if err := remarshal(raw, &result); err != nil {
		return nil, err
	}
	return &CreateDomainResponse{Domain: &result}, nil
}

func (s *DomainsService) List(ctx context.Context, params ListDomainsParams) ([]Domain, error) {
	query := url.Values{}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset != 0 {
		query.Set("offset", strconv.Itoa(params.Offset))
	}
	path := "/domains"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var result []Domain
	if err := s.client.request(ctx, http.MethodGet, path, nil, &result, nil); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *DomainsService) Get(ctx context.Context, id string) (*Domain, error) {
	var result Domain
	if err := s.client.request(ctx, http.MethodGet, "/domains/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *DomainsService) Update(ctx context.Context, id string, params UpdateDomainParams) (*Domain, error) {
	var result Domain
	if err := s.client.request(ctx, http.MethodPatch, "/domains/"+url.PathEscape(id), params, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *DomainsService) Verify(ctx context.Context, id string) (*Domain, error) {
	var result Domain
	if err := s.client.request(ctx, http.MethodPost, "/domains/"+url.PathEscape(id)+"/verify", nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *DomainsService) Delete(ctx context.Context, id string) (*Domain, error) {
	var result Domain
	if err := s.client.request(ctx, http.MethodDelete, "/domains/"+url.PathEscape(id), nil, &result, nil); err != nil {
		return nil, err
	}
	return &result, nil
}
