package appid

import (
	"context"
	"fmt"
)

type ApplicationService service

type Application struct {
	ClientID          string  `json:"clientId"`
	TenantID          string  `json:"tenantId"`
	Secret            *string `json:"secret,omitempty"`
	Name              string  `json:"name"`
	OAuthServerURL    string  `json:"oAuthServerUrl"`
	ProfilesURL       string  `json:"profilesURL"`
	DiscoveryEndpoint string  `json:"discoveryEndpoint"`
	Type              string  `json:"type"`
}

type CreateApplicationInput struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (s *ApplicationService) GetApplication(ctx context.Context, tenantID string, clientID string) (*Application, error) {
	path := fmt.Sprintf("/management/v4/%s/applications/%s", tenantID, clientID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp := &Application{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ApplicationService) GetApplicationScopes(ctx context.Context, tenantID string, clientID string) ([]string, error) {
	path := fmt.Sprintf("/management/v4/%s/applications/%s/scopes", tenantID, clientID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp := &struct {
		Scopes []string `json:"scopes"`
	}{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp.Scopes, nil
}

func (s *ApplicationService) CreateApplication(ctx context.Context, tenantID string, input *CreateApplicationInput) (*Application, error) {
	path := fmt.Sprintf("/management/v4/%s/applications", tenantID)

	req, err := s.client.NewRequest("POST", path, input)
	if err != nil {
		return nil, err
	}

	resp := &Application{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ApplicationService) UpdateApplication(ctx context.Context, tenantID string, clientID string, name string) (*Application, error) {
	path := fmt.Sprintf("/management/v4/%s/applications/%s", tenantID, clientID)

	input := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	req, err := s.client.NewRequest("PUT", path, input)
	if err != nil {
		return nil, err
	}

	resp := &Application{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ApplicationService) DeleteApplication(ctx context.Context, tenantID string, clientID string) error {
	path := fmt.Sprintf("/management/v4/%s/applications/%s", tenantID, clientID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)

	if err != nil {
		return err
	}

	return nil
}
