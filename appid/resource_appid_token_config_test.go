package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppIDTokenConfig_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppIDTokenConfigCreate(testTenantID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("appid_token_config.test_config", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "access_token_expires_in", "7200"),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "anonymous_access_enabled", "false"),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "anonymous_token_expires_in", "7200"),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "refresh_token_enabled", "true"),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "refresh_token_expires_in", "7200"),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "access_token_claim.#", "2"),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "id_token_claim.#", "0"),
					// the order here is deterministic: https://github.com/hashicorp/terraform-plugin-sdk/blob/main/helper/schema/set.go#L268
					resource.TestCheckResourceAttr("appid_token_config.test_config", "access_token_claim.0.destination_claim", "employeeId"),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "access_token_claim.0.source", "appid_custom"),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "access_token_claim.0.source_claim", "employeeId"),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "access_token_claim.1.destination_claim", "groupIds"),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "access_token_claim.1.source", "roles"),
					resource.TestCheckResourceAttr("appid_token_config.test_config", "access_token_claim.1.source_claim", ""),
				),
			},
		},
	})
}

func testAccCheckAppIDTokenConfigCreate(tenantID string) string {
	return fmt.Sprintf(`
		resource "appid_token_config" "test_config" {
			tenant_id = "%s"
			access_token_expires_in = 7200    
			anonymous_access_enabled = false
			anonymous_token_expires_in = 7200
			refresh_token_enabled = true
			refresh_token_expires_in = 7200
			
			access_token_claim {
				source = "roles"
				destination_claim = "groupIds"
			}

			access_token_claim {
				source = "appid_custom"
				source_claim = "employeeId"
				destination_claim = "employeeId"
			}			
		}
	`, tenantID)
}
