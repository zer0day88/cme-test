variable "kind_cluster_name" {
  type        = string
  description = "The name of the cluster."
  default     = "demo-local"
}

variable "kind_cluster_config_path" {
  type        = string
  description = "The location where this cluster's kubeconfig will be saved to."
  default     = "~/.kube/config"
}

variable "ingress_nginx_helm_version" {
  type        = string
  description = "The Helm version for the nginx ingress controller."
  default     = "4.7.1"
}

variable "ingress_nginx_namespace" {
  type        = string
  description = "The nginx ingress namespace (it will be created if needed)."
  default     = "ingress-nginx"
}

variable "redis_helm_version" {
  type        = string
  description = "The redis version."
  default     = "19.6.4"
}

variable "redis_namespace" {
  type        = string
  description = "The redis namespace (it will be created if needed)."
  default     = "app"
}

variable "cassandra_helm_version" {
  type        = string
  description = "The cassandra version."
  default     = "11.3.11"
}

variable "cassandra_namespace" {
  type        = string
  description = "The cassandra namespace (it will be created if needed)."
  default     = "app"
}
