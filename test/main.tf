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
  name        = "shawn_test"
  use_existing   = true
}
