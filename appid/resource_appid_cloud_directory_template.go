package appid

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"log"
	"strings"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppIDCloudDirectoryTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDCloudDirectoryTemplateCreate,
		ReadContext:   resourceAppIDCloudDirectoryTemplateRead,
		DeleteContext: resourceAppIDCloudDirectoryTemplateDelete,
		UpdateContext: resourceAppIDCloudDirectoryTemplateUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"template_name": {
				Description:  "The type of email template. This can be `USER_VERIFICATION`, `WELCOME`, `PASSWORD_CHANGED`, `RESET_PASSWORD` or `MFA_VERIFICATION`",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(supportedTemplates, false),
				ForceNew:     true,
			},
			"language": {
				Description: "Preferred language for resource. Format as described at RFC5646. According to the configured languages codes returned from the `GET /management/v4/{tenantId}/config/ui/languages API`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "en",
				ForceNew:    true,
			},
			"subject": {
				Description: "The subject of the email",
				Type:        schema.TypeString,
				Required:    true,
			},
			"html_body": {
				Description: "The HTML body of the email",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"base64_encoded_html_body": {
				Description: "The HTML body of the email encoded in Base64",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"plain_text_body": {
				Description: "The text body of the email.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceAppIDCloudDirectoryTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := d.Id()
	idParts := strings.Split(id, "/")

	tenantID := idParts[0]
	templateName := idParts[1]
	language := idParts[2]

	c := m.(*appid.AppIDManagementV4)

	template, _, err := c.GetTemplateWithContext(ctx, &appid.GetTemplateOptions{
		TenantID:     getStringPtr(tenantID),
		TemplateName: getStringPtr(templateName),
		Language:     getStringPtr(language),
	})

	if err != nil {
		return diag.Errorf("Error loading Cloud Directory template: %s", err)
	}

	if template.Subject != nil {
		d.Set("subject", *template.Subject)
	}

	if template.HTMLBody != nil {
		d.Set("html_body", *template.HTMLBody)
	}

	if template.Base64EncodedHTMLBody != nil {
		d.Set("base64_encoded_html_body", *template.Base64EncodedHTMLBody)
	}

	if template.PlainTextBody != nil {
		d.Set("plain_text_body", *template.PlainTextBody)
	}

	d.Set("tenant_id", tenantID)
	d.Set("template_name", templateName)
	d.Set("language", language)

	return diags
}

func resourceAppIDCloudDirectoryTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	templateName := d.Get("template_name").(string)
	language := d.Get("language").(string)

	input := &appid.UpdateTemplateOptions{
		TenantID:     getStringPtr(tenantID),
		TemplateName: getStringPtr(templateName),
		Language:     getStringPtr(language),
		Subject:      getStringPtr(d.Get("subject").(string)),
	}

	c := m.(*appid.AppIDManagementV4)

	if htmlBody, ok := d.GetOk("html_body"); ok {
		// don't want to set HTMLBody here otherwise might run into issues with Cloudfare filtering
		input.Base64EncodedHTMLBody = getStringPtr(b64.StdEncoding.EncodeToString([]byte(htmlBody.(string))))
	}

	if textBody, ok := d.GetOk("plain_text_body"); ok {
		input.PlainTextBody = getStringPtr(textBody.(string))
	}

	log.Printf("[DEBUG] Updating CD Email Template: %+v", input)
	_, _, err := c.UpdateTemplateWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error updating Cloud Directory email template: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", tenantID, templateName, language))

	return resourceAppIDCloudDirectoryTemplateRead(ctx, d, m)
}

func resourceAppIDCloudDirectoryTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	templateName := d.Get("template_name").(string)
	language := d.Get("language").(string)

	c := m.(*appid.AppIDManagementV4)

	log.Printf("[DEBUG] Deleting CD Email Template: %s", d.Id())

	_, err := c.DeleteTemplateWithContext(ctx, &appid.DeleteTemplateOptions{
		TenantID:     getStringPtr(tenantID),
		TemplateName: getStringPtr(templateName),
		Language:     getStringPtr(language),
	})

	if err != nil {
		return diag.Errorf("Error deleting Cloud Directory email template: %s", err)
	}

	d.SetId("")

	return diags
}

func resourceAppIDCloudDirectoryTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// this is just a configuration, can reuse create method
	return resourceAppIDCloudDirectoryTemplateCreate(ctx, d, m)
}
