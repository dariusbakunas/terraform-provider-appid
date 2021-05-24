resource "appid_redirect_urls" "urls" {
    tenant_id = "<your appid tenant_id>"
    urls = [
        "https://localhost:3000"
    ]
}