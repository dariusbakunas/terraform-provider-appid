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
    appid_base_url = "us-south"
}

data "appid_idp_saml" "saml" {
  tenant_id = var.tenant_id  
}

output "saml" {
  value = data.appid_idp_saml.saml
}