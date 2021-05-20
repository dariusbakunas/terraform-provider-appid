package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppIDActionURLDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: setupActionURLConfig(testTenantID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_action_url.url", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("data.appid_action_url.url", "action", "on_user_verified"),
					resource.TestCheckResourceAttr("data.appid_action_url.url", "url", "https://www.example.com/?user=verified"),
				),
			},
		},
	})
}

func setupActionURLConfig(tenantID string) string {
	return fmt.Sprintf(`
		resource "appid_action_url" "url" {
			tenant_id = "%s"
			action = "on_user_verified"
			url = "https://www.example.com/?user=verified"
		}

		data "appid_action_url" "url" {
			tenant_id = appid_action_url.url
			action = "on_user_verified"

			depends_on = [
				appid_action_url.url
			]
		}
	`, tenantID)
}
