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
