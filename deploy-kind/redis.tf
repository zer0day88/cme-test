

resource "helm_release" "redis" {
  name       = "redis"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "redis"
  version    = var.redis_helm_version

  namespace        = var.redis_namespace
  create_namespace = true

  values = [file("redis_values.yaml")]

  depends_on = [kind_cluster.default]
}

resource "null_resource" "wait_for_redis" {
  triggers = {
    key = uuid()
  }

  provisioner "local-exec" {
    command = <<EOF
      printf "\nWaiting for redis...\n"
      kubectl wait --namespace ${helm_release.redis.namespace} \
        --for=condition=ready pod \
        --selector=app.kubernetes.io/name=redis \
        --timeout=300s
    EOF
  }

  depends_on = [helm_release.redis]
}
