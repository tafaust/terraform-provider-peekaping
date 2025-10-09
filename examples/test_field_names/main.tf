terraform {
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.0.3"
    }
  }
}

provider "peekaping" {
  endpoint = var.endpoint
  email    = var.email
  password = var.password
}

# Test all field names and types for monitors
resource "peekaping_monitor" "field_test" {
  name = "Field Test Monitor"
  type = "http"
  config = jsonencode({
    url                   = "https://example.com"
    method                = "GET"
    headers               = { "User-Agent" = "Test" }
    body                  = "test body"
    accepted_status_codes = [200, 201, 202]
    timeout               = 30
    follow_redirects      = true
    ignore_tls_errors     = false
    expected_body         = "expected content"
    expected_body_regex   = ".*success.*"
    max_redirects         = 5
    client_cert           = "cert content"
    client_key            = "key content"
    client_cert_password  = "password"
    basic_auth_username   = "user"
    basic_auth_password   = "pass"
    digest_auth_username  = "user"
    digest_auth_password  = "pass"
    ntlm_domain           = "domain"
    ntlm_username         = "user"
    ntlm_password         = "pass"
    oauth2_client_id      = "client_id"
    oauth2_client_secret  = "client_secret"
    oauth2_token_url      = "https://oauth.example.com/token"
    oauth2_scopes         = ["read", "write"]
    oauth2_audience       = "api.example.com"
    oauth2_resource       = "api.example.com"
  })
  interval         = 60
  timeout          = 30
  max_retries      = 3
  retry_interval   = 60
  resend_interval  = 10
  active           = true
  notification_ids = []
  tag_ids          = []
  proxy_id         = null
}
