package appid

import (
	"context"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDMedia() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDMediaRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"logo_url": {
				Description: "AppID Login logo URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAppIDMediaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	media, _, err := c.GetMediaWithContext(ctx, &appid.GetMediaOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID media: %s", err)
	}

	if media.Image != nil {
		d.Set("logo_url", *media.Image)
	}

	d.SetId(tenantID)
	return diags
}
