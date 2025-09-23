package provider

import (
	"context"

	"terraform-provider-devops/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DevOpsDataSource{}

func NewDevOpsDataSource() datasource.DataSource { return &DevOpsDataSource{} }

type DevOpsDataSource struct{ client *client.Client }

func (d *DevOpsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devops"
}

func (d *DevOpsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"devops": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Computed: true},
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
				},
			},
		},
	}
}

func (d *DevOpsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil { return }
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", "Provider data was not the expected *client.Client")
		return
	}
	d.client = c
}

func (d *DevOpsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DevOpsDataSourceModel
	items, err := d.client.GetDevOps()
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read DevOps", err.Error())
		return
	}
	for _, v := range items {
		devopsState := devopsModel{ID: types.StringValue(v.ID)}
		for _, o := range v.Ops {
			opsState := opsModel{ID: types.StringValue(o.ID), Name: types.StringValue(o.Name)}
			for _, e := range o.Engineers {
				opsState.Engineers = append(opsState.Engineers, engineerNested{ID: types.StringValue(e.ID), Name: types.StringValue(e.Name), Email: types.StringValue(e.Email)})
			}
			devopsState.Ops = append(devopsState.Ops, opsState)
		}
		for _, dv := range v.Devs {
			devState := devModel{ID: types.StringValue(dv.ID), Name: types.StringValue(dv.Name)}
			for _, e := range dv.Engineers {
				devState.Engineers = append(devState.Engineers, engineerNested{ID: types.StringValue(e.ID), Name: types.StringValue(e.Name), Email: types.StringValue(e.Email)})
			}
			devopsState.Devs = append(devopsState.Devs, devState)
		}
		state.DevOps = append(state.DevOps, devopsState)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type DevOpsDataSourceModel struct { DevOps []devopsModel `tfsdk:"devops"` }

type devopsModel struct {
	ID  types.String `tfsdk:"id"`
	Ops []opsModel   `tfsdk:"ops"`
	Devs []devModel  `tfsdk:"devs"`
}
