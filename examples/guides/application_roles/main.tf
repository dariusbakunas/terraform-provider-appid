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