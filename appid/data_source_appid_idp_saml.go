package appid

import (
	"context"
	"fmt"
	"log"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDIDPSAML() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDIDPSAMLRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
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
						"authn_context": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"class": {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed: true,
									},
									"comparison": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAppIDIDPSAMLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	saml, _, err := c.GetSAMLIDPWithContext(ctx, &appid.GetSAMLIDPOptions{
		TenantID: getStringPtr(tenantID),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Got SAML IDP config: %+v", saml)

	d.Set("is_active", *saml.IsActive)

	if saml.Config != nil {
		if err := d.Set("config", flattenSAMLConfig(saml.Config)); err != nil {
			return diag.Errorf("failed setting config: %s", err)
		}
	}

	d.SetId(fmt.Sprintf("%s/idp/saml", tenantID))

	return diags
}

func flattenSAMLConfig(config *appid.SAMLConfigParams) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	mConfig := map[string]interface{}{}
	mConfig["entity_id"] = *config.EntityID
	mConfig["sign_in_url"] = *config.SignInURL
	mConfig["certificates"] = flattenStringList(config.Certificates)

	if config.DisplayName != nil {
		mConfig["display_name"] = *config.DisplayName
	}

	if config.SignRequest != nil {
		mConfig["sign_request"] = *config.SignRequest
	}

	if config.EncryptResponse != nil {
		mConfig["encrypt_response"] = *config.EncryptResponse
	}

	if config.IncludeScoping != nil {
		mConfig["include_scoping"] = *config.IncludeScoping
	}

	if config.AuthnContext != nil {
		mConfig["authn_context"] = flattenAuthNContext(config.AuthnContext)
	}

	return []interface{}{mConfig}
}

func flattenAuthNContext(context *appid.SAMLConfigParamsAuthnContext) []interface{} {
	if context == nil {
		return []interface{}{}
	}

	mContext := map[string]interface{}{}

	if context.Class != nil {
		class := []interface{}{}

		for _, c := range context.Class {
			class = append(class, c)
		}

		mContext["class"] = class
	}

	if context.Comparison != nil {
		mContext["comparison"] = *context.Comparison
	}

	return []interface{}{mContext}
}
