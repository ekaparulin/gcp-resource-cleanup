resource "google_pubsub_topic" "topic" {
  name = var.topic

  labels = {
    category = "maintenance-tools"
  }
}