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

resource "appid_redirect_urls" "urls" {
    tenant_id = var.tenant_id
    urls = [
        "https://localhost:3000"
    ]
}