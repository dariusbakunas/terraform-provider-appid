package appid

import (
	"context"
	"fmt"
	"log"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDApplicationRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_id": {
				Description: "The `client_id` is a public identifier for applications",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The application name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"secret": {
				Description: "The `secret` is a secret known only to the application and the authorization server",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"oauth_server_url": {
				Description: "Base URL for common OAuth endpoints, like `/authorization`, `/token` and `/publickeys`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"profiles_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"discovery_endpoint": {
				Description: "This URL returns OAuth Authorization Server Metadata",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The type of application to be registered. Allowed types are `regularwebapp` and `singlepageapp`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"scopes": {
				Description: "A `scope` is a runtime action in your application that you register with IBM Cloud App ID to create an access permission",
				Type:        schema.TypeList,
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

	c := m.(*appid.AppIDManagementV4)

	app, _, err := c.GetApplicationWithContext(ctx, &appid.GetApplicationOptions{
		TenantID: getStringPtr(tenantID),
		ClientID: getStringPtr(clientID),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Read application: %+v", app)

	scopes, _, err := c.GetApplicationScopesWithContext(ctx, &appid.GetApplicationScopesOptions{
		TenantID: getStringPtr(tenantID),
		ClientID: getStringPtr(clientID),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Read application scopes: %v", scopes)

	d.Set("name", *app.Name)

	if app.Secret != nil {
		d.Set("secret", *app.Secret)
	}

	d.Set("oauth_server_url", *app.OAuthServerURL)
	d.Set("profiles_url", *app.ProfilesURL)
	d.Set("discovery_endpoint", *app.DiscoveryEndpoint)
	d.Set("type", *app.Type)

	if err := d.Set("scopes", scopes.Scopes); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, clientID))
	return diags
}
