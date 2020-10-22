resource "google_cloud_scheduler_job" "job" {
  name        = "resource-cleanup-job"
  description = "GCP resource cleanup job"
  schedule    = var.schedule

  pubsub_target {
    topic_name = google_pubsub_topic.topic.id
    data       = base64encode("instance templates")
  }
}