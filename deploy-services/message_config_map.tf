resource "kubernetes_config_map" "message-service" {
  metadata {
    name      = "message-config"
    namespace = "app"
  }

  data = {
    "app-config.yaml" = "${file("message-config.yaml")}"
  }

}
