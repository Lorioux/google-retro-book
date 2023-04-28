
locals {
  states    = jsondecode(file("${path.root}/state.json"))["resources"]
  instances = [for resource in local.states : merge(resource["instances"]...) if resource["name"] == "example"]
}


resource "google_compute_instance" "example" {
  for_each     = { for k, vm in local.instances : vm["index_key"] => vm["attributes"] }
  name         = each.value["name"]
  zone         = each.value["zone"]
  project      = each.value["project"]
  machine_type = each.value["machine_type"]
  # id           = each.value["id"]
  boot_disk {
    initialize_params {
      image = "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-11-bullseye-v20230206"
      # "labels" = {}
      size = 10
      type = "pd-standard"
    }
  }
  dynamic network_interface {
    for_each = each.value.network_interface
    content {
      dynamic access_config {
        for_each = network_interface.value["access_config"]
        content {
          nat_ip = access_config.value.nat_ip
          network_tier = access_config.value.network_tier
        }
      }
    }
  }

}