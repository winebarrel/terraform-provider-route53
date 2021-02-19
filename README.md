# terraform-provider-route53

## Usage

```terraform
provider "route53" {
}

data "route53_record" "my_record" {
  hosted_zone_id = "..."
  name           = "db.example.com"
  type           = "CNAME"
}

output "records" {
  value = join(",", data.route53_record.my_record.records)
}
```
