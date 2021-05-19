package appid

import (
	"context"
	b64 "encoding/base64"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppIDPasswordRegex() *schema.Resource {
	return &schema.Resource{
		Description:   "The regular expression used by App ID for password strength validation",
		CreateContext: resourceAppIDPasswordRegexCreate,
		ReadContext:   dataSourceAppIDPasswordRegexRead,
		DeleteContext: resourceAppIDPasswordRegexDelete,
		UpdateContext: resourceAppIDPasswordRegexUpdate,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"base64_encoded_regex": {
				Description: "The regex expression rule for acceptable password encoded in base64",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"error_message": {
				Description: "Custom error message",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"regex": {
				Description: "The escaped regex expression rule for acceptable password",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceAppIDPasswordRegexCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	regex := d.Get("regex").(string)

	c := m.(*appid.AppIDManagementV4)

	input := &appid.SetCloudDirectoryPasswordRegexOptions{
		TenantID:           getStringPtr(tenantID),
		Base64EncodedRegex: getStringPtr(b64.StdEncoding.EncodeToString([]byte(regex))),
	}

	if msg, ok := d.GetOk("error_message"); ok {
		input.ErrorMessage = getStringPtr(msg.(string))
	}

	_, _, err := c.SetCloudDirectoryPasswordRegexWithContext(ctx, input)

	if err != nil {
		return diag.FromErr(err)
	}

	return dataSourceAppIDPasswordRegexRead(ctx, d, m)
}

func resourceAppIDPasswordRegexUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceAppIDPasswordRegexCreate(ctx, d, m)
}

func resourceAppIDPasswordRegexDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	input := &appid.SetCloudDirectoryPasswordRegexOptions{
		TenantID:           getStringPtr(tenantID),
		Base64EncodedRegex: getStringPtr(""),
	}

	_, _, err := c.SetCloudDirectoryPasswordRegexWithContext(ctx, input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
