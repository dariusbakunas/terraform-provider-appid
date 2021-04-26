package appid

import (
	"context"
	"log"

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
					},
				},
			},
		},
	}
}

func resourceAppIDIDPSAMLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	c := m.(*Client)

	config := &SAML{
		IsActive: isActive,
	}

	if isActive {
		if cfg, ok := d.GetOk("config"); ok {
			config.Config = expandSAMLConfig(cfg.([]interface{}))
		}
	}

	log.Printf("[DEBUG] Applying SAML config: %+v", config)
	err := c.IDPService.UpdateSAMLConfig(ctx, tenantID, config)

	if err != nil {
		return diag.FromErr(err)
	}

	return dataSourceAppIDIDPSAMLRead(ctx, d, m)
}

func expandSAMLConfig(cfg []interface{}) *SAMLConfig {
	config := &SAMLConfig{}

	if len(cfg) == 0 || cfg[0] == nil {
		return config
	}

	mCfg := cfg[0].(map[string]interface{})

	config.EntityID = mCfg["entity_id"].(string)
	config.SignInURL = mCfg["sign_in_url"].(string)
	config.DisplayName = mCfg["display_name"].(string)

	if encResponse, ok := mCfg["encrypt_response"]; ok {
		config.EncryptResponse = getBoolPtr(encResponse.(bool))
	}

	if signRequest, ok := mCfg["sign_request"]; ok {
		config.SignRequest = getBoolPtr(signRequest.(bool))
	}

	if includeScoping, ok := mCfg["include_scoping"]; ok {
		config.IncludeScoping = getBoolPtr(includeScoping.(bool))
	}

	if certificates, ok := mCfg["certificates"].([]interface{}); ok && len(certificates) > 0 {
		config.Certificates = []string{}

		for _, cert := range certificates {
			config.Certificates = append(config.Certificates, cert.(string))
		}
	}

	return config
}

func samlConfigDefaults() *SAML {
	return &SAML{
		IsActive: false,
	}
}

func resourceAppIDIDPSAMLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)
	tenantID := d.Get("tenant_id").(string)
	config := samlConfigDefaults()

	log.Printf("[DEBUG] Resetting SAML config: %v", config)
	err := c.IDPService.UpdateSAMLConfig(ctx, tenantID, config)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceAppIDIDPSAMLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: add saml idp update
	return dataSourceAppIDIDPSAMLRead(ctx, d, m)
}
