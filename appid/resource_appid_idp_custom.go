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
		ReadContext:   resourceAppIDIDPCustomRead,
		DeleteContext: resourceAppIDIDPCustomDelete,
		UpdateContext: resourceAppIDIDPCustomUpdate,
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
			"is_active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"public_key": {
				Description: "This is the public key used to validate your signed JWT. It is required to be a PEM in the RS256 or greater format.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceAppIDIDPCustomRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Id()
	c := m.(*appid.AppIDManagementV4)

	config, _, err := c.GetCustomIDPWithContext(ctx, &appid.GetCustomIDPOptions{
		TenantID: getStringPtr(tenantID),
	})

	if err != nil {
		return diag.Errorf("Error loading custom IDP: %s", err)
	}

	log.Printf("[DEBUG] Got Custom IDP config: %+v", config)

	d.Set("is_active", *config.IsActive)

	if config.Config != nil && config.Config.PublicKey != nil {
		if err := d.Set("public_key", *config.Config.PublicKey); err != nil {
			return diag.Errorf("failed setting config: %s", err)
		}
	}

	d.Set("tenant_id", tenantID)

	return diags
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
		return diag.Errorf("Error applying custom IDP configuration: %s", err)
	}

	d.SetId(tenantID)

	return resourceAppIDIDPCustomRead(ctx, d, m)
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
		return diag.Errorf("Error resetting custom IDP configuration: %s", err)
	}

	d.SetId("")

	return diags
}

func resourceAppIDIDPCustomUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// since this is configuration we can reuse create method
	return resourceAppIDIDPCustomCreate(ctx, d, m)
}
