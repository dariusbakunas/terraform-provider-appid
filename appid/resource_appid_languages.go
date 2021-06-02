package appid

import (
	"context"
	"log"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDLanguages() *schema.Resource {
	return &schema.Resource{
		Description:   "User localization configuration",
		CreateContext: resourceAppIDLanguagesCreate,
		ReadContext:   resourceAppIDLanguagesRead,
		DeleteContext: resourceAppIDLanguagesDelete,
		UpdateContext: resourceAppIDLanguagesCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"languages": {
				Description: "The list of languages that can be used to customize email templates for Cloud Directory",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceAppIDLanguagesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Id()
	c := m.(*appid.AppIDManagementV4)

	langs, _, err := c.GetLocalizationWithContext(ctx, &appid.GetLocalizationOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting languages: %s", err)
	}

	d.Set("languages", langs.Languages)
	d.Set("tenant_id", tenantID)

	return diags
}

func resourceAppIDLanguagesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	languages := expandStringList(d.Get("languages").([]interface{}))

	c := m.(*appid.AppIDManagementV4)

	input := &appid.UpdateLocalizationOptions{
		TenantID:  &tenantID,
		Languages: languages,
	}

	log.Printf("[DEBUG] Updating languages: %+v", input)
	_, err := c.UpdateLocalizationWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error updating languages: %s", err)
	}

	d.SetId(tenantID)

	return resourceAppIDLanguagesRead(ctx, d, m)
}

func resourceAppIDLanguagesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*appid.AppIDManagementV4)
	tenantID := d.Get("tenant_id").(string)

	input := &appid.UpdateLocalizationOptions{
		TenantID:  &tenantID,
		Languages: []string{"en"}, // AppID default
	}

	log.Printf("[DEBUG] Resetting AppID languages: %+v", input)
	_, err := c.UpdateLocalizationWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error resetting AppID languages: %s", err)
	}

	d.SetId("")

	return diags
}
