package appid

import (
	"context"
	"fmt"
)

type IDPService service

type AuthNContext struct {
	Class      []string `json:"class,omitempty"`
	Comparison string   `json:"comparison,omitempty"`
}

type SAMLConfig struct {
	EntityID        string        `json:"entityID"`
	DisplayName     string        `json:"displayName"`
	SignInURL       string        `json:"signInUrl"`
	Certificates    []string      `json:"certificates"`
	AuthNContext    *AuthNContext `json:"authnContext"`
	SignRequest     *bool         `json:"signRequest"`
	EncryptResponse *bool         `json:"encryptResponse"`
	IncludeScoping  *bool         `json:"includeScoping"`
}

type SAMLResponse struct {
	IsActive bool        `json:"isActive"`
	Config   *SAMLConfig `json:"config,omitempty"`
}

func (s *IDPService) GetSAMLConfig(ctx context.Context, tenantID string) (*SAMLResponse, error) {
	path := fmt.Sprintf("/management/v4/%s/config/idps/saml", tenantID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp := &SAMLResponse{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
