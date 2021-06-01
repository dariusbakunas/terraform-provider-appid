package appid

import (
	"context"
	"log"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDIDPGoogle() *schema.Resource {
	return &schema.Resource{
		Description: "Returns the Google identity provider configuration.",
		ReadContext: dataSourceAppIDIDPGoogleRead,
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
							Description: "Google application id",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"application_secret": {
							Description: "Google application secret",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"redirect_url": {
				Description: "Paste the URI into the into the Authorized redirect URIs field in the Google Developer Console",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAppIDIDPGoogleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	googleIDP, _, err := c.GetGoogleIDPWithContext(ctx, &appid.GetGoogleIDPOptions{
		TenantID: getStringPtr(tenantID),
	})

	if err != nil {
		return diag.Errorf("Error loading Google IDP: %s", err)
	}

	log.Printf("[DEBUG] Got Google IDP config: %+v", googleIDP)

	d.Set("is_active", *googleIDP.IsActive)

	if googleIDP.RedirectURL != nil {
		d.Set("redirect_url", *googleIDP.RedirectURL)
	}

	if googleIDP.Config != nil {
		if err := d.Set("config", flattenGoogleIDPConfig(googleIDP.Config)); err != nil {
			return diag.Errorf("failed setting config: %s", err)
		}
	}

	d.SetId(tenantID)

	return diags
}

func flattenGoogleIDPConfig(config *appid.GoogleConfigParamsConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	mConfig := map[string]interface{}{}
	mConfig["application_id"] = *config.IDPID
	mConfig["application_secret"] = *config.Secret

	return []interface{}{mConfig}
}
