package appid

import (
	"context"
	"fmt"
	"log"
	"strings"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDRole() *schema.Resource {
	return &schema.Resource{
		Description:   "A role is a collection of `scopes` that allow varying permissions to different types of app users",
		CreateContext: resourceAppIDRoleCreate,
		ReadContext:   resourceAppIDRoleRead,
		DeleteContext: resourceAppIDRoleDelete,
		UpdateContext: resourceAppIDRoleUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Role ID",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service `tenantId`",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique role name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional role description",
			},
			"access": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Application `client_id`",
						},
						"scopes": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAppIDRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*appid.AppIDManagementV4)

	id := d.Id()
	idParts := strings.Split(id, "/")

	tenantID := idParts[0]
	roleID := idParts[1]

	role, _, err := c.GetRoleWithContext(ctx, &appid.GetRoleOptions{
		RoleID:   &roleID,
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error loading AppID role: %s", err)
	}

	d.Set("name", *role.Name)

	if role.Description != nil {
		d.Set("description", *role.Description)
	}

	d.Set("access", flattenRoleAccess(role.Access))

	d.Set("tenant_id", tenantID)
	d.Set("role_id", roleID)

	return diags
}

func resourceAppIDRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)

	input := &appid.CreateRoleOptions{
		Name:     getStringPtr(d.Get("name").(string)),
		TenantID: getStringPtr(tenantID),
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = getStringPtr(description.(string))
	}

	c := m.(*appid.AppIDManagementV4)

	input.Access = expandRoleAccess(d.Get("access").(*schema.Set).List())

	log.Printf("[DEBUG] Creating AppID role: %+v", input)
	role, _, err := c.CreateRoleWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error creating Cloud Directory role: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, *role.ID))

	return resourceAppIDRoleRead(ctx, d, m)
}

func resourceAppIDRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*appid.AppIDManagementV4)

	roleID := d.Id()
	tenantID := d.Get("tenant_id").(string)

	log.Printf("[DEBUG] Deleting AppID role: %s", roleID)

	_, err := c.DeleteRoleWithContext(ctx, &appid.DeleteRoleOptions{
		TenantID: getStringPtr(tenantID),
		RoleID:   getStringPtr(roleID),
	})

	if err != nil {
		return diag.Errorf("Error deleting Cloud Directory role: %s", err)
	}

	d.SetId("")

	log.Printf("[DEBUG] Finished deleting AppID role: %s", d.Id())

	return diags
}

func expandRoleAccess(l []interface{}) []appid.RoleAccessItem {
	if len(l) == 0 {
		return []appid.RoleAccessItem{}
	}

	result := make([]appid.RoleAccessItem, len(l))

	for i, item := range l {
		aMap := item.(map[string]interface{})

		access := &appid.RoleAccessItem{
			ApplicationID: getStringPtr(aMap["application_id"].(string)),
		}

		if scopes, ok := aMap["scopes"].([]interface{}); ok && len(scopes) > 0 {
			access.Scopes = []string{}

			for _, s := range scopes {
				access.Scopes = append(access.Scopes, s.(string))
			}
		}

		result[i] = *access
	}

	return result
}

func resourceAppIDRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	roleID := d.Id()

	input := &appid.UpdateRoleOptions{
		TenantID: getStringPtr(tenantID),
		RoleID:   getStringPtr(roleID),
		Name:     getStringPtr(d.Get("name").(string)),
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = getStringPtr(description.(string))
	}

	input.Access = expandRoleAccess(d.Get("access").(*schema.Set).List())

	c := m.(*appid.AppIDManagementV4)

	log.Printf("[DEBUG] Updating AppID role: %+v", input)
	_, _, err := c.UpdateRoleWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error updating Cloud Directory role: %s", err)
	}

	return dataSourceAppIDRoleRead(ctx, d, m)
}
