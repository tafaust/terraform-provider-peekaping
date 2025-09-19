// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/tafaust/terraform-provider-peekaping/internal/peekaping"
)

var _ resource.Resource = &NotificationResource{}
var _ resource.ResourceWithImportState = &NotificationResource{}

type NotificationResource struct {
	client *peekaping.Client
}

func NewNotificationResource() resource.Resource { return &NotificationResource{} }

// notificationNameValidator validates notification name constraints.
type notificationNameValidator struct{}

func (v notificationNameValidator) Description(_ context.Context) string {
	return "Notification name must be between 1 and 100 characters"
}

func (v notificationNameValidator) MarkdownDescription(_ context.Context) string {
	return "Notification name must be between 1 and 100 characters"
}

func (v notificationNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	name := req.ConfigValue.ValueString()
	if len(name) < 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Name",
			"Notification name must be at least 1 character long",
		)
		return
	}
	if len(name) > 100 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Name",
			"Notification name must be at most 100 characters long",
		)
		return
	}
}

// notificationTypeValidator validates notification type.
type notificationTypeValidator struct{}

func (v notificationTypeValidator) Description(_ context.Context) string {
	return "Notification type must be one of the supported types"
}

func (v notificationTypeValidator) MarkdownDescription(_ context.Context) string {
	return "Notification type must be one of: `smtp`, `telegram`, `webhook`, `slack`, `ntfy`, `pagerduty`, `opsgenie`, `google_chat`, `grafana_oncall`, `signal`, `gotify`, `pushover`, `mattermost`, `matrix`, `discord`, `wecom`, `whatsapp`, `twilio`, `sendgrid`, `pushbullet`, `pagertree`"
}

func (v notificationTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	notificationType := req.ConfigValue.ValueString()
	supportedTypes := []string{
		"smtp", "telegram", "webhook", "slack", "ntfy", "pagerduty", "opsgenie",
		"google_chat", "grafana_oncall", "signal", "gotify", "pushover", "mattermost",
		"matrix", "discord", "wecom", "whatsapp", "twilio", "sendgrid", "pushbullet", "pagertree",
	}

	for _, supportedType := range supportedTypes {
		if notificationType == supportedType {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Notification Type",
		fmt.Sprintf("Notification type '%s' is not supported. Supported types are: %s", notificationType, strings.Join(supportedTypes, ", ")),
	)
}

// notificationConfigValidator validates notification configuration JSON.
type notificationConfigValidator struct{}

func (v notificationConfigValidator) Description(_ context.Context) string {
	return "Notification configuration must be valid JSON"
}

func (v notificationConfigValidator) MarkdownDescription(_ context.Context) string {
	return "Notification configuration must be valid JSON"
}

func (v notificationConfigValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	config := req.ConfigValue.ValueString()
	if config == "" {
		return // Empty config is allowed for some notification types
	}

	// Validate JSON format
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(config), &jsonData); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSON Configuration",
			fmt.Sprintf("Notification configuration must be valid JSON: %s", err.Error()),
		)
		return
	}

	// Basic validation - more specific validation would require the notification type context
	// which is not available in this validator. The API will perform type-specific validation.
}

type notificationResourceModel struct {
	ID        types.String         `tfsdk:"id"`
	Name      types.String         `tfsdk:"name"`
	Type      types.String         `tfsdk:"type"`
	Config    jsontypes.Normalized `tfsdk:"config"`
	Active    types.Bool           `tfsdk:"active"`
	IsDefault types.Bool           `tfsdk:"is_default"`
	CreatedAt types.String         `tfsdk:"created_at"`
	UpdatedAt types.String         `tfsdk:"updated_at"`
}

func (r *NotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification"
}

func (r *NotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Notification channel name",
				Validators: []validator.String{
					notificationNameValidator{},
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Notification type (smtp, webhook, slack, discord, etc.)",
				Validators: []validator.String{
					notificationTypeValidator{},
				},
			},
			"config": schema.StringAttribute{
				Required:    true,
				Description: "Notification configuration (JSON string)",
				CustomType:  jsontypes.NormalizedType{},
				Validators: []validator.String{
					notificationConfigValidator{},
				},
			},
			"active": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the notification channel is active",
			},
			"is_default": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this is the default notification channel",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Last update timestamp",
			},
		},
	}
}

func (r *NotificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = client
}

func (r *NotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// active and is_default are computed by the API, not set by user

	// Log the notification we're creating for debugging
	tflog.Info(ctx, "Creating notification", map[string]interface{}{
		"name":   plan.Name.ValueString(),
		"type":   plan.Type.ValueString(),
		"config": plan.Config.ValueString(),
	})

	in := peekaping.NotificationCreate{
		Name:   plan.Name.ValueString(),
		Type:   plan.Type.ValueString(),
		Config: plan.Config.ValueString(),
	}

	n, err := r.client.CreateNotification(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError("create notification failed", err.Error())
		return
	}
	setModelFromNotification(&plan, n)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	n, err := r.client.GetNotification(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read notification failed", err.Error())
		return
	}

	// Log what the API returned for debugging
	tflog.Info(ctx, "API returned notification", map[string]interface{}{
		"id":         n.ID,
		"name":       n.Name,
		"type":       n.Type,
		"active":     n.Active,
		"is_default": n.IsDefault,
	})

	// Use direct field mapping for Read operations
	setModelFromNotification(&state, n)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to get the ID
	var state notificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upd := peekaping.NotificationUpdate{}
	if !plan.Name.IsNull() {
		v := plan.Name.ValueString()
		upd.Name = &v
	}
	if !plan.Type.IsNull() {
		v := plan.Type.ValueString()
		upd.Type = &v
	}
	if !plan.Config.IsNull() {
		v := plan.Config.ValueString()
		upd.Config = &v
	}
	// active and is_default are computed by the API, not updated by user

	// Use state.ID instead of plan.ID
	n, err := r.client.UpdateNotification(ctx, state.ID.ValueString(), upd)
	if err != nil {
		resp.Diagnostics.AddError("update notification failed", err.Error())
		return
	}
	setModelFromNotification(&plan, n)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteNotification(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("delete notification failed", err.Error())
		return
	}
}

func (r *NotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state notificationResourceModel
	state.ID = types.StringValue(req.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func setModelFromNotification(m *notificationResourceModel, from *peekaping.Notification) {
	// Required fields - always present
	m.ID = types.StringValue(from.ID)
	m.Name = types.StringValue(from.Name)
	m.Type = types.StringValue(from.Type)

	// Config field - normalize JSON for consistency
	if from.Config != "" {
		m.Config = jsontypes.NewNormalizedValue(from.Config)
	} else {
		m.Config = jsontypes.NewNormalizedValue("{}")
	}

	// Active field - always set from API response
	m.Active = types.BoolValue(from.Active)

	// IsDefault field - always set from API response
	m.IsDefault = types.BoolValue(from.IsDefault)

	// Timestamp fields
	if from.CreatedAt != "" {
		m.CreatedAt = types.StringValue(from.CreatedAt)
	} else {
		m.CreatedAt = types.StringNull()
	}
	if from.UpdatedAt != "" {
		m.UpdatedAt = types.StringValue(from.UpdatedAt)
	} else {
		m.UpdatedAt = types.StringNull()
	}
}

// setModelFromNotificationWithState handles field mapping with state comparison to resolve API inconsistencies
