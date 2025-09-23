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

data "devops-bootcamp_devops" "all" {}

output "devops" {
  value = data.devops-bootcamp_devops.all.devops
}
