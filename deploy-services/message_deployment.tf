resource "kubernetes_deployment" "message-service" {
  metadata {
    name      = "message-app"
    namespace = "app"
  }

  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "message-app"
      }
    }

    template {
      metadata {
        labels = {
          app = "message-app"
        }
      }

      spec {
        container {
          image = "zer0day88/cme-message:latest"
          name  = "message-deployment"

          volume_mount {
            mount_path = "/root/config"
            name       = "message-config"
          }

          port {
            container_port = 7070
          }
        }

        volume {
          name = "message-config"
          config_map {
            name = kubernetes_config_map.message-service.metadata[0].name
          }

        }
      }
    }
  }
}
