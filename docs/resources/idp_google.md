---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "appid_idp_google Resource - terraform-provider-appid"
subcategory: ""
description: |-
  Update Google identity provider configuration.
---

# appid_idp_google (Resource)

Update Google identity provider configuration.

## Example Usage

```terraform
resource "appid_idp_google" "gg" {
  tenant_id = "<your appid tenant_id>"
  is_active = true
  
  config {
    application_id 		= "test_id"
    application_secret 	= "test_secret"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **is_active** (Boolean)
- **tenant_id** (String) The service `tenantId`

### Optional

- **config** (Block List, Max: 1) (see [below for nested schema](#nestedblock--config))
- **id** (String) The ID of this resource.

### Read-Only

- **redirect_url** (String) Paste the URI into the into the Authorized redirect URIs field in the Google Developer Console

<a id="nestedblock--config"></a>
### Nested Schema for `config`

Required:

- **application_id** (String) Google application id
- **application_secret** (String) Google application secret

