package appid

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDRole() *schema.Resource {
	return &schema.Resource{
		Description: "A role is a collection of `scopes` that allow varying permissions to different types of app users",
		ReadContext: dataSourceAppIDRoleRead,
		Schema: map[string]*schema.Schema{
			"role_id": {
				Description: "Role ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service `tenantId`",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique role name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional role description",
			},
			"access": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scopes": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAppIDRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*appid.AppIDManagementV4)

	tenantID := d.Get("tenant_id").(string)
	id := d.Get("role_id").(string)

	role, _, err := c.GetRoleWithContext(ctx, &appid.GetRoleOptions{
		RoleID:   getStringPtr(id),
		TenantID: getStringPtr(tenantID),
	})

	if err != nil {
		return diag.Errorf("Error loading AppID role: %s", err)
	}

	d.Set("name", *role.Name)

	if role.Description != nil {
		d.Set("description", *role.Description)
	}

	d.Set("access", flattenRoleAccess(role.Access))

	d.SetId(fmt.Sprintf("%s/%s", tenantID, *role.ID))

	return diags
}

func flattenRoleAccess(ra []appid.RoleAccessItem) []interface{} {
	var result []interface{}

	for _, a := range ra {
		access := map[string]interface{}{
			"scopes": flattenStringList(a.Scopes),
		}

		if a.ApplicationID != nil {
			access["application_id"] = *a.ApplicationID
		}

		result = append(result, access)
	}

	return result
}
