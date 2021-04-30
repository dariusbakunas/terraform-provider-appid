package api

import (
	"context"
	"fmt"
)

type CloudDirectoryService service

type EmailTemplate struct {
	Subject     string `json:"subject"`
	HTMLBody    string `json:"html_body,omitempty"`
	B64HTMLBody string `json:"base64_encoded_html_body,omitempty"`
	TextBody    string `json:"plain_text_body,omitempty"`
}

func (s *CloudDirectoryService) GetEmailTemplate(ctx context.Context, tenantID string, templateName string, language string) (*EmailTemplate, error) {
	path := fmt.Sprintf("/management/v4/%s/config/cloud_directory/templates/%s/%s", tenantID, templateName, language)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp := &EmailTemplate{}

	_, err = s.client.Do(ctx, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *CloudDirectoryService) UpdateEmailTemplate(ctx context.Context, tenantID string, templateName string, language string, template *EmailTemplate) error {
	path := fmt.Sprintf("/management/v4/%s/config/cloud_directory/templates/%s/%s", tenantID, templateName, language)

	req, err := s.client.NewRequest("PUT", path, template)

	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)

	return err
}

func (s *CloudDirectoryService) DeleteEmailTemplate(ctx context.Context, tenantID string, templateName string, language string) error {
	path := fmt.Sprintf("/management/v4/%s/config/cloud_directory/templates/%s/%s", tenantID, templateName, language)

	req, err := s.client.NewRequest("DELETE", path, nil)

	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)

	return err
}
