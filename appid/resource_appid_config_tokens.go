package appid

import (
	"context"

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
			},
		},
	}
}

func resourceAppIDConfigTokensCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)

	c := m.(*Client)

	input := &TokenConfig{
		Access:          &AccessTokenConfig{},
		Refresh:         &RefreshTokenConfig{},
		AnonymousAccess: &AnonymusAccessConfig{},
	}

	if expiresIn, ok := d.GetOk("access_token_expires_in"); ok {
		input.Access.ExpiresIn = expiresIn.(int)
	}

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

	return diags
}

func resourceAppIDConfigTokensUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)

	c := m.(*Client)

	input := &TokenConfig{}

	if d.HasChange("access_token_expires_in") {
		expiresIn := d.Get("access_token_expires_in").(int)

		input.Access = &AccessTokenConfig{
			ExpiresIn: expiresIn,
		}
	}

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
