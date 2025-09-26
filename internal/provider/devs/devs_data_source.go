package devs

import (
	"context"

	"terraform-provider-devops/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &devsDataSource{}
	_ datasource.DataSourceWithConfigure = &devsDataSource{}
)

// NewDevsDataSource is a helper function to simplify the provider implementation.
func NewDevsDataSource() datasource.DataSource {
	return &devsDataSource{}
}

// devsDataSource is the data source implementation.
type devsDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *devsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev"
}

// Schema defines the schema for the data source.
func (d *devsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"devs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
						"engineers": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

// Configure accepts provider configured data to set the API client.
func (d *devsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			"Expected *client.Client for provider data, but received a different type.",
		)
		return
	}

	d.client = c
}

// Read refreshes the Terraform state with the latest data.
func (d *devsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DevsDataSourceModel

	devs, err := d.client.GetDevs()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Dev groups",
			err.Error(),
		)
		return
	}

	for _, dv := range devs {
		dvm := devsModel{
			ID:   types.StringValue(dv.ID),
			Name: types.StringValue(dv.Name),
		}

		// Convert engineers (objects) to a list of engineer IDs
		engineerIDs := make([]string, 0, len(dv.Engineers))
		for _, eng := range dv.Engineers {
			engineerIDs = append(engineerIDs, eng.ID)
		}
		engList, diags := types.ListValueFrom(ctx, types.StringType, engineerIDs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		dvm.Engineers = engList

		state.Devs = append(state.Devs, dvm)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// DevsDataSourceModel maps the data source schema data.
type DevsDataSourceModel struct {
	Devs []devsModel `tfsdk:"devs"`
}

// devsModel maps Devs schema data.
type devsModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Engineers types.List   `tfsdk:"engineers"`
}

// devsInfoModel maps info data
type DevInfoModel struct {
	ID types.String `tfsdk:"id"`
}
