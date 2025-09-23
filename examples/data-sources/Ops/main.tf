terraform {
  required_providers {
    devops-bootcamp = {
      source = "liatr.io/terraform/devops-bootcamp"
    }
  }
}

provider "devops-bootcamp" {
  endpoint = "http://localhost:8080"
}

data "devops-bootcamp_ops" "all" {}

output "ops" {
  value = data.devops-bootcamp_ops.all.ops
}
