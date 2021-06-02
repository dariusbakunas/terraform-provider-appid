---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "appid_languages Data Source - terraform-provider-appid"
subcategory: ""
description: |-
  User localization configuration
---

# appid_languages (Data Source)

User localization configuration

## Example Usage

```terraform
data "appid_languages" "languages" {
    tenant_id = "<your tenant id>"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **tenant_id** (String) The service `tenantId`

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **languages** (List of String) The list of languages that can be used to customize email templates for Cloud Directory

