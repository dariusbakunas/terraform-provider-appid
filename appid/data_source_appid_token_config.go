package appid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDTokenConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDTokenConfigRead,
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
			"source": v.Source,
		}

		if v.SourceClaim != nil {
			claim["source_claim"] = *v.SourceClaim
		}

		if v.DestinationClaim != nil {
			claim["destination_claim"] = *v.DestinationClaim
		}

		s = append(s, claim)
	}

	return s
}

func dataSourceAppIDTokenConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)

	c := m.(*Client)

	tokenConfig, err := c.ConfigAPI.GetTokens(ctx, tenantID)

	if err != nil {
		return diag.FromErr(err)
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

	if tokenConfig.Access != nil {
		d.Set("anonymous_token_expires_in", tokenConfig.AnonymousAccess.ExpiresIn)
	}

	if tokenConfig.Refresh != nil {
		d.Set("refresh_token_expires_in", tokenConfig.Refresh.ExpiresIn)

		if tokenConfig.Refresh.Enabled != nil {
			d.Set("refresh_token_expires_in", tokenConfig.Refresh.ExpiresIn)
		}
	}

	if tokenConfig.AnonymousAccess != nil {
		d.Set("anonymous_token_expires_in", tokenConfig.AnonymousAccess.ExpiresIn)

		if tokenConfig.AnonymousAccess.Enabled != nil {
			d.Set("anonymous_access_enabled", *tokenConfig.AnonymousAccess.Enabled)
		}
	}

	d.SetId(tenantID)

	return diags
}
