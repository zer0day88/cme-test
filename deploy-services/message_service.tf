resource "kubernetes_service" "message-service" {
  metadata {
    name      = "message-app"
    namespace = "app"
  }

  spec {
    selector = {
      app = kubernetes_deployment.message-service.metadata[0].name
    }

    port {
      name        = "http"
      protocol    = "TCP"
      port        = 7070
      target_port = 7070
    }
  }

}
