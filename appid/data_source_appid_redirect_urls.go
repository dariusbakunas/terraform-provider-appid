package appid

import (
	"context"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDRedirectURLs() *schema.Resource {
	return &schema.Resource{
		Description: "Redirect URIs that can be used as callbacks of App ID authentication flow",
		ReadContext: dataSourceAppIDRedirectURLsRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service `tenantId`",
			},
			"urls": {
				Description: "A list of redirect URLs",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceAppIDRedirectURLsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*appid.AppIDManagementV4)

	tenantID := d.Get("tenant_id").(string)

	urls, _, err := c.GetRedirectUrisWithContext(ctx, &appid.GetRedirectUrisOptions{
		TenantID: getStringPtr(tenantID),
	})
	if err != nil {
		return diag.Errorf("Error loading redirect urls: %s", err)
	}

	if err := d.Set("urls", urls.RedirectUris); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tenantID)

	return diags
}
