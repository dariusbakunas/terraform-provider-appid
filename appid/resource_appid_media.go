package appid

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
)

func resourceAppIDMedia() *schema.Resource {
	return &schema.Resource{
		Description:   "Custom logo image of the login widget. *Note:* Currently there is no supported way of deleting the image",
		ReadContext:   resourceAppIDMediaRead,
		CreateContext: resourceAppIDMediaCreate,
		UpdateContext: resourceAppIDMediaCreate,
		DeleteContext: resourceAppIDMediaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"logo_url": {
				Description: "AppID Login logo URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"source": {
				Description: "Path to logo image",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceAppIDMediaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	tenantID := d.Id()
	c := m.(*appid.AppIDManagementV4)

	media, _, err := c.GetMediaWithContext(ctx, &appid.GetMediaOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID media: %s", err)
	}

	if media.Image != nil {
		d.Set("logo_url", *media.Image)
	}

	d.Set("tenant_id", tenantID)
	return diags
}

func resourceAppIDMediaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	source := d.Get("source").(string)
	c := m.(*appid.AppIDManagementV4)

	path, err := homedir.Expand(source)

	if err != nil {
		return diag.Errorf("Error parsing source: %s", err)
	}

	file, err := os.Open(path)

	if err != nil {
		return diag.Errorf("Error opening AppID media source file: %s", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("[WARN] Error closing AppID media source file: %s", err)
		}
	}()

	fileContentType, err := detectFileContentType(file)

	if err != nil {
		return diag.Errorf("Error detecting source content type: %s", err)
	}

	if fileContentType != "image/png" && fileContentType != "image/jpg" {
		return diag.Errorf("Only PNG and JPG images are supported, detected type: %s", fileContentType)
	}

	log.Printf("[DEBUG] Uploading AppID media: %s", source)
	_, err = c.PostMediaWithContext(ctx, &appid.PostMediaOptions{
		TenantID:        &tenantID,
		MediaType:       getStringPtr("logo"),
		File:            file,
		FileContentType: &fileContentType,
	})

	if err != nil {
		return diag.Errorf("Error uploading AppID media: %s", err)
	}

	d.SetId(tenantID)

	return resourceAppIDMediaRead(ctx, d, m)
}

func resourceAppIDMediaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Printf("[WARN] AppID Management APIs does not yet support media deletion, skipping")
	d.SetId("")
	return diags
}

func detectFileContentType(r io.ReadSeeker) (string, error) {
	var fileHeader [512]byte

	_, err := io.ReadFull(r, fileHeader[:])

	if err != nil {
		return "", err
	}

	_, err = r.Seek(0, io.SeekStart)

	if err != nil {
		return "", err
	}

	fileContentType := http.DetectContentType(fileHeader[:])
	return fileContentType, nil
}
