package provider

import (
	"context"
	"fmt"
	"terraform-provider-gotify/provider/internal"

	"github.com/gotify/go-api-client/v2/client/application"
	"github.com/gotify/go-api-client/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	// Ensure implementation satisfies expected interfaces - compilition fails here otherwise.
	_ resource.Resource              = &ApplicationResource{}
	_ resource.ResourceWithConfigure = &ApplicationResource{}
)

type ApplicationResource struct {
	gotify *internal.AuthedGotifyClient
}

func NewApplicationResource() resource.Resource {
	return &ApplicationResource{}
}

type ApplicationResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	// Read-only after apply
	Id    types.Int64  `tfsdk:"id"`
	Token types.String `tfsdk:"token"`
}

func (r *ApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *ApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "An application is used to publish messages to Gotify from a specific App. Each app receives it's own channel where all its messages end up in.",
		MarkdownDescription: "After applying the resource, use the `token` to send messages from your application.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "Numeric identifier of this specific Application.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the application sending messages. This is also the message channel name in the UI.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the application sending messages. Will show up in the Apps list.",
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The Token to both identify the sending application AND authenticate it against the server.",
			},
		},
	}
}

func (r *ApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// IMPORTANT: This method is called MULTIPLE times. An initial call might not have configured the Provider yet, so we need
		// to handle this gracefully. It will eventually be called with a configured provider.
		return
	}

	client, ok := req.ProviderData.(*internal.AuthedGotifyClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *AuthedGotifyClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.gotify = client
}

func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApplicationResourceModel

	// Read planned initial values from the Terraform Plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Send the request
	params := application.NewCreateAppParams()
	params.Body = &models.Application{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}
	app, err := r.gotify.Client.Application.CreateApp(params, r.gotify.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
		return
	}

	// Update model with computed information
	data.Id = types.Int64Value(int64(app.Payload.ID))
	data.Token = types.StringValue(app.Payload.Token)
	data.Name = types.StringValue(app.Payload.Name)
	data.Description = types.StringValue(app.Payload.Description)

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
	app_list, err := r.gotify.Client.Application.GetApps(params, r.gotify.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
		return
	}

	// Find this application and it's data
	var found *models.Application
	for _, app := range app_list.Payload {
		if app.ID == uint(state.Id.ValueInt64()) {
			found = app
			break
		}
	}
	if found != nil {
		// Update information on state
		state.Name = types.StringValue(found.Name)
		state.Description = types.StringValue(found.Description)
		state.Token = types.StringValue(found.Token)

		// Write new information to tf-state
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	} else {
		// The Application is no longer there, remove it and let terraform re-create it later.
		// https://discuss.hashicorp.com/t/how-should-read-signal-that-a-resource-has-vanished-from-the-api-server/40833/2
		resp.State.RemoveResource(ctx)
	}
}

func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApplicationResourceModel

	// Read planned changes from the Terraform Plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create API request
	params := application.NewUpdateApplicationParams()
	params.ID = data.Id.ValueInt64()
	params.Body = &models.Application{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}
	app, err := r.gotify.Client.Application.UpdateApplication(params, r.gotify.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
		return
	}

	// Update model with updated information
	data.Name = types.StringValue(app.Payload.Name)
	data.Description = types.StringValue(app.Payload.Description)
	data.Token = types.StringValue(app.Payload.Token)

	// Write new data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApplicationResourceModel

	// Read data from state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Send DELETE request
	params := application.NewDeleteAppParams()
	params.ID = state.Id.ValueInt64()
	_, err := r.gotify.Client.Application.DeleteApp(params, r.gotify.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
		return
	}
}
