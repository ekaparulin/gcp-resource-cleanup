# Local test

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