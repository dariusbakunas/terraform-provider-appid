package appid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
)

func resourceAppIDAuditStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDAuditStatusCreate,
		ReadContext:   dataSourceAppIDAuditStatusRead, // reusing data source read, same schema
		DeleteContext: resourceAppIDAuditStatusDelete,
		UpdateContext: resourceAppIDAuditStatusUpdate,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service `tenantId`",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "The auditing status of the tenant.",
			},
		},
	}
}

func resourceAppIDAuditStatusCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)
	c := m.(*api.Client)

	err := c.ConfigAPI.SetAuditStatuts(ctx, tenantID, &api.AuditStatus{IsActive: isActive})

	if err != nil {
		return diag.Errorf("error setting audit status: %s", err)
	}

	return dataSourceAppIDAuditStatusRead(ctx, d, m)
}

func resourceAppIDAuditStatusDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	tenantID := d.Get("tenant_id").(string)
	c := m.(*api.Client)

	err := c.ConfigAPI.SetAuditStatuts(ctx, tenantID, &api.AuditStatus{IsActive: false})

	if err != nil {
		return diag.Errorf("error setting audit status: %s", err)
	}

	d.SetId("")
	return diags
}

func resourceAppIDAuditStatusUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceAppIDAuditStatusCreate(ctx, d, m)
}
