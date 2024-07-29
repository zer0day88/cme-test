

resource "helm_release" "cassandra" {
  name       = "cassandra"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "cassandra"
  version    = var.cassandra_helm_version

  namespace        = var.cassandra_namespace
  create_namespace = true

  values = [file("cassandra_values.yaml")]

  depends_on = [kind_cluster.default]

}

resource "null_resource" "wait_for_cassandra" {
  triggers = {
    key = uuid()
  }

  provisioner "local-exec" {
    command = <<EOF
      printf "\nWaiting for cassandra...\n"
      kubectl wait --namespace ${helm_release.cassandra.namespace} \
        --for=condition=ready pod \
        --selector=app.kubernetes.io/name=cassandra \
        --timeout=300s
      
    EOF
  }

  depends_on = [helm_release.cassandra]
}

