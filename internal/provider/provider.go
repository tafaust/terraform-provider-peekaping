// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/tafaust/terraform-provider-peekaping/internal/peekaping"
)

var _ provider.Provider = &PeekapingProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PeekapingProvider{
			version: version,
		}
	}
}

type PeekapingProvider struct {
	client  *peekaping.Client
	version string
}

type providerModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Email    types.String `tfsdk:"email"`
	Password types.String `tfsdk:"password"`
	Token    types.String `tfsdk:"token"`
}

func (p *PeekapingProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "peekaping"
	resp.Version = p.version
}

func (p *PeekapingProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "Base URL of the Peekaping server (e.g. http://localhost:8034)",
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Description: "Email for login.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password for login.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "2FA token for login (if 2FA is enabled).",
			},
		},
	}
}

func (p *PeekapingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("PEEKAPING_ENDPOINT")
	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}
	if endpoint == "" {
		endpoint = "http://localhost:8034"
	}

	email := os.Getenv("PEEKAPING_EMAIL")
	password := os.Getenv("PEEKAPING_PASSWORD")
	token := os.Getenv("PEEKAPING_TOKEN")
	if !config.Email.IsNull() {
		email = config.Email.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	client := peekaping.New(endpoint, peekaping.WithCredentials(email, password), peekaping.WithToken(token))
	if email != "" && password != "" {
		if err := client.Login(ctx); err != nil {
			resp.Diagnostics.AddError("login failed", err.Error())
			return
		}
		tflog.Info(ctx, "logged in to peekaping")
	}

	p.client = client
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PeekapingProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMonitorResource,
		NewNotificationResource,
		NewTagResource,
		NewMaintenanceResource,
		NewStatusPageResource,
		NewProxyResource,
	}
}

func (p *PeekapingProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMonitorDataSource,
		NewNotificationDataSource,
		NewTagDataSource,
		NewMaintenanceDataSource,
		NewStatusPageDataSource,
		NewProxyDataSource,
	}
}
