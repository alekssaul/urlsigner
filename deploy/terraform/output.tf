output "keyset" {
    value = google_network_services_edge_cache_keyset.default.name
    description = "Media CDN keyset"
}

output "media_cdn_ip" {
    description = "Media CDN IP Address"
    value = google_network_services_edge_cache_service.instance.ipv4_addresses
}