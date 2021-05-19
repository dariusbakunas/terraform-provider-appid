package appid

import (
	"context"
	b64 "encoding/base64"
	"log"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
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
				ForceNew: true,
			},
			"template_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(supportedTemplates, false),
				ForceNew:     true,
			},
			"language": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "en",
				ForceNew: true,
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
		return diag.FromErr(err)
	}

	return dataSourceAppIDCloudDirectoryTemplateRead(ctx, d, m)
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
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceAppIDCloudDirectoryTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// this is just a configuration, can reuse create method
	return resourceAppIDCloudDirectoryTemplateCreate(ctx, d, m)
}
