package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigTokensDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckConfigTokensDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_config_tokens.test_config", "tenant_id", tenantID),
					resource.TestCheckResourceAttrSet(
						"data.appid_config_tokens.test_config", "id"),
				),
			},
		},
	})
}

func testAccCheckConfigTokensDataSource() string {
	return fmt.Sprintf(`
data "appid_config_tokens" "test_config" {
	tenant_id = "%s"
}
	`, tenantID)
}
