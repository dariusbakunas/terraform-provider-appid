package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTokenConfigDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTokenConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_token_config.test_config", "tenant_id", tenantID),
					resource.TestCheckResourceAttrSet(
						"data.appid_token_config.test_config", "id"),
				),
			},
		},
	})
}

func testAccCheckTokenConfigDataSource() string {
	return fmt.Sprintf(`
data "appid_token_config" "test_config" {
	tenant_id = "%s"
}
	`, tenantID)
}
