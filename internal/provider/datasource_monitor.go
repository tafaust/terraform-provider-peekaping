// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tafaust/terraform-provider-peekaping/internal/peekaping"
)

var _ datasource.DataSource = &MonitorDataSource{}

type MonitorDataSource struct {
	client *peekaping.Client
}

func NewMonitorDataSource() datasource.DataSource { return &MonitorDataSource{} }

type monitorDataSourceModel struct {
	ID     types.String         `tfsdk:"id"`
	Name   types.String         `tfsdk:"name"`
	Type   types.String         `tfsdk:"type"`
	Config jsontypes.Normalized `tfsdk:"config"`
}

func (d *MonitorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (d *MonitorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Optional: true, Description: "Monitor ID"},
			"name": schema.StringAttribute{Optional: true, Description: "Monitor name"},
			"type": schema.StringAttribute{Computed: true, Description: "Monitor type"},
			"config": schema.StringAttribute{
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "Monitor configuration",
			},
		},
	}
}

func (d *MonitorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*peekaping.Client)
}

func (d *MonitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data monitorDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		m, err := d.client.GetMonitor(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("lookup by id failed", err.Error())
			return
		}
		data.ID = types.StringValue(m.ID)
		data.Name = types.StringValue(m.Name)
		data.Type = types.StringValue(string(m.Type))
		data.Config = jsontypes.NewNormalizedValue(m.Config)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	list, err := d.client.ListMonitors(ctx)
	if err != nil {
		resp.Diagnostics.AddError("list monitors failed", err.Error())
		return
	}
	name := strings.ToLower(data.Name.ValueString())
	for _, m := range list.Items {
		if strings.ToLower(m.Name) == name || (name != "" && strings.Contains(strings.ToLower(m.Name), name)) {
			data.ID = types.StringValue(m.ID)
			data.Name = types.StringValue(m.Name)
			data.Type = types.StringValue(string(m.Type))
			data.Config = jsontypes.NewNormalizedValue(m.Config)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
	}
	resp.Diagnostics.AddError("not found", "no monitor matched the given criteria")
}
