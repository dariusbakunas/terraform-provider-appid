resource "appid_apm" "apm" {
    tenant_id = "<your tenant id>"
    enabled = true
    prevent_password_with_username = true

    password_reuse {
        enabled = true
        max_password_reuse = 4
    }

    password_expiration {
        enabled = true
        days_to_expire = 25
    }

    lockout_policy {
        enabled = true
        lockout_time_sec = 2600
        num_of_attempts = 4
    }

    min_password_change_interval {
        enabled = true
        min_hours_to_change_password = 1
    }
}