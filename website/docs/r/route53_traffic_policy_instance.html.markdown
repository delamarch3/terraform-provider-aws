---
subcategory: "Route 53"
layout: "aws"
page_title: "AWS: aws_route53_traffic_policy_instance"
description: |-
    Provides a Route53 traffic policy instance resource.
---

# Resource: aws_route53_traffic_policy_instance

Provides a Route53 traffic policy instance resource.

## Example Usage

```terraform
resource "aws_route53_traffic_policy_instance" "test" {
  name                   = "test.example.com"
  traffic_policy_id      = "b3gb108f-ea6f-45a5-baab-9d112d8b4037"
  traffic_policy_version = 1
  hosted_zone_id         = "Z033120931TAQO548OGJC"
  ttl                    = 360
}
```

## Argument Reference

The following arguments are required:

* `name` - (Required) Domain name for which Amazon Route 53 responds to DNS queries by using the resource record sets that Route 53 creates for this traffic policy instance.
* `traffic_policy_id` - (Required) ID of the traffic policy that you want to use to create resource record sets in the specified hosted zone.
* `traffic_policy_version` - (Required) Version of the traffic policy
* `hosted_zone_id` - (Required) ID of the hosted zone that you want Amazon Route 53 to create resource record sets in by using the configuration in a traffic policy.
* `ttl` - (Required) TTL that you want Amazon Route 53 to assign to all the resource record sets that it creates in the specified hosted zone.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `arn` - Amazon Resource Name (ARN) of the traffic policy instance.
* `id` - ID of traffic policy instance.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Route53 traffic policy instance using its id. For example:

```terraform
import {
  to = aws_route53_traffic_policy_instance.test
  id = "df579d9a-6396-410e-ac22-e7ad60cf9e7e"
}
```

Using `terraform import`, import Route53 traffic policy instance using its id. For example:

```console
% terraform import aws_route53_traffic_policy_instance.test df579d9a-6396-410e-ac22-e7ad60cf9e7e
```
