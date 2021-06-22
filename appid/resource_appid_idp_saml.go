package appid

import (
	"context"
	"log"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDIDPSaml() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDIDPSAMLCreate,
		ReadContext:   dataSourceAppIDIDPSAMLRead,
		DeleteContext: resourceAppIDIDPSAMLDelete,
		UpdateContext: resourceAppIDIDPSAMLUpdate,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sign_in_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"certificates": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MaxItems: 2,
							Required: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypt_response": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"sign_request": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"include_scoping": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"authn_context": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"class": {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional: true,
									},
									"comparison": {
										Type:     schema.TypeString,
										Optional: true,
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

func resourceAppIDIDPSAMLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	c := m.(*appid.AppIDManagementV4)

	config := &appid.SetSAMLIDPOptions{
		TenantID: getStringPtr(tenantID),
		IsActive: getBoolPtr(isActive),
	}

	if isActive {
		if cfg, ok := d.GetOk("config"); ok {
			config.Config = expandSAMLConfig(cfg.([]interface{}))
		}
	}

	log.Printf("[DEBUG] Applying SAML config: %v", config)
	_, _, err := c.SetSAMLIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error applying SAML IDP configuration: %s", err)
	}

	return dataSourceAppIDIDPSAMLRead(ctx, d, m)
}

func expandAuthNContext(ctx []interface{}) *appid.SAMLConfigParamsAuthnContext {
	context := &appid.SAMLConfigParamsAuthnContext{}

	if len(ctx) == 0 || ctx[0] == nil {
		return nil
	}

	mContext := ctx[0].(map[string]interface{})

	if comparison, ok := mContext["comparison"]; ok {
		context.Comparison = getStringPtr(comparison.(string))
	}

	if class, ok := mContext["class"]; ok {
		context.Class = expandStringList(class.([]interface{}))
	}

	return context
}

func expandSAMLConfig(cfg []interface{}) *appid.SAMLConfigParams {
	config := &appid.SAMLConfigParams{}

	if len(cfg) == 0 || cfg[0] == nil {
		return nil
	}

	mCfg := cfg[0].(map[string]interface{})

	config.EntityID = getStringPtr(mCfg["entity_id"].(string))
	config.SignInURL = getStringPtr(mCfg["sign_in_url"].(string))

	if dispName, ok := mCfg["display_name"]; ok {
		config.DisplayName = getStringPtr(dispName.(string))
	}

	if encResponse, ok := mCfg["encrypt_response"]; ok {
		config.EncryptResponse = getBoolPtr(encResponse.(bool))
	}

	if signRequest, ok := mCfg["sign_request"]; ok {
		config.SignRequest = getBoolPtr(signRequest.(bool))
	}

	if includeScoping, ok := mCfg["include_scoping"]; ok {
		config.IncludeScoping = getBoolPtr(includeScoping.(bool))
	}

	if certificates, ok := mCfg["certificates"]; ok {
		config.Certificates = []string{}

		for _, cert := range certificates.([]interface{}) {
			if cert != nil {
				config.Certificates = append(config.Certificates, cert.(string))
			}
		}
	}

	if context, ok := mCfg["authn_context"]; ok {
		config.AuthnContext = expandAuthNContext(context.([]interface{}))
	}

	return config
}

func samlConfigDefaults(tenantID string) *appid.SetSAMLIDPOptions {
	return &appid.SetSAMLIDPOptions{
		IsActive: getBoolPtr(false),
		TenantID: getStringPtr(tenantID),
	}
}

func resourceAppIDIDPSAMLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*appid.AppIDManagementV4)
	tenantID := d.Get("tenant_id").(string)
	config := samlConfigDefaults(tenantID)

	log.Printf("[DEBUG] Resetting SAML config: %v", config)
	_, _, err := c.SetSAMLIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error resetting SAML IDP configuration: %s", err)
	}

	d.SetId("")

	return diags
}

func resourceAppIDIDPSAMLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// since this is configuration we can reuse create method
	return resourceAppIDIDPSAMLCreate(ctx, d, m)
}
