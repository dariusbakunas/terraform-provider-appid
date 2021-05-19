---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "appid_password_regex Data Source - terraform-provider-appid"
subcategory: ""
description: |-
  
---

# appid_password_regex (Data Source)



## Example Usage

```terraform
data "appid_password_regex" "rgx" {
    tenant_id = "<your appid tenant_id>"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **tenant_id** (String) The service `tenantId`

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **base64_encoded_regex** (String) The regex expression rule for acceptable password encoded in base64
- **error_message** (String) Custom error message
- **regex** (String) The escaped regex expression rule for acceptable password

