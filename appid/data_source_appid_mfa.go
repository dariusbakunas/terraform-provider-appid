package appid

import (
	"context"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDMFA() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDMFARead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_active": {
				Description: "`true` if MFA is active",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func dataSourceAppIDMFARead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	mfa, _, err := c.GetMFAConfigWithContext(ctx, &appid.GetMFAConfigOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID MFA configuration: %s", err)
	}

	if mfa.IsActive != nil {
		d.Set("is_active", *mfa.IsActive)
	}

	d.SetId(tenantID)

	return diags
}
