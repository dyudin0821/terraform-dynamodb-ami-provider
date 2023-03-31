package amidynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"strings"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &amiDynamoDataSource{}
	_ datasource.DataSourceWithConfigure = &amiDynamoDataSource{}
)

func NewAmiDynamoDataSource() datasource.DataSource {
	return &amiDynamoDataSource{}
}

// amiDynamoDataSource is the data source implementation.
type amiDynamoDataSource struct {
	dynamodb *dynamodb.DynamoDB
}

type amiDynamoDataSourceModel struct {
	TableName            types.String `tfsdk:"table_name"`
	FilterExpression     types.String `tfsdk:"filter_expression"`
	ExpressionAttrValues types.Map    `tfsdk:"expression_attribute_values"`
	AmiID                types.String `tfsdk:"ami_id"`
}

func (a *amiDynamoDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dynamodb.DynamoDB)

	if !ok {
		tflog.Error(ctx, "Unable to prepare AWS DynamoDB client")
		return
	}

	a.dynamodb = client
}

func (a *amiDynamoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_get_item"
}

func (a *amiDynamoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"table_name": schema.StringAttribute{
				Description: "AWS DynamoDB table name",
				Required:    true,
			},
			"filter_expression": schema.StringAttribute{
				Description: "AWS DynamoDB the condition(s) an attribute(s) must meet",
				Required:    true,
			},
			"expression_attribute_values": schema.MapAttribute{
				ElementType: types.StringType,
				Description: "One or more substitution tokens for attribute names in an expression.",
				Required:    true,
			},
			"ami_id": schema.StringAttribute{
				Description: "The return value of the AWS AMI ID from the DynamoDB table",
				Computed:    true,
			},
		},
		Blocks:      map[string]schema.Block{},
		Description: "Interface with amidynamodb API",
	}
}

func (a *amiDynamoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Preparing to read item data source")
	var state amiDynamoDataSourceModel
	var expressionAttributeNames map[string]*string
	var projectionExpression *string

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(state.TableName.ValueString()),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: convertExpressionAttributeValues(ctx, state.ExpressionAttrValues.Elements()),
		FilterExpression:          aws.String(state.FilterExpression.ValueString()),
		ProjectionExpression:      projectionExpression,
	}

	output, err := a.dynamodb.Scan(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Item",
			err.Error(),
		)
		return
	}

	if len(output.Items) > 0 {
		amiID := output.Items[0]["ami"].S
		state.AmiID = types.StringValue(*amiID)
	} else {
		resp.Diagnostics.AddWarning("No matching items found in the DynamoDB table.", "")
		return
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func convertExpressionAttributeValues(ctx context.Context, values map[string]attr.Value) map[string]*dynamodb.AttributeValue {
	result := make(map[string]*dynamodb.AttributeValue)
	for key, value := range values {
		v := strings.Replace(value.String(), `"`, "", -1)
		if boolVal, _ := strconv.ParseBool(v); boolVal {
			result[key] = &dynamodb.AttributeValue{
				BOOL: aws.Bool(boolVal),
			}
		} else {

			result[key] = &dynamodb.AttributeValue{
				S: aws.String(v),
			}
		}
	}
	return result
}
