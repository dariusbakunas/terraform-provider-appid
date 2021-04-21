package appid

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppIDTokenConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDTokenConfigCreate,
		ReadContext:   resourceAppIDTokenConfigRead,
		UpdateContext: resourceAppIDTokenConfigUpdate,
		DeleteContext: resourceAppIDTokenConfigDelete,
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
			"anonymous_token_expires_in": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"anonymous_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"refresh_token_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"access_token_claim": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"saml", "cloud_directory", "appid_custom", "facebook", "google", "ibmid", "attributes", "roles"}, false),
						},
						"source_claim": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_claim": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"id_token_claim": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"saml", "cloud_directory", "appid_custom", "facebook", "google", "ibmid", "attributes", "roles"}, false),
						},
						"source_claim": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_claim": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAppIDTokenConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)

	c := m.(*Client)

	input := expandTokenConfig(d)

	log.Printf("[DEBUG] Applying AppID token config: %v", input)
	err := c.ConfigAPI.UpdateTokens(ctx, tenantID, input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tenantID)

	return resourceAppIDTokenConfigRead(ctx, d, m)
}

func resourceAppIDTokenConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)

	c := m.(*Client)

	tokenConfig, err := c.ConfigAPI.GetTokens(ctx, tenantID)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Received AppID token config: %v", tokenConfig)

	if tokenConfig.Access != nil {
		if err := d.Set("access_token_expires_in", tokenConfig.Access.ExpiresIn); err != nil {
			return diag.FromErr(err)
		}
	}

	if tokenConfig.Refresh != nil {
		if err := d.Set("refresh_token_enabled", *tokenConfig.Refresh.Enabled); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("refresh_token_expires_in", tokenConfig.Refresh.ExpiresIn); err != nil {
			return diag.FromErr(err)
		}
	}

	if tokenConfig.AnonymousAccess != nil {
		if err := d.Set("anonymous_access_enabled", *tokenConfig.AnonymousAccess.Enabled); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("anonymous_token_expires_in", tokenConfig.AnonymousAccess.ExpiresIn); err != nil {
			return diag.FromErr(err)
		}
	}

	if tokenConfig.AccessTokenClaims != nil {
		if err := d.Set("access_token_claim", flattenTokenClaims(tokenConfig.AccessTokenClaims)); err != nil {
			return diag.FromErr(err)
		}
	}

	if tokenConfig.IDTokenClaims != nil {
		if err := d.Set("id_token_claim", flattenTokenClaims(tokenConfig.IDTokenClaims)); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func expandTokenClaims(l []interface{}) []TokenClaim {
	if len(l) == 0 {
		return nil
	}

	result := make([]TokenClaim, len(l))

	for i, item := range l {
		cMap := item.(map[string]interface{})

		claim := TokenClaim{
			Source: cMap["source"].(string),
		}

		// source_claim and destination_claim are optional
		if sClaim, ok := cMap["source_claim"]; ok {
			claim.SourceClaim = getStringPtr(sClaim.(string))
		}

		if dClaim, ok := cMap["destination_claim"]; ok {
			claim.DestinationClaim = getStringPtr(dClaim.(string))
		}

		result[i] = claim
	}

	return result
}

func expandTokenConfig(d *schema.ResourceData) *TokenConfig {
	config := &TokenConfig{}

	if accessExpiresIn, ok := d.GetOk("access_token_expires_in"); ok {
		config.Access = &AccessTokenConfig{
			ExpiresIn: accessExpiresIn.(int),
		}
	}

	if anonymousExpiresIn, ok := d.GetOk("anonymous_token_expires_in"); ok {
		config.AnonymousAccess = &AnonymusAccessConfig{
			ExpiresIn: anonymousExpiresIn.(int),
		}
	}

	if refreshExpiresIn, ok := d.GetOk("refresh_token_expires_in"); ok {
		config.Refresh = &RefreshTokenConfig{
			ExpiresIn: refreshExpiresIn.(int),
		}
	}

	// can't really use GetOk with bool
	anonymousAccessEnabled := d.Get("anonymous_access_enabled")

	if anonymousAccessEnabled != nil {
		if config.AnonymousAccess == nil {
			config.AnonymousAccess = &AnonymusAccessConfig{}
		}

		config.AnonymousAccess.Enabled = getBoolPtr(anonymousAccessEnabled.(bool))
	}

	refreshTokenEnabled := d.Get("refresh_token_enabled")

	if refreshTokenEnabled != nil {
		if config.Refresh == nil {
			config.Refresh = &RefreshTokenConfig{}
		}

		config.Refresh.Enabled = getBoolPtr(refreshTokenEnabled.(bool))
	}

	if accessClaims, ok := d.GetOk("access_token_claim"); ok {
		config.AccessTokenClaims = expandTokenClaims(accessClaims.(*schema.Set).List())
	}

	if idClaims, ok := d.GetOk("id_token_claim"); ok {
		config.IDTokenClaims = expandTokenClaims(idClaims.(*schema.Set).List())
	}

	return config
}

func resourceAppIDTokenConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)

	c := m.(*Client)

	// AppID resets value to default if it is not provided, so we can't do partial updates
	input := expandTokenConfig(d)

	log.Printf("[DEBUG] Updating AppID token config: %v", input)
	err := c.ConfigAPI.UpdateTokens(ctx, tenantID, input)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAppIDTokenConfigRead(ctx, d, m)
}

func resourceAppIDTokenConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")

	return diags
}
