package appid

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppIDApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDApplicationCreate,
		ReadContext:   dataSourceAppIDApplicationRead, // reusing data source read, same schema
		DeleteContext: resourceAppIDApplicationDelete,
		UpdateContext: resourceAppIDApplicationUpdate,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
			},
			"type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "regularwebapp",
				ValidateFunc: validation.StringInSlice([]string{"regularwebapp", "singlepageapp"}, false),
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
			"scopes": {
				Type: schema.TypeList,
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

	c := m.(*Client)

	input := &CreateApplicationInput{
		Name: name,
		Type: appType,
	}

	log.Printf("[DEBUG] Creating AppID application: %+v", input)
	app, err := c.ApplicationAPI.CreateApplication(ctx, tenantID, input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, app.ClientID))

	if err := d.Set("client_id", app.ClientID); err != nil {
		return diag.FromErr(err)
	}

	return dataSourceAppIDApplicationRead(ctx, d, m)
}

func resourceAppIDApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)
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
	c := m.(*Client)
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
