terraform {
  required_providers {
    appid = {
      source = "dariusbakunas/appid"
      version = "0.2.0"
    } 
  }
}

provider "appid" {  
    iam_access_token = var.iam_access_token      
    region = "us-south"
}

data "appid_applications" "source" {
  tenant_id = var.source_tenant_id
}

resource "appid_application" "destination" {
  count = length(data.appid_applications.source.applications)
  tenant_id = var.destination_tenant_id
  name = data.appid_applications.source.applications[count.index].name
  type = data.appid_applications.source.applications[count.index].type
  scopes = data.appid_applications.source.applications[count.index].scopes
}
