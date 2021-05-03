---
subcategory: ""
page_title: "Manage Applications/Roles using variables"
description: |-
    An example of using variable to create applications and associated roles.
---

Given variable file:

```terraform
variable "tenant_id" {
  type = string  
  description = "AppID tenant ID"
}

variable "iam_access_token" {
    type = string
}

variable "applications" {
    type = map
    default = {
        "Web App #1": {
            type = "regularwebapp"
            roles = [
                { name: "admin", scopes: ["create_posts"]},
                { name: "regular_user", scopes: ["read_posts", "fav_posts"]}
            ],            
        },
        "Web App #2": {
            type = "singlepageapp"
            roles = [
                { name: "admin", scopes: ["create_posts"]},
                { name: "regular_user", scopes: ["read_posts"]},
                { name: "guest", scopes: ["list_posts"]},
            ],            
        }
    }        
}
```

You could create AppID applications and associate roles:

```terraform
terraform {
  required_providers {
    appid = {
      version = "0.1"
      source  = "us.ibm.com/watson-health/appid"
    }
  }
}

provider "appid" {  
    iam_access_token = var.iam_access_token  
    iam_base_url = "https://iam.cloud.ibm.com"
    appid_base_url = "https://us-south.appid.cloud.ibm.com"
}

locals {
    applications = {
        for app_name, app in var.applications:
        app_name => {
            "type" = app.type
            "scopes" = toset(flatten(app.roles[*].scopes))
        }
    }
    all_roles = flatten([
        for app_name, app in var.applications: [
            for role in app.roles:
            {
                name = role.name
                access = [
                    {
                        application = app_name
                        scopes = role.scopes
                    }                    
                ]
            }
        ] 
    ])

    roles = {
        for role_name in toset(local.all_roles[*].name):
        role_name => {            
            access = flatten(matchkeys(local.all_roles[*].access, local.all_roles[*].name, [role_name]))
        }
    }
}

resource "appid_application" "apps" {
    for_each = local.applications
    tenant_id = var.tenant_id
    name = each.key
    type = each.value.type
    scopes = each.value.scopes
}

resource "appid_role" "roles" {
    for_each = local.roles
    tenant_id = var.tenant_id
    name = each.key

    dynamic "access" {
        for_each = each.value.access
        content {
            application_id = appid_application.apps[access.value.application].client_id
            scopes = access.value.scopes
        }        
    }
}

output applications {
    value = local.applications
}

output roles {
    value = local.roles
}
```