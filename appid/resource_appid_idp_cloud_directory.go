package appid

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
)

func resourceAppIDIDPCloudDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDIDPCloudDirectoryCreate,
		ReadContext:   dataSourceAppIDIDPCloudDirectoryRead,
		DeleteContext: resourceAppIDIDPCloudDirectoryDelete,
		UpdateContext: resourceAppIDIDPCloudDirectoryUpdate,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"identity_confirmation": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "FULL",
							ValidateFunc: validation.StringInSlice([]string{"FULL", "RESTRICTIVE", "OFF"}, false),
						},
						"methods": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
					},
				},
			},
			"identity_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAppIDIDPCloudDirectoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	c := m.(*api.Client)

	config := &api.CloudDirectoryIDP{
		IsActive: isActive,
	}

	if isActive {
		config.Config = &api.CloudDirectoryConfig{
			SelfServiceEnabled: d.Get("self_service_enabled").(bool),
			SignupEnabled:      getBoolPtr(d.Get("signup_enabled").(bool)),
			Interactions: &api.CloudDirectoryInteractions{
				WelcomeEnabled:                   d.Get("welcome_enabled").(bool),
				ResetPasswordEnabled:             d.Get("reset_password_enabled").(bool),
				ResetPasswordNotificationEnabled: d.Get("reset_password_notification_enabled").(bool),
				IdentityConfirmation:             expandIdentityConfirmation(d.Get("identity_confirmation").([]interface{})),
			},
			IdentityField: d.Get("identity_field").(string),
		}
	}

	log.Printf("[DEBUG] Applying Cloud Directory IDP config: %v", config)
	err := c.IDPAPI.UpdateCloudDirectoryConfig(ctx, tenantID, config)

	if err != nil {
		return diag.FromErr(err)
	}

	return dataSourceAppIDIDPCloudDirectoryRead(ctx, d, m)
}

func resourceAppIDIDPCloudDirectoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// since this is configuration we can reuse create method
	return resourceAppIDIDPCloudDirectoryCreate(ctx, d, m)
}

func resourceAppIDIDPCloudDirectoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*api.Client)
	tenantID := d.Get("tenant_id").(string)
	config := cloudDirectoryDefaults()

	log.Printf("[DEBUG] Resetting Cloud Directory IDP config: %v", config)
	err := c.IDPAPI.UpdateCloudDirectoryConfig(ctx, tenantID, config)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func expandIdentityConfirmation(ic []interface{}) *api.IdentityConfirmation {
	confirmation := &api.IdentityConfirmation{}

	if len(ic) == 0 || ic[0] == nil {
		return nil
	}

	mIc := ic[0].(map[string]interface{})

	confirmation.AccessMode = mIc["access_mode"].(string)
	if methods, ok := mIc["methods"]; ok {
		confirmation.Methods = expandStringList(methods.([]interface{}))
	}

	return confirmation
}

func cloudDirectoryDefaults() *api.CloudDirectoryIDP {
	return &api.CloudDirectoryIDP{
		IsActive: false,
		Config: &api.CloudDirectoryConfig{
			SelfServiceEnabled: true,
			Interactions: &api.CloudDirectoryInteractions{
				IdentityConfirmation: &api.IdentityConfirmation{
					AccessMode: "FULL",
					Methods:    []string{"email"},
				},
				WelcomeEnabled:                   true,
				ResetPasswordEnabled:             true,
				ResetPasswordNotificationEnabled: true,
			},
		},
	}
}