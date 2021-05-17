package appid

import (
	"context"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDRedirectURLs() *schema.Resource {
	return &schema.Resource{
		Description:   "Redirect URIs that can be used as callbacks of App ID authentication flow",
		CreateContext: resourceAppIDRedirectURLsCreate,
		ReadContext:   dataSourceAppIDRedirectURLsRead,
		UpdateContext: resourceAppIDRedirectURLsUpdate,
		DeleteContext: resourceAppIDRedirectURLsDelete,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service `tenantId`",
			},
			"urls": {
				Description: "A list of redirect URLs",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceAppIDRedirectURLsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	if urls, ok := d.GetOk("urls"); ok {
		redirectURLs := expandStringList(urls.([]interface{}))
		_, err := c.UpdateRedirectUrisWithContext(ctx, &appid.UpdateRedirectUrisOptions{
			TenantID: getStringPtr(tenantID),
			RedirectUrisArray: &appid.RedirectURIConfig{
				RedirectUris: redirectURLs,
			},
		})

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(tenantID)
	return nil
}

func resourceAppIDRedirectURLsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceAppIDRedirectURLsCreate(ctx, d, m)
}

func resourceAppIDRedirectURLsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*appid.AppIDManagementV4)
	tenantID := d.Get("tenant_id").(string)

	_, err := c.UpdateRedirectUrisWithContext(ctx, &appid.UpdateRedirectUrisOptions{
		TenantID: getStringPtr(tenantID),
		RedirectUrisArray: &appid.RedirectURIConfig{
			RedirectUris: []string{},
		},
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
