resource "appid_idp_custom" "idp" {
    tenant_id = "<your tenant_id>"
    is_active = true
    public_key = "<PUBLIC KEY>"
}