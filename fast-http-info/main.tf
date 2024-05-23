/**
 * Copyright 2023 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 4.34.0"
    }
  }
}
locals {
  region     = "us-central1"
  project_id = "alpha0809"
}
provider "google" {
  region  = local.region
  project = local.project_id
}
resource "random_id" "default" {
  byte_length = 8
}

resource "google_storage_bucket" "source_bucket" {
  name                        = "${random_id.default.hex}-gcf-source-bucket"
  location                    = local.region
  uniform_bucket_level_access = true
}

data "archive_file" "default" {
  type        = "zip"
  output_path = "/tmp/function-source.zip"
  source_dir  = "app/"
}

resource "google_storage_bucket_object" "default" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.source_bucket.name
  source = data.archive_file.default.output_path # Path to the zipped function source code
}

resource "google_service_account" "account" {
  project      = local.project_id
  account_id   = "gcf-sa-http-info"
  display_name = "Service Account - used for both the cloud function"
}

data "google_cloud_run_service" "function_service" {
  name     = google_cloudfunctions2_function.default.name
  location = google_cloudfunctions2_function.default.location
}

resource "google_cloud_run_service_iam_member" "member" {
  location = data.google_cloud_run_service.function_service.location
  service  = data.google_cloud_run_service.function_service.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_project_iam_member" "artifactregistry_reader" {
  project    = local.project_id
  role       = "roles/artifactregistry.reader"
  member     = "serviceAccount:${google_service_account.account.email}"
  depends_on = [google_service_account.account]
}

resource "google_cloudfunctions2_function" "default" {
  depends_on = [
    google_project_iam_member.artifactregistry_reader,
  ]
  name        = "gcf-go-fast-http-info"
  location    = local.region
  description = "fast http golang function"

  build_config {
    runtime     = "go122"
    entry_point = "entryPoint" # Set the entry point in the code
    source {
      storage_source {
        bucket = google_storage_bucket.source_bucket.name
        object = google_storage_bucket_object.default.name
      }
    }
  }

  service_config {
    available_memory               = "256M"
    timeout_seconds                = 60
    ingress_settings               = "ALLOW_ALL"
    all_traffic_on_latest_revision = true
    service_account_email          = google_service_account.account.email
  }
}
output "trigger_url" {
  value = google_cloudfunctions2_function.default.service_config[0].uri
}