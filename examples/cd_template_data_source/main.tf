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

data "appid_cloud_directory_template" "template" {
  tenant_id = var.tenant_id
  template_name = "USER_VERIFICATION"  
}

output "template" {
  value = data.appid_cloud_directory_template.template
}