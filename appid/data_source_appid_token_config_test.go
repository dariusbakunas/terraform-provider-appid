package appid

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
)

func TestAccTokenConfigDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTokenConfigDataSource(testTenantID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_token_config.test_config", "tenant_id", testTenantID),
					resource.TestCheckResourceAttrSet(
						"data.appid_token_config.test_config", "id"),
				),
			},
		},
	})
}

func TestFlattenTokenClaims(t *testing.T) {
	testcases := []struct {
		claims   []api.TokenClaim
		expected []interface{}
	}{
		{
			claims: []api.TokenClaim{
				{Source: "appid_custom", SourceClaim: getStringPtr("sClaim"), DestinationClaim: getStringPtr("dClaim")},
				{Source: "appid_custom", DestinationClaim: getStringPtr("dClaim")},
			},
			expected: []interface{}{
				map[string]interface{}{"source": "appid_custom", "source_claim": "sClaim", "destination_claim": "dClaim"},
				map[string]interface{}{"source": "appid_custom", "destination_claim": "dClaim"},
			},
		},
	}

	for _, c := range testcases {
		actual := flattenTokenClaims(c.claims)
		assert.Equal(t, actual, c.expected)
	}
}

func testAccCheckTokenConfigDataSource(tenantID string) string {
	return fmt.Sprintf(`
data "appid_token_config" "test_config" {
	tenant_id = "%s"
}
	`, tenantID)
}
