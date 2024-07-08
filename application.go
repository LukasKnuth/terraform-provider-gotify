package main

import (
	"context"
	"fmt"

	"github.com/gotify/go-api-client/v2/client/application"
	"github.com/gotify/go-api-client/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ApplicationResource struct {
	client *AuthedGotifyClient
}

func NewApplicationResource() resource.Resource {
	return &ApplicationResource{}
}

type ApplicationResourceModel struct {
	name        types.String `tfsdk:"name"`
	description types.String `tfsdk:"description"`
	// Read-only after apply
	id    types.Int64  `tfsdk:"id"`
	token types.String `tfsdk:"token"`
}

func (r *ApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *ApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: false,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"token": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *ApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := req.ProviderData.(*AuthedGotifyClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApplicationResourceModel

	// Read Plan data first
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Send the request
	params := application.NewCreateAppParams()
	params.Body = &models.Application{
		Name:        data.name.ValueString(),
		Description: data.description.ValueString(),
	}
	app, err := r.client.client.Application.CreateApp(params, r.client.auth)
	if err != nil {
		resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
		return
	}

	// Update model with computed information
	data.id = types.Int64Value(int64(app.Payload.ID))
	data.token = types.StringValue(app.Payload.Token)

	// Write new data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApplicationResourceModel

	// Read current state from Terraform state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read all apps
	params := application.NewGetAppsParams()
	app_list, err := r.client.client.Application.GetApps(params, r.client.auth)
	if err != nil {
		resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
		return
	}

	// Find this application and it's data
	var found *models.Application
	for _, app := range app_list.Payload {
		if app.ID == uint(state.id.ValueInt64()) {
			found = app
			break
		}
	}
	if found != nil {
		// Update information on state
		state.name = types.StringValue(found.Name)
		state.description = types.StringValue(found.Description)
		state.token = types.StringValue(found.Token)

		// Write new information to tf-state
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	} else {
		// The Application is no longer there, remove it and let terraform re-create it later.
		// https://discuss.hashicorp.com/t/how-should-read-signal-that-a-resource-has-vanished-from-the-api-server/40833/2
		resp.State.RemoveResource(ctx)
	}
}

func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
