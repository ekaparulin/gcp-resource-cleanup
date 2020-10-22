terraform {
    # Optionally uncomment the backend settings for terraform state and set to an existing bucket
    # backend "gcs" {
    #     bucket  = ""
    #     prefix  = "resource-cleanup"
    # }
}

module "resource-cleanup" {
    source = "./module"

    # set set GCP project name
    gcp_project = ""

    # set GCP region where cloud functions are supported
    gcp_region = "europe-west3"

    # set/change topic name
    topic = "resource-cleanup"

    # set/change schedule
    schedule = "0 9 * * 1"
}