package appid

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
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

	c := m.(*api.Client)

	config := &api.CustomIDP{
		IsActive: isActive,
	}

	if isActive {
		config.Config = &api.CustomIDPConfig{}

		if pKey, ok := d.GetOk("public_key"); ok {
			config.Config.PublicKey = pKey.(string)
		}
	}

	log.Printf("[DEBUG] Applying custom IDP config: %v", config)
	err := c.IDPAPI.UpdateCustomIDPConfig(ctx, tenantID, config)

	if err != nil {
		return diag.FromErr(err)
	}

	return dataSourceAppIDIDPCustomRead(ctx, d, m)
}

func customIDPDefaults() *api.CustomIDP {
	return &api.CustomIDP{
		IsActive: false,
	}
}

func resourceAppIDIDPCustomDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*api.Client)
	tenantID := d.Get("tenant_id").(string)
	config := customIDPDefaults()

	log.Printf("[DEBUG] Resetting custom IDP config: %v", config)
	err := c.IDPAPI.UpdateCustomIDPConfig(ctx, tenantID, config)

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
