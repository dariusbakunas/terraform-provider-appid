package appid

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDApplicationIDs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDApplicationIDsRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_ids": {
				Type: schema.TypeList,
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
	c := m.(*Client)

	apps, err := c.ApplicationAPI.ListApplications(ctx, tenantID)

	if err != nil {
		return diag.FromErr(err)
	}

	ids := make([]string, 0)

	for _, app := range apps {
		ids = append(ids, app.ClientID)
	}

	d.SetId(fmt.Sprintf("%s/ids", tenantID))

	if err := d.Set("client_ids", ids); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
