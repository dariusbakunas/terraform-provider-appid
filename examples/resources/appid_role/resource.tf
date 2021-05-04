resource "appid_role" "role" {
  tenant_id = var.tenant_id
  name = "test-tf-role"    
  
  access {
    application_id = "<application id that contains the scopes listed below>"
    scopes = [      
      "scope_1",
      "scope_2"
    ]
  }

  access {
    application_id = "<application id #2>"
    scopes = [      
      "scope_1",
      "scope_2"
    ]
  }

  access {
    application_id = "<application id #3>"
    scopes = [      
      "scope_1",
      "scope_2"
    ]
  }
}
