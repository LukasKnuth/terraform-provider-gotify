package provider

import (
	"context"
	"fmt"
	"terraform-provider-gotify/provider/internal"

	"github.com/go-openapi/runtime"
	"github.com/gotify/go-api-client/v2/client/plugin"
	"github.com/gotify/go-api-client/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	// Ensure implementation satisfies expected interfaces - compilition fails here otherwise.
	_ resource.Resource              = &PluginResource{}
	_ resource.ResourceWithConfigure = &PluginResource{}
)

type PluginResource struct {
	gotify *internal.AuthedGotifyClient
}

func NewPluginResource() resource.Resource {
	return &PluginResource{}
}

type PluginResourceModel struct {
	ModulePath types.String `tfsdk:"module_path"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	// Read-only after apply
	Token       types.String `tfsdk:"token"`
	WebhookPath types.String `tfsdk:"webhook_path"`
}

func (r *PluginResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plugin"
}

func (r *PluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Configures a plugin installed on the Gotify server.",
		MarkdownDescription: "The plugin must already be on the server. It must be compatible as well, you can check this manually by navigating to \"Plugins\" in the Web interface.",
		Attributes: map[string]schema.Attribute{
			"module_path": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of this plugin, chosen by the author. Check \"Plugins\" in the Web interface to find this out manually.",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Sets the desired plugin status.",
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The token generated for this plugin. Mainly used for Webhooks.",
			},
			"webhook_path": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Optional:            true,
				Description:         "Generates the webhook base path. If the plugin registers a webhook, this is where it'll be available at.",
				MarkdownDescription: "You are responsible for setting the host/port AND the sub-path the plugin sets itself. Usually, the plugin description has more information, check \"Plugins\" in the Web interface.\n\nFor example, if the full plugin webhook path is `https://localhost:8080/plugin/1/custom/t0k3n/slack_message` then this field will contain `/plugin/1/custom/t0k3n`\n\nNOTE: The path **does** include a leading slash but **not** a trailing slash.",
			},
		},
	}
}

func (r *PluginResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// TODO this is the same everywhere. Make function and reuse
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

func (r *PluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PluginResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 1. Find plugin ID
	found, err := r.findPlugin(data.ModulePath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Could not fetch plugin list", err.Error())
		return
	} else if found == nil {
		resp.Diagnostics.AddError("Could not find plugin %s, is it installed?", data.ModulePath.String())
		return
	}

	// 2. Enable/Disable the plugin
	if found.Enabled != data.Enabled.ValueBool() {
		err = r.applyPluginState(int64(found.ID), data.Enabled.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError("Could not enable/disable plugin", err.Error())
			return
		}
	}

	// Store state info
	data.Token = types.StringValue(found.Token)
	data.WebhookPath = toWebhookPath(found.ID, found.Token)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PluginResource) applyPluginState(id int64, enable bool) error {
	var err error
	if enable {
		params := plugin.NewEnablePluginParams()
		params.ID = id
		_, err = r.gotify.Client.Plugin.EnablePlugin(params, r.gotify.Auth)
	} else {
		params := plugin.NewDisablePluginParams()
		params.ID = id
		_, err = r.gotify.Client.Plugin.DisablePlugin(params, r.gotify.Auth)
	}
	if err == nil {
		return nil
	} else {
		apiErr, ok := err.(*runtime.APIError)
		if ok && apiErr.Code == 400 {
			// Plugin is already enabled/disabled, ignore and continue
			return nil
		} else {
			return err
		}
	}
}

func (r *PluginResource) findPlugin(modulePath string) (*models.PluginConfExternal, error) {
	params := plugin.NewGetPluginsParams()
	plugins, err := r.gotify.Client.Plugin.GetPlugins(params, r.gotify.Auth)
	if err != nil {
		return nil, err
	} else {
		for _, plugin := range plugins.Payload {
			if plugin.ModulePath == modulePath {
				return plugin, nil
			}
		}
		return nil, nil
	}
}

func toWebhookPath(id uint, token string) basetypes.StringValue {
	return types.StringValue(fmt.Sprintf("/plugin/%v/custom/%v", id, token))
}

func (r *PluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PluginResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Find this application and it's data
	found, err := r.findPlugin(state.ModulePath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
	} else if found != nil {
		// Update information on state
		state.Token = types.StringValue(found.Token)
		state.Enabled = types.BoolValue(found.Enabled)
		state.Token = types.StringValue(found.Token)
		state.WebhookPath = toWebhookPath(found.ID, found.Token)

		// Write new information to tf-state
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	} else {
		// The Application is no longer there, remove it and let terraform re-create it later.
		// https://discuss.hashicorp.com/t/how-should-read-signal-that-a-resource-has-vanished-from-the-api-server/40833/2
		resp.State.RemoveResource(ctx)
	}
}

func (r *PluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PluginResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	found, err := r.findPlugin(plan.ModulePath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Gotify API request failed", err.Error())
		return
	} else if found == nil {
		resp.Diagnostics.AddError("Could not find plugin %s, is it installed?", plan.ModulePath.ValueString())
		return
	}

	if !plan.Enabled.Equal(state.Enabled) {
		err := r.applyPluginState(int64(found.ID), plan.Enabled.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError("Gotify API Request failed", err.Error())
			return
		}
	}

	plan.Token = types.StringValue(found.Token)
	plan.WebhookPath = toWebhookPath(found.ID, found.Token)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Plugins can not be deleted",
		"This will only remove the plugin from your Terraform State. It will not disable the plugin, nor will it uninstall it. If you want to disable the plugin, add the `gotify_plugin` resource with `enabled = false`. If you want to uninstall the plugin, remove the so file from the Gotify Plugin path.",
	)
}
