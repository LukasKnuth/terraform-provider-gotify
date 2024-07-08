package main

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GotifyProvider struct {
	version string
}

// Map Terraform HCL schema to Go types
type GotifyProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *GotifyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gotify"
	resp.Version = p.version
}

// What can be configured through HCL for this provider
func (p *GotifyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional: false,
			},
			"username": schema.StringAttribute{
				Optional: false,
			},
			"password": schema.StringAttribute{
				Optional:  false,
				Sensitive: true,
			},
		},
	}
}

// The actual code of taking HCL values and creating a provider instance from them
func (p *GotifyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var model GotifyProviderModel

	// Parse into model and add any errors...
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default to ENV variables but override with explicit config
	endpoint := os.Getenv("GOTIFY_ENDPOINT")
	username := os.Getenv("GOTIFY_USERNAME")
	password := os.Getenv("GOTIFY_PASSWORD")

	if !model.Endpoint.IsNull() {
		endpoint = model.Endpoint.ValueString()
	}
	if !model.Username.IsNull() {
		username = model.Username.ValueString()
	}
	if !model.Password.IsNull() {
		password = model.Password.ValueString()
	}

	// Verify we have values for everything
	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Endpoint configuration",
			"Configure the endpoint to reach the Gotify API, either via the `GOTIFY_ENDPOINT` environment variable, or configuration.",
		)
	}
	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Username configuration",
			"Configure the username to authenticate against the Gotify API, either via `GOTIFY_USERNAME` environment variable, or configuration.",
		)
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Password configuration",
			"Configure the password to authenticate against the Gotify API, either via `GOTIFY_PASSWORD` environment variable, or configuration.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := NewAuthedClient(endpoint, username, password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed while constructing client", err.Error(),
		)
		return
	}

	// Make client available to data/resource
	resp.DataSourceData = client
	resp.ResourceData = client
}

// All Resources this provider offers
func (p *GotifyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// All DataSources (read) this provider offers
func (p *GotifyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// All custom functions this provider offers
func (p *GotifyProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

// Create new instance of this provider
func NewProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GotifyProvider{
			version: version,
		}
	}
}
