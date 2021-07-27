package appid

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDRoles() *schema.Resource {
	return &schema.Resource{
		Description: "A list of AppID roles",
		ReadContext: dataSourceAppIDRolesRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Description: "Role ID",
							Type:        schema.TypeString,
							Computed:    true,
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
				},
			},
		},
	}
}

func dataSourceAppIDRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)

	c := m.(*appid.AppIDManagementV4)
	roles, _, err := c.ListRolesWithContext(ctx, &appid.ListRolesOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error listing AppID roles: %s", err)
	}

	roleList := make([]interface{}, 0)

	for _, role := range roles.Roles {
		rMap := map[string]interface{}{}
		rMap["role_id"] = *role.ID
		rMap["name"] = *role.Name

		if role.Description != nil {
			rMap["description"] = *role.Description
		}

		rMap["access"] = flattenRoleAccess(role.Access)
		roleList = append(roleList, rMap)
	}

	if err := d.Set("roles", roleList); err != nil {
		return diag.Errorf("Error setting roles: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/roles", tenantID))
	return nil
}
