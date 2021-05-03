package api

import (
	"context"
	"fmt"
)

type RolesService service

type RoleAccess struct {
	ApplicationID string   `json:"application_id"`
	Scopes        []string `json:"scopes"`
}

type RoleInput struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Access      []RoleAccess `json:"access"`
}
type Role struct {
	RoleInput
	ID string `json:"id"`
}

func (s *RolesService) GetRole(ctx context.Context, tenantID string, roleID string) (*Role, error) {
	path := fmt.Sprintf("/management/v4/%s/roles/%s", tenantID, roleID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp := &Role{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *RolesService) CreateRole(ctx context.Context, tenantID string, input *RoleInput) (*Role, error) {
	path := fmt.Sprintf("/management/v4/%s/roles", tenantID)

	req, err := s.client.NewRequest("POST", path, input)
	if err != nil {
		return nil, err
	}

	resp := &Role{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *RolesService) DeleteRole(ctx context.Context, tenantID string, roleID string) error {
	path := fmt.Sprintf("/management/v4/%s/roles/%s", tenantID, roleID)

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
