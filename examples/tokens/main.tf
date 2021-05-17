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

resource "appid_token_config" "tokens" {
    tenant_id = var.tenant_id
    access_token_expires_in = 7200    
    anonymous_access_enabled = true
    anonymous_token_expires_in = 3200    
    refresh_token_enabled = false 
    access_token_claim {
      source = "roles"
      destination_claim = "groupIds"
    }

    access_token_claim {
      source = "appid_custom"
      source_claim = "employeeId"
      destination_claim = "employeeId"
    }

    access_token_claim {
      source = "saml"
      source_claim = "attributes.uid"
      destination_claim = "employeeId"
    }

    access_token_claim {
      source = "attributes"
      source_claim = "employeeId"
      destination_claim = "employeeId"
    }

    access_token_claim {
      source = "attributes"
      source_claim = "tenantId"
      destination_claim = "tenantId"
    }

    access_token_claim {
      source = "attributes"
      source_claim = "userType"
      destination_claim = "userType"
    }
}

data "appid_token_config" "tokens" {
  tenant_id = appid_token_config.tokens.tenant_id

  depends_on = [
    appid_token_config.tokens
  ]
}