package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppIDIDPGoogleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: setupGoogleIDPConfig(testTenantID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_idp_google.gg", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("data.appid_idp_google.gg", "config.0.application_id", "test_id"),
					resource.TestCheckResourceAttr("data.appid_idp_google.gg", "config.0.application_secret", "test_secret"),
				),
			},
		},
	})
}

func setupGoogleIDPConfig(tenantID string) string {
	return fmt.Sprintf(`
		resource "appid_idp_google" "gg" {
			tenant_id = "%s"
			is_active = true
			
			config {
				application_id 		= "test_id"
				application_secret 	= "test_secret"
			}
		}

		data "appid_idp_google" "gg" {
			tenant_id = appid_idp_google.gg.tenant_id

			depends_on = [
				appid_idp_google.gg
			]
		}
	`, tenantID)
}
