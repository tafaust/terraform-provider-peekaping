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

var _ datasource.DataSource = &NotificationDataSource{}

type NotificationDataSource struct {
	client *peekaping.Client
}

func NewNotificationDataSource() datasource.DataSource { return &NotificationDataSource{} }

type notificationDataSourceModel struct {
	ID        types.String         `tfsdk:"id"`
	Name      types.String         `tfsdk:"name"`
	Type      types.String         `tfsdk:"type"`
	Config    jsontypes.Normalized `tfsdk:"config"`
	Active    types.Bool           `tfsdk:"active"`
	IsDefault types.Bool           `tfsdk:"is_default"`
}

func (d *NotificationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification"
}

func (d *NotificationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Optional: true, Description: "Notification ID"},
			"name": schema.StringAttribute{Optional: true, Description: "Notification name"},
			"type": schema.StringAttribute{Computed: true, Description: "Notification type"},
			"config": schema.StringAttribute{
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "Notification configuration",
			},
			"active":     schema.BoolAttribute{Computed: true, Description: "Whether the notification channel is active"},
			"is_default": schema.BoolAttribute{Computed: true, Description: "Whether this is the default notification channel"},
		},
	}
}

func (d *NotificationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*peekaping.Client)
}

func (d *NotificationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data notificationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		n, err := d.client.GetNotification(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("lookup by id failed", err.Error())
			return
		}
		data.ID = types.StringValue(n.ID)
		data.Name = types.StringValue(n.Name)
		data.Type = types.StringValue(n.Type)
		data.Config = jsontypes.NewNormalizedValue(n.Config)
		data.Active = types.BoolValue(n.Active)
		data.IsDefault = types.BoolValue(n.IsDefault)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	list, err := d.client.ListNotifications(ctx)
	if err != nil {
		resp.Diagnostics.AddError("list notifications failed", err.Error())
		return
	}
	name := strings.ToLower(data.Name.ValueString())
	for _, n := range list.Items {
		if strings.ToLower(n.Name) == name || (name != "" && strings.Contains(strings.ToLower(n.Name), name)) {
			data.ID = types.StringValue(n.ID)
			data.Name = types.StringValue(n.Name)
			data.Type = types.StringValue(n.Type)
			data.Config = jsontypes.NewNormalizedValue(n.Config)
			data.Active = types.BoolValue(n.Active)
			data.IsDefault = types.BoolValue(n.IsDefault)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
	}
	resp.Diagnostics.AddError("not found", "no notification matched the given criteria")
}
