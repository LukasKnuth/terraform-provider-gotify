package provider

import (
	"context"
	"fmt"
	"terraform-provider-gotify/provider/internal"

	"github.com/gotify/go-api-client/v2/client/client"
	"github.com/gotify/go-api-client/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	// Ensure implementation satisfies expected interfaces - compilition fails here otherwise.
	_ resource.Resource              = &ClientResource{}
	_ resource.ResourceWithConfigure = &ClientResource{}
)

type ClientResource struct {
	gotify *internal.AuthedGotifyClient
}

func NewClientResource() resource.Resource {
	return &ClientResource{}
}

type ClientResourceModel struct {
	Name types.String `tfsdk:"name"`
	// Read-only after apply
	Id    types.Int64  `tfsdk:"id"`
	Token types.String `tfsdk:"token"`
}

func (r *ClientResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client"
}

func (r *ClientResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A Gotify client is used to fetch messages from the server to display them in a client application.",
		MarkdownDescription: "After applying the client, use the `token` to fetch messages in your client.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "Numerical identifier of this specific client.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name for the client. Will show up in the clients list UI.",
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The Token to both identify the reading client AND authenticate it against the server.",
			},
		},
	}
}

func (r *ClientResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClientResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := client.NewCreateClientParams()
	params.Body = &models.Client{
		Name: data.Name.ValueString(),
	}
	new_client, err := r.gotify.Client.Client.CreateClient(params, r.gotify.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
		return
	}

	data.Id = types.Int64Value(int64(new_client.Payload.ID))
	data.Token = types.StringValue(new_client.Payload.Token)
	data.Name = types.StringValue(new_client.Payload.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClientResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := client.NewGetClientsParams()
	client_list, err := r.gotify.Client.Client.GetClients(params, r.gotify.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
		return
	}

	// Find this application and it's data
	var found *models.Client
	for _, client := range client_list.Payload {
		if client.ID == uint(state.Id.ValueInt64()) {
			found = client
			break
		}
	}
	if found != nil {
		// Update information on state
		state.Name = types.StringValue(found.Name)
		state.Token = types.StringValue(found.Token)

		// Write new information to tf-state
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	} else {
		// The Application is no longer there, remove it and let terraform re-create it later.
		// https://discuss.hashicorp.com/t/how-should-read-signal-that-a-resource-has-vanished-from-the-api-server/40833/2
		resp.State.RemoveResource(ctx)
	}
}

func (r *ClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClientResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := client.NewUpdateClientParams()
	params.ID = data.Id.ValueInt64()
	params.Body = &models.Client{
		Name: data.Name.ValueString(),
	}
	updated_client, err := r.gotify.Client.Client.UpdateClient(params, r.gotify.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
		return
	}

	data.Name = types.StringValue(updated_client.Payload.Name)
	data.Token = types.StringValue(updated_client.Payload.Token)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ClientResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := client.NewDeleteClientParams()
	params.ID = state.Id.ValueInt64()
	_, err := r.gotify.Client.Client.DeleteClient(params, r.gotify.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
		return
	}
}
