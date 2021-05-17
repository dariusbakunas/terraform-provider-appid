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

data "appid_token_config" "source" {
  tenant_id = var.source_tenant_id
}

resource "appid_token_config" "destination" {    
    tenant_id = var.destination_tenant_id
    access_token_expires_in = data.appid_token_config.source.access_token_expires_in
    anonymous_access_enabled = data.appid_token_config.source.anonymous_access_enabled
    anonymous_token_expires_in = data.appid_token_config.source.anonymous_token_expires_in
    refresh_token_enabled = data.appid_token_config.source.refresh_token_enabled    
    refresh_token_expires_in = data.appid_token_config.source.refresh_token_expires_in    

    dynamic "access_token_claim" {
        for_each = data.appid_token_config.source.access_token_claim
        content {
            source = access_token_claim.value["source"]
            source_claim = access_token_claim.value["source_claim"]
            destination_claim = access_token_claim.value["destination_claim"]
        }
    }
}