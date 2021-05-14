package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type ConfigService service
type AccessTokenConfig struct {
	ExpiresIn int `json:"expires_in,omitempty"`
}

type RefreshTokenConfig struct {
	Enabled   *bool `json:"enabled,omitempty"`
	ExpiresIn int   `json:"expires_in,omitempty"`
}

type AnonymusAccessConfig struct {
	Enabled   *bool `json:"enabled,omitempty"`
	ExpiresIn int   `json:"expires_in,omitempty"`
}
type TokenClaim struct {
	Source           string  `json:"source"`
	SourceClaim      *string `json:"sourceClaim,omitempty"`
	DestinationClaim *string `json:"destinationClaim,omitempty"`
}

type TokenConfig struct {
	Access            *AccessTokenConfig    `json:"access,omitempty"`
	Refresh           *RefreshTokenConfig   `json:"refresh,omitempty"`
	AnonymousAccess   *AnonymusAccessConfig `json:"anonymousAccess,omitempty"`
	IDTokenClaims     []TokenClaim          `json:"idTokenClaims,omitempty"`
	AccessTokenClaims []TokenClaim          `json:"accessTokenClaims,omitempty"`
}

func (c *TokenConfig) String() string {
	str, _ := json.Marshal(c)
	return string(str)
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

	_, err = s.client.Do(ctx, req, nil)

	return err
}

func (s *ConfigService) ListRedirectURLs(ctx context.Context, tenantID string) ([]string, error) {
	path := fmt.Sprintf("/management/v4/%s/config/redirect_uris", tenantID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp := &struct {
		RedirectURIs []string `json:"redirectUris"`
	}{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp.RedirectURIs, nil
}

func (s *ConfigService) UpdateRedirectURLs(ctx context.Context, tenantID string, urls []string) error {
	path := fmt.Sprintf("/management/v4/%s/config/redirect_uris", tenantID)

	input := struct {
		RedirectURIs []string `json:"redirectUris"`
	}{
		RedirectURIs: urls,
	}

	req, err := s.client.NewRequest("PUT", path, input)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}
