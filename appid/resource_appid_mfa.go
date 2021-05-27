package appid

import (
	"context"
	"log"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDMFA() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAppIDMFARead,
		CreateContext: resourceAppIDMFACreate,
		UpdateContext: resourceAppIDMFACreate,
		DeleteContext: resourceAppIDMFADelete,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_active": {
				Description: "`true` if MFA is active",
				Type:        schema.TypeBool,
				Required:    true,
			},
		},
	}
}

func resourceAppIDMFACreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)
	c := m.(*appid.AppIDManagementV4)

	input := &appid.UpdateMFAConfigOptions{
		TenantID: &tenantID,
		IsActive: &isActive,
	}

	log.Printf("[DEBUG] Applying AppID MFA configuration: %+v", input)
	_, _, err := c.UpdateMFAConfigWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error updating MFA configuration: %s", err)
	}

	return dataSourceAppIDMFARead(ctx, d, m)
}

func resourceAppIDMFADelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	input := &appid.UpdateMFAConfigOptions{
		TenantID: &tenantID,
		IsActive: getBoolPtr(false),
	}

	log.Printf("[DEBUG] Resetting AppID MFA configuration: %+v", input)
	_, _, err := c.UpdateMFAConfigWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error resetting MFA configuration: %s", err)
	}

	d.SetId("")
	return diags
}
