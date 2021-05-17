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

data "appid_application" "app" {
  tenant_id = var.tenant_id
  client_id = var.client_id
}

output "application" {
  value = data.appid_application.app
}