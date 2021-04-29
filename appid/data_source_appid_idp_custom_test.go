package appid

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppIDIDPCustomDataSource_basic(t *testing.T) {
	publicKey := `-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzb19EC2vJfTLaJTs3/4F\ndmoHnpYHJo4Q5SJYJK2YfclwRJc49zs1juoNGvXsUOsEi58PHarot3aAUpzBk8g9\n1RdDoovQDKBhMbT7BXP291qp5WQsvrv5W6xPoTbNONYPmAWTN75e3AvvvQElgv9N\n4BBkXZ962bf/OM1Ccm786laop9fC03D7vmUUypISPMZ61O6aA3dRI2JSvHh+VL4s\nEtXkZvLR7DvvWl4sl4oA5EvpYqw5/qbXTp4bnllfiQuCuwgYz/MH1mQA4qGWEVTN\nE4z3b0jsHNHVAzsPfB3Bnok/Zvgtxc3cjVlm3el+bie9O3vW1jFQf1JCke/qusj7\neQIDAQAB\n-----END PUBLIC KEY-----\n`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: setupCustomIDPConfig(testTenantID, publicKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_idp_custom.idp", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("data.appid_idp_custom.idp", "is_active", "true"),
					resource.TestCheckResourceAttr("data.appid_idp_custom.idp", "public_key", strings.Replace(publicKey, "\\n", "\n", -1)),
				),
			},
		},
	})
}

func setupCustomIDPConfig(tenantID string, publicKey string) string {
	return fmt.Sprintf(`
	resource "appid_idp_custom" "idp" {
		tenant_id = "%s"
		is_active = true
		public_key = "%s"
	}

	data "appid_idp_custom"  "idp" {
		tenant_id = appid_idp_custom.idp.tenant_id

		depends_on = [
			appid_idp_custom.idp
		]
	}
	`, tenantID, publicKey)
}
