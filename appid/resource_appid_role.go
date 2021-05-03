package appid

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
)

func resourceAppIDRole() *schema.Resource {
	return &schema.Resource{
		Description:   "A role is a collection of `scopes` that allow varying permissions to different types of app users",
		CreateContext: resourceAppIDRoleCreate,
		ReadContext:   dataSourceAppIDRoleRead,
		DeleteContext: resourceAppIDRoleDelete,
		UpdateContext: resourceAppIDRoleUpdate,
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

func resourceAppIDRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	input := &api.RoleInput{
		Name: d.Get("name").(string),
	}

	tenantID := d.Get("tenant_id").(string)

	if description, ok := d.GetOk("description"); ok {
		input.Description = description.(string)
	}

	c := m.(*api.Client)

	input.Access = expandRoleAccess(d.Get("access").(*schema.Set).List())

	log.Printf("[DEBUG] Creating AppID role: %+v", input)
	role, err := c.RolesAPI.CreateRole(ctx, tenantID, input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(role.ID)
	d.Set("role_id", role.ID)

	return dataSourceAppIDRoleRead(ctx, d, m)
}

func resourceAppIDRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)

	roleID := d.Id()
	tenantID := d.Get("tenant_id").(string)

	log.Printf("[DEBUG] Deleting AppID role: %s", roleID)

	err := c.RolesAPI.DeleteRole(ctx, tenantID, roleID)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Finished deleting AppID role: %s", d.Id())

	return diags
}

func expandRoleAccess(l []interface{}) []api.RoleAccess {
	if len(l) == 0 {
		return []api.RoleAccess{}
	}

	result := make([]api.RoleAccess, len(l))

	for i, item := range l {
		aMap := item.(map[string]interface{})

		access := &api.RoleAccess{
			ApplicationID: aMap["application_id"].(string),
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
	// TODO: implement role update
	return dataSourceAppIDRoleRead(ctx, d, m)
}
