package appid

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAppIDActionURL() *schema.Resource {
	return &schema.Resource{
		Description: "The custom url to redirect to when Cloud Directory action is executed.",
		ReadContext: dataSourceAppIDActionURLRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"action": {
				Description:  "The type of the action: `on_user_verified` - the URL of your custom user verified page, `on_reset_password` - the URL of your custom reset password page",
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"on_user_verified", "on_reset_password"}, false),
				Required:     true,
			},
			"url": {
				Description: "The action URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAppIDActionURLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	action := d.Get("action").(string)

	c := m.(*appid.AppIDManagementV4)

	resp, _, err := c.GetCloudDirectoryActionURLWithContext(ctx, &appid.GetCloudDirectoryActionURLOptions{
		TenantID: getStringPtr(tenantID),
		Action:   getStringPtr(action),
	})

	if err != nil {
		return diag.Errorf("Error getting actionURL: %s", err)
	}

	if resp.ActionURL != nil {
		d.Set("url", *resp.ActionURL)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, action))

	return diags
}
