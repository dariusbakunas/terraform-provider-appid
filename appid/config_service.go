package appid

import (
	"context"
	"fmt"
)

type ConfigService service
type AccessTokenConfig struct {
	ExpiresIn int `json:"expires_in"`
}

type RefreshTokenConfig struct {
	Enabled   *bool `json:"enabled,omitempty"`
	ExpiresIn int   `json:"expires_in"`
}

type AnonymusAccessConfig struct {
	Enabled   *bool `json:"enabled,omitempty"`
	ExpiresIn int   `json:"expires_in"`
}
type TokenClaim struct {
	Source           string `json:"source"`
	SourceClaim      string `json:"sourceClaim,omitempty"`
	DestinationClaim string `json:"destinationClaim,omitempty"`
}

type TokenConfig struct {
	Access            *AccessTokenConfig    `json:"access,omitempty"`
	Refresh           *RefreshTokenConfig   `json:"refresh,omitempty"`
	AnonymousAccess   *AnonymusAccessConfig `json:"anonymousAccess,omitempty"`
	IDTokenClaims     []TokenClaim          `json:"idTokenClaims"`
	AccessTokenClaims []TokenClaim          `json:"accessTokenClaims"`
}

func (s *ConfigService) GetTokens(ctx context.Context, tenantID string) (*TokenConfig, error) {
	path := fmt.Sprintf("/management/v4/%s/config/tokens", tenantID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp := &TokenConfig{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ConfigService) UpdateTokens(ctx context.Context, tenantID string, config *TokenConfig) error {
	path := fmt.Sprintf("/management/v4/%s/config/tokens", tenantID)

	req, err := s.client.NewRequest("PUT", path, config)

	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, config)

	return err
}
