package appid

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
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
			"scopes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceAppIDApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	c := m.(*api.Client)

	app, err := c.ApplicationAPI.GetApplication(ctx, tenantID, clientID)

	log.Printf("[DEBUG] Read application: %+v", app)

	if err != nil {
		return diag.FromErr(err)
	}

	scopes, err := c.ApplicationAPI.GetApplicationScopes(ctx, tenantID, clientID)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Read application scopes: %v", scopes)

	d.Set("name", app.Name)

	if app.Secret != nil {
		d.Set("secret", *app.Secret)
	}

	d.Set("oauth_server_url", app.OAuthServerURL)
	d.Set("profiles_url", app.ProfilesURL)
	d.Set("discovery_endpoint", app.DiscoveryEndpoint)
	d.Set("type", app.Type)

	if err := d.Set("scopes", scopes); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, clientID))
	return diags
}
