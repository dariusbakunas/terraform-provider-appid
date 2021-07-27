package appid

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDApplicationIDs() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceAppIDApplicationIDsRead,
		DeprecationMessage: "This datasource will be removed in next release, use appid_applications instead",
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_ids": {
				Description: "A List of application client IDs for current applications in AppID instance",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceAppIDApplicationIDsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	apps, _, err := c.ListApplicationsWithContext(ctx, &appid.ListApplicationsOptions{
		TenantID: getStringPtr(tenantID),
	})

	if err != nil {
		return diag.Errorf("Error getting application IDs: %s", err)
	}

	ids := make([]string, 0)

	for _, app := range apps.Applications {
		ids = append(ids, *app.ClientID)
	}

	d.SetId(fmt.Sprintf("%s/ids", tenantID))

	if err := d.Set("client_ids", ids); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
