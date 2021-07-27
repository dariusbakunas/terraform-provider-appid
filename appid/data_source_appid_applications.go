package appid

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDApplicationsRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func dataSourceAppIDApplicationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)

	c := m.(*appid.AppIDManagementV4)

	apps, _, err := c.ListApplicationsWithContext(ctx, &appid.ListApplicationsOptions{
		TenantID: getStringPtr(tenantID),
	})

	if err != nil {
		return diag.Errorf("Error listing AppID applications: %s", err)
	}

	applicationList := make([]interface{}, 0)

	for _, app := range apps.Applications {
		application := map[string]interface{}{}
		application["client_id"] = *app.ClientID
		application["name"] = *app.Name

		if app.Secret != nil {
			application["secret"] = *app.Secret
		}

		if app.OAuthServerURL != nil {
			application["oauth_server_url"] = *app.OAuthServerURL
		}

		if app.ProfilesURL != nil {
			application["profiles_url"] = *app.ProfilesURL
		}

		if app.DiscoveryEndpoint != nil {
			application["discovery_endpoint"] = *app.DiscoveryEndpoint
		}

		if app.Type != nil {
			application["type"] = *app.Type
		}

		scopes, _, err := c.GetApplicationScopesWithContext(ctx, &appid.GetApplicationScopesOptions{
			TenantID: getStringPtr(tenantID),
			ClientID: app.ClientID,
		})

		if err != nil {
			return diag.Errorf("Error getting AppID application scopes: %s", err)
		}

		application["scopes"] = flattenStringList(scopes.Scopes)
		applicationList = append(applicationList, application)
	}

	if err := d.Set("applications", applicationList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/applications", tenantID))
	return nil
}
