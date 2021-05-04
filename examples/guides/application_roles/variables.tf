
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
