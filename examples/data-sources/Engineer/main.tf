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

data "devops-bootcamp_engineer" "all" {}

output "engineers" {
  value = data.devops-bootcamp_engineer.all.engineers
}
