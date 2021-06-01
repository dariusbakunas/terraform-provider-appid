package appid

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppIDLanguagesDataSource_basic(t *testing.T) {
	languages := []string{"en", "es-ES", "fr-FR"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: setupLanguagesConfig(testTenantID, languages),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.appid_languages.lang", "tenant_id", testTenantID),
					resource.TestCheckResourceAttr("data.appid_languages.lang", "languages.#", strconv.Itoa(len(languages))),
				),
			},
		},
	})
}

func setupLanguagesConfig(tenantID string, languages []string) string {
	langs := strings.Replace(fmt.Sprintf("%q", languages), " ", ", ", -1)

	return fmt.Sprintf(`
		resource "appid_languages" "lang" {
			tenant_id = "%s"
			languages = %s
		}
		data "appid_languages" "lang" {
			tenant_id = appid_languages.lang.tenant_id

			depends_on = [
				appid_languages.lang
			]
		}
	`, tenantID, langs)
}
