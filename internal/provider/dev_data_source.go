package provider

import (
	"context"

	"terraform-provider-devops/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DevDataSource{}

func NewDevDataSource() datasource.DataSource { return &DevDataSource{} }

type DevDataSource struct{ client *client.Client }

func (d *DevDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev"
}

func (d *DevDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"devs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
						"engineers": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{Computed: true},
									"name": schema.StringAttribute{Computed: true},
									"email": schema.StringAttribute{Computed: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *DevDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil { return }
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", "Provider data was not the expected *client.Client")
		return
	}
	d.client = c
}

func (d *DevDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DevDataSourceModel
	items, err := d.client.GetDevs()
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Devs", err.Error())
		return
	}
	for _, v := range items {
		devState := devModel{
			ID:   types.StringValue(v.ID),
			Name: types.StringValue(v.Name),
		}
		for _, e := range v.Engineers {
			devState.Engineers = append(devState.Engineers, engineerNested{
				ID: types.StringValue(e.ID), Name: types.StringValue(e.Name), Email: types.StringValue(e.Email),
			})
		}
		state.Devs = append(state.Devs, devState)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type DevDataSourceModel struct { Devs []devModel `tfsdk:"devs"` }

type devModel struct {
	ID        types.String    `tfsdk:"id"`
	Name      types.String    `tfsdk:"name"`
	Engineers []engineerNested `tfsdk:"engineers"`
}

type engineerNested struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}
