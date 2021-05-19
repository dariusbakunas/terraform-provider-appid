resource "appid_token_config" "tc" {
  tenant_id = "<your appid tenant_id>"  
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
}