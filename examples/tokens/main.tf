terraform {
  required_providers {
    appid = {
      version = "0.1"
      source  = "us.ibm.com/watson-health/appid"
    }
  }
}

provider "appid" {
    iam_base_url = "https://iam.cloud.ibm.com"
    appid_base_url = "https://us-south.appid.cloud.ibm.com"
}

resource "appid_config_tokens" "tokens" {
    tenant_id = var.tenant_id
    access_token_expires_in = 7200    
    anonymous_access_enabled = false    
    refresh_token_enabled = false
}
