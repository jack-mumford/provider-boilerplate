package devs

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-devops/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &DevsResource{}
	_ resource.ResourceWithConfigure = &DevsResource{}
)

// NewDevsResource is a helper function to simplify the provider implementation.
func NewDevsResource() resource.Resource {
	return &DevsResource{}
}

// DevsResource is the resource implementation.
type DevsResource struct {
	client *client.Client
}

// Metadata returns the resource type name.
func (r *DevsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev"
}

// Schema defines the schema for the resource.
func (r *DevsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"engineers": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *DevsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DevsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Getting the engineer IDs from the plan
	var EngineerID []string
	diags = plan.Engineers.ElementsAs(ctx, &EngineerID, false)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	engineers := make([]client.Engineer, len(EngineerID))
	for i, EngineerID := range EngineerID {
		engineers[i] = client.Engineer{ID: EngineerID}
	}

	dev := client.Dev{
		Name:      plan.Name.ValueString(),
		Engineers: engineers,
	}

	createdDev, err := r.client.CreateDev(dev)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Dev",
			"Could not create Dev, unexpected error: "+err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(createdDev.ID)
	plan.Name = types.StringValue(createdDev.Name)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *DevsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DevsResourceModel

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dev, err := r.client.GetDev(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Dev",
			"Could not read Dev: "+err.Error(),
		)

		return
	}

	if dev == nil {
		// Treat as not found; remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(dev.ID)
	state.Name = types.StringValue(dev.Name)

	// Set engineers list (IDs)
	engineerIDs := make([]string, 0, len(dev.Engineers))
	for _, eng := range dev.Engineers {
		engineerIDs = append(engineerIDs, eng.ID)
	}
	engList, ldiags := types.ListValueFrom(ctx, types.StringType, engineerIDs)
	resp.Diagnostics.Append(ldiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Engineers = engList

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *DevsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan DevsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var reqDev = client.Dev{
		Name: plan.Name.ValueString(),
	}
	reqDev.ID = plan.ID.ValueString()

	// Update existing Dev
	_, err := r.client.UpdateDev(plan.ID.ValueString(), reqDev)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Dev",
			"Could not update Dev, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated Dev after update
	dev, err := r.client.GetDev(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Dev",
			"Could not read Dev ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.ID = types.StringValue(dev.ID)
	plan.Name = types.StringValue(dev.Name)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DevsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state DevsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Dev
	err := r.client.DeleteDev(state.ID.ValueString())
	if err != nil {
		// If backend returns 404, treat as already deleted
		if strings.Contains(err.Error(), "status: 404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Dev: "+state.ID.ValueString(),
			"Could not delete Dev, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *DevsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// devResourceModel maps the resource schema data.
type DevsResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Engineers types.List   `tfsdk:"engineers"`
}
