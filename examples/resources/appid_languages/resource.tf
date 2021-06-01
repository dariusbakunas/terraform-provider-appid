resource "appid_languages" "langs" {
    tenant_id = "<your appid tenant_id>"
    languages = [
        "en",
        "es-ES",
        "fr-FR",
    ]
}