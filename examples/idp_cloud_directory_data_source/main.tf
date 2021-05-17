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

data "appid_idp_cloud_directory" "cd" {
  tenant_id = var.tenant_id  
}

output "cd" {
  value = data.appid_idp_cloud_directory.cd
}