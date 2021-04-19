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
			"access_token_claims": {
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
			"id_token_claims": {
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
			"source_claim":      v.SourceClaim,
			"destination_claim": v.DestinationClaim,
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
		if err := d.Set("access_token_claims", flattenTokenClaims(tokenConfig.AccessTokenClaims)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("id_token_claims", flattenTokenClaims(tokenConfig.IDTokenClaims)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
