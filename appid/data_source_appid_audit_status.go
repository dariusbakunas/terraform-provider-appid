package appid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
)

func dataSourceAppIDAuditStatus() *schema.Resource {
	return &schema.Resource{
		Description: "Tenant audit status",
		ReadContext: dataSourceAppIDAuditStatusRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service `tenantId`",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The auditing status of the tenant.",
			},
		},
	}
}

func dataSourceAppIDAuditStatusRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)

	c := m.(*api.Client)

	auditStatus, err := c.ConfigAPI.GetAuditStatus(ctx, tenantID)

	if err != nil {
		return diag.Errorf("error getting audit status: %s", err)
	}

	d.Set("is_active", auditStatus.IsActive)
	d.SetId("auditStatus")

	return diags
}
