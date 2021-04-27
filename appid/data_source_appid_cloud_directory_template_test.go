package appid

import (
	b64 "encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppIDCloudDirectoryTemplateDataSource_basic(t *testing.T) {
	htmlBody := "<HTML><HEAD><TITLE>Test title</TITLE></HEAD><BODY>test</BODY></HTML>"
	b64HTML := b64.StdEncoding.EncodeToString([]byte(htmlBody))
	textBody := "This is the test"
	subject := "Please Verify Your Email Address %%{user.displayName} TEST"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: setupTemplateConfig(testTenantID, subject, htmlBody, textBody),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_cloud_directory_template.test_tpl", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("data.appid_cloud_directory_template.test_tpl", "template_name", "USER_VERIFICATION"),
					resource.TestCheckResourceAttr("data.appid_cloud_directory_template.test_tpl", "subject", strings.Replace(subject, "%%", "%", 1)),
					resource.TestCheckResourceAttr("data.appid_cloud_directory_template.test_tpl", "html_body", htmlBody),
					resource.TestCheckResourceAttr("data.appid_cloud_directory_template.test_tpl", "base64_encoded_html_body", b64HTML),
					resource.TestCheckResourceAttr("data.appid_cloud_directory_template.test_tpl", "plain_text_body", textBody),
				),
			},
		},
	})
}

func setupTemplateConfig(tenantID string, subject string, htmlBody string, textBody string) string {
	return fmt.Sprintf(`
		resource "appid_cloud_directory_template" "test_tpl" {
			tenant_id = "%s"
			template_name = "USER_VERIFICATION"
			subject = "%s"
			html_body = "%s"
			plain_text_body = "%s"
		}

		data "appid_cloud_directory_template" "test_tpl" {
			tenant_id = appid_cloud_directory_template.test_tpl.tenant_id
			template_name = "USER_VERIFICATION"
			depends_on = [appid_cloud_directory_template.test_tpl]
		}
	`, tenantID, subject, htmlBody, textBody)
}
