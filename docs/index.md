# SumoLogic Extra Provider
The SumoLogic Extra provider extends the official SumoLogic provider with `use_existing` flags to make it easier
to import existing resources without having to run `terraform import`.
## Example Usage
```hcl
terraform {
  required_providers {
    sumologicextra = {
      source  = "scastria/sumologicextra"
      version = "~> 0.1.0"
    }
  }
}

# Configure the SumoLogic Extra Provider
provider "sumologicextra" {
  num_retries = 3
  retry_delay = 30
}
```
## Argument Reference
* `num_retries` - **(Optional, Integer)** Number of retries for each SumoLogic API call in case of 429-Too Many Requests or any 5XX status code. Can be specified via env variable `SUMOLOGIC_NUM_RETRIES`. Default: 3.
* `retry_delay` - **(Optional, Integer)** How long to wait (in seconds) in between retries. Can be specified via env variable `SUMOLOGIC_RETRY_DELAY`. Default: 30.
