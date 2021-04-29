package appid

import (
	"context"
	"encoding/json"
	"fmt"
)

type IDPService service

type AuthNContext struct {
	Class      []string `json:"class,omitempty"`
	Comparison string   `json:"comparison,omitempty"`
}

type SAMLConfig struct {
	EntityID        string        `json:"entityID"`
	DisplayName     string        `json:"displayName,omitempty"`
	SignInURL       string        `json:"signInUrl"`
	Certificates    []string      `json:"certificates"`
	AuthNContext    *AuthNContext `json:"authnContext,omitempty"`
	SignRequest     *bool         `json:"signRequest,omitempty"`
	EncryptResponse *bool         `json:"encryptResponse,omitempty"`
	IncludeScoping  *bool         `json:"includeScoping,omitempty"`
}

func (s *SAMLConfig) String() string {
	str, _ := json.MarshalIndent(s, "", "\t")
	return string(str)
}

type SAMLIDP struct {
	IsActive bool        `json:"isActive"`
	Config   *SAMLConfig `json:"config,omitempty"`
}

type CloudDirectoryConfig struct {
}
type CloudDirectoryIDP struct {
	IsActive bool                  `json:"isActive"`
	Config   *CloudDirectoryConfig `json:"config,omitempty"`
}

func (s *IDPService) GetCloudDirectoryConfig(ctx context.Context, tenantID string) (*CloudDirectoryIDP, error) {
	path := fmt.Sprintf("/management/v4/%s/config/idps/cloud_directory", tenantID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp := &CloudDirectoryIDP{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *IDPService) GetSAMLConfig(ctx context.Context, tenantID string) (*SAMLIDP, error) {
	path := fmt.Sprintf("/management/v4/%s/config/idps/saml", tenantID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp := &SAMLIDP{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *IDPService) UpdateSAMLConfig(ctx context.Context, tenantID string, config *SAMLIDP) error {
	path := fmt.Sprintf("/management/v4/%s/config/idps/saml", tenantID)

	req, err := s.client.NewRequest("PUT", path, config)

	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, config)

	return err
}
