package appid

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const testResourcePrefix = "tf-acc-test-"

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var tenantID string
var appIDBaseURL string
var iamBaseURL string

func init() {
	tenantID = os.Getenv("APPID_TENANT_ID")
	if tenantID == "" {
		tenantID = "bac33a56-501f-493c-8b1e-bfda921f4a3e"
		log.Printf("[INFO] Set the environment variable APPID_TENANT_ID for testing AppID resources else it is set to default '%s'", tenantID)
	}

	appIDBaseURL = os.Getenv("APPID_BASE_URL")
	if appIDBaseURL == "" {
		appIDBaseURL = "https://us-south.appid.cloud.ibm.com"
		log.Printf("[INFO] Set the environment variable APPID_BASE_URL for testing AppID resources else it is set to default '%s'", appIDBaseURL)
	}

	iamBaseURL = os.Getenv("IAM_BASE_URL")
	if iamBaseURL == "" {
		iamBaseURL = "https://iam.cloud.ibm.com"
		log.Printf("[INFO] Set the environment variable IAM_BASE_URL for testing AppID resources else it is set to default '%s'", iamBaseURL)
	}
}

func init() {
	testAccProvider = Provider()

	testAccProviders = map[string]*schema.Provider{
		"appid": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("IAM_API_KEY"); v == "" {
		t.Fatal("IAM_API_KEY must be set for acceptance tests")
	}
}
