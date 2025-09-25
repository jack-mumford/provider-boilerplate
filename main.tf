terraform {
  required_providers {
    dob = {
      source = "liatr.io/terraform/devops-bootcamp"
    }
  }
}

provider "dob" {
  endpoint = "http://localhost:8080"
}

resource "dob_engineer" "Madi" {
  name  = "Madi"
  email = "madi@liatrio.com"
}

resource "dob_engineer" "Colin" {
  name  = "Colin"
  email = "colin@liatrio.com"
}

resource "dob_engineer" "Angel" {
  name  = "Angel"
  email = "angel@liatrio.com"
}

resource "dob_engineer" "Austin" {
  name  = "Austin"
  email = "austin@liatrio.com"
}

resource "dob_engineer" "Jack" {
  name  = "Jack"
  email = "jack@liatrio.com"
}
