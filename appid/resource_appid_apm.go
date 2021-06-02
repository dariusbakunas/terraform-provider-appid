package appid

import (
	"context"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDAPM() *schema.Resource {
	return &schema.Resource{
		Description:   "AppID advanced password management configuration",
		ReadContext:   resourceAppIDAPMRead,
		CreateContext: resourceAppIDAPMCreate,
		UpdateContext: resourceAppIDAPMCreate,
		DeleteContext: resourceAppIDAPMDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "`true` if APM is enabled",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"prevent_password_with_username": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"password_reuse": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"max_password_reuse": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  8,
						},
					},
				},
			},
			"password_expiration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"days_to_expire": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
						},
					},
				},
			},
			"lockout_policy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"lockout_time_sec": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1800,
						},
						"num_of_attempts": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
			},
			"min_password_change_interval": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"min_hours_to_change_password": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
		},
	}
}

func resourceAppIDAPMRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	tenantID := d.Id()
	c := m.(*appid.AppIDManagementV4)

	apm, _, err := c.GetCloudDirectoryAdvancedPasswordManagementWithContext(ctx, &appid.GetCloudDirectoryAdvancedPasswordManagementOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID APM configuration: %s", err)
	}

	if apm.AdvancedPasswordManagement != nil {
		d.Set("enabled", *apm.AdvancedPasswordManagement.Enabled)

		if err := d.Set("password_reuse", flattenPasswordReuse(apm.AdvancedPasswordManagement.PasswordReuse)); err != nil {
			return diag.Errorf("Failed setting password_reuse: %s", err)
		}

		if apm.AdvancedPasswordManagement.PreventPasswordWithUsername != nil {
			d.Set("prevent_password_with_username", *apm.AdvancedPasswordManagement.PreventPasswordWithUsername.Enabled)
		}

		if err := d.Set("password_expiration", flattenPasswordExpiration(apm.AdvancedPasswordManagement.PasswordExpiration)); err != nil {
			return diag.Errorf("Failed setting password_expiration: %s", err)
		}

		if err := d.Set("lockout_policy", flattenLockoutPolicy(apm.AdvancedPasswordManagement.LockOutPolicy)); err != nil {
			return diag.Errorf("Failed setting lockout_policy: %s", err)
		}
		if err := d.Set("min_password_change_interval", flattenPasswordChangeInterval(apm.AdvancedPasswordManagement.MinPasswordChangeInterval)); err != nil {
			return diag.Errorf("Failed setting min_password_change_interval: %s", err)
		}

	}

	d.Set("tenant_id", tenantID)
	return diags
}

func resourceAppIDAPMCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	enabled := d.Get("enabled").(bool)

	c := m.(*appid.AppIDManagementV4)

	config := &appid.SetCloudDirectoryAdvancedPasswordManagementOptions{
		TenantID: &tenantID,
		AdvancedPasswordManagement: &appid.ApmSchemaAdvancedPasswordManagement{
			Enabled:                   &enabled,
			PasswordReuse:             expandPasswordReuse(d.Get("password_reuse").([]interface{})),
			PasswordExpiration:        expandPasswordExpiration(d.Get("password_expiration").([]interface{})),
			LockOutPolicy:             expandLockoutPolicy(d.Get("lockout_policy").([]interface{})),
			MinPasswordChangeInterval: expandMinPasswordChangeInterval(d.Get("min_password_change_interval").([]interface{})),
			PreventPasswordWithUsername: &appid.ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername{
				Enabled: getBoolPtr(d.Get("prevent_password_with_username").(bool)),
			},
		},
	}

	_, _, err := c.SetCloudDirectoryAdvancedPasswordManagementWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error updating AppID APM configuration: %s", err)
	}

	d.SetId(tenantID)
	return resourceAppIDAPMRead(ctx, d, m)
}

func expandPasswordReuse(reuse []interface{}) *appid.ApmSchemaAdvancedPasswordManagementPasswordReuse {
	if len(reuse) == 0 || reuse[0] == nil {
		return nil
	}

	mReuse := reuse[0].(map[string]interface{})

	result := &appid.ApmSchemaAdvancedPasswordManagementPasswordReuse{
		Enabled: getBoolPtr(mReuse["enabled"].(bool)),
		Config: &appid.ApmSchemaAdvancedPasswordManagementPasswordReuseConfig{
			MaxPasswordReuse: getInt64Ptr(int64(mReuse["max_password_reuse"].(int))),
		},
	}

	return result
}

func expandPasswordExpiration(exp []interface{}) *appid.ApmSchemaAdvancedPasswordManagementPasswordExpiration {
	if len(exp) == 0 || exp[0] == nil {
		return nil
	}

	mExp := exp[0].(map[string]interface{})

	result := &appid.ApmSchemaAdvancedPasswordManagementPasswordExpiration{
		Enabled: getBoolPtr(mExp["enabled"].(bool)),
		Config: &appid.ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig{
			DaysToExpire: getInt64Ptr(int64(mExp["days_to_expire"].(int))),
		},
	}

	return result
}

func expandLockoutPolicy(loc []interface{}) *appid.ApmSchemaAdvancedPasswordManagementLockOutPolicy {
	if len(loc) == 0 || loc[0] == nil {
		return nil
	}

	mLock := loc[0].(map[string]interface{})

	result := &appid.ApmSchemaAdvancedPasswordManagementLockOutPolicy{
		Enabled: getBoolPtr(mLock["enabled"].(bool)),
		Config: &appid.ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig{
			LockOutTimeSec: getInt64Ptr(int64(mLock["lockout_time_sec"].(int))),
			NumOfAttempts:  getInt64Ptr(int64(mLock["num_of_attempts"].(int))),
		},
	}

	return result
}

func expandMinPasswordChangeInterval(chg []interface{}) *appid.ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval {
	if len(chg) == 0 || chg[0] == nil {
		return nil
	}

	mChg := chg[0].(map[string]interface{})

	result := &appid.ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval{
		Enabled: getBoolPtr(mChg["enabled"].(bool)),
		Config: &appid.ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig{
			MinHoursToChangePassword: getInt64Ptr(int64(mChg["min_hours_to_change_password"].(int))),
		},
	}

	return result
}

func resourceAppIDAPMDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	tenantID := d.Get("tenant_id").(string)
	config := getDefaultAPMConfig()
	c := m.(*appid.AppIDManagementV4)

	_, _, err := c.SetCloudDirectoryAdvancedPasswordManagementWithContext(ctx, &appid.SetCloudDirectoryAdvancedPasswordManagementOptions{
		TenantID:                   &tenantID,
		AdvancedPasswordManagement: config,
	})

	if err != nil {
		return diag.Errorf("Error resetting AppID APM configuration: %s", err)
	}

	d.SetId("")

	return diags
}

func getDefaultAPMConfig() *appid.ApmSchemaAdvancedPasswordManagement {
	return &appid.ApmSchemaAdvancedPasswordManagement{
		Enabled: getBoolPtr(false),
		PasswordReuse: &appid.ApmSchemaAdvancedPasswordManagementPasswordReuse{
			Enabled: getBoolPtr(false),
			Config: &appid.ApmSchemaAdvancedPasswordManagementPasswordReuseConfig{
				MaxPasswordReuse: getInt64Ptr(8),
			},
		},
		PasswordExpiration: &appid.ApmSchemaAdvancedPasswordManagementPasswordExpiration{
			Enabled: getBoolPtr(false),
			Config: &appid.ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig{
				DaysToExpire: getInt64Ptr(30),
			},
		},
		MinPasswordChangeInterval: &appid.ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval{
			Enabled: getBoolPtr(false),
			Config: &appid.ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig{
				MinHoursToChangePassword: getInt64Ptr(0),
			},
		},
		LockOutPolicy: &appid.ApmSchemaAdvancedPasswordManagementLockOutPolicy{
			Enabled: getBoolPtr(false),
			Config: &appid.ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig{
				LockOutTimeSec: getInt64Ptr(1800),
				NumOfAttempts:  getInt64Ptr(3),
			},
		},
		PreventPasswordWithUsername: &appid.ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername{
			Enabled: getBoolPtr(false),
		},
	}
}
