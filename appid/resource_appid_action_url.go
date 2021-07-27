package appid

import (
	"context"
	"fmt"
	"log"
	"strings"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppIDActionURL() *schema.Resource {
	return &schema.Resource{
		Description:   "The custom url to redirect to when Cloud Directory action is executed.",
		CreateContext: resourceAppIDActionURLCreate,
		ReadContext:   resourceAppIDActionURLRead,
		DeleteContext: resourceAppIDActionURLDelete,
		UpdateContext: resourceAppIDActionURLUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"action": {
				Description:  "The type of the action: `on_user_verified` - the URL of your custom user verified page, `on_reset_password` - the URL of your custom reset password page",
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"on_user_verified", "on_reset_password"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"url": {
				Description: "The action URL",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceAppIDActionURLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := d.Id()
	idParts := strings.Split(id, "/")

	tenantID := idParts[0]
	action := idParts[1]

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

	d.Set("tenant_id", tenantID)
	d.Set("action", action)

	return diags
}

func resourceAppIDActionURLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	action := d.Get("action").(string)
	actionURL := d.Get("url").(string)

	c := m.(*appid.AppIDManagementV4)

	input := &appid.SetCloudDirectoryActionOptions{
		TenantID:  getStringPtr(tenantID),
		Action:    getStringPtr(action),
		ActionURL: getStringPtr(actionURL),
	}

	log.Printf("[DEBUG] Setting Cloud Directory action URL: %+v", input)

	_, _, err := c.SetCloudDirectoryActionWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error setting Cloud Directory action URL: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, action))

	return resourceAppIDActionURLRead(ctx, d, m)
}

func resourceAppIDActionURLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*appid.AppIDManagementV4)
	tenantID := d.Get("tenant_id").(string)
	action := d.Get("action").(string)

	log.Printf("[DEBUG] Deleting Cloud Directory action URL: %s", d.Id())

	_, err := c.DeleteActionURLWithContext(ctx, &appid.DeleteActionURLOptions{
		TenantID: getStringPtr(tenantID),
		Action:   getStringPtr(action),
	})

	if err != nil {
		return diag.Errorf("Error deleting Cloud Directory action URL: %s", err)
	}

	d.SetId("")

	return diags
}

func resourceAppIDActionURLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceAppIDActionURLCreate(ctx, d, m)
}
