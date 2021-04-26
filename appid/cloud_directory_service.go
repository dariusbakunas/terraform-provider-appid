package appid

import (
	"context"
	"fmt"
)

type CloudDirectoryService service

type EmailTeamplate struct {
	Subject  string `json:"subject"`
	HTMLBody string `json:"html_body"`
}

func (s *CloudDirectoryService) GetEmailTemplate(ctx context.Context, tenantID string, templateName string, language string) (*EmailTeamplate, error) {
	path := fmt.Sprintf("/management/v4/{%s}/config/cloud_directory/templates/%s/%s", tenantID, templateName, language)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp := &EmailTeamplate{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
