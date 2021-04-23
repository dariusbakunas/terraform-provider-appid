package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppIDApplicationDataSource_basic(t *testing.T) {
	appName := fmt.Sprintf("%s_app_%d", testResourcePrefix, acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: setupApplicationConfig(testTenantID, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_application.test_app", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("data.appid_application.test_app", "name", appName),
					resource.TestCheckResourceAttr("data.appid_application.test_app", "type", "singlepageapp"),
					resource.TestCheckResourceAttrSet("data.appid_application.test_app", "client_id"),
					resource.TestCheckResourceAttr("data.appid_application.test_app", "scopes.#", "3"),
					resource.TestCheckResourceAttr("data.appid_application.test_app", "scopes.0", "test_scope_1"),
					resource.TestCheckResourceAttr("data.appid_application.test_app", "scopes.1", "test_scope_2"),
					resource.TestCheckResourceAttr("data.appid_application.test_app", "scopes.2", "test_scope_3"),
				),
			},
		},
	})
}

func setupApplicationConfig(tenantID string, name string) string {
	return fmt.Sprintf(`
		resource "appid_application" "test_app" {
			tenant_id = "%s"
			name = "%s"  
			type = "singlepageapp"
			scopes = ["test_scope_1", "test_scope_2", "test_scope_3"]
		}
		data "appid_application" "test_app" {
			tenant_id = "%s"
			client_id = appid_application.test_app.client_id
		}
	`, tenantID, name, tenantID)
}
