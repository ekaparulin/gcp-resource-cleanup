resource "google_storage_bucket" "bucket" {
  name = "${var.gcp_project}-gcf-sources-cleanup"
  location = var.gcp_region
  labels = {
    category = "maintenance-tools"
    purpose = "gcf-source-archives"
  }
}

data "archive_file" "sources" {
  type        = "zip"
  output_path = "${path.module}/files/function-cleanup.zip"

  source {
    content  = file("${path.module}/../../src/cleanup.go")
    filename = "cleanup.go"
  }

  source {
    content  = file("${path.module}/../../src/go.mod")
    filename = "go.mod"
  }
}

resource "google_storage_bucket_object" "archive" {
  name   = "function-cleanup-${filesha256(data.archive_file.sources.output_path)}.zip"
  bucket = google_storage_bucket.bucket.name
  source = data.archive_file.sources.output_path
}

resource "google_cloudfunctions_function" "function" {
  name        = "infra-cleanup"
  description = "Function to cleanup unused GCP resources"
  runtime     = "go113"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  ingress_settings = "ALLOW_INTERNAL_ONLY"
  max_instances = 1

  event_trigger {
      event_type = "google.pubsub.topic.publish"
      resource = "projects/${var.gcp_project}/topics/${google_pubsub_topic.topic.name}"
      failure_policy {
          retry = false
      }
  }
  timeout               = 60
  entry_point           = "Cleanup"
  labels = {
    category = "maintenance-tools"
    purpose = "cleanup"
  }

  environment_variables = {
    GCP_PROJECT = var.gcp_project
    DELETE_OLDER_DAYS = 14
  }
}
