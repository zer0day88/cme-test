resource "kubernetes_service" "user-service" {
  metadata {
    name      = "user-app"
    namespace = "app"
  }

  spec {
    selector = {
      app = kubernetes_deployment.user-service.metadata[0].name
    }

    port {
      name        = "http"
      protocol    = "TCP"
      port        = 8080
      target_port = 8080
    }
  }

}
