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

resource "appid_idp_saml" "saml" {
  tenant_id = var.tenant_id
  is_active = true

  config {
    entity_id = "https://test-saml-idp"
    sign_in_url = "https://test-saml-idp/login"
    certificates = [
      <<EOT
      MIIFmjCCA4ICCQDsTVT6SQ82GTANBgkqhkiG9w0BAQsFADCBjjELMAkGA1UEBhMC
VVMxDTALBgNVBAgMBE9ISU8xEjAQBgNVBAcMCUNsZXZlbGFuZDEQMA4GA1UECgwH
S3luZHJ5bDEUMBIGA1UECwwLQmFzZW1lbnQgSVQxFDASBgNVBAMMC2t5bmRyeWwu
Y29tMR4wHAYJKoZIhvcNAQkBFg9reW5kcnlsQGlibS5jb20wHhcNMjEwNDI2MTUz
NTIwWhcNMjIwNDI2MTUzNTIwWjCBjjELMAkGA1UEBhMCVVMxDTALBgNVBAgMBE9I
SU8xEjAQBgNVBAcMCUNsZXZlbGFuZDEQMA4GA1UECgwHS3luZHJ5bDEUMBIGA1UE
CwwLQmFzZW1lbnQgSVQxFDASBgNVBAMMC2t5bmRyeWwuY29tMR4wHAYJKoZIhvcN
AQkBFg9reW5kcnlsQGlibS5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIK
AoICAQDeQZGeFHQ6rkqPeaZYLuGAE0O0x7nAinCivs9i8ZrM/At6UyS98OBCXfk3
tyLFzDg8nmUaFnqgJHUVRRFU5C+MatdvqWtpNHaqMLDxQmcy0w+kPOL4W9ECoOUe
08xakQfXmXIhbt0RMN0dBgyXVsUb0mrFeEOh9gw5O5xsz6EoowtnJHhqk2/dKMlo
R/Cx9tASIFVCpcsbianPSy5zf6KDmDW7f9Piay9ibAc7yEvqlfv3DFw4x+/V/rZa
KYbkVvh+0T9PKwbQkeEOFLJv5KMkoDG6YSLWbm1ho/28uW8i7SpciRUpmBux4ARV
lHtWNNN4PHf9ZfaCP8/3hWenOQqb6Fqqk/sZsHAkBkdao13dz3DUXWW+3c32LDuc
Lo3+9Uv5pXiSjHzxC+dpVff324WFynC9toc63IX8orn7ZNNHyskiQ3nbZll0aa3T
qB2dGPC7AxEdFJ7ZZORNm3TfK+PtI8GptCtBupSSHFq1r2F/Y4arVfzpNvC2Y5UW
s8kdfHN3+DT8+WTbnPHu6/+WFlNHnGH8B9DCri4yLCTjYbK1grcYunbEvaYmSVZX
7kjSxWsSB4l2dHmXEP7f8tZP1pm+t9J7rEZ3PxqLb1suFxbrL9TDQm71v+HEJ82b
ue8mTE6N3BQbfhb1ggZseM2UCQ8PqVZNWp6iAIm3t3w9/NTgcQIDAQABMA0GCSqG
SIb3DQEBCwUAA4ICAQAjLG1KbdY1tPa272bpW/36V4AqUmXbsXTPqm/wnRbu0nMM
sYo8oX5XFkNXJnGisgXYGsds+uPpgsYG1OpPHtiF6fP3bv78bB6AeNo/OHSs1xgo
+3cVRVB6UydpWBxVYhrrU+tSSOQH41MOlU10a/iJCTD2ftBKV2wIIto9Xx57gfUP
bjc+CeWHgx0rWNmnVFxVCx0Q0W1mfi85YHo5e+oT5b3V+YroEXa1vdkFjlm6/aY3
JAnyOwkHXB48e0aSd4PE46RaRI/wOG4WA8iA9sbzCqoATn8ZBYa6MJBTxBcJpzrE
N9sB9dAAUWj2hH8YDka6doZb8TJsn7/hNXoGDaHeTeVMWhECSgUu61FPWRHnAV6C
6pq2Z4wR7M4vp9gc3GLdrAC3DlLW43Wsxbe7aeLlYaZZ5Saezr+1MjEyI27/cGC6
6Tc7lo7mqXCcfUSHFMVRXGAhQ7hShVl5jOP8NPxC8yeDlwRHPdQGQ7GLfC70vsGl
q195EOxxLJ5fyJOdvlEbCw9WnOPLJ9sf4C0Lg8dbOnsxUKufRC6zJR3P2vNWQ2Z1
5oKtRnb/s6hOMdDoXMVen/1tS4NIi22L80OARJEmHSB3bdxlJe/TDkVXmQxbvi7+
VrWz2D2R2MUEAyw8m/J1d5k+agb/BmTguAa/pdhI4w6S2Gg0h67eU48Omdr+fQ==
      EOT
    ]
    display_name = "Test IDP"
    encrypt_response = false
    sign_request = false
  }
}

output "saml" {
  value = appid_idp_saml.saml
}