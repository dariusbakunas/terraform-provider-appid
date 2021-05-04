package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppIDRole_basic(t *testing.T) {
	roleName := fmt.Sprintf("%s_role_%d", testResourcePrefix, acctest.RandIntRange(10, 100))
	appName := fmt.Sprintf("%s_app_%d", testResourcePrefix, acctest.RandIntRange(10, 100))
	description := "test role"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: setupRoleConfig(testTenantID, appName, roleName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_role.role", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("data.appid_role.role", "name", roleName),
					resource.TestCheckResourceAttr("data.appid_role.role", "description", description),
				),
			},
		},
	})
}

func setupRoleConfig(tenantID string, appName string, roleName string, description string) string {
	return fmt.Sprintf(`
		resource "appid_application" "app" {
			tenant_id = "%s"
			name = "%s"  
			type = "singlepageapp"
			scopes = ["pancakes", "cartoons"]
	  	}

		resource "appid_role" "role" {
			tenant_id = appid_application.app.tenant_id
			name = "%s"
			description = "%s"
			access {
				application_id = appid_application.app.client_id
				scopes = [
					"pancakes",
				]
			}
		}

		data "appid_role" "role" {
			tenant_id = appid_role.role.tenant_id
			role_id = appid_role.role.role_id
		}
	`, tenantID, appName, roleName, description)
}
