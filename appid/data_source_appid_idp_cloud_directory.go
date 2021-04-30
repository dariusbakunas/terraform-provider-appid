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
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"self_service_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"signup_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"interactions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
								},
							},
						},
						"identity_field": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
		if err := d.Set("config", flattenCloudDirectoryConfig(config.Config)); err != nil {
			return diag.Errorf("failed setting config: %s", err)
		}
	}

	d.SetId(fmt.Sprintf("%s/idp/cloud_directory", tenantID))

	return diags
}

func flattenCloudDirectoryConfig(config *api.CloudDirectoryConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	mConfig := map[string]interface{}{}
	mConfig["self_service_enabled"] = config.SelfServiceEnabled

	if config.SignupEnabled != nil {
		mConfig["signup_enabled"] = *config.SignupEnabled
	}

	if config.IdentityField != "" {
		mConfig["identity_field"] = config.IdentityField
	}

	mConfig["interactions"] = flattenCloudDirectoryConfigInteractions(&config.Interactions)

	return []interface{}{mConfig}
}

func flattenCloudDirectoryConfigInteractions(interactions *api.CloudDirectoryInteractions) []interface{} {
	if interactions == nil {
		return []interface{}{}
	}

	mInteractions := map[string]interface{}{}

	mInteractions["welcome_enabled"] = interactions.WelcomeEnabled
	mInteractions["reset_password_enabled"] = interactions.ResetPasswordEnabled
	mInteractions["reset_password_notification_enabled"] = interactions.ResetPasswordNotificationEnabled
	mInteractions["identity_confirmation"] = flattenIdentityConfirmation(&interactions.IdentityConfirmation)

	return []interface{}{mInteractions}
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
