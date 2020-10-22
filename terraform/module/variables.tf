variable gcp_project {
    description = "GCP project name"
}

variable gcp_region {
    description = "GCP Region to run cloud function, e.g. europe-west3"
}

variable schedule {
    description = "Cloud Scheduler schedule, e.g. 0 9 * * 1'"
}

variable topic {
    description = "PubSub topic to use, e.g. infra-cleanup"
}