terraform {
  required_providers {
    google = {}

    google-beta = {}
  }
}



provider "google" {
  impersonate_service_account = "cw-prod-resman-000@cw-prod-iac-core-000.iam.gserviceaccount.com"
  # access_token = ""
}