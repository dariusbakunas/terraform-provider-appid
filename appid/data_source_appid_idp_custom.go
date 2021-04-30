package appid

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
)

func dataSourceAppIDIDPCustom() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDIDPCustomRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppIDIDPCustomRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	c := m.(*api.Client)

	config, err := c.IDPAPI.GetCustomIDPConfig(ctx, tenantID)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Got Custom IDP config: %+v", config)

	d.Set("is_active", config.IsActive)

	if config.Config != nil && config.Config.PublicKey != "" {
		if err := d.Set("public_key", config.Config.PublicKey); err != nil {
			return diag.Errorf("failed setting config: %s", err)
		}
	}

	d.SetId(fmt.Sprintf("%s/idp/custom_idp", tenantID))

	return diags
}
