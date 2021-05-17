package appid

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const testResourcePrefix = "tf_acc_test"

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var testTenantID string
var appIDBaseURL string
var iamBaseURL string

func init() {
	testTenantID = os.Getenv("APPID_TENANT_ID")
	if testTenantID == "" {
		testTenantID = "24565a1c-2dac-409b-a60c-0ff130c6943c"
		log.Printf("[INFO] Set the environment variable APPID_TENANT_ID for testing AppID resources else it is set to default '%s'", testTenantID)
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
	apiKey := os.Getenv("IAM_API_KEY")
	accessToken := os.Getenv("IAM_ACCESS_TOKEN")

	if apiKey == "" && accessToken == "" {
		t.Fatal("IAM_API_KEY or IAM_ACCESS_TOKEN env must be set for acceptance tests")
	}
}
