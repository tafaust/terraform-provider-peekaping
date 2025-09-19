// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tafaust/terraform-provider-peekaping/internal/peekaping"
)

var _ datasource.DataSource = &MaintenanceDataSource{}

type MaintenanceDataSource struct {
	client *peekaping.Client
}

func NewMaintenanceDataSource() datasource.DataSource { return &MaintenanceDataSource{} }

type maintenanceDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Title    types.String `tfsdk:"title"`
	Strategy types.String `tfsdk:"strategy"`
	Active   types.Bool   `tfsdk:"active"`
}

func (d *MaintenanceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_maintenance"
}

func (d *MaintenanceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Optional: true, Description: "Maintenance ID"},
			"title":    schema.StringAttribute{Optional: true, Description: "Maintenance title"},
			"strategy": schema.StringAttribute{Computed: true, Description: "Maintenance strategy"},
			"active":   schema.BoolAttribute{Computed: true, Description: "Whether maintenance is active"},
		},
	}
}

func (d *MaintenanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*peekaping.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *peekaping.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *MaintenanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data maintenanceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		m, err := d.client.GetMaintenance(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("lookup by id failed", err.Error())
			return
		}
		data.ID = types.StringValue(m.ID)
		data.Title = types.StringValue(m.Title)
		data.Strategy = types.StringValue(m.Strategy)
		data.Active = types.BoolValue(m.Active)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	list, err := d.client.ListMaintenance(ctx)
	if err != nil {
		resp.Diagnostics.AddError("list maintenance failed", err.Error())
		return
	}
	title := strings.ToLower(data.Title.ValueString())
	for _, m := range list.Items {
		if strings.ToLower(m.Title) == title || (title != "" && strings.Contains(strings.ToLower(m.Title), title)) {
			data.ID = types.StringValue(m.ID)
			data.Title = types.StringValue(m.Title)
			data.Strategy = types.StringValue(m.Strategy)
			data.Active = types.BoolValue(m.Active)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
	}
	resp.Diagnostics.AddError("not found", "no maintenance window matched the given criteria")
}
