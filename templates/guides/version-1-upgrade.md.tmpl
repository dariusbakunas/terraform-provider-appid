# Terraform AppID Provider Version 1 Upgrade Guide

Version 1.0.0 of the AppID provider for Terraform is a major release and includes some breaking changes.
Few resources had their ID format changed, to support new Terraform import functionality.

Resources affected:

* **appid_apm** (id changed from `<tenant_id>/apm` to `<tenant_id>)
* **appid_audit_status** (id changed from `<tenant_id>/auditStatus` to `<tenant_id>`)
* **appid_idp_cloud_directory** (id changed from `<tenant_id>/idp/cloud_directory` to `<tenant_id>`)
* **appid_idp_custom** (id changed from `<tenant_id>/idp/custom_idp` to `<tenant_id>`)
* **appid_idp_saml** (id changed from `<tenant_id>/idp/saml` to `<tenant_id>`)
* **appid_languages** (id changed from `<tenant_id>/languages` to `<tenant_id>`)
* **appid_media** (id changed from `<tenant_id>/media` to `<tenant_id>`)
* **appid_mfa** (id changed from `<tenant_id>/mfa` to `<tenant_id>`)
* **appid_role** (id changed from `<role_id>` to `<tenant_id>/<role_id>`)
* **appid_theme_color** (id changed from `<tenant_id>/themeColors` to `<tenant_id>`)

If you run `terraform plan` while Terraform state contains any of these resources that were created using AppID provider v0.x.x you may get this error:

```bash
2021/06/03 10:51:25 [INFO] backend/local: plan operation completed
panic: runtime error: index out of range [1] with length 1
2021-06-03T10:51:25.336-0400 [DEBUG] plugin.terraform-provider-appid: 
2021-06-03T10:51:25.336-0400 [DEBUG] plugin.terraform-provider-appid: goroutine 121 [running]:
2021-06-03T10:51:25.336-0400 [DEBUG] plugin.terraform-provider-appid: github.ibm.com/dbakuna/terraform-provider-appid/appid.resourceAppIDRoleRead(0x1b6f0c8, 0xc0006399e0, 0xc0001f7400, 0x1a41f00, 0xc000010200, 0xc0007a4610, 0xc000317908, 0x100f818)
<...>
```

One way of solving this, is to remove affected resource from state (make sure to note current resource id, for example `role_id`):

```bash
terraform state rm appid_role.role
```

And then re-import that resource:

```bash
terraform import appid_role.role <tenant_id>/<role_id>
```