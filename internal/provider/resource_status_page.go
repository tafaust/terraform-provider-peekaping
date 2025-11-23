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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/tafaust/terraform-provider-peekaping/internal/peekaping"
)

var _ resource.Resource = &StatusPageResource{}
var _ resource.ResourceWithImportState = &StatusPageResource{}

// normalizeMonitorIDsPlanModifier uses API's order from state for updates.
type normalizeMonitorIDsPlanModifier struct{}

func (m normalizeMonitorIDsPlanModifier) Description(_ context.Context) string {
	return "Use API's monitor order from state for updates to ensure consistency"
}

func (m normalizeMonitorIDsPlanModifier) MarkdownDescription(_ context.Context) string {
	return "Use API's monitor order from state for updates to ensure consistency"
}

func (m normalizeMonitorIDsPlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// Only modify if we have a planned value and it's not unknown
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	// Debug logging
	tflog.Debug(ctx, "normalizeMonitorIDsPlanModifier called", map[string]interface{}{
		"has_state": !req.StateValue.IsNull() && !req.StateValue.IsUnknown(),
		"state_len": len(req.StateValue.Elements()),
		"plan_len":  len(req.PlanValue.Elements()),
	})

	// If we have a state value (existing resource) and it's not empty, use that order as the source of truth
	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() && len(req.StateValue.Elements()) > 0 {
		// Get the state monitor IDs
		stateIDs := make([]string, 0, len(req.StateValue.Elements()))
		for _, elem := range req.StateValue.Elements() {
			if str, ok := elem.(types.String); ok && !str.IsNull() && !str.IsUnknown() {
				stateIDs = append(stateIDs, str.ValueString())
			}
		}

		// Get the planned monitor IDs
		plannedIDs := make([]string, 0, len(req.PlanValue.Elements()))
		for _, elem := range req.PlanValue.Elements() {
			if str, ok := elem.(types.String); ok && !str.IsNull() && !str.IsUnknown() {
				plannedIDs = append(plannedIDs, str.ValueString())
			}
		}

		// Check if the planned IDs are the same as the state IDs (ignoring order)
		if len(plannedIDs) == len(stateIDs) {
			stateIDSet := make(map[string]bool)
			for _, id := range stateIDs {
				stateIDSet[id] = true
			}
			plannedIDSet := make(map[string]bool)
			for _, id := range plannedIDs {
				plannedIDSet[id] = true
			}

			// If the sets are the same, use the state order
			sameSet := true
			for id := range stateIDSet {
				if !plannedIDSet[id] {
					sameSet = false
					break
				}
			}
			for id := range plannedIDSet {
				if !stateIDSet[id] {
					sameSet = false
					break
				}
			}

			if sameSet {
				// Use the state order as the source of truth
				tflog.Debug(ctx, "Using state order for monitor_ids", map[string]interface{}{
					"state_order":   stateIDs,
					"planned_order": plannedIDs,
				})
				resp.PlanValue = req.StateValue
				return
			}
		}
	}

	// For new resources or when IDs have changed, preserve user order
	// Don't modify the plan - let user order be preserved
}

// statusPageTitleValidator validates status page title constraints.
type statusPageTitleValidator struct{}

func (v statusPageTitleValidator) Description(_ context.Context) string {
	return "Status page title must be between 1 and 100 characters"
}

func (v statusPageTitleValidator) MarkdownDescription(_ context.Context) string {
	return "Status page title must be between 1 and 100 characters"
}

func (v statusPageTitleValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	title := req.ConfigValue.ValueString()
	if len(title) < 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Title",
			"Status page title must be at least 1 character long",
		)
		return
	}
	if len(title) > 100 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Title",
			"Status page title must be at most 100 characters long",
		)
		return
	}
}

// statusPageSlugValidator validates status page slug format.
type statusPageSlugValidator struct{}

func (v statusPageSlugValidator) Description(_ context.Context) string {
	return "Status page slug must be a valid URL slug (lowercase letters, numbers, and hyphens only)"
}

func (v statusPageSlugValidator) MarkdownDescription(_ context.Context) string {
	return "Status page slug must be a valid URL slug (lowercase letters, numbers, and hyphens only)"
}

func (v statusPageSlugValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	slug := req.ConfigValue.ValueString()
	if slug == "" {
		return // Empty slug is allowed
	}

	// Check if it's a valid URL slug
	slugPattern := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !slugPattern.MatchString(slug) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Slug",
			"Status page slug must contain only lowercase letters, numbers, and hyphens",
		)
		return
	}

	// Check length
	if len(slug) > 50 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Slug",
			"Status page slug must be at most 50 characters long",
		)
		return
	}
}

// statusPageThemeValidator validates status page theme.
type statusPageThemeValidator struct{}

func (v statusPageThemeValidator) Description(_ context.Context) string {
	return "Status page theme must be one of: light, dark, auto"
}

func (v statusPageThemeValidator) MarkdownDescription(_ context.Context) string {
	return "Status page theme must be one of: `light`, `dark`, `auto`"
}

func (v statusPageThemeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	theme := req.ConfigValue.ValueString()
	if theme == "" {
		return // Empty theme is allowed
	}

	supportedThemes := []string{"light", "dark", "auto"}
	for _, supportedTheme := range supportedThemes {
		if theme == supportedTheme {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Theme",
		fmt.Sprintf("Status page theme '%s' is not supported. Supported themes are: %s", theme, strings.Join(supportedThemes, ", ")),
	)
}

// statusPageDescriptionValidator validates status page description constraints.
type statusPageDescriptionValidator struct{}

func (v statusPageDescriptionValidator) Description(_ context.Context) string {
	return "Status page description must be at most 500 characters"
}

func (v statusPageDescriptionValidator) MarkdownDescription(_ context.Context) string {
	return "Status page description must be at most 500 characters"
}

func (v statusPageDescriptionValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	description := req.ConfigValue.ValueString()
	if len(description) > 500 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Description",
			"Status page description must be at most 500 characters long",
		)
		return
	}
}

// statusPageFooterTextValidator validates status page footer text constraints.
type statusPageFooterTextValidator struct{}

func (v statusPageFooterTextValidator) Description(_ context.Context) string {
	return "Status page footer text must be at most 200 characters"
}

func (v statusPageFooterTextValidator) MarkdownDescription(_ context.Context) string {
	return "Status page footer text must be at most 200 characters"
}

func (v statusPageFooterTextValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	footerText := req.ConfigValue.ValueString()
	if len(footerText) > 200 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Footer Text",
			"Status page footer text must be at most 200 characters long",
		)
		return
	}
}

type StatusPageResource struct {
	client *peekaping.Client
}

func NewStatusPageResource() resource.Resource { return &StatusPageResource{} }

type statusPageResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	Title                 types.String   `tfsdk:"title"`
	Description           types.String   `tfsdk:"description"`
	Slug                  types.String   `tfsdk:"slug"`
	Domains               []types.String `tfsdk:"domains"`
	MonitorIDs            types.List     `tfsdk:"monitor_ids"`
	Published             types.Bool     `tfsdk:"published"`
	Theme                 types.String   `tfsdk:"theme"`
	Icon                  types.String   `tfsdk:"icon"`
	FooterText            types.String   `tfsdk:"footer_text"`
	CustomCSS             types.String   `tfsdk:"custom_css"`
	GoogleAnalyticsTagID  types.String   `tfsdk:"google_analytics_tag_id"`
	AutoRefreshInterval   types.Int64    `tfsdk:"auto_refresh_interval"`
	SearchEngineIndex     types.Bool     `tfsdk:"search_engine_index"`
	ShowCertificateExpiry types.Bool     `tfsdk:"show_certificate_expiry"`
	ShowPoweredBy         types.Bool     `tfsdk:"show_powered_by"`
	ShowTags              types.Bool     `tfsdk:"show_tags"`
	Password              types.String   `tfsdk:"password"`
	CreatedAt             types.String   `tfsdk:"created_at"`
	UpdatedAt             types.String   `tfsdk:"updated_at"`
}

func (r *StatusPageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (r *StatusPageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description: "Status page title",
				Validators: []validator.String{
					statusPageTitleValidator{},
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Status page description",
				Validators: []validator.String{
					statusPageDescriptionValidator{},
				},
			},
			"slug": schema.StringAttribute{
				Optional:    true,
				Description: "URL slug for the status page",
				Validators: []validator.String{
					statusPageSlugValidator{},
				},
			},
			"domains": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of custom domains for the status page",
			},
			"monitor_ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of monitor IDs to display on the status page",
				PlanModifiers: []planmodifier.List{
					normalizeMonitorIDsPlanModifier{},
				},
			},
			"published": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the status page is published",
			},
			"theme": schema.StringAttribute{
				Optional:    true,
				Description: "Status page theme",
				Validators: []validator.String{
					statusPageThemeValidator{},
				},
			},
			"icon": schema.StringAttribute{
				Optional:    true,
				Description: "Status page icon",
			},
			"footer_text": schema.StringAttribute{
				Optional:    true,
				Description: "Footer text for the status page",
				Validators: []validator.String{
					statusPageFooterTextValidator{},
				},
			},
			"custom_css": schema.StringAttribute{
				Computed:    true,
				Description: "Custom CSS for the status page",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"google_analytics_tag_id": schema.StringAttribute{
				Computed:    true,
				Description: "Google Analytics tag ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_refresh_interval": schema.Int64Attribute{
				Computed:    true,
				Description: "Auto refresh interval in seconds",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"search_engine_index": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to allow search engine indexing",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"show_certificate_expiry": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to show certificate expiry information",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"show_powered_by": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to show 'Powered by' text",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"show_tags": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to show tags on the status page",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Password for the status page (if protected)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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

func (r *StatusPageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StatusPageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan statusPageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the status page we're creating for debugging
	tflog.Info(ctx, "Creating status page", map[string]interface{}{
		"title":     plan.Title.ValueString(),
		"slug":      plan.Slug.ValueString(),
		"published": plan.Published.ValueBool(),
		"theme":     plan.Theme.ValueString(),
	})

	in := peekaping.StatusPageCreate{
		Title:       plan.Title.ValueString(),
		Description: plan.Description.ValueString(),
		Slug:        plan.Slug.ValueString(),
		Domains:     toStrSliceFromStringSlice(plan.Domains),
		MonitorIDs:  toStrSliceFromList(plan.MonitorIDs),
		Published:   plan.Published.ValueBool(),
		Theme:       plan.Theme.ValueString(),
		Icon:        plan.Icon.ValueString(),
		FooterText:  plan.FooterText.ValueString(),
	}

	sp, err := r.client.CreateStatusPage(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError("create status page failed", err.Error())
		return
	}
	setModelFromStatusPageWithState(&plan, sp, &plan)
	// Preserve the plan's monitor_ids, custom_css, google_analytics_tag_id, password
	// and boolean flags to maintain Terraform state consistency
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state statusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sp, err := r.client.GetStatusPage(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read status page failed", err.Error())
		return
	}
	setModelFromStatusPageWithState(&state, sp, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StatusPageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan statusPageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to get the ID
	var state statusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upd := peekaping.StatusPageUpdate{
		Domains:    toStrSliceFromStringSlice(plan.Domains),
		MonitorIDs: toStrSliceFromList(plan.MonitorIDs),
	}
	if !plan.Title.IsNull() {
		v := plan.Title.ValueString()
		upd.Title = &v
	}
	if !plan.Description.IsNull() {
		v := plan.Description.ValueString()
		upd.Description = &v
	}
	if !plan.Slug.IsNull() {
		v := plan.Slug.ValueString()
		upd.Slug = &v
	}
	if !plan.Published.IsNull() {
		v := plan.Published.ValueBool()
		upd.Published = &v
	}
	if !plan.Theme.IsNull() {
		v := plan.Theme.ValueString()
		upd.Theme = &v
	}
	if !plan.Icon.IsNull() {
		v := plan.Icon.ValueString()
		upd.Icon = &v
	}
	if !plan.FooterText.IsNull() {
		v := plan.FooterText.ValueString()
		upd.FooterText = &v
	}
	// custom_css, google_analytics_tag_id, monitor_ids, password, and boolean flags
	// are computed by the API, not updated by user

	// Use state.ID instead of plan.ID
	_, err := r.client.UpdateStatusPage(ctx, state.ID.ValueString(), upd)
	if err != nil {
		resp.Diagnostics.AddError("update status page failed", err.Error())
		return
	}

	// The API doesn't return monitor_ids in the PATCH response, so we need to fetch the full status page
	// to get the updated monitor_ids
	fullSp, err := r.client.GetStatusPage(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to fetch updated status page", err.Error())
		return
	}

	setModelFromStatusPageWithState(&plan, fullSp, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StatusPageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state statusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteStatusPage(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("delete status page failed", err.Error())
		return
	}
}

func (r *StatusPageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state statusPageResourceModel
	state.ID = types.StringValue(req.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func setModelFromStatusPageWithState(m *statusPageResourceModel, from *peekaping.StatusPage, currentState *statusPageResourceModel) {
	m.ID = types.StringValue(from.ID)
	m.Title = types.StringValue(from.Title)
	if from.Description != "" {
		m.Description = types.StringValue(from.Description)
	} else {
		m.Description = types.StringNull()
	}
	if from.Slug != "" {
		m.Slug = types.StringValue(from.Slug)
	} else {
		m.Slug = types.StringNull()
	}
	m.Published = types.BoolValue(from.Published)
	if from.Theme != "" {
		m.Theme = types.StringValue(from.Theme)
	} else {
		m.Theme = types.StringNull()
	}
	if from.Icon != "" {
		m.Icon = types.StringValue(from.Icon)
	} else {
		m.Icon = types.StringNull()
	}
	if from.FooterText != "" {
		m.FooterText = types.StringValue(from.FooterText)
	} else {
		m.FooterText = types.StringNull()
	}
	// Always set computed fields from API response
	m.CustomCSS = types.StringValue(from.CustomCSS)
	m.GoogleAnalyticsTagID = types.StringValue(from.GoogleAnalyticsTagID)
	m.AutoRefreshInterval = types.Int64Value(int64(from.AutoRefreshInterval))
	// Always set computed boolean fields from API response
	m.SearchEngineIndex = types.BoolValue(from.SearchEngineIndex)
	m.ShowCertificateExpiry = types.BoolValue(from.ShowCertificateExpiry)
	m.ShowPoweredBy = types.BoolValue(from.ShowPoweredBy)
	m.ShowTags = types.BoolValue(from.ShowTags)
	// Always set computed password field from API response
	if from.Password != "" {
		m.Password = types.StringValue(from.Password)
	} else {
		m.Password = types.StringNull()
	}
	// Always set computed timestamp fields from API response
	m.CreatedAt = types.StringValue(from.CreatedAt)
	m.UpdatedAt = types.StringValue(from.UpdatedAt)

	// Handle Domains - preserve state when API doesn't return them
	// The API may not return domains in its response, so we preserve the user's configured value
	if len(from.Domains) > 0 {
		domains := make([]types.String, 0, len(from.Domains))
		for _, domain := range from.Domains {
			domains = append(domains, types.StringValue(domain))
		}
		m.Domains = domains
	} else if currentState != nil && len(currentState.Domains) > 0 {
		// API didn't return domains, preserve state
		m.Domains = currentState.Domains
	} else {
		m.Domains = nil
	}

	// Handle MonitorIDs - use state as ground truth when available
	// The API returns monitors ordered by created_at DESC, but we want to preserve the user's order
	if currentState != nil && !currentState.MonitorIDs.IsNull() && !currentState.MonitorIDs.IsUnknown() {
		// Use state order as ground truth
		m.MonitorIDs = currentState.MonitorIDs
	} else if len(from.MonitorIDs) > 0 {
		// No state available, use API response
		ids := make([]attr.Value, 0, len(from.MonitorIDs))
		for _, id := range from.MonitorIDs {
			ids = append(ids, types.StringValue(id))
		}
		m.MonitorIDs = types.ListValueMust(types.StringType, ids)
	} else {
		// If API returns null or empty list, set to null (not empty list)
		m.MonitorIDs = types.ListNull(types.StringType)
	}
}

// Helper function to convert Terraform list to string slice.
func toStrSliceFromList(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var result []string
	for _, elem := range list.Elements() {
		if str, ok := elem.(types.String); ok && !str.IsNull() && !str.IsUnknown() {
			result = append(result, str.ValueString())
		}
	}
	return result
}

// Helper function to convert []types.String to string slice.
func toStrSliceFromStringSlice(strSlice []types.String) []string {
	if strSlice == nil {
		return nil
	}

	var result []string
	for _, str := range strSlice {
		if !str.IsNull() && !str.IsUnknown() {
			result = append(result, str.ValueString())
		}
	}
	return result
}
