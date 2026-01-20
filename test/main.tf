terraform {
  required_providers {
    sumologicextra = {
      source = "github.com/scastria/sumologicextra"
    }
  }
}

provider "sumologicextra" {
}

resource "sumologicextra_collector" "collector" {
  name        = "my_collector"
  use_existing   = true
}
