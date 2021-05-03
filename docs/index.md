---
layout: ""
page_title: "Provider: AppID"
description: |-
  The AppID provider provides resources to interact with a IBM AppID API.
---

# AppID Provider

The AppID provider provides resources to interact with a IBM AppID API.

## Example Usage

```terraform
provider "appid" {  
    iam_access_token = var.iam_access_token
    region = "us-south"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **appid_base_url** (String) AppID API base URL (for example 'https://us-south.appid.cloud.ibm.com')
- **iam_access_token** (String, Sensitive) The IBM Cloud Identity and Access Management token used to access AppID APIs
- **iam_api_key** (String, Sensitive) The IBM Cloud IAM api key used to retrieve IAM access token if `iam_access_token` is not specified
- **iam_base_url** (String) IBM IAM base URL
- **region** (String) The IBM cloud Region (for example 'us-south').