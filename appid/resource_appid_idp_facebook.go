package appid

import (
	"context"
	"log"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDIDPFacebook() *schema.Resource {
	return &schema.Resource{
		Description:   "Update Facebook identity provider configuration.",
		CreateContext: resourceAppIDIDPFacebookCreate,
		ReadContext:   resourceAppIDIDPFacebookRead,
		DeleteContext: resourceAppIDIDPFacebookDelete,
		UpdateContext: resourceAppIDIDPFacebookUpdate,
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
							Description: "Facebook application id",
							Type:        schema.TypeString,
							Required:    true,
						},
						"application_secret": {
							Description: "Facebook application secret",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
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

func resourceAppIDIDPFacebookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Id()
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

	d.Set("tenant_id", tenantID)

	return diags
}

func resourceAppIDIDPFacebookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	c := m.(*appid.AppIDManagementV4)

	config := &appid.SetFacebookIDPOptions{
		TenantID: getStringPtr(tenantID),
		IDP: &appid.FacebookGoogleConfigParams{
			IsActive: getBoolPtr(isActive),
		},
	}

	if isActive {
		config.IDP.Config = expandFBConfig(d.Get("config").([]interface{}))
	}

	log.Printf("[DEBUG] Applying Facebook IDP config: %v", config)
	_, _, err := c.SetFacebookIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error applying Facebook IDP configuration: %s", err)
	}

	d.SetId(tenantID)

	return resourceAppIDIDPFacebookRead(ctx, d, m)
}

func resourceAppIDIDPFacebookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*appid.AppIDManagementV4)
	tenantID := d.Get("tenant_id").(string)
	config := facebookIDPConfigDefaults(tenantID)

	log.Printf("[DEBUG] Resetting Facebook IDP config: %v", config)
	_, _, err := c.SetFacebookIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error resetting Facebook IDP configuration: %s", err)
	}

	d.SetId("")

	return diags
}

func resourceAppIDIDPFacebookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// since this is configuration we can reuse create method
	return resourceAppIDIDPFacebookCreate(ctx, d, m)
}

func expandFBConfig(cfg []interface{}) *appid.FacebookGoogleConfigParamsConfig {
	config := &appid.FacebookGoogleConfigParamsConfig{}

	if len(cfg) == 0 || cfg[0] == nil {
		return nil
	}

	mCfg := cfg[0].(map[string]interface{})

	config.IDPID = getStringPtr(mCfg["application_id"].(string))
	config.Secret = getStringPtr(mCfg["application_secret"].(string))

	return config
}

func facebookIDPConfigDefaults(tenantID string) *appid.SetFacebookIDPOptions {
	return &appid.SetFacebookIDPOptions{
		TenantID: getStringPtr(tenantID),
		IDP: &appid.FacebookGoogleConfigParams{
			IsActive: getBoolPtr(false),
		},
	}
}
