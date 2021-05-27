package appid

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	c := m.(*appid.AppIDManagementV4)

	auditStatus, _, err := c.GetAuditStatusWithContext(ctx, &appid.GetAuditStatusOptions{
		TenantID: getStringPtr(tenantID),
	})

	if err != nil {
		return diag.Errorf("error getting audit status: %s", err)
	}

	d.Set("is_active", *auditStatus.IsActive)
	d.SetId(fmt.Sprintf("%s/auditStatus", tenantID))

	return diags
}
