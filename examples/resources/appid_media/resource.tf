resource "appid_media" "media" {
    tenant_id = var.tenant_id
    source = "~/Downloads/logo.png"
}

# Alternative usage:

resource "appid_media" "media" {
    tenant_id = var.tenant_id
    source_content = filebase64("~/Downloads/logo.png")
}
