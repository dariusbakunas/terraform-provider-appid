package appid

import (
	"context"
	"log"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const defaultHeaderColor = "#EEF2F5" // AppID default

func resourceAppIDThemeColor() *schema.Resource {
	return &schema.Resource{
		Description:   "Colors of the App ID login widget",
		CreateContext: resourceAppIDThemeColorCreate,
		UpdateContext: resourceAppIDThemeColorCreate,
		ReadContext:   resourceAppIDThemeColorRead,
		DeleteContext: resourceAppIDThemeColorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service `tenantId`",
			},
			"header_color": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAppIDThemeColorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*appid.AppIDManagementV4)

	tenantID := d.Id()

	colors, _, err := c.GetThemeColorWithContext(ctx, &appid.GetThemeColorOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID theme colors: %s", err)
	}

	if colors.HeaderColor != nil {
		d.Set("header_color", *colors.HeaderColor)
	}

	d.Set("tenant_id", tenantID)

	return diags
}

func resourceAppIDThemeColorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)

	c := m.(*appid.AppIDManagementV4)

	input := &appid.PostThemeColorOptions{
		TenantID:    &tenantID,
		HeaderColor: getStringPtr(d.Get("header_color").(string)),
	}

	log.Printf("[DEBUG] Applying AppID theme color: %+v", input)

	_, err := c.PostThemeColorWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error setting AppID theme color: %s", err)
	}

	d.SetId(tenantID)

	return resourceAppIDThemeColorRead(ctx, d, m)
}

func resourceAppIDThemeColorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	input := &appid.PostThemeColorOptions{
		TenantID:    &tenantID,
		HeaderColor: getStringPtr(defaultHeaderColor),
	}

	log.Printf("[DEBUG] Resetting AppID theme color: %+v", input)

	_, err := c.PostThemeColorWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error resetting AppID theme color: %s", err)
	}

	d.SetId("")

	return diags
}
