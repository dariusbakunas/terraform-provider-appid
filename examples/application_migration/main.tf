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
