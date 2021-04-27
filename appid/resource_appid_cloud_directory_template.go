package appid

import (
	"context"
	b64 "encoding/base64"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppIDCloudDirectoryTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDCloudDirectoryTemplateCreate,
		ReadContext:   dataSourceAppIDCloudDirectoryTemplateRead,
		DeleteContext: resourceAppIDCloudDirectoryTemplateDelete,
		UpdateContext: resourceAppIDCloudDirectoryTemplateUpdate,
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
				Required: true,
			},
			"html_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"base64_encoded_html_body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plain_text_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAppIDCloudDirectoryTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	templateName := d.Get("template_name").(string)
	language := d.Get("language").(string)

	input := &EmailTemplate{
		Subject: d.Get("subject").(string),
	}

	c := m.(*Client)

	if htmlBody, ok := d.GetOk("html_body"); ok {
		// don't want to set HTMLBody here otherwise might run into issues with Cloudfare filtering
		input.B64HTMLBody = b64.StdEncoding.EncodeToString([]byte(htmlBody.(string)))
	}

	if textBody, ok := d.GetOk("plain_text_body"); ok {
		input.TextBody = textBody.(string)
	}

	err := c.CloudDirectoryAPI.UpdateEmailTemplate(ctx, tenantID, templateName, language, input)

	if err != nil {
		return diag.FromErr(err)
	}

	return dataSourceAppIDCloudDirectoryTemplateRead(ctx, d, m)
}

func resourceAppIDCloudDirectoryTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	templateName := d.Get("template_name").(string)
	language := d.Get("language").(string)

	c := m.(*Client)

	log.Printf("[DEBUG] Deleting CD Email Template: %s", d.Id())

	err := c.CloudDirectoryAPI.DeleteEmailTemplate(ctx, tenantID, templateName, language)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceAppIDCloudDirectoryTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: implement template update
	return dataSourceAppIDCloudDirectoryTemplateRead(ctx, d, m)
}
