---
subcategory: "EC2 Image Builder"
layout: "aws"
page_title: "AWS: aws_imagebuilder_distribution_configurations"
description: |-
    Get information on Image Builder Distribution Configurations.
---


<!-- Please do not edit this file, it is generated. -->
# Data Source: aws_imagebuilder_distribution_configurations

Use this data source to get the ARNs and names of Image Builder Distribution Configurations matching the specified criteria.

## Example Usage

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.data_aws_imagebuilder_distribution_configurations import DataAwsImagebuilderDistributionConfigurations
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        DataAwsImagebuilderDistributionConfigurations(self, "example",
            filter=[DataAwsImagebuilderDistributionConfigurationsFilter(
                name="name",
                values=["example"]
            )
            ]
        )
```

## Argument Reference

This data source supports the following arguments:

* `filter` - (Optional) Configuration block(s) for filtering. Detailed below.

## filter Configuration Block

The `filter` configuration block supports the following arguments:

* `name` - (Required) Name of the filter field. Valid values can be found in the [Image Builder ListDistributionConfigurations API Reference](https://docs.aws.amazon.com/imagebuilder/latest/APIReference/API_ListDistributionConfigurations.html).
* `values` - (Required) Set of values that are accepted for the given filter field. Results will be selected if any given value matches.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `arns` - Set of ARNs of the matched Image Builder Distribution Configurations.
* `names` - Set of names of the matched Image Builder Distribution Configurations.

<!-- cache-key: cdktf-0.20.8 input-2b1ffe9a927d65041ade7f031f4f84f9b05ecfa961af5a27d1b6007c6b171cf3 -->