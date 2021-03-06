package appid

import (
	"context"
	"log"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDIDPCustom() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDIDPCustomRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_key": {
				Description: "This is the public key used to validate your signed JWT. It is required to be a PEM in the RS256 or greater format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAppIDIDPCustomRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
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

	d.SetId(tenantID)

	return diags
}
