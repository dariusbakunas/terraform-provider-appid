package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppIDPasswordRefexDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: setupPasswordRegexConfig(testTenantID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_password_regex.rgx", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("data.appid_password_regex.rgx", "regex", "^(?:(?=.*\\d)(?=.*[a-z])(?=.*[A-Z]).*)$"),
					resource.TestCheckResourceAttr("data.appid_password_regex.rgx", "error_message", "test error"),
				),
			},
		},
	})
}

func setupPasswordRegexConfig(tenantID string) string {
	return fmt.Sprintf(`
		resource "appid_password_regex" "rgx" {
			tenant_id = "%s"
			regex = "^(?:(?=.*\\d)(?=.*[a-z])(?=.*[A-Z]).*)$"
			error_message = "test error"
		}

		data "appid_password_regex" "rgx" {
			tenant_id = appid_password_regex.rgx.tenant_id

			depends_on = [
				appid_password_regex.rgx
			]
		}
	`, tenantID)
}
