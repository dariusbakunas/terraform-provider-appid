package appid

import (
	"context"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDAuditStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppIDAuditStatusCreate,
		ReadContext:   resourceAppIDAuditStatusRead, // reusing data source read, same schema
		DeleteContext: resourceAppIDAuditStatusDelete,
		UpdateContext: resourceAppIDAuditStatusUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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

func resourceAppIDAuditStatusRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Id()

	c := m.(*appid.AppIDManagementV4)

	auditStatus, _, err := c.GetAuditStatusWithContext(ctx, &appid.GetAuditStatusOptions{
		TenantID: getStringPtr(tenantID),
	})

	if err != nil {
		return diag.Errorf("error getting audit status: %s", err)
	}

	d.Set("is_active", *auditStatus.IsActive)
	d.Set("tenant_id", tenantID)

	return diags
}

func resourceAppIDAuditStatusCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)
	c := m.(*appid.AppIDManagementV4)

	_, err := c.SetAuditStatusWithContext(ctx, &appid.SetAuditStatusOptions{
		TenantID: getStringPtr(tenantID),
		IsActive: getBoolPtr(isActive),
	})

	if err != nil {
		return diag.Errorf("error setting audit status: %s", err)
	}

	d.SetId(tenantID)
	return resourceAppIDAuditStatusRead(ctx, d, m)
}

func resourceAppIDAuditStatusDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	_, err := c.SetAuditStatusWithContext(ctx, &appid.SetAuditStatusOptions{
		TenantID: getStringPtr(tenantID),
		IsActive: getBoolPtr(false),
	})

	if err != nil {
		return diag.Errorf("error setting audit status: %s", err)
	}

	d.SetId("")
	return diags
}

func resourceAppIDAuditStatusUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceAppIDAuditStatusCreate(ctx, d, m)
}
