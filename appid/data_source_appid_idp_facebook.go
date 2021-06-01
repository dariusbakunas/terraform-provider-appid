package appid

import (
	"context"
	"log"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDIDPFacebook() *schema.Resource {
	return &schema.Resource{
		Description: "Returns the Facebook identity provider configuration.",
		ReadContext: dataSourceAppIDIDPFacebookRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Description: "Facebook application id",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"application_secret": {
							Description: "Facebook application secret",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"redirect_url": {
				Description: "Paste the URI into the Valid OAuth redirect URIs field in the Facebook Login section of the Facebook Developers Portal",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAppIDIDPFacebookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	fb, _, err := c.GetFacebookIDPWithContext(ctx, &appid.GetFacebookIDPOptions{
		TenantID: getStringPtr(tenantID),
	})

	if err != nil {
		return diag.Errorf("Error loading Facebook IDP: %s", err)
	}

	log.Printf("[DEBUG] Got Facebook IDP config: %+v", fb)

	d.Set("is_active", *fb.IsActive)

	if fb.RedirectURL != nil {
		d.Set("redirect_url", *fb.RedirectURL)
	}

	if fb.Config != nil {
		if err := d.Set("config", flattenFacebookIDPConfig(fb.Config)); err != nil {
			return diag.Errorf("failed setting config: %s", err)
		}
	}

	d.SetId(tenantID)

	return diags
}

func flattenFacebookIDPConfig(config *appid.FacebookConfigParamsConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	mConfig := map[string]interface{}{}
	mConfig["application_id"] = *config.IDPID
	mConfig["application_secret"] = *config.Secret

	return []interface{}{mConfig}
}
