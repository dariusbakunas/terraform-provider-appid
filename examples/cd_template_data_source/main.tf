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

data "appid_cloud_directory_template" "template" {
  tenant_id = var.tenant_id
  template_name = "USER_VERIFICATION"  
}

output "template" {
  value = data.appid_cloud_directory_template.template
}