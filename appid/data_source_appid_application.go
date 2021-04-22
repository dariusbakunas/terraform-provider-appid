package appid

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDApplicationRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"oauth_server_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profiles_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"discovery_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppIDApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	c := m.(*Client)

	app, err := c.ConfigAPI.GetApplication(ctx, tenantID, clientID)

	if err := d.Set("name", app.Name); err != nil {
		return diag.FromErr(err)
	}

	if app.Secret != nil {
		if err := d.Set("secret", *app.Secret); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("oauth_server_url", app.OAuthServerURL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("profiles_url", app.ProfilesURL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("discovery_endpoint", app.DiscoveryEndpoint); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", app.Type); err != nil {
		return diag.FromErr(err)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, clientID))
	return diags
}
