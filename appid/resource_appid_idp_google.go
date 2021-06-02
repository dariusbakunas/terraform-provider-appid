package appid

import (
	"context"
	"log"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDIDPGoogle() *schema.Resource {
	return &schema.Resource{
		Description:   "Update Google identity provider configuration.",
		CreateContext: resourceAppIDIDPGoogleCreate,
		ReadContext:   resourceAppIDIDPGoogleRead,
		DeleteContext: resourceAppIDIDPGoogleDelete,
		UpdateContext: resourceAppIDIDPGoogleUpdate,
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
			"is_active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Description: "Google application id",
							Type:        schema.TypeString,
							Required:    true,
						},
						"application_secret": {
							Description: "Google application secret",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
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

func resourceAppIDIDPGoogleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Id()
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

	d.Set("tenant_id", tenantID)

	return diags
}

func resourceAppIDIDPGoogleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	c := m.(*appid.AppIDManagementV4)

	config := &appid.SetGoogleIDPOptions{
		TenantID: getStringPtr(tenantID),
		IDP: &appid.FacebookGoogleConfigParams{
			IsActive: getBoolPtr(isActive),
		},
	}

	if isActive {
		config.IDP.Config = expandGoogleIDPConfig(d.Get("config").([]interface{}))
	}

	log.Printf("[DEBUG] Applying Google IDP config: %v", config)
	_, _, err := c.SetGoogleIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error applying Google IDP configuration: %s", err)
	}

	d.SetId(tenantID)

	return resourceAppIDIDPGoogleRead(ctx, d, m)
}

func resourceAppIDIDPGoogleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*appid.AppIDManagementV4)
	tenantID := d.Get("tenant_id").(string)
	config := googleIDPConfigDefaults(tenantID)

	log.Printf("[DEBUG] Resetting Google IDP config: %v", config)
	_, _, err := c.SetGoogleIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error resetting Google IDP configuration: %s", err)
	}

	d.SetId("")

	return diags
}

func resourceAppIDIDPGoogleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// since this is configuration we can reuse create method
	return resourceAppIDIDPGoogleCreate(ctx, d, m)
}

func expandGoogleIDPConfig(cfg []interface{}) *appid.FacebookGoogleConfigParamsConfig {
	config := &appid.FacebookGoogleConfigParamsConfig{}

	if len(cfg) == 0 || cfg[0] == nil {
		return nil
	}

	mCfg := cfg[0].(map[string]interface{})

	config.IDPID = getStringPtr(mCfg["application_id"].(string))
	config.Secret = getStringPtr(mCfg["application_secret"].(string))

	return config
}

func googleIDPConfigDefaults(tenantID string) *appid.SetGoogleIDPOptions {
	return &appid.SetGoogleIDPOptions{
		TenantID: getStringPtr(tenantID),
		IDP: &appid.FacebookGoogleConfigParams{
			IsActive: getBoolPtr(false),
		},
	}
}
