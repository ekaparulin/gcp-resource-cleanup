# General

This setup cleans up unused GCP resources on a schedule.

## Limitations

Currently only GCE instance templates are cleaned up which are not used by any managed instance group and are older than specified age.

## Architecture

GCP Cloud scheduler runs on specified schedule and sends a message to a PubSub topic, which triggers an execution of a cloud Function which does the resource lookup and cleanup


# Local testing

- go to 'src' directory
- export the following variables:

    ```
    # path to service account json file
    export GOOGLE_APPLICATION_CREDENTIALS=...
    # GCP project name 
    export GCP_PROJECT=...
    # Days to preserve unused templates
    export DELETE_OLDER_DAYS=...
    ```

- run the test

    ```
    go test
    ```

# Deployment

Deployment is performed using Terraform 0.13.4 or newer. Please download it for your platform.

## Preparations

1. Enter the terraform directory
1. Set values in main.tf:

- optionally uncomment the backend settings for terraform state and set to an existing bucket
- set GCP project name
- set GCP region where cloud functions are supported
- set/change topic name
- set/change schedule


2. Run

terraform init

## Deploy/redeploy

terraform apply

## Destroy

terraform destroy