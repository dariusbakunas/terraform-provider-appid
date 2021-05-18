resource "appid_idp_google" "gg" {
  tenant_id = "<your appid tenant_id>"
  is_active = true
  
  config {
    application_id 		= "test_id"
    application_secret 	= "test_secret"
  }
}