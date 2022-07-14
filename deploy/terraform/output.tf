output "keyset" {
  value       = google_network_services_edge_cache_keyset.default.name
  description = "Media CDN keyset"
}

output "media_cdn_ip" {
  description = "Media CDN IP Address"
  value       = google_network_services_edge_cache_service.instance.ipv4_addresses
}

output "keyset_primary_private" {
    description = "Version of Private Keyset secret in Secret Manager"
    value = format("%s:%s", google_secret_manager_secret.keyset_primary_private.secret_id , data.google_secret_manager_secret_version.keyset_primary_private_version.version)
}