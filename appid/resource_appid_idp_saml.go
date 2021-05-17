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
		config.Config = expandSAMLConfig(d.Get("config").([]interface{}))
	}

	log.Printf("[DEBUG] Applying SAML config: %v", config)
	_, _, err := c.SetSAMLIDPWithContext(ctx, config)

	if err != nil {
		return diag.FromErr(err)
	}

	return dataSourceAppIDIDPSAMLRead(ctx, d, m)
}

func expandAuthNContext(ctx []interface{}) *appid.SAMLConfigParamsAuthnContext {
	context := &appid.SAMLConfigParamsAuthnContext{}

	if len(ctx) == 0 || ctx[0] == nil {
		return nil
	}

	mContext := ctx[0].(map[string]interface{})

	context.Comparison = getStringPtr(mContext["comparison"].(string))

	if class, ok := mContext["class"].([]interface{}); ok && len(class) > 0 {
		context.Class = expandStringList(class)
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
	config.DisplayName = getStringPtr(mCfg["display_name"].(string))

	if encResponse, ok := mCfg["encrypt_response"].(bool); ok {
		config.EncryptResponse = getBoolPtr(encResponse)
	}

	if signRequest, ok := mCfg["sign_request"].(bool); ok {
		config.SignRequest = getBoolPtr(signRequest)
	}

	if includeScoping, ok := mCfg["include_scoping"].(bool); ok {
		config.IncludeScoping = getBoolPtr(includeScoping)
	}

	if certificates, ok := mCfg["certificates"].([]interface{}); ok && len(certificates) > 0 {
		config.Certificates = []string{}

		for _, cert := range certificates {
			config.Certificates = append(config.Certificates, cert.(string))
		}
	}

	config.AuthnContext = expandAuthNContext(mCfg["authn_context"].([]interface{}))

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
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceAppIDIDPSAMLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// since this is configuration we can reuse create method
	return resourceAppIDIDPSAMLCreate(ctx, d, m)
}
