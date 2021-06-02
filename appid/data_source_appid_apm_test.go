package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppIDAPMDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: setupAPMConfig(testTenantID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_apm.apm", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("data.appid_apm.apm", "enabled", "true"),
					resource.TestCheckResourceAttr("data.appid_apm.apm", "prevent_password_with_username", "true"),
					resource.TestCheckResourceAttr("data.appid_apm.apm", "password_reuse.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.appid_apm.apm", "password_reuse.0.max_password_reuse", "4"),
					resource.TestCheckResourceAttr("data.appid_apm.apm", "password_expiration.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.appid_apm.apm", "password_expiration.0.days_to_expire", "25"),
					resource.TestCheckResourceAttr("data.appid_apm.apm", "lockout_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.appid_apm.apm", "lockout_policy.0.lockout_time_sec", "2600"),
					resource.TestCheckResourceAttr("data.appid_apm.apm", "lockout_policy.0.num_of_attempts", "4"),
					resource.TestCheckResourceAttr("data.appid_apm.apm", "min_password_change_interval.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.appid_apm.apm", "min_password_change_interval.0.min_hours_to_change_password", "1"),
				),
			},
		},
	})
}

func setupAPMConfig(tenantID string) string {
	return fmt.Sprintf(`
		resource "appid_apm" "apm" {
			tenant_id = "%s"
			enabled = true
			prevent_password_with_username = true
		
			password_reuse {
				enabled = true
				max_password_reuse = 4
			}
		
			password_expiration {
				enabled = true
				days_to_expire = 25
			}
		
			lockout_policy {
				enabled = true
				lockout_time_sec = 2600
				num_of_attempts = 4
			}
		
			min_password_change_interval {
				enabled = true
				min_hours_to_change_password = 1
			}
		}

		data "appid_apm" "apm" {
			tenant_id = appid_apm.apm.tenant_id

			depends_on = [
				appid_apm.apm
			]
		}
	`, tenantID)
}
