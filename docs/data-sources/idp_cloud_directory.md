---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "appid_idp_cloud_directory Data Source - terraform-provider-appid"
subcategory: ""
description: |-
  
---

# appid_idp_cloud_directory (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **tenant_id** (String)

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **config** (List of Object) (see [below for nested schema](#nestedatt--config))
- **is_active** (Boolean)

<a id="nestedatt--config"></a>
### Nested Schema for `config`

Read-Only:

- **identity_field** (String)
- **interactions** (List of Object) (see [below for nested schema](#nestedobjatt--config--interactions))
- **self_service_enabled** (Boolean)
- **signup_enabled** (Boolean)

<a id="nestedobjatt--config--interactions"></a>
### Nested Schema for `config.interactions`

Read-Only:

- **identity_confirmation** (List of Object) (see [below for nested schema](#nestedobjatt--config--interactions--identity_confirmation))
- **reset_password_enabled** (Boolean)
- **reset_password_notification_enabled** (Boolean)
- **welcome_enabled** (Boolean)

<a id="nestedobjatt--config--interactions--identity_confirmation"></a>
### Nested Schema for `config.interactions.identity_confirmation`

Read-Only:

- **access_mode** (String)
- **methods** (List of String)

