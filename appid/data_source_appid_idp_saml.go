package appid

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDIDPSAML() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDIDPSAMLRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"entity_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sign_in_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificates": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypt_response": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sign_request": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"include_scoping": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceAppIDIDPSAMLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	c := m.(*Client)

	saml, err := c.IDPService.GetSAMLConfig(ctx, tenantID)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Got SAML IDP config: %+v", saml)

	d.Set("is_active", saml.IsActive)

	if saml.Config != nil {
		d.Set("entity_id", saml.Config.EntityID)
		d.Set("sign_in_url", saml.Config.SignInURL)
		d.Set("certificates", flattenStringList(saml.Config.Certificates))
		d.Set("display_name", saml.Config.DisplayName)

		if saml.Config.SignRequest != nil {
			d.Set("sign_request", saml.Config.SignRequest)
		}

		if saml.Config.EncryptResponse != nil {
			d.Set("encrypt_Response", saml.Config.EncryptResponse)
		}

		if saml.Config.IncludeScoping != nil {
			d.Set("include_scoping", saml.Config.IncludeScoping)
		}
	}

	d.SetId(fmt.Sprintf("%s/idp/saml", tenantID))

	return diags
}
