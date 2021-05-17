package appid

import (
	"context"
	"log"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDIDPCustom() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDIDPCustomCreate,
		ReadContext:   dataSourceAppIDIDPCustomRead,
		DeleteContext: resourceAppIDIDPCustomDelete,
		UpdateContext: resourceAppIDIDPCustomUpdate,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAppIDIDPCustomCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	c := m.(*appid.AppIDManagementV4)

	config := &appid.SetCustomIDPOptions{
		TenantID: getStringPtr(tenantID),
		IsActive: getBoolPtr(isActive),
	}

	if isActive {
		config.Config = &appid.CustomIDPConfigParamsConfig{}

		if pKey, ok := d.GetOk("public_key"); ok {
			config.Config.PublicKey = getStringPtr(pKey.(string))
		}
	}

	log.Printf("[DEBUG] Applying custom IDP config: %v", config)
	_, _, err := c.SetCustomIDPWithContext(ctx, config)

	if err != nil {
		return diag.FromErr(err)
	}

	return dataSourceAppIDIDPCustomRead(ctx, d, m)
}

func customIDPDefaults(tenantID string) *appid.SetCustomIDPOptions {
	return &appid.SetCustomIDPOptions{
		TenantID: getStringPtr(tenantID),
		IsActive: getBoolPtr(false),
	}
}

func resourceAppIDIDPCustomDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*appid.AppIDManagementV4)
	tenantID := d.Get("tenant_id").(string)
	config := customIDPDefaults(tenantID)

	log.Printf("[DEBUG] Resetting custom IDP config: %v", config)
	_, _, err := c.SetCustomIDPWithContext(ctx, config)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceAppIDIDPCustomUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// since this is configuration we can reuse create method
	return resourceAppIDIDPCustomCreate(ctx, d, m)
}
