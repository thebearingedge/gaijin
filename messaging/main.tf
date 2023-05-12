provider "google" {
  project = "pub-sub-golang"
}

resource "google_pubsub_topic" "bussin_tx" {
  name = "bussin-tx"

  message_retention_duration = "86600s"
}


resource "google_pubsub_subscription" "bussin_rx" {
  name  = "bussin-rx"
  topic = google_pubsub_topic.bussin_tx.name

  message_retention_duration = "1200s"

  ack_deadline_seconds = 20

  expiration_policy {
    ttl = "300000.5s"
  }

  retry_policy {
    minimum_backoff = "10s"
  }
}
