resource "kubernetes_config_map" "user-service" {
  metadata {
    name      = "user-config"
    namespace = "app"
  }

  data = {
    "app-config.yaml" = "${file("user-config.yaml")}"
  }

}
