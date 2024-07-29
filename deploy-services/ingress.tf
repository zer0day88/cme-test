
resource "kubernetes_ingress_v1" "ingress-cme" {
  wait_for_load_balancer = true
  metadata {
    name = "ingress-cme"
    annotations = {
      "nginx.ingress.kubernetes.io/rewrite-target" : "/$2"
    }
    namespace = "app"
  }
  spec {

    rule {
      http {
        path {
          path      = "/user(/|$)(.*)"
          path_type = "Prefix"
          backend {
            service {
              name = kubernetes_service.user-service.metadata.0.name
              port {
                number = 8080
              }
            }

          }
        }

        path {
          path      = "/message(/|$)(.*)"
          path_type = "Prefix"
          backend {
            service {
              name = kubernetes_service.message-service.metadata.0.name
              port {
                number = 7070
              }
            }

          }
        }
      }
    }
  }
}

