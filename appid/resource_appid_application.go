package appid

import (
	"context"
	"fmt"
	"log"
	"strings"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppIDApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDApplicationCreate,
		ReadContext:   resourceAppIDApplicationRead, // reusing data source read, same schema
		DeleteContext: resourceAppIDApplicationDelete,
		UpdateContext: resourceAppIDApplicationUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"client_id": {
				Description: "The `client_id` is a public identifier for applications",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "The application name to be registered. Application name cannot exceed 50 characters.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
			},
			"type": {
				Description:  "The type of application to be registered. Allowed types are `regularwebapp` and `singlepageapp`, default is `regularwebapp`.",
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "regularwebapp",
				ValidateFunc: validation.StringInSlice([]string{"regularwebapp", "singlepageapp"}, false),
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
			"scopes": {
				Description: "A `scope` is a runtime action in your application that you register with IBM Cloud App ID to create an access permission",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceAppIDApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := d.Id()
	idParts := strings.Split(id, "/")

	tenantID := idParts[0]
	clientID := idParts[1]

	c := m.(*appid.AppIDManagementV4)

	app, _, err := c.GetApplicationWithContext(ctx, &appid.GetApplicationOptions{
		TenantID: getStringPtr(tenantID),
		ClientID: getStringPtr(clientID),
	})

	if err != nil {
		return diag.Errorf("Error getting AppID application: %s", err)
	}

	log.Printf("[DEBUG] Read application: %+v", app)

	scopes, _, err := c.GetApplicationScopesWithContext(ctx, &appid.GetApplicationScopesOptions{
		TenantID: getStringPtr(tenantID),
		ClientID: getStringPtr(clientID),
	})

	if err != nil {
		return diag.Errorf("Error getting AppID application scopes: %s", err)
	}

	log.Printf("[DEBUG] Read application scopes: %v", scopes)

	if app.Name != nil {
		d.Set("name", *app.Name)
	}

	if app.Secret != nil {
		d.Set("secret", *app.Secret)
	}

	if app.OAuthServerURL != nil {
		d.Set("oauth_server_url", *app.OAuthServerURL)
	}

	if app.ProfilesURL != nil {
		d.Set("profiles_url", *app.ProfilesURL)
	}

	if app.DiscoveryEndpoint != nil {
		d.Set("discovery_endpoint", *app.DiscoveryEndpoint)
	}

	if app.Type != nil {
		d.Set("type", *app.Type)
	}

	if err := d.Set("scopes", scopes.Scopes); err != nil {
		return diag.FromErr(err)
	}

	d.Set("tenant_id", tenantID)
	d.Set("client_id", clientID)

	return diags
}

func resourceAppIDApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	name := d.Get("name").(string)
	appType := d.Get("type").(string)

	scopes := expandStringList(d.Get("scopes").([]interface{}))

	c := m.(*appid.AppIDManagementV4)

	input := &appid.RegisterApplicationOptions{
		TenantID: getStringPtr(tenantID),
		Name:     getStringPtr(name),
		Type:     getStringPtr(appType),
	}

	log.Printf("[DEBUG] Creating AppID application: %+v", input)
	app, _, err := c.RegisterApplicationWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error creating AppID application: %s", err)
	}

	if len(scopes) != 0 {
		scopeOpts := &appid.PutApplicationsScopesOptions{
			TenantID: getStringPtr(tenantID),
			ClientID: getStringPtr(*app.ClientID),
			Scopes:   scopes,
		}

		_, _, err := c.PutApplicationsScopesWithContext(ctx, scopeOpts)

		if err != nil {
			// this is not ideal, but we have to delete created app otherwise next apply will fail
			// another option would be adding separate application_scopes resource
			_, deleteErr := c.DeleteApplicationWithContext(ctx, &appid.DeleteApplicationOptions{
				TenantID: getStringPtr(tenantID),
				ClientID: getStringPtr(*app.ClientID),
			})
			diags := diag.FromErr(err)

			if deleteErr != nil {
				log.Printf("[WARN] Failed cleaning up partially created application: %s/%s", *app.TenantID, *app.ClientID)
				diags = append(diags, diag.FromErr(deleteErr)...)
			}

			return diags
		}
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, *app.ClientID))

	return resourceAppIDApplicationRead(ctx, d, m)
}

func resourceAppIDApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*appid.AppIDManagementV4)
	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	log.Printf("[DEBUG] Deleting AppID application: %s", d.Id())

	_, err := c.DeleteApplicationWithContext(ctx, &appid.DeleteApplicationOptions{
		TenantID: getStringPtr(tenantID),
		ClientID: getStringPtr(clientID),
	})

	if err != nil {
		return diag.Errorf("Error deleting AppID application: %s", err)
	}

	d.SetId("")

	log.Printf("[DEBUG] Finished deleting AppID application: %s", d.Id())

	return diags
}

func resourceAppIDApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*appid.AppIDManagementV4)
	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	if d.HasChange("name") {
		name := d.Get("name").(string)

		log.Printf("[DEBUG] Updating AppID application: %s", d.Id())
		_, _, err := c.UpdateApplicationWithContext(ctx, &appid.UpdateApplicationOptions{
			TenantID: getStringPtr(tenantID),
			Name:     getStringPtr(name),
			ClientID: getStringPtr(clientID),
		})

		if err != nil {
			return diag.Errorf("Error updating AppID application: %s", err)
		}
	}

	if d.HasChange("scopes") {
		scopes := expandStringList(d.Get("scopes").([]interface{}))

		scopeOpts := &appid.PutApplicationsScopesOptions{
			TenantID: &tenantID,
			ClientID: &clientID,
			Scopes:   scopes,
		}

		_, _, err := c.PutApplicationsScopesWithContext(ctx, scopeOpts)

		if err != nil {
			return diag.Errorf("Error updating application scopes: %s", err)
		}
	}

	log.Printf("[DEBUG] Finished updating AppID application: %s", d.Id())
	return resourceAppIDApplicationRead(ctx, d, m)
}
