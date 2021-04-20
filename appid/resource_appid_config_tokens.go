package appid

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDConfigTokens() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDConfigTokensCreate,
		ReadContext:   resourceAppIDConfigTokensRead,
		UpdateContext: resourceAppIDConfigTokensUpdate,
		DeleteContext: resourceAppIDConfigTokensDelete,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_token_expires_in": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"refresh_token_expires_in": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAppIDConfigTokensCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)

	c := m.(*Client)

	input := expandTokenConfig(d)

	log.Printf("[DEBUG] Applying AppID token config: %+v", input)
	err := c.ConfigAPI.UpdateTokens(ctx, tenantID, input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tenantID)

	return diags
}

func resourceAppIDConfigTokensRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)

	c := m.(*Client)

	tokenConfig, err := c.ConfigAPI.GetTokens(ctx, tenantID)

	if err != nil {
		return diag.FromErr(err)
	}

	if tokenConfig.Access != nil {
		if err := d.Set("access_token_expires_in", tokenConfig.Access.ExpiresIn); err != nil {
			return diag.FromErr(err)
		}
	}

	if tokenConfig.Refresh != nil {
		if err := d.Set("refresh_token_expires_in", tokenConfig.Refresh.ExpiresIn); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func expandTokenConfig(d *schema.ResourceData) *TokenConfig {
	config := &TokenConfig{}

	if accessExpiresIn, ok := d.GetOk("access_token_expires_in"); ok {
		config.Access = &AccessTokenConfig{
			ExpiresIn: accessExpiresIn.(int),
		}
	}

	if refreshExpiresIn, ok := d.GetOk("refresh_token_expires_in"); ok {
		config.Refresh = &RefreshTokenConfig{
			ExpiresIn: refreshExpiresIn.(int),
		}
	}

	return config
}

func resourceAppIDConfigTokensUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)

	c := m.(*Client)

	// AppID resets value to default if it is not provided, so we can't do partial updates
	input := expandTokenConfig(d)

	log.Printf("[DEBUG] Updating AppID token config: %+v", input)
	err := c.ConfigAPI.UpdateTokens(ctx, tenantID, input)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAppIDConfigTokensRead(ctx, d, m)
}

func resourceAppIDConfigTokensDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")

	return diags
}
