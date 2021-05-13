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
    region = "us-south"
}

resource "appid_redirect_urls" "urls" {
    tenant_id = var.tenant_id
    urls = [
        "https://localhost:3000"
    ]
}