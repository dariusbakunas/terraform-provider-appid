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