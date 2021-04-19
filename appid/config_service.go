package appid

import (
	"context"
	"fmt"
)

type ConfigService service

type IDTokenClaim struct {
	Source      string `json:"source"`
	SourceClaim string `json:"sourceClaim,omitempty"`
}

type TokensResponse struct {
	IDTokenClaims []IDTokenClaim `json:"idTokenClaims"`
}

func (s *ConfigService) GetTokens(ctx context.Context, tenantID string) (*TokensResponse, error) {
	path := fmt.Sprintf("/management/v4/%s/config/tokens", tenantID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp := &TokensResponse{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
