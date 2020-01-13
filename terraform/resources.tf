# ---------------------------------------------------------------------------------------- #
# Cloud Storage Resource Bucket: event_bucket
# ---------------------------------------------------------------------------------------- #
resource "google_storage_bucket" "event_bucket" {
  name     = "${var.event_bucket}"
  location = "${var.location}"
  force_destroy = true
}

# ---------------------------------------------------------------------------------------- #
# Cloud Storage Resource Bucket: code_bucket
# ---------------------------------------------------------------------------------------- #
resource "google_storage_bucket" "code_bucket" {
  name     = "${var.code_bucket}"
  location = "${var.location}"
  force_destroy = true
}

# ---------------------------------------------------------------------------------------- #
# Cloud Functions Resource Function: gcs-to-bigquery
# ---------------------------------------------------------------------------------------- #
data "archive_file" "gcs_to_bigquery" {
  type        = "zip"
  output_path = "../output/gcs-to-bigquery.zip"
  source_dir = "../gcs-to-bigquery"
}
 
resource "google_storage_bucket_object" "gcs_to_bigquery" {
  name   = "gcs-to-bigquery.zip"
  bucket = "${google_storage_bucket.code_bucket.name}"
  source = "../output/gcs-to-bigquery.zip"
  depends_on = ["data.archive_file.gcs_to_bigquery"]
}
 
resource "google_cloudfunctions_function" "gcs_to_bigquery" {
  name        = "gcs-to-bigquery"
  description = "CloudFunctions GCS Sample Application."
  runtime     = "go111"

  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.code_bucket.name}"
  source_archive_object = "${google_storage_bucket_object.gcs_to_bigquery.name}"
  event_trigger = {
    event_type = "google.storage.object.finalize"
    resource = "${google_storage_bucket.event_bucket.name}"
  }
  entry_point           = "GcfEventHandler"

  environment_variables = {
    PROJECT_ID = "${var.project_id}"
    SRC_BUCKET = "${google_storage_bucket.event_bucket.name}"
    BIGQUERY_DATASET_NAME = "${var.dataset_id}"
    BIGQUERY_TABLE_NAME = "${var.table_id}"
  }
}

# ---------------------------------------------------------------------------------------- #
# Cloud Functions Resource Function: httptrigger-to-gcs
# ---------------------------------------------------------------------------------------- #
data "archive_file" "httptrigger_to_gcs" {
  type        = "zip"
  output_path = "../output/httptrigger-to-gcs.zip"
  source_dir = "../httptrigger-to-gcs"
}

resource "google_storage_bucket_object" "httptrigger_to_gcs" {
  name   = "httptrigger-to-gcs.zip"
  bucket = "${google_storage_bucket.code_bucket.name}"
  source = "../output/httptrigger-to-gcs.zip"
  depends_on = ["data.archive_file.httptrigger_to_gcs"]
}

resource "google_cloudfunctions_function" "httptrigger_to_gcs" {
  name        = "httptrigger-to-gcs"
  description = "CloudFunctions Sample Application."
  runtime     = "go111"

  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.code_bucket.name}"
  source_archive_object = "${google_storage_bucket_object.httptrigger_to_gcs.name}"
  trigger_http          = true
  entry_point           = "GcfWebhookHandler"

  environment_variables = {
    DEST_BUCKET = "${google_storage_bucket.event_bucket.name}"
  }
}

resource "google_cloudfunctions_function_iam_member" "httptrigger_to_gcs" {
  project        = "${google_cloudfunctions_function.httptrigger_to_gcs.project}"
  region         = "${google_cloudfunctions_function.httptrigger_to_gcs.region}"
  cloud_function = "${google_cloudfunctions_function.httptrigger_to_gcs.name}"

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}

# ---------------------------------------------------------------------------------------- #
# Big Query Resource Dataset / Table: sample
# ---------------------------------------------------------------------------------------- #
resource "google_bigquery_dataset" "gcs_to_bigquery" {
  dataset_id  = "${var.dataset_id}"
  description = "Sample dataset"
  location    = "${var.location}"
}

resource "google_bigquery_table" "gcs_to_bigquery" {
  dataset_id = "${google_bigquery_dataset.gcs_to_bigquery.dataset_id}"
  table_id = "${var.table_id}"
  description = "Sample table"
  schema     = "${file("bigquery/schema.json")}"

  time_partitioning {
    type  = "DAY"
  }
}

# ---------------------------------------------------------------------------------------- #
# Output
# ---------------------------------------------------------------------------------------- #
output "https_trigger_url" {
  value = "${google_cloudfunctions_function.httptrigger_to_gcs.https_trigger_url}"
}
