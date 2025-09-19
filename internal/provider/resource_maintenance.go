// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/tafaust/terraform-provider-peekaping/internal/peekaping"
)

var _ resource.Resource = &MaintenanceResource{}
var _ resource.ResourceWithImportState = &MaintenanceResource{}

type MaintenanceResource struct {
	client *peekaping.Client
}

func NewMaintenanceResource() resource.Resource { return &MaintenanceResource{} }

// maintenanceTitleValidator validates maintenance title constraints.
type maintenanceTitleValidator struct{}

func (v maintenanceTitleValidator) Description(_ context.Context) string {
	return "Maintenance title must be between 1 and 100 characters"
}

func (v maintenanceTitleValidator) MarkdownDescription(_ context.Context) string {
	return "Maintenance title must be between 1 and 100 characters"
}

func (v maintenanceTitleValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	title := req.ConfigValue.ValueString()
	if len(title) < 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Title",
			"Maintenance title must be at least 1 character long",
		)
		return
	}
	if len(title) > 100 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Title",
			"Maintenance title must be at most 100 characters long",
		)
		return
	}
}

// maintenanceStrategyValidator validates maintenance strategy.
type maintenanceStrategyValidator struct{}

func (v maintenanceStrategyValidator) Description(_ context.Context) string {
	return "Maintenance strategy must be one of: once, recurring, cron"
}

func (v maintenanceStrategyValidator) MarkdownDescription(_ context.Context) string {
	return "Maintenance strategy must be one of: `once`, `recurring`, `cron`"
}

func (v maintenanceStrategyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	strategy := req.ConfigValue.ValueString()
	supportedStrategies := []string{"once", "recurring", "cron"}

	for _, supportedStrategy := range supportedStrategies {
		if strategy == supportedStrategy {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Strategy",
		fmt.Sprintf("Maintenance strategy '%s' is not supported. Supported strategies are: %s", strategy, strings.Join(supportedStrategies, ", ")),
	)
}

// maintenanceCronValidator validates cron expression format.
type maintenanceCronValidator struct{}

func (v maintenanceCronValidator) Description(_ context.Context) string {
	return "Cron expression must be a valid 5-field cron format"
}

func (v maintenanceCronValidator) MarkdownDescription(_ context.Context) string {
	return "Cron expression must be a valid 5-field cron format (minute hour day month weekday)"
}

func (v maintenanceCronValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	cron := req.ConfigValue.ValueString()
	if cron == "" {
		return // Empty cron is allowed
	}

	// Basic cron validation - check for 5 fields separated by spaces
	cronPattern := regexp.MustCompile(`^(\S+\s+){4}\S+$`)
	if !cronPattern.MatchString(cron) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Cron Expression",
			"Cron expression must be in 5-field format: minute hour day month weekday",
		)
		return
	}
}

// maintenanceTimeValidator validates time format (HH:MM).
type maintenanceTimeValidator struct{}

func (v maintenanceTimeValidator) Description(_ context.Context) string {
	return "Time must be in HH:MM format (24-hour)"
}

func (v maintenanceTimeValidator) MarkdownDescription(_ context.Context) string {
	return "Time must be in HH:MM format (24-hour)"
}

func (v maintenanceTimeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	timeStr := req.ConfigValue.ValueString()
	if timeStr == "" {
		return // Empty time is allowed
	}

	// Validate HH:MM format
	timePattern := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)
	if !timePattern.MatchString(timeStr) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Time Format",
			"Time must be in HH:MM format (24-hour), e.g., 14:30",
		)
		return
	}
}

// maintenanceDescriptionValidator validates maintenance description constraints.
type maintenanceDescriptionValidator struct{}

func (v maintenanceDescriptionValidator) Description(_ context.Context) string {
	return "Maintenance description must be at most 500 characters"
}

func (v maintenanceDescriptionValidator) MarkdownDescription(_ context.Context) string {
	return "Maintenance description must be at most 500 characters"
}

func (v maintenanceDescriptionValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	description := req.ConfigValue.ValueString()
	if len(description) > 500 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Description",
			"Maintenance description must be at most 500 characters long",
		)
		return
	}
}

type maintenanceResourceModel struct {
	ID            types.String  `tfsdk:"id"`
	Title         types.String  `tfsdk:"title"`
	Description   types.String  `tfsdk:"description"`
	Strategy      types.String  `tfsdk:"strategy"`
	Active        types.Bool    `tfsdk:"active"`
	MonitorIDs    types.List    `tfsdk:"monitor_ids"`
	StartDateTime types.String  `tfsdk:"start_date_time"`
	EndDateTime   types.String  `tfsdk:"end_date_time"`
	Duration      types.Int64   `tfsdk:"duration"`
	Timezone      types.String  `tfsdk:"timezone"`
	Cron          types.String  `tfsdk:"cron"`
	Weekdays      []types.Int64 `tfsdk:"weekdays"`
	DaysOfMonth   []types.Int64 `tfsdk:"days_of_month"`
	IntervalDay   types.Int64   `tfsdk:"interval_day"`
	StartTime     types.String  `tfsdk:"start_time"`
	EndTime       types.String  `tfsdk:"end_time"`
	CreatedAt     types.String  `tfsdk:"created_at"`
	UpdatedAt     types.String  `tfsdk:"updated_at"`
}

func (r *MaintenanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_maintenance"
}

func (r *MaintenanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Maintenance window title",
				Validators: []validator.String{
					maintenanceTitleValidator{},
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Maintenance window description",
				Validators: []validator.String{
					maintenanceDescriptionValidator{},
				},
			},
			"strategy": schema.StringAttribute{
				Required:    true,
				Description: "Maintenance strategy (once, recurring, cron)",
				Validators: []validator.String{
					maintenanceStrategyValidator{},
				},
			},
			"active": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the maintenance window is active",
			},
			"monitor_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of monitor IDs to include in maintenance",
			},
			"start_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "Start date and time (ISO 8601 format)",
			},
			"end_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "End date and time (ISO 8601 format)",
			},
			"duration": schema.Int64Attribute{
				Computed:    true,
				Description: "Duration in minutes",
			},
			"timezone": schema.StringAttribute{
				Optional:    true,
				Description: "Timezone for the maintenance window",
			},
			"cron": schema.StringAttribute{
				Optional:    true,
				Description: "Cron expression for recurring maintenance",
				Validators: []validator.String{
					maintenanceCronValidator{},
				},
			},
			"weekdays": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Days of the week (0=Sunday, 1=Monday, etc.)",
			},
			"days_of_month": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Days of the month (1-31)",
			},
			"interval_day": schema.Int64Attribute{
				Optional:    true,
				Description: "Interval in days for recurring maintenance",
			},
			"start_time": schema.StringAttribute{
				Optional:    true,
				Description: "Start time (HH:MM format)",
				Validators: []validator.String{
					maintenanceTimeValidator{},
				},
			},
			"end_time": schema.StringAttribute{
				Optional:    true,
				Description: "End time (HH:MM format)",
				Validators: []validator.String{
					maintenanceTimeValidator{},
				},
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

func (r *MaintenanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MaintenanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan maintenanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the maintenance we're creating for debugging
	tflog.Info(ctx, "Creating maintenance", map[string]interface{}{
		"title":       plan.Title.ValueString(),
		"description": plan.Description.ValueString(),
		"strategy":    plan.Strategy.ValueString(),
		"timezone":    plan.Timezone.ValueString(),
	})

	in := peekaping.MaintenanceCreate{
		Title:       plan.Title.ValueString(),
		Description: plan.Description.ValueString(),
		Strategy:    plan.Strategy.ValueString(),
		Timezone:    plan.Timezone.ValueString(),
		Cron:        plan.Cron.ValueString(),
		Weekdays:    toIntSlice(plan.Weekdays),
		DaysOfMonth: toIntSlice(plan.DaysOfMonth),
		IntervalDay: int(plan.IntervalDay.ValueInt64()),
		StartTime:   plan.StartTime.ValueString(),
		EndTime:     plan.EndTime.ValueString(),
	}

	m, err := r.client.CreateMaintenance(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError("create maintenance failed", err.Error())
		return
	}
	setModelFromMaintenanceWithState(&plan, m)
	// Preserve the plan's monitor_ids and duration to maintain Terraform state consistency
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MaintenanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state maintenanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	m, err := r.client.GetMaintenance(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read maintenance failed", err.Error())
		return
	}
	setModelFromMaintenance(&state, m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MaintenanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan maintenanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to get the ID
	var state maintenanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upd := peekaping.MaintenanceUpdate{
		Weekdays:    toIntSlice(plan.Weekdays),
		DaysOfMonth: toIntSlice(plan.DaysOfMonth),
	}
	if !plan.Title.IsNull() {
		v := plan.Title.ValueString()
		upd.Title = &v
	}
	if !plan.Description.IsNull() {
		v := plan.Description.ValueString()
		upd.Description = &v
	}
	if !plan.Strategy.IsNull() {
		v := plan.Strategy.ValueString()
		upd.Strategy = &v
	}
	// active is computed by the API, not updated by user
	if !plan.Duration.IsNull() {
		v := int(plan.Duration.ValueInt64())
		upd.Duration = &v
	}
	if !plan.Timezone.IsNull() {
		v := plan.Timezone.ValueString()
		upd.Timezone = &v
	}
	if !plan.Cron.IsNull() {
		v := plan.Cron.ValueString()
		upd.Cron = &v
	}
	if !plan.IntervalDay.IsNull() {
		v := int(plan.IntervalDay.ValueInt64())
		upd.IntervalDay = &v
	}
	if !plan.StartTime.IsNull() {
		v := plan.StartTime.ValueString()
		upd.StartTime = &v
	}
	if !plan.EndTime.IsNull() {
		v := plan.EndTime.ValueString()
		upd.EndTime = &v
	}

	// Use state.ID instead of plan.ID
	m, err := r.client.UpdateMaintenance(ctx, state.ID.ValueString(), upd)
	if err != nil {
		resp.Diagnostics.AddError("update maintenance failed", err.Error())
		return
	}
	setModelFromMaintenance(&plan, m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MaintenanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state maintenanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteMaintenance(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("delete maintenance failed", err.Error())
		return
	}
}

func (r *MaintenanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state maintenanceResourceModel
	state.ID = types.StringValue(req.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func toIntSlice(xs []types.Int64) []int {
	out := make([]int, 0, len(xs))
	for _, x := range xs {
		if !x.IsNull() {
			out = append(out, int(x.ValueInt64()))
		}
	}
	return out
}

// formatDateTime formats a datetime string to the format expected by the API

func setModelFromMaintenanceWithState(m *maintenanceResourceModel, from *peekaping.Maintenance) {
	// For computed fields, always use API response values
	// For user-configurable fields, preserve plan values if they exist

	// Set all fields from API response
	setModelFromMaintenance(m, from)

	// For computed fields, use API values (not plan values)
	// These fields are marked as Computed: true in the schema
	// so they should always reflect the server state
}

func setModelFromMaintenance(m *maintenanceResourceModel, from *peekaping.Maintenance) {
	m.ID = types.StringValue(from.ID)
	m.Title = types.StringValue(from.Title)
	if from.Description != "" {
		m.Description = types.StringValue(from.Description)
	} else {
		m.Description = types.StringNull()
	}
	m.Strategy = types.StringValue(from.Strategy)
	m.Active = types.BoolValue(from.Active)
	if from.StartDateTime != "" {
		m.StartDateTime = types.StringValue(from.StartDateTime)
	} else {
		m.StartDateTime = types.StringNull()
	}
	if from.EndDateTime != "" {
		m.EndDateTime = types.StringValue(from.EndDateTime)
	} else {
		m.EndDateTime = types.StringNull()
	}
	if from.Duration != 0 {
		m.Duration = types.Int64Value(int64(from.Duration))
	} else {
		m.Duration = types.Int64Null()
	}
	if from.Timezone != "" {
		m.Timezone = types.StringValue(from.Timezone)
	} else {
		m.Timezone = types.StringNull()
	}
	if from.Cron != "" {
		m.Cron = types.StringValue(from.Cron)
	} else {
		m.Cron = types.StringNull()
	}
	if from.IntervalDay != 0 {
		m.IntervalDay = types.Int64Value(int64(from.IntervalDay))
	} else {
		m.IntervalDay = types.Int64Null()
	}
	if from.StartTime != "" {
		m.StartTime = types.StringValue(from.StartTime)
	} else {
		m.StartTime = types.StringNull()
	}
	if from.EndTime != "" {
		m.EndTime = types.StringValue(from.EndTime)
	} else {
		m.EndTime = types.StringNull()
	}
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

	// Handle MonitorIDs - always set from API response since it's computed
	ids := make([]attr.Value, 0, len(from.MonitorIDs))
	for _, id := range from.MonitorIDs {
		ids = append(ids, types.StringValue(id))
	}
	m.MonitorIDs = types.ListValueMust(types.StringType, ids)
	if len(from.Weekdays) > 0 {
		days := make([]types.Int64, 0, len(from.Weekdays))
		for _, d := range from.Weekdays {
			days = append(days, types.Int64Value(int64(d)))
		}
		m.Weekdays = days
	} else {
		m.Weekdays = nil
	}
	if len(from.DaysOfMonth) > 0 {
		days := make([]types.Int64, 0, len(from.DaysOfMonth))
		for _, d := range from.DaysOfMonth {
			days = append(days, types.Int64Value(int64(d)))
		}
		m.DaysOfMonth = days
	} else {
		m.DaysOfMonth = nil
	}
}
