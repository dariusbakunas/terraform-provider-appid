resource "appid_application" "app" {
  tenant_id = "<your appid tenant_id>"
  name = "<application name, must not exceed 50 characters>"  
  type = "<singlepageapp or regularwebapp>"
  scopes = ["test_scope_1", "test_scope_2", "test_scope_3"]
}
