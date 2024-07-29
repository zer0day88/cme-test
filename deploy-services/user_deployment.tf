resource "kubernetes_deployment" "user-service" {
  metadata {
    name      = "user-app"
    namespace = "app"
  }

  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "user-app"
      }
    }

    template {
      metadata {
        labels = {
          app = "user-app"
        }
      }

      spec {
        container {
          image = "zer0day88/cme-user:latest"
          name  = "user-deployment"

          volume_mount {
            mount_path = "/root/config"
            name       = "user-config"
          }

          port {
            container_port = 8080
          }
        }

        volume {
          name = "user-config"
          config_map {
            name = kubernetes_config_map.user-service.metadata[0].name
          }

        }
      }
    }
  }
}
