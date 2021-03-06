package appid

import (
	"context"
	"log"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppIDIDPCloudDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDIDPCloudDirectoryCreate,
		ReadContext:   resourceAppIDIDPCloudDirectoryRead,
		DeleteContext: resourceAppIDIDPCloudDirectoryDelete,
		UpdateContext: resourceAppIDIDPCloudDirectoryUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"self_service_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"signup_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"welcome_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"reset_password_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"reset_password_notification_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"identity_confirm_access_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "FULL",
				ValidateFunc: validation.StringInSlice([]string{"FULL", "RESTRICTIVE", "OFF"}, false),
			},
			"identity_confirm_methods": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"identity_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAppIDIDPCloudDirectoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Id()
	c := m.(*appid.AppIDManagementV4)

	config, _, err := c.GetCloudDirectoryIDPWithContext(ctx, &appid.GetCloudDirectoryIDPOptions{
		TenantID: getStringPtr(tenantID),
	})

	if err != nil {
		return diag.Errorf("Error loading Cloud Directory IDP: %s", err)
	}

	log.Printf("[DEBUG] Got CloudDirectory IDP config: %+v", config)

	d.Set("is_active", *config.IsActive)

	if config.Config != nil {
		d.Set("self_service_enabled", *config.Config.SelfServiceEnabled)

		if config.Config.SignupEnabled != nil {
			d.Set("signup_enabled", *config.Config.SignupEnabled)
		}

		if config.Config.IdentityField != nil {
			d.Set("identity_field", *config.Config.IdentityField)
		}

		if config.Config.Interactions != nil {
			d.Set("welcome_enabled", *config.Config.Interactions.WelcomeEnabled)
			d.Set("reset_password_enabled", *config.Config.Interactions.ResetPasswordEnabled)
			d.Set("reset_password_notification_enabled", *config.Config.Interactions.ResetPasswordNotificationEnable)
			d.Set("identity_confirm_access_mode", *config.Config.Interactions.IdentityConfirmation.AccessMode)
			d.Set("identity_confirm_methods", config.Config.Interactions.IdentityConfirmation.Methods)
		}
	}

	d.Set("tenant_id", tenantID)

	return diags
}

func resourceAppIDIDPCloudDirectoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	c := m.(*appid.AppIDManagementV4)

	config := &appid.SetCloudDirectoryIDPOptions{
		TenantID: getStringPtr(tenantID),
		IsActive: getBoolPtr(isActive),
		Config: &appid.CloudDirectoryConfigParams{
			SelfServiceEnabled: getBoolPtr(d.Get("self_service_enabled").(bool)),
			SignupEnabled:      getBoolPtr(d.Get("signup_enabled").(bool)),
			Interactions: &appid.CloudDirectoryConfigParamsInteractions{
				WelcomeEnabled:                  getBoolPtr(d.Get("welcome_enabled").(bool)),
				ResetPasswordEnabled:            getBoolPtr(d.Get("reset_password_enabled").(bool)),
				ResetPasswordNotificationEnable: getBoolPtr(d.Get("reset_password_notification_enabled").(bool)),
				IdentityConfirmation: &appid.CloudDirectoryConfigParamsInteractionsIdentityConfirmation{
					AccessMode: getStringPtr(d.Get("identity_confirm_access_mode").(string)),
				},
			},
		},
	}

	if idField, ok := d.GetOk("identity_field"); ok {
		config.Config.IdentityField = getStringPtr(idField.(string))
	}

	if methods, ok := d.GetOk("identity_confirm_methods"); ok {
		config.Config.Interactions.IdentityConfirmation.Methods = expandStringList(methods.([]interface{}))
	}

	log.Printf("[DEBUG] Applying Cloud Directory IDP config: %+v", config)
	_, _, err := c.SetCloudDirectoryIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error applying Cloud Directory IDP configuration: %s", err)
	}

	d.SetId(tenantID)

	return resourceAppIDIDPCloudDirectoryRead(ctx, d, m)
}

func resourceAppIDIDPCloudDirectoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// since this is configuration we can reuse create method
	return resourceAppIDIDPCloudDirectoryCreate(ctx, d, m)
}

func resourceAppIDIDPCloudDirectoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*appid.AppIDManagementV4)
	tenantID := d.Get("tenant_id").(string)
	config := cloudDirectoryDefaults(tenantID)

	log.Printf("[DEBUG] Resetting Cloud Directory IDP config: %v", config)
	_, _, err := c.SetCloudDirectoryIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error resetting Cloud Directory IDP configuration: %s", err)
	}

	d.SetId("")

	return diags
}

func cloudDirectoryDefaults(tenantID string) *appid.SetCloudDirectoryIDPOptions {
	return &appid.SetCloudDirectoryIDPOptions{
		TenantID: getStringPtr(tenantID),
		IsActive: getBoolPtr(false),
		Config: &appid.CloudDirectoryConfigParams{
			SelfServiceEnabled: getBoolPtr(true),
			Interactions: &appid.CloudDirectoryConfigParamsInteractions{
				IdentityConfirmation: &appid.CloudDirectoryConfigParamsInteractionsIdentityConfirmation{
					AccessMode: getStringPtr("FULL"),
					Methods:    []string{"email"},
				},
				WelcomeEnabled:                  getBoolPtr(true),
				ResetPasswordEnabled:            getBoolPtr(true),
				ResetPasswordNotificationEnable: getBoolPtr(true),
			},
		},
	}
}
