variable "project_id" {
  type        = string
  description = "GCP Project ID"
}

variable "region" {
  type        = string
  description = "GCP Region"
  default     = "us-east1"
}

variable "zone" {
  type        = string
  description = "GCP Zone"
  default     = "us-east1-b"

}

variable "gcs_bucket_name" {
  type        = string
  description = "Name of GCS Bucket to create"
}

variable "domain_name" {
  type        = string
  description = "Domain name to be used for CME setup"
}