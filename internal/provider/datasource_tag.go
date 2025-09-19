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

var _ datasource.DataSource = &TagDataSource{}

type TagDataSource struct {
	client *peekaping.Client
}

func NewTagDataSource() datasource.DataSource { return &TagDataSource{} }

type tagDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Color       types.String `tfsdk:"color"`
	Description types.String `tfsdk:"description"`
}

func (d *TagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (d *TagDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Optional: true, Description: "Tag ID"},
			"name":        schema.StringAttribute{Optional: true, Description: "Tag name"},
			"color":       schema.StringAttribute{Computed: true, Description: "Tag color"},
			"description": schema.StringAttribute{Computed: true, Description: "Tag description"},
		},
	}
}

func (d *TagDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data tagDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		t, err := d.client.GetTag(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("lookup by id failed", err.Error())
			return
		}
		data.ID = types.StringValue(t.ID)
		data.Name = types.StringValue(t.Name)
		if t.Color != "" {
			data.Color = types.StringValue(t.Color)
		} else {
			data.Color = types.StringNull()
		}
		if t.Description != "" {
			data.Description = types.StringValue(t.Description)
		} else {
			data.Description = types.StringNull()
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	list, err := d.client.ListTags(ctx)
	if err != nil {
		resp.Diagnostics.AddError("list tags failed", err.Error())
		return
	}
	name := strings.ToLower(data.Name.ValueString())
	for _, t := range list.Items {
		if strings.ToLower(t.Name) == name || (name != "" && strings.Contains(strings.ToLower(t.Name), name)) {
			data.ID = types.StringValue(t.ID)
			data.Name = types.StringValue(t.Name)
			if t.Color != "" {
				data.Color = types.StringValue(t.Color)
			} else {
				data.Color = types.StringNull()
			}
			if t.Description != "" {
				data.Description = types.StringValue(t.Description)
			} else {
				data.Description = types.StringNull()
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
	}
	resp.Diagnostics.AddError("not found", "no tag matched the given criteria")
}
