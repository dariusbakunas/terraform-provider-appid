package appid

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var supportedTemplates = []string{"USER_VERIFICATION", "RESET_PASSWORD", "WELCOME", "PASSWORD_CHANGED", "MFA_VERIFICATION"}

func dataSourceAppIDCloudDirectoryTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDCloudDirectoryTemplateRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(supportedTemplates, false),
			},
			"language": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "en",
			},
			"subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"html_body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"base64_encoded_html_body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plain_text_body": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppIDCloudDirectoryTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	templateName := d.Get("template_name").(string)
	language := d.Get("language").(string)

	c := m.(*Client)

	template, err := c.CloudDirectoryAPI.GetEmailTemplate(ctx, tenantID, templateName, language)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("subject", template.Subject)
	d.Set("html_body", template.HTMLBody)
	d.Set("base64_encoded_html_body", template.B64HTMLBody)
	d.Set("plain_text_body", template.TextBody)

	d.SetId(fmt.Sprintf("%s/%s/%s", tenantID, templateName, language))

	return diags
}
