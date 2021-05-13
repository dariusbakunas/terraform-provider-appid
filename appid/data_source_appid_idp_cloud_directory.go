package appid

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
)

func dataSourceAppIDIDPCloudDirectory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDIDPCloudDirectoryRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"self_service_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"signup_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"welcome_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"reset_password_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"reset_password_notification_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"identity_confirmation": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"methods": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
					},
				},
			},
			"identity_field": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppIDIDPCloudDirectoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	c := m.(*api.Client)

	config, err := c.IDPAPI.GetCloudDirectoryConfig(ctx, tenantID)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Got CloudDirectory IDP config: %+v", config)

	d.Set("is_active", config.IsActive)

	if config.Config != nil {
		d.Set("self_service_enabled", config.Config.SelfServiceEnabled)
		d.Set("signup_enabled", config.Config.SignupEnabled)

		if config.Config.IdentityField != "" {
			d.Set("identity_field", config.Config.IdentityField)
		}

		if config.Config.Interactions != nil {
			d.Set("welcome_enabled", config.Config.Interactions.WelcomeEnabled)
			d.Set("reset_password_enabled", config.Config.Interactions.ResetPasswordEnabled)
			d.Set("reset_password_notification_enabled", config.Config.Interactions.ResetPasswordNotificationEnabled)
			if config.Config.Interactions.IdentityConfirmation != nil {
				d.Set("identity_confirmation", flattenIdentityConfirmation(config.Config.Interactions.IdentityConfirmation))
			}
		}
	}

	d.SetId(fmt.Sprintf("%s/idp/cloud_directory", tenantID))

	return diags
}

func flattenIdentityConfirmation(confirmation *api.IdentityConfirmation) []interface{} {
	if confirmation == nil {
		return []interface{}{}
	}

	mConfirmation := map[string]interface{}{}
	mConfirmation["access_mode"] = confirmation.AccessMode
	mConfirmation["methods"] = flattenStringList(confirmation.Methods)

	return []interface{}{mConfirmation}
}
