# Resource: sumologicextra_collector
Represents a collector
## Example usage
```hcl
resource "sumologicextra_collector" "example" {
  name = "My Collector"
}
```
## Argument Reference
* `name` - **(Required, String)** The name of the repository.
* `use_existing` - **(Optional, Boolean, IgnoreDiffs)** During a CREATE only, look for an existing collector with the same `name`.  Prevents the need for an import. Default: `false`
## Attribute Reference
* `id` - **(String)** The ID of the collector.
## Import
Collectors can be imported using a proper value of `id` as described above
