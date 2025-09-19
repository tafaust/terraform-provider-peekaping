// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tafaust/terraform-provider-peekaping/internal/peekaping"
)

var _ datasource.DataSource = &ProxyDataSource{}

type ProxyDataSource struct {
	client *peekaping.Client
}

func NewProxyDataSource() datasource.DataSource { return &ProxyDataSource{} }

type proxyDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Host     types.String `tfsdk:"host"`
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Auth     types.Bool   `tfsdk:"auth"`
	Username types.String `tfsdk:"username"`
}

func (d *ProxyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxy"
}

func (d *ProxyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Optional: true, Description: "Proxy ID"},
			"host":     schema.StringAttribute{Optional: true, Description: "Proxy host"},
			"port":     schema.Int64Attribute{Computed: true, Description: "Proxy port"},
			"protocol": schema.StringAttribute{Computed: true, Description: "Proxy protocol"},
			"auth":     schema.BoolAttribute{Computed: true, Description: "Whether authentication is required"},
			"username": schema.StringAttribute{Computed: true, Description: "Username for authentication"},
		},
	}
}

func (d *ProxyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*peekaping.Client)
}

func (d *ProxyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data proxyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		p, err := d.client.GetProxy(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("lookup by id failed", err.Error())
			return
		}
		data.ID = types.StringValue(p.ID)
		data.Host = types.StringValue(p.Host)
		data.Port = types.Int64Value(int64(p.Port))
		data.Protocol = types.StringValue(string(p.Protocol))
		data.Auth = types.BoolValue(p.Auth)
		if p.Username != "" {
			data.Username = types.StringValue(p.Username)
		} else {
			data.Username = types.StringNull()
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	list, err := d.client.ListProxies(ctx)
	if err != nil {
		resp.Diagnostics.AddError("list proxies failed", err.Error())
		return
	}
	host := strings.ToLower(data.Host.ValueString())
	for _, p := range list.Items {
		if strings.ToLower(p.Host) == host || (host != "" && strings.Contains(strings.ToLower(p.Host), host)) {
			data.ID = types.StringValue(p.ID)
			data.Host = types.StringValue(p.Host)
			data.Port = types.Int64Value(int64(p.Port))
			data.Protocol = types.StringValue(string(p.Protocol))
			data.Auth = types.BoolValue(p.Auth)
			if p.Username != "" {
				data.Username = types.StringValue(p.Username)
			} else {
				data.Username = types.StringNull()
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
	}
	resp.Diagnostics.AddError("not found", "no proxy matched the given criteria")
}
