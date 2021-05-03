package appid

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
)

func resourceAppIDApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDApplicationCreate,
		ReadContext:   dataSourceAppIDApplicationRead, // reusing data source read, same schema
		DeleteContext: resourceAppIDApplicationDelete,
		UpdateContext: resourceAppIDApplicationUpdate,
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

func resourceAppIDApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	name := d.Get("name").(string)
	appType := d.Get("type").(string)

	var scopes = make([]string, 0)
	if data, ok := d.GetOk("scopes"); ok {
		for _, scope := range data.([]interface{}) {
			scopes = append(scopes, scope.(string))
		}
	}

	c := m.(*api.Client)

	input := &api.CreateApplicationInput{
		Name: name,
		Type: appType,
	}

	log.Printf("[DEBUG] Creating AppID application: %+v", input)
	app, err := c.ApplicationAPI.CreateApplication(ctx, tenantID, input)

	if err != nil {
		return diag.FromErr(err)
	}

	if len(scopes) != 0 {
		_, err := c.ApplicationAPI.SetApplicationScopes(ctx, tenantID, app.ClientID, scopes)

		if err != nil {
			// this is not ideal, but we have to delete created app otherwise next apply will fail
			// another option would be adding separate application_scopes resource
			deleteErr := c.ApplicationAPI.DeleteApplication(ctx, tenantID, app.ClientID)
			diags := diag.FromErr(err)

			if deleteErr != nil {
				log.Printf("[WARN] Failed cleaning up partially created application: %s/%s", app.TenantID, app.ClientID)
				diags = append(diags, diag.FromErr(deleteErr)...)
			}

			return diags
		}
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, app.ClientID))
	d.Set("client_id", app.ClientID)

	return dataSourceAppIDApplicationRead(ctx, d, m)
}

func resourceAppIDApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*api.Client)
	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	log.Printf("[DEBUG] Deleting AppID application: %s", d.Id())

	err := c.ApplicationAPI.DeleteApplication(ctx, tenantID, clientID)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Finished deleting AppID application: %s", d.Id())

	return diags
}

func resourceAppIDApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*api.Client)
	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	if d.HasChange("name") {
		name := d.Get("name").(string)

		log.Printf("[DEBUG] Updating AppID application: %s", d.Id())
		_, err := c.ApplicationAPI.UpdateApplication(ctx, tenantID, clientID, name)

		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Finished updating AppID application: %s", d.Id())
	}

	return dataSourceAppIDApplicationRead(ctx, d, m)
}
