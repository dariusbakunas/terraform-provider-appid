resource "appid_idp_cloud_directory" "cd" {  
  tenant_id = "<your tenant id>"
  is_active = true
  reset_password_enabled = true
  reset_password_notification_enabled = true
  self_service_enabled = true
  signup_enabled = true
  welcome_enabled = true
}