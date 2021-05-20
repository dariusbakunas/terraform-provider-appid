resource "appid_cloud_directory_template" "tpl" {
    tenant_id = "<appid tenant id>"
    template_name = "WELCOME"
    html_body = file("template.html")
    language = "en"    
    subject = "Welcome %%{user.displayName}"
}