data "amidynamodb_get_item" "example" {
  table_name = "ami-table"
  expression = "latest = :latest"
    expression_attribute_values = {
    ":latest" = true
  }
}