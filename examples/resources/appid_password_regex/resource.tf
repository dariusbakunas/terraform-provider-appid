resource "appid_password_regex" "rgx" {
  tenant_id = "<your appid tenant_id>"
  regex = "^(?:(?=.*\\d)(?=.*[a-z])(?=.*[A-Z]).*)$"
  error_message = "Must have one number, one lowercase letter, and one capital letter."
}