
terraform {
  required_providers {
    google = {}

    google-beta ={}
  }
}



provider "google" {
  impersonate_service_account = "cw-prod-resman-000@cw-prod-iac-core-000.iam.gserviceaccount.com"
  access_token = "ya29.a0Ael9sCM8RDF2l5IRB8gefad4Jm6o0K2x3N4KJA3G46XNOakUmtFXSsBQyBHS2QuH5tHOKj1KMm8Rs2n5RzSIkMwmKEQ5dDu5JfM2BgVOib5z2VOhB9VNeRqj17mbhOQCq1mHLry_25FzUPQYkUdYd_s8QalXudHIr_3fU-EaCgYKAT4SARISFQF4udJh-wcrPTw7inwyyla0O7hLNA0174"
}