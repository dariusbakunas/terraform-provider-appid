package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccThemeColorDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: setupThemeColorConfig(testTenantID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_theme_color.color", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("data.appid_theme_color.color", "header_color", "#000000"),
				),
			},
		},
	})
}

func setupThemeColorConfig(tenantID string) string {
	return fmt.Sprintf(`
		resource "appid_theme_color" "color" {
			tenant_id = "%s"
			header_color = "#000000"
		}

		data "appid_theme_color" "color" {
			tenant_id = appid_theme_color.color.tenant_id

			depends_on = [
				appid_theme_color.color
			]
		}
	`, tenantID)
}
