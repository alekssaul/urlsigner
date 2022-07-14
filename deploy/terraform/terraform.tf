data "google_project" "project" {
  project_id = var.project_id
}

# Enable APIs for Media Cache
resource "google_project_service" "networkservices" {
  project            = data.google_project.project.id
  service            = "networkservices.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "certificatemanager" {
  project            = data.google_project.project.id
  service            = "certificatemanager.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "edgecache" {
  project            = data.google_project.project.id
  service            = "edgecache.googleapis.com"
  disable_on_destroy = false
}

# Media CDN Keyset
resource "random_string" "key_id" {
  length  = 16
  special = false
  number = false
  upper = false
}

resource "google_network_services_edge_cache_keyset" "default" {
  name        = "default"
  description = "The default keyset"
  public_key {
    id    = random_string.key_id.result
    value = file("${path.module}/assets/public.key")
  }
  depends_on = [
    google_project_service.edgecache,
    google_project_service.certificatemanager,
    google_project_service.networkservices
  ]
}

resource "google_storage_bucket" "origin" {
  name                        = var.gcs_bucket_name
  location                    = "US"
  force_destroy               = false
  uniform_bucket_level_access = true
}

resource "google_network_services_edge_cache_origin" "instance" {
  name           = "my-origin"
  origin_address = google_storage_bucket.origin.url
  description    = "The default bucket for media edge test"
  max_attempts   = 2
  timeout {
    connect_timeout = "10s"
  }
  depends_on = [
    google_project_service.edgecache,
    google_project_service.certificatemanager,
    google_project_service.networkservices
  ]
}

resource "google_network_services_edge_cache_service" "instance" {
  name        = "my-service"
  description = "some description"
  depends_on = [
    google_project_service.edgecache,
    google_project_service.certificatemanager,
    google_project_service.networkservices
  ]
  edge_ssl_certificates = ["projects/ce-alekssaul-340117/locations/global/certificates/stream-example-com"]
  routing {
    host_rule {
      description  = "host rule description"
      hosts        = [var.domain_name]
      path_matcher = "routes"
    }
    path_matcher {
      name = "routes"
      route_rule {
        description = "a route rule to match against"
        priority    = 1
        match_rule {
          prefix_match = "/"
        }
        origin = google_network_services_edge_cache_origin.instance.name
        route_action {
          cdn_policy {
            cache_mode            = "CACHE_ALL_STATIC"
            default_ttl           = "3600s"
            signed_request_mode   = "REQUIRE_SIGNATURES"
            signed_request_keyset = google_network_services_edge_cache_keyset.default.name
          }
        }
        header_action {
          response_header_to_add {
            header_name  = "x-cache-status"
            header_value = "{cdn_cache_status}"
          }
        }
      }
    }
  }
}

# Test File
resource "local_file" "helloworld" {
  content  = "Hello World"
  filename = "${path.module}/helloworld.txt"
}

resource "google_storage_bucket_object" "picture" {
  name   = "helloworld.txt"
  source = "${path.module}/helloworld.txt"
  bucket = google_storage_bucket.origin.name
}

resource "google_storage_bucket_iam_member" "mediaedge_acl" {
  bucket = google_storage_bucket.origin.name
  role   = "roles/storage.objectViewer"
  member = format("serviceAccount:service-%s@gcp-sa-mediaedgefill.iam.gserviceaccount.com", data.google_project.project.number)
}