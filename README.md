# Terraform AWS AMI ID Provider

This Terraform provider retrieves the Amazon Machine Image (AMI) ID value from an AWS DynamoDB table using a conditional expression. It is written in Golang and uses the `terraform-plugin-framework` package.

## Requirements

- Terraform v1.0+
- Golang 1.16+ (for building the provider)

## Installation

1. Clone the repository and build the provider:

```bash
git clone https://github.com/yourusername/terraform-dynamodb-ami-provider.git
cd terraform-provider-amidynamodb
go build -o terraform-provider-amidynamodb

mkdir -p ~/.terraform.d/plugins/yourusername/example/1.0.0/linux_amd64
mv terraform-provider-amidynamodb ~/.terraform.d/plugins/yourusername/example/1.0.0/linux_amd64
```

or 

dowload from Terraform Registry: [https://registry.terraform.io/providers/dyudin0821/amidynamodb](https://registry.terraform.io/providers/dyudin0821/amidynamodb/latest)


## Usage

### Provider Configuration

To configure the provider, you need to specify the AWS region:

```hcl
provider "amidynamodb" {
  aws_region = "us-west-2"
}
```

### Data Source: amidynamodb_get_item
This data source retrieves the AMI ID from a DynamoDB table based on the provided conditional expression.

Arguments:
* `table_name` - (Required) The name of the DynamoDB table.
* `filter_expression` - (Required) The conditional expression used to filter values from the DynamoDB table.
* `expression_attribute_values` - (Optional) A map of attribute values used in the conditional expression.


Attributes:  
* `ami_id` - The AMI ID retrieved from the DynamoDB table.

Example:  
```hcl
data "amidynamodb_get_item" "example" {
  table_name = "ami-table"
  filter_expression = "latest = :latest"
  expression_attribute_values = {
    ":latest" = true
  }
}
```