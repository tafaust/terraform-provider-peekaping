// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/tafaust/terraform-provider-peekaping/internal/peekaping"
)

var _ resource.Resource = &TagResource{}
var _ resource.ResourceWithImportState = &TagResource{}

type TagResource struct {
	client *peekaping.Client
}

func NewTagResource() resource.Resource { return &TagResource{} }

// tagNameValidator validates tag name constraints.
type tagNameValidator struct{}

func (v tagNameValidator) Description(_ context.Context) string {
	return "Tag name must be between 1 and 50 characters"
}

func (v tagNameValidator) MarkdownDescription(_ context.Context) string {
	return "Tag name must be between 1 and 50 characters"
}

func (v tagNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	name := req.ConfigValue.ValueString()
	if len(name) < 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Name",
			"Tag name must be at least 1 character long",
		)
		return
	}
	if len(name) > 50 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Name",
			"Tag name must be at most 50 characters long",
		)
		return
	}
}

// tagColorValidator validates tag color hex format.
type tagColorValidator struct{}

func (v tagColorValidator) Description(_ context.Context) string {
	return "Tag color must be a valid hex color code (e.g., #FF0000)"
}

func (v tagColorValidator) MarkdownDescription(_ context.Context) string {
	return "Tag color must be a valid hex color code (e.g., `#FF0000`)"
}

func (v tagColorValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	color := req.ConfigValue.ValueString()
	if color == "" {
		return // Empty color is allowed
	}

	// Check if it's a valid hex color code
	hexPattern := regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
	if !hexPattern.MatchString(color) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Color",
			"Tag color must be a valid hex color code (e.g., #FF0000 or #F00)",
		)
		return
	}
}

// tagDescriptionValidator validates tag description constraints.
type tagDescriptionValidator struct{}

func (v tagDescriptionValidator) Description(_ context.Context) string {
	return "Tag description must be at most 200 characters"
}

func (v tagDescriptionValidator) MarkdownDescription(_ context.Context) string {
	return "Tag description must be at most 200 characters"
}

func (v tagDescriptionValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	description := req.ConfigValue.ValueString()
	if len(description) > 200 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Description",
			"Tag description must be at most 200 characters long",
		)
		return
	}
}

type tagResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Color       types.String `tfsdk:"color"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func (r *TagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *TagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description: "Tag name",
				Validators: []validator.String{
					tagNameValidator{},
				},
			},
			"color": schema.StringAttribute{
				Optional:    true,
				Description: "Tag color (hex code)",
				Validators: []validator.String{
					tagColorValidator{},
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Tag description",
				Validators: []validator.String{
					tagDescriptionValidator{},
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

func (r *TagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the tag we're creating for debugging
	tflog.Info(ctx, "Creating tag", map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"color":       plan.Color.ValueString(),
		"description": plan.Description.ValueString(),
	})

	in := peekaping.TagCreate{
		Name:        plan.Name.ValueString(),
		Color:       plan.Color.ValueString(),
		Description: plan.Description.ValueString(),
	}

	t, err := r.client.CreateTag(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError("create tag failed", err.Error())
		return
	}
	setModelFromTag(&plan, t)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	t, err := r.client.GetTag(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read tag failed", err.Error())
		return
	}

	// Log what the API returned for debugging
	tflog.Info(ctx, "API returned tag", map[string]interface{}{
		"id":          t.ID,
		"name":        t.Name,
		"color":       t.Color,
		"description": t.Description,
	})

	// Use state-aware field mapping to handle API inconsistencies
	setModelFromTagWithState(&state, t, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to get the ID
	var state tagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upd := peekaping.TagUpdate{}
	if !plan.Name.IsNull() {
		v := plan.Name.ValueString()
		upd.Name = &v
	}
	if !plan.Color.IsNull() {
		v := plan.Color.ValueString()
		upd.Color = &v
	}
	if !plan.Description.IsNull() {
		v := plan.Description.ValueString()
		upd.Description = &v
	}

	// Use state.ID instead of plan.ID
	t, err := r.client.UpdateTag(ctx, state.ID.ValueString(), upd)
	if err != nil {
		resp.Diagnostics.AddError("update tag failed", err.Error())
		return
	}
	setModelFromTag(&plan, t)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteTag(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("delete tag failed", err.Error())
		return
	}
}

func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state tagResourceModel
	state.ID = types.StringValue(req.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func setModelFromTag(m *tagResourceModel, from *peekaping.Tag) {
	// Required fields - always present
	m.ID = types.StringValue(from.ID)
	m.Name = types.StringValue(from.Name)

	// Optional string fields
	if from.Color != "" {
		m.Color = types.StringValue(from.Color)
	} else {
		m.Color = types.StringNull()
	}
	if from.Description != "" {
		m.Description = types.StringValue(from.Description)
	} else {
		m.Description = types.StringNull()
	}

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

// setModelFromTagWithState handles field mapping with state comparison to resolve API inconsistencies.
func setModelFromTagWithState(m *tagResourceModel, from *peekaping.Tag, currentState *tagResourceModel) {
	// Required fields - always present
	m.ID = types.StringValue(from.ID)
	m.Name = types.StringValue(from.Name)

	// Optional string fields - preserve current state if API doesn't return them
	if from.Color != "" {
		m.Color = types.StringValue(from.Color)
	} else if currentState != nil && !currentState.Color.IsNull() {
		// If API doesn't return color, preserve current state
		m.Color = currentState.Color
	} else {
		m.Color = types.StringNull()
	}

	if from.Description != "" {
		m.Description = types.StringValue(from.Description)
	} else if currentState != nil && !currentState.Description.IsNull() {
		// If API doesn't return description, preserve current state
		m.Description = currentState.Description
	} else {
		m.Description = types.StringNull()
	}

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
