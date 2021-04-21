package appid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDConfigTokens() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDConfigTokensRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_token_expires_in": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"refresh_token_expires_in": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"anonymous_token_expires_in": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"anonymous_access_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"refresh_token_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"access_token_claim": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_claim": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"destination_claim": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"id_token_claim": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_claim": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_claim": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenTokenClaims(c []TokenClaim) []interface{} {
	var s []interface{}

	for _, v := range c {
		claim := map[string]interface{}{
			"source":            v.Source,
			"destination_claim": v.DestinationClaim,
		}

		if v.SourceClaim != nil {
			claim["source_claim"] = *v.SourceClaim
		}

		s = append(s, claim)
	}

	return s
}

func dataSourceAppIDConfigTokensRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)

	c := m.(*Client)

	tokenConfig, err := c.ConfigAPI.GetTokens(ctx, tenantID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tenantID)

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

	if tokenConfig.Access != nil {
		if err := d.Set("access_token_expires_in", tokenConfig.Access.ExpiresIn); err != nil {
			return diag.FromErr(err)
		}
	}

	if tokenConfig.Refresh != nil {
		if err := d.Set("refresh_token_expires_in", tokenConfig.Refresh.ExpiresIn); err != nil {
			return diag.FromErr(err)
		}

		if tokenConfig.Refresh.Enabled != nil {
			if err := d.Set("refresh_token_enabled", *tokenConfig.Refresh.Enabled); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if tokenConfig.AnonymousAccess != nil {
		if err := d.Set("anonymous_token_expires_in", tokenConfig.AnonymousAccess.ExpiresIn); err != nil {
			return diag.FromErr(err)
		}

		if tokenConfig.AnonymousAccess.Enabled != nil {
			if err := d.Set("anonymous_access_enabled", *tokenConfig.AnonymousAccess.Enabled); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return diags
}
