terraform {
  required_providers {
    sumologicextra = {
      source = "scastria/sumologicextra"
    }
  }
}

provider "sumologicextra" {
}

resource "sumologicextra_collector" "collector" {
  name        = "my-collector"
  description = "Test"
}
