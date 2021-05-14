package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppIDIDPCloudDirectoryDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: setupCloudDirectoryIDPConfig(testTenantID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_idp_cloud_directory.idp", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("data.appid_idp_cloud_directory.idp", "is_active", "true"),
					resource.TestCheckResourceAttr("data.appid_idp_cloud_directory.idp", "self_service_enabled", "false"),
					resource.TestCheckResourceAttr("data.appid_idp_cloud_directory.idp", "signup_enabled", "false"),
					resource.TestCheckResourceAttr("data.appid_idp_cloud_directory.idp", "welcome_enabled", "true"),
					resource.TestCheckResourceAttr("data.appid_idp_cloud_directory.idp", "reset_password_enabled", "false"),
					resource.TestCheckResourceAttr("data.appid_idp_cloud_directory.idp", "reset_password_notification_enabled", "false"),
					resource.TestCheckResourceAttr("data.appid_idp_cloud_directory.idp", "identity_confirm_access_mode", "FULL"),
				),
			},
		},
	})
}

func setupCloudDirectoryIDPConfig(tenantID string) string {
	return fmt.Sprintf(`
		resource "appid_idp_cloud_directory" "idp" {
			tenant_id = "%s"
			is_active = true
			self_service_enabled = false
			signup_enabled = false
			welcome_enabled = true
			reset_password_enabled = false
			reset_password_notification_enabled = false			
		}
		data "appid_idp_cloud_directory" "idp" {
			tenant_id = appid_idp_cloud_directory.idp.tenant_id
			depends_on = [
				appid_idp_cloud_directory.idp
			]
		}
	`, tenantID)
}
