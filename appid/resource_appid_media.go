package appid

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
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
				Description:   "Path to logo image",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source_content"},
			},
			"source_content": {
				Description:   "base64 encoded logo image contents",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
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
	var contentFile io.ReadSeekCloser
	var contentType string

	tenantID := d.Get("tenant_id").(string)
	c := m.(*appid.AppIDManagementV4)

	source, sourceOK := d.GetOk("source")
	sourcePath := source.(string)
	content, contentOK := d.GetOk("source_content")

	if !sourceOK && !contentOK {
		return diag.Errorf("AppID media `source` or `source_content` must be specified")
	}

	if contentOK {
		decoded, err := base64.StdEncoding.DecodeString(content.(string))

		if err != nil {
			// if this is not base64 encoded string use it as is
			decoded = []byte(content.(string))
		}

		contentFile, err := os.CreateTemp(os.TempDir(), "appid-media-")
		if err != nil {
			return diag.Errorf("Error creating AppID media temporary file: %s", err)
		}

		defer func() {
			err = contentFile.Close()
			if err != nil {
				log.Printf("[WARN] Error closing AppID media temporary file: %s", err)
			}

			err = os.Remove(contentFile.Name())
			if err != nil {
				log.Printf("[WARN] Error removing temporary AppID media source file: %s", err)
			}
		}()

		if _, err = contentFile.Write(decoded); err != nil {
			return diag.Errorf("Error writing source content to temporary file: %s", err)
		}

		sourcePath = contentFile.Name()
	}

	path, err := homedir.Expand(sourcePath)

	if err != nil {
		return diag.Errorf("Error parsing source: %s", err)
	}

	contentFile, err = os.Open(path)

	if err != nil {
		return diag.Errorf("Error opening AppID media source file: %s", err)
	}

	defer func() {
		err := contentFile.Close()
		if err != nil {
			log.Printf("[WARN] Error closing AppID media source file: %s", err)
		}
	}()

	contentType, err = detectFileContentType(contentFile)

	if err != nil {
		return diag.Errorf("Error detecting source content type: %s", err)
	}

	if contentType != "image/png" && contentType != "image/jpg" {
		return diag.Errorf("Only PNG and JPG images are supported, detected type: %s", contentType)
	}

	_, err = c.PostMediaWithContext(ctx, &appid.PostMediaOptions{
		TenantID:        &tenantID,
		MediaType:       getStringPtr("logo"),
		File:            contentFile,
		FileContentType: &contentType,
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
