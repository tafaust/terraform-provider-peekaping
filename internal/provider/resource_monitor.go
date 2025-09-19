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

var _ resource.Resource = &MonitorResource{}
var _ resource.ResourceWithImportState = &MonitorResource{}

// monitorTypeValidator validates that the monitor type is supported.
type monitorTypeValidator struct{}

func (v monitorTypeValidator) Description(_ context.Context) string {
	return "Monitor type must be one of the supported types"
}

func (v monitorTypeValidator) MarkdownDescription(_ context.Context) string {
	return "Monitor type must be one of: `http`, `http-keyword`, `http-json-query`, `push`, `tcp`, `ping`, `dns`, `docker`, `grpc-keyword`, `snmp`, `mongodb`, `mysql`, `postgres`, `sqlserver`, `redis`, `mqtt`, `rabbitmq`, `kafka-producer`"
}

func (v monitorTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	monitorType := req.ConfigValue.ValueString()
	supportedTypes := []string{
		"http", "http-keyword", "http-json-query", "push", "tcp", "ping", "dns", "docker",
		"grpc-keyword", "snmp", "mongodb", "mysql", "postgres", "sqlserver", "redis",
		"mqtt", "rabbitmq", "kafka-producer",
	}

	for _, supportedType := range supportedTypes {
		if monitorType == supportedType {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Monitor Type",
		fmt.Sprintf("Monitor type '%s' is not supported. Supported types are: %s", monitorType, strings.Join(supportedTypes, ", ")),
	)
}

// monitorConfigValidator validates monitor configuration based on type.
type monitorConfigValidator struct{}

func (v monitorConfigValidator) Description(_ context.Context) string {
	return "Monitor configuration must be valid JSON and appropriate for the monitor type"
}

func (v monitorConfigValidator) MarkdownDescription(_ context.Context) string {
	return "Monitor configuration must be valid JSON and appropriate for the monitor type"
}

func (v monitorConfigValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	config := req.ConfigValue.ValueString()
	if config == "" {
		return // Empty config is allowed for some monitor types
	}

	// Validate JSON format
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(config), &jsonData); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSON Configuration",
			fmt.Sprintf("Monitor configuration must be valid JSON: %s", err.Error()),
		)
		return
	}

	// Basic validation - more specific validation would require the monitor type context
	// which is not available in this validator. The API will perform type-specific validation.
}

// monitorNameValidator validates monitor name constraints.
type monitorNameValidator struct{}

func (v monitorNameValidator) Description(_ context.Context) string {
	return "Monitor name must be between 3 and 100 characters"
}

func (v monitorNameValidator) MarkdownDescription(_ context.Context) string {
	return "Monitor name must be between 3 and 100 characters"
}

func (v monitorNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	name := req.ConfigValue.ValueString()
	if len(name) < 3 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Name",
			"Monitor name must be at least 3 characters long",
		)
		return
	}
	if len(name) > 100 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Name",
			"Monitor name must be at most 100 characters long",
		)
		return
	}
}

// monitorTimeoutValidator validates timeout constraints.
type monitorTimeoutValidator struct{}

func (v monitorTimeoutValidator) Description(_ context.Context) string {
	return "Monitor timeout must be at least 16 seconds and less than 80% of interval"
}

func (v monitorTimeoutValidator) MarkdownDescription(_ context.Context) string {
	return "Monitor timeout must be at least 16 seconds and less than 80% of interval"
}

func (v monitorTimeoutValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	timeout := req.ConfigValue.ValueInt64()
	if timeout < 16 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Timeout",
			"Monitor timeout must be at least 16 seconds",
		)
		return
	}

	// Note: The constraint that timeout must be less than 80% of interval
	// would require access to the interval value, which is not available in this validator.
	// This validation is handled by the API.
}

// monitorIntervalValidator validates interval constraints.
type monitorIntervalValidator struct{}

func (v monitorIntervalValidator) Description(_ context.Context) string {
	return "Monitor interval must be at least 20 seconds"
}

func (v monitorIntervalValidator) MarkdownDescription(_ context.Context) string {
	return "Monitor interval must be at least 20 seconds"
}

func (v monitorIntervalValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	interval := req.ConfigValue.ValueInt64()
	if interval < 20 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Interval",
			"Monitor interval must be at least 20 seconds",
		)
		return
	}
}

// monitorRetryIntervalValidator validates retry interval constraints.
type monitorRetryIntervalValidator struct{}

func (v monitorRetryIntervalValidator) Description(_ context.Context) string {
	return "Monitor retry interval must be at least 20 seconds"
}

func (v monitorRetryIntervalValidator) MarkdownDescription(_ context.Context) string {
	return "Monitor retry interval must be at least 20 seconds"
}

func (v monitorRetryIntervalValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	retryInterval := req.ConfigValue.ValueInt64()
	if retryInterval < 20 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Retry Interval",
			"Monitor retry interval must be at least 20 seconds",
		)
		return
	}
}

// monitorMaxRetriesValidator validates max retries constraints.
type monitorMaxRetriesValidator struct{}

func (v monitorMaxRetriesValidator) Description(_ context.Context) string {
	return "Monitor max retries must be at least 0"
}

func (v monitorMaxRetriesValidator) MarkdownDescription(_ context.Context) string {
	return "Monitor max retries must be at least 0"
}

func (v monitorMaxRetriesValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	maxRetries := req.ConfigValue.ValueInt64()
	if maxRetries < 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Max Retries",
			"Monitor max retries must be at least 0",
		)
		return
	}
}

// monitorResendIntervalValidator validates resend interval constraints.
type monitorResendIntervalValidator struct{}

func (v monitorResendIntervalValidator) Description(_ context.Context) string {
	return "Monitor resend interval must be at least 0"
}

func (v monitorResendIntervalValidator) MarkdownDescription(_ context.Context) string {
	return "Monitor resend interval must be at least 0"
}

func (v monitorResendIntervalValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	resendInterval := req.ConfigValue.ValueInt64()
	if resendInterval < 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Resend Interval",
			"Monitor resend interval must be at least 0",
		)
		return
	}
}

type MonitorResource struct {
	client *peekaping.Client
}

func NewMonitorResource() resource.Resource { return &MonitorResource{} }

type monitorResourceModel struct {
	ID              types.String         `tfsdk:"id"`
	Name            types.String         `tfsdk:"name"`
	Type            types.String         `tfsdk:"type"`
	Config          jsontypes.Normalized `tfsdk:"config"`
	Interval        types.Int64          `tfsdk:"interval"`
	Active          types.Bool           `tfsdk:"active"`
	Timeout         types.Int64          `tfsdk:"timeout"`
	MaxRetries      types.Int64          `tfsdk:"max_retries"`
	RetryInterval   types.Int64          `tfsdk:"retry_interval"`
	ResendInterval  types.Int64          `tfsdk:"resend_interval"`
	ProxyID         types.String         `tfsdk:"proxy_id"`
	PushToken       types.String         `tfsdk:"push_token"`
	NotificationIDs []types.String       `tfsdk:"notification_ids"`
	TagIDs          []types.String       `tfsdk:"tag_ids"`
	Status          types.Int64          `tfsdk:"status"`
	CreatedAt       types.String         `tfsdk:"created_at"`
	UpdatedAt       types.String         `tfsdk:"updated_at"`
}

func (r *MonitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *MonitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description: "Monitor name",
				Validators: []validator.String{
					monitorNameValidator{},
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Monitor type (http, tcp, ping, dns, push, grpc, etc.)",
				Validators: []validator.String{
					monitorTypeValidator{},
				},
			},
			"config": schema.StringAttribute{
				Required:    true,
				Description: "Monitor configuration (URL for http, host:port for tcp, etc.)",
				CustomType:  jsontypes.NormalizedType{},
				Validators: []validator.String{
					monitorConfigValidator{},
				},
			},
			"interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Monitor interval in seconds (minimum 20)",
				Validators: []validator.Int64{
					monitorIntervalValidator{},
				},
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the monitor is active",
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Monitor timeout in seconds (minimum 16)",
				Validators: []validator.Int64{
					monitorTimeoutValidator{},
				},
			},
			"max_retries": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum retries before marking as down",
				Validators: []validator.Int64{
					monitorMaxRetriesValidator{},
				},
			},
			"retry_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Retry interval in seconds (minimum 20)",
				Validators: []validator.Int64{
					monitorRetryIntervalValidator{},
				},
			},
			"resend_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Resend notification if down X times consecutively",
				Validators: []validator.Int64{
					monitorResendIntervalValidator{},
				},
			},
			"proxy_id": schema.StringAttribute{
				Optional:    true,
				Description: "Proxy ID for the monitor",
			},
			"push_token": schema.StringAttribute{
				Optional:    true,
				Description: "Push token for push monitors",
			},
			"notification_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of notification channel IDs",
			},
			"tag_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of tag IDs",
			},
			"status": schema.Int64Attribute{
				Computed:    true,
				Description: "Monitor status (0=down, 1=up, 2=pending, 3=maintenance)",
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

func (r *MonitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan monitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	notificationIDs := toStrSlice(plan.NotificationIDs)
	if notificationIDs == nil {
		notificationIDs = []string{}
	}

	// Log the config we're sending for debugging
	tflog.Info(ctx, "Creating monitor", map[string]interface{}{
		"name":    plan.Name.ValueString(),
		"type":    plan.Type.ValueString(),
		"config":  plan.Config.ValueString(),
		"tag_ids": toStrSlice(plan.TagIDs),
	})

	// Handle active field - use plan value or default to true
	active := true // Default value
	if !plan.Active.IsNull() {
		active = plan.Active.ValueBool()
	}

	in := peekaping.MonitorCreate{
		Name:            plan.Name.ValueString(),
		Type:            peekaping.MonitorType(plan.Type.ValueString()),
		Config:          plan.Config.ValueString(), // jsontypes.Normalized handles normalization automatically
		Interval:        plan.Interval.ValueInt64(),
		Active:          active,
		NotificationIDs: notificationIDs,         // Always send, even if empty (API requires it)
		TagIDs:          toStrSlice(plan.TagIDs), // Always send, even if empty (API requires it)
	}
	if !plan.Timeout.IsNull() {
		in.Timeout = plan.Timeout.ValueInt64()
	}
	if !plan.MaxRetries.IsNull() {
		in.MaxRetries = plan.MaxRetries.ValueInt64()
	}
	if !plan.RetryInterval.IsNull() {
		in.RetryInterval = plan.RetryInterval.ValueInt64()
	}
	if !plan.ResendInterval.IsNull() {
		in.ResendInterval = plan.ResendInterval.ValueInt64()
	}
	if !plan.ProxyID.IsNull() {
		in.ProxyID = plan.ProxyID.ValueString()
	}
	if !plan.PushToken.IsNull() {
		in.PushToken = plan.PushToken.ValueString()
	}

	m, err := r.client.CreateMonitor(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError("create monitor failed", err.Error())
		return
	}
	setModelFromMonitor(ctx, &plan, m)

	// Preserve the plan's active value to maintain Terraform state consistency
	// The API may return different defaults than what the plan specifies
	if !plan.Active.IsNull() {
		plan.Active = types.BoolValue(active)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state monitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	m, err := r.client.GetMonitor(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read monitor failed", err.Error())
		return
	}

	// Log what the API returned for debugging
	tflog.Info(ctx, "API returned monitor", map[string]interface{}{
		"id":               m.ID,
		"name":             m.Name,
		"tag_ids":          m.TagIDs,
		"notification_ids": m.NotificationIDs,
		"active":           m.Active,
	})

	// Use regular field mapping but don't touch tag_ids and notification_ids
	// since the API doesn't return these fields and we want to preserve current state
	setModelFromMonitor(ctx, &state, m)

	// Don't modify tag_ids and notification_ids - let Terraform preserve them from current state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan monitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to get the ID
	var state monitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upd := peekaping.MonitorUpdate{
		NotificationIDs: toStrSlice(plan.NotificationIDs), // Always send, even if empty (API requires it)
		TagIDs:          toStrSlice(plan.TagIDs),          // Always send, even if empty (API requires it)
	}
	if !plan.Name.IsNull() {
		v := plan.Name.ValueString()
		upd.Name = &v
	}
	if !plan.Type.IsNull() {
		v := plan.Type.ValueString()
		mt := peekaping.MonitorType(v)
		upd.Type = &mt
	}
	if !plan.Config.IsNull() {
		config := plan.Config.ValueString()
		upd.Config = &config
	}
	if !plan.Interval.IsNull() {
		v := plan.Interval.ValueInt64()
		upd.Interval = &v
	}
	if !plan.Active.IsNull() {
		v := plan.Active.ValueBool()
		upd.Active = &v
	}
	if !plan.Timeout.IsNull() {
		v := plan.Timeout.ValueInt64()
		upd.Timeout = &v
	}
	if !plan.MaxRetries.IsNull() {
		v := plan.MaxRetries.ValueInt64()
		upd.MaxRetries = &v
	}
	if !plan.RetryInterval.IsNull() {
		v := plan.RetryInterval.ValueInt64()
		upd.RetryInterval = &v
	}
	if !plan.ResendInterval.IsNull() {
		v := plan.ResendInterval.ValueInt64()
		upd.ResendInterval = &v
	}
	if !plan.ProxyID.IsNull() {
		v := plan.ProxyID.ValueString()
		upd.ProxyID = &v
	}
	if !plan.PushToken.IsNull() {
		v := plan.PushToken.ValueString()
		upd.PushToken = &v
	}

	// Use state.ID instead of plan.ID
	_, err := r.client.UpdateMonitor(ctx, state.ID.ValueString(), upd)
	if err != nil {
		resp.Diagnostics.AddError("update monitor failed", err.Error())
		return
	}

	// Fetch the monitor again to get the current status and other computed fields
	// The update response might not reflect the current monitoring status
	fullMonitor, err := r.client.GetMonitor(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to fetch updated monitor", err.Error())
		return
	}

	// For computed fields, populate the plan with current state values before setting from API
	// This prevents Terraform from seeing computed field changes as inconsistencies
	// Note: We don't set Status here as it can legitimately change during updates
	// Note: We don't set CreatedAt/UpdatedAt here as they can change during updates

	setModelFromMonitorWithState(&plan, fullMonitor, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state monitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteMonitor(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("delete monitor failed", err.Error())
		return
	}
}

func (r *MonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state monitorResourceModel
	state.ID = types.StringValue(req.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func toStrSlice(xs []types.String) []string {
	out := make([]string, 0, len(xs))
	for _, s := range xs {
		if !s.IsNull() {
			out = append(out, s.ValueString())
		}
	}
	return out
}

func setModelFromMonitor(ctx context.Context, m *monitorResourceModel, from *peekaping.Monitor) {
	// Required fields - always present
	m.ID = types.StringValue(from.ID)
	m.Name = types.StringValue(from.Name)
	m.Type = types.StringValue(string(from.Type))

	// Config field - use jsontypes.Normalized for automatic JSON normalization
	if from.Config != "" {
		m.Config = jsontypes.NewNormalizedValue(from.Config)
	} else {
		m.Config = jsontypes.NewNormalizedValue("{}")
	}

	// Optional numeric fields with proper null handling
	if from.Interval > 0 {
		m.Interval = types.Int64Value(from.Interval)
	} else {
		m.Interval = types.Int64Null()
	}

	// Active field - preserve existing value (set by plan), don't override with API response
	// The API may return different defaults than what the plan specifies

	if from.Timeout > 0 {
		m.Timeout = types.Int64Value(from.Timeout)
	} else {
		m.Timeout = types.Int64Null()
	}

	if from.MaxRetries > 0 {
		m.MaxRetries = types.Int64Value(from.MaxRetries)
	} else {
		m.MaxRetries = types.Int64Null()
	}

	if from.RetryInterval > 0 {
		m.RetryInterval = types.Int64Value(from.RetryInterval)
	} else {
		m.RetryInterval = types.Int64Null()
	}

	if from.ResendInterval > 0 {
		m.ResendInterval = types.Int64Value(from.ResendInterval)
	} else {
		m.ResendInterval = types.Int64Null()
	}

	// Optional string fields
	if from.ProxyID != "" {
		m.ProxyID = types.StringValue(from.ProxyID)
	} else {
		m.ProxyID = types.StringNull()
	}

	if from.PushToken != "" {
		m.PushToken = types.StringValue(from.PushToken)
	} else {
		m.PushToken = types.StringNull()
	}

	// Status field - handle API bug where it returns wrong status
	// The API sometimes returns status=1 instead of the correct value
	// We'll set it to 0 (down) as a safe default, but this should be handled by state comparison
	m.Status = types.Int64Value(int64(from.Status))

	// Timestamp fields - handle API bug where it returns zero time
	tflog.Debug(ctx, "Setting timestamps from API response", map[string]interface{}{
		"created_at": from.CreatedAt,
		"updated_at": from.UpdatedAt,
	})

	// Handle API bug where it returns zero time instead of actual timestamp
	if from.CreatedAt == "" || from.CreatedAt == "0001-01-01T00:00:00Z" {
		// Use null value to indicate unknown timestamp
		m.CreatedAt = types.StringNull()
	} else {
		m.CreatedAt = types.StringValue(from.CreatedAt)
	}

	if from.UpdatedAt == "" || from.UpdatedAt == "0001-01-01T00:00:00Z" {
		// Use null value to indicate unknown timestamp
		m.UpdatedAt = types.StringNull()
	} else {
		m.UpdatedAt = types.StringValue(from.UpdatedAt)
	}

	// Handle TagIDs and NotificationIDs - preserve plan values, don't override with API response
	// The API response may not match the plan, so we preserve the plan's values
	// This ensures Terraform state consistency
}

// setModelFromMonitorWithState handles field mapping with state comparison to resolve API inconsistencies.
func setModelFromMonitorWithState(m *monitorResourceModel, from *peekaping.Monitor, currentState *monitorResourceModel) {
	// Required fields - always present
	m.ID = types.StringValue(from.ID)
	m.Name = types.StringValue(from.Name)
	m.Type = types.StringValue(string(from.Type))

	// Config field - use jsontypes.Normalized for automatic JSON normalization
	if from.Config != "" {
		m.Config = jsontypes.NewNormalizedValue(from.Config)
	} else {
		m.Config = jsontypes.NewNormalizedValue("{}")
	}

	// Optional numeric fields with proper null handling
	if from.Interval > 0 {
		m.Interval = types.Int64Value(from.Interval)
	} else {
		m.Interval = types.Int64Null()
	}

	// Active field - use API value directly as this is user-configurable
	// The state-based logic should only apply to computed fields that the API returns incorrectly
	m.Active = types.BoolValue(from.Active)

	if from.Timeout > 0 {
		m.Timeout = types.Int64Value(from.Timeout)
	} else {
		m.Timeout = types.Int64Null()
	}

	if from.MaxRetries > 0 {
		m.MaxRetries = types.Int64Value(from.MaxRetries)
	} else {
		m.MaxRetries = types.Int64Null()
	}

	if from.RetryInterval > 0 {
		m.RetryInterval = types.Int64Value(from.RetryInterval)
	} else {
		m.RetryInterval = types.Int64Null()
	}

	if from.ResendInterval > 0 {
		m.ResendInterval = types.Int64Value(from.ResendInterval)
	} else {
		m.ResendInterval = types.Int64Null()
	}

	// Optional string fields
	if from.ProxyID != "" {
		m.ProxyID = types.StringValue(from.ProxyID)
	} else {
		m.ProxyID = types.StringNull()
	}

	if from.PushToken != "" {
		m.PushToken = types.StringValue(from.PushToken)
	} else {
		m.PushToken = types.StringNull()
	}

	// Status field - always use API value as this is a computed field
	// Valid status values: 0=down, 1=up, 2=pending, 3=maintenance
	// Status is computed and should always reflect the current API state
	apiStatus := int64(from.Status)
	if apiStatus < 0 || apiStatus > 3 {
		// Invalid status from API, default to down
		apiStatus = 0
	}
	m.Status = types.Int64Value(apiStatus)

	// Timestamp fields - use state as ground truth when API returns invalid values
	tflog.Debug(context.Background(), "Setting timestamps from API response in setModelFromMonitorWithState", map[string]interface{}{
		"created_at": from.CreatedAt,
		"updated_at": from.UpdatedAt,
	})

	// Handle created_at - use state as ground truth when API returns invalid values
	if (from.CreatedAt == "" || from.CreatedAt == "0001-01-01T00:00:00Z") && currentState != nil && !currentState.CreatedAt.IsNull() {
		// API returned invalid timestamp, preserve state value
		m.CreatedAt = currentState.CreatedAt
	} else if from.CreatedAt == "" || from.CreatedAt == "0001-01-01T00:00:00Z" {
		// No state value available, use null
		m.CreatedAt = types.StringNull()
	} else {
		// Use API value
		m.CreatedAt = types.StringValue(from.CreatedAt)
	}

	// Handle updated_at - use state as ground truth when API returns invalid values
	if (from.UpdatedAt == "" || from.UpdatedAt == "0001-01-01T00:00:00Z") && currentState != nil && !currentState.UpdatedAt.IsNull() {
		// API returned invalid timestamp, preserve state value
		m.UpdatedAt = currentState.UpdatedAt
	} else if from.UpdatedAt == "" || from.UpdatedAt == "0001-01-01T00:00:00Z" {
		// No state value available, use null
		m.UpdatedAt = types.StringNull()
	} else {
		// Use API value
		m.UpdatedAt = types.StringValue(from.UpdatedAt)
	}

	// Handle TagIDs - API doesn't return this field, so always preserve current state
	// These fields are provider-managed only and should never change based on API response
	if currentState != nil && currentState.TagIDs != nil {
		m.TagIDs = currentState.TagIDs
	} else {
		m.TagIDs = []types.String{}
	}

	// Handle NotificationIDs - API doesn't return this field, so always preserve current state
	// These fields are provider-managed only and should never change based on API response
	if currentState != nil && currentState.NotificationIDs != nil {
		m.NotificationIDs = currentState.NotificationIDs
	} else {
		m.NotificationIDs = []types.String{}
	}
}
