package provider

import (
	"context"

	"terraform-provider-devops/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &OpsDataSource{}

func NewOpsDataSource() datasource.DataSource { return &OpsDataSource{} }

type OpsDataSource struct{ client *client.Client }

func (d *OpsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ops"
}

func (d *OpsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ops": schema.ListNestedAttribute{
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

func (d *OpsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil { return }
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", "Provider data was not the expected *client.Client")
		return
	}
	d.client = c
}

func (d *OpsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state OpsDataSourceModel
	items, err := d.client.GetOps()
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Ops", err.Error())
		return
	}
	for _, v := range items {
		opsState := opsModel{
			ID:   types.StringValue(v.ID),
			Name: types.StringValue(v.Name),
		}
		for _, e := range v.Engineers {
			opsState.Engineers = append(opsState.Engineers, engineerNested{
				ID: types.StringValue(e.ID), Name: types.StringValue(e.Name), Email: types.StringValue(e.Email),
			})
		}
		state.Ops = append(state.Ops, opsState)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type OpsDataSourceModel struct { Ops []opsModel `tfsdk:"ops"` }

type opsModel struct {
	ID        types.String    `tfsdk:"id"`
	Name      types.String    `tfsdk:"name"`
	Engineers []engineerNested `tfsdk:"engineers"`
}
