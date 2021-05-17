terraform {
  required_providers {
    appid = {
      source = "dariusbakunas/appid"
      version = "0.2.0"
    } 
  }
}

provider "appid" {  
    iam_access_token = var.iam_access_token      
    region = "us-south"
}

resource "appid_application" "app" {
  tenant_id = var.tenant_id
  name = "test-tf-application"  
  type = "singlepageapp"
  scopes = ["test_scope_1", "test_scope_2", "test_scope_3"]
}

data "appid_application" "app" {
  tenant_id = var.tenant_id
  client_id = appid_application.app.client_id
}

output "application" {
  value = data.appid_application.app
}
