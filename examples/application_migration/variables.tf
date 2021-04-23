variable "source_tenant_id" {
  type = string
  default = "1a14fe5e-e258-4323-b5fb-cd35ab7d69da"
  description = "Source AppID tenant ID"
}

variable "destination_tenant_id" {
  type = string
  default = "bac33a56-501f-493c-8b1e-bfda921f4a3e"
  description = "Destination AppID tenant ID"
}

variable "iam_access_token" {
    type = string
}