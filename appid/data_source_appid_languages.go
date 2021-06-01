package appid

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppIDLanguages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppIDLanguagesRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"languages": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceAppIDLanguagesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	langs, _, err := c.GetLocalizationWithContext(ctx, &appid.GetLocalizationOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting languages: %s", err)
	}

	d.Set("languages", langs.Languages)
	d.SetId(fmt.Sprintf("%s/languages", tenantID))

	return diags
}
