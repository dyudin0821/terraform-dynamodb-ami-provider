package amidynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &amiDynamoProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &amiDynamoProvider{
			version: version,
		}
	}
}

// amiDynamoProvider is the provider implementation.
type amiDynamoProvider struct {
	version string
}

type amiPromoteProviderModel struct {
	Region types.String `tfsdk:"region"`
}

// Metadata returns the provider type name.
func (p *amiDynamoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "amidynamodb"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *amiDynamoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required:    true,
				Description: "AWS Region",
			},
		},
		Blocks:      map[string]schema.Block{},
		Description: "AMI DynamoDB provider that implements obtaining the AWS AMI ID value from AWS DynamoDB tables for further deployment of EC2 instances.",
	}
}

func (p *amiDynamoProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring amidynamodb client")
	var config amiPromoteProviderModel
	diags := request.Config.Get(ctx, &config)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	var awsRegion string = config.Region.ValueString()

	if awsRegion == "" {
		response.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Missing AWS Region",
			"The provider cannot create the AWS API client as there is a missing or empty value for the amidynamodb AWS Region.",
		)
	}

	tflog.Info(ctx, fmt.Sprintf("region: %s", awsRegion))

	if response.Diagnostics.HasError() {
		return
	}

	awsConfig := aws.Config{
		Region: aws.String(awsRegion),
	}

	sess, err := session.NewSession(&awsConfig)
	if err != nil {
		response.Diagnostics.AddError("AWS Session сreation error", fmt.Sprintf("An error occurred creating an AWS session. Error Details: %s", err))
	}

	dynamoDBClient := dynamodb.New(sess)
	response.DataSourceData = dynamoDBClient
	response.ResourceData = dynamoDBClient

	tflog.Info(ctx, "Configured amidynamodb client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *amiDynamoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAmiDynamoDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *amiDynamoProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
