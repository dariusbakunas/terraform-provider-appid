resource "appid_action_url" "url" {
  tenant_id = "<your appid tenant_id>"
  action = "on_reset_password"
  url = "https://your-domain.com/?action=pw-reset"
}
