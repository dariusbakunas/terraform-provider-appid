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
  scopes = ["pancakes", "cartoons"]
}

resource "appid_role" "role" {
  tenant_id = var.tenant_id
  name = "test-tf-role"    
  access {
    application_id = appid_application.app.client_id
    scopes = [      
      "pancakes"
    ]
  }
}


data "appid_role" "role" {
  tenant_id = appid_role.role.tenant_id
  role_id = appid_role.role.role_id

  depends_on = [
    appid_role.role
  ]
}

output "role" {
  value = data.appid_role.role
}