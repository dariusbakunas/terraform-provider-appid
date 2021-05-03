package appid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
)

func dataSourceAppIDRole() *schema.Resource {
	return &schema.Resource{
		Description: "A role is a group of scopes that apply to the user.",
		ReadContext: dataSourceAppIDRoleRead,
		Schema: map[string]*schema.Schema{
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

	c := m.(*api.Client)

	tenantID := d.Get("tenant_id").(string)
	id := d.Id()

	role, err := c.RolesAPI.GetRole(ctx, tenantID, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(role.ID)
	d.Set("id", role.ID)
	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("access", flattenRoleAccess(role.Access))

	return diags
}

func flattenRoleAccess(ra []api.RoleAccess) []interface{} {
	var result []interface{}

	for _, a := range ra {
		access := map[string]interface{}{
			"application_id": a.ApplicationID,
			"scopes":         flattenStringList(a.Scopes),
		}

		result = append(result, access)
	}

	return result
}
