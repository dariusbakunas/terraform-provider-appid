terraform {
  required_providers {
    appid = {
      version = "0.1"
      source  = "us.ibm.com/watson-health/appid"
    }
  }
}

provider "appid" {  
    iam_access_token = var.iam_access_token  
    iam_base_url = "https://iam.cloud.ibm.com"
    appid_base_url = "https://us-south.appid.cloud.ibm.com"
}

resource "appid_application" "app" {
  tenant_id = var.tenant_id
  name = "test-tf-application"  
  type = "singlepageapp"
  scopes = ["test_scope_1", "test_scope_2", "test_scope_3"]
}

output "application" {
  value = appid_application.app
}