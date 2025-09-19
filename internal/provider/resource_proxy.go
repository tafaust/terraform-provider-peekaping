// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/tafaust/terraform-provider-peekaping/internal/peekaping"
)

var _ resource.Resource = &ProxyResource{}
var _ resource.ResourceWithImportState = &ProxyResource{}

type ProxyResource struct {
	client *peekaping.Client
}

func NewProxyResource() resource.Resource { return &ProxyResource{} }

// proxyHostValidator validates proxy host format.
type proxyHostValidator struct{}

func (v proxyHostValidator) Description(_ context.Context) string {
	return "Proxy host must be a valid hostname or IP address"
}

func (v proxyHostValidator) MarkdownDescription(_ context.Context) string {
	return "Proxy host must be a valid hostname or IP address"
}

func (v proxyHostValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	host := req.ConfigValue.ValueString()
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Host",
			"Proxy host cannot be empty",
		)
		return
	}

	// Basic hostname/IP validation
	hostPattern := regexp.MustCompile(`^[a-zA-Z0-9.-]+$`)
	if !hostPattern.MatchString(host) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Host",
			"Proxy host must be a valid hostname or IP address",
		)
		return
	}
}

// proxyPortValidator validates proxy port range.
type proxyPortValidator struct{}

func (v proxyPortValidator) Description(_ context.Context) string {
	return "Proxy port must be between 1 and 65535"
}

func (v proxyPortValidator) MarkdownDescription(_ context.Context) string {
	return "Proxy port must be between 1 and 65535"
}

func (v proxyPortValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	port := req.ConfigValue.ValueInt64()
	if port < 1 || port > 65535 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Port",
			"Proxy port must be between 1 and 65535",
		)
		return
	}
}

// proxyProtocolValidator validates proxy protocol.
type proxyProtocolValidator struct{}

func (v proxyProtocolValidator) Description(_ context.Context) string {
	return "Proxy protocol must be one of: http, https, socks, socks4, socks5, socks5h"
}

func (v proxyProtocolValidator) MarkdownDescription(_ context.Context) string {
	return "Proxy protocol must be one of: `http`, `https`, `socks`, `socks4`, `socks5`, `socks5h`"
}

func (v proxyProtocolValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	protocol := req.ConfigValue.ValueString()
	supportedProtocols := []string{
		"http", "https", "socks", "socks4", "socks5", "socks5h",
	}

	for _, supportedProtocol := range supportedProtocols {
		if protocol == supportedProtocol {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Protocol",
		fmt.Sprintf("Proxy protocol '%s' is not supported. Supported protocols are: %s", protocol, strings.Join(supportedProtocols, ", ")),
	)
}

// proxyUsernameValidator validates proxy username constraints.
type proxyUsernameValidator struct{}

func (v proxyUsernameValidator) Description(_ context.Context) string {
	return "Proxy username must be between 1 and 100 characters"
}

func (v proxyUsernameValidator) MarkdownDescription(_ context.Context) string {
	return "Proxy username must be between 1 and 100 characters"
}

func (v proxyUsernameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	username := req.ConfigValue.ValueString()
	if len(username) < 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Username",
			"Proxy username must be at least 1 character long",
		)
		return
	}
	if len(username) > 100 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Username",
			"Proxy username must be at most 100 characters long",
		)
		return
	}
}

// proxyPasswordValidator validates proxy password constraints.
type proxyPasswordValidator struct{}

func (v proxyPasswordValidator) Description(_ context.Context) string {
	return "Proxy password must be between 1 and 100 characters"
}

func (v proxyPasswordValidator) MarkdownDescription(_ context.Context) string {
	return "Proxy password must be between 1 and 100 characters"
}

func (v proxyPasswordValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	password := req.ConfigValue.ValueString()
	if len(password) < 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Password",
			"Proxy password must be at least 1 character long",
		)
		return
	}
	if len(password) > 100 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Password",
			"Proxy password must be at most 100 characters long",
		)
		return
	}
}

type proxyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Host        types.String `tfsdk:"host"`
	Port        types.Int64  `tfsdk:"port"`
	Protocol    types.String `tfsdk:"protocol"`
	Auth        types.Bool   `tfsdk:"auth"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	CreatedDate types.String `tfsdk:"created_date"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func (r *ProxyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxy"
}

func (r *ProxyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Required:    true,
				Description: "Proxy host",
				Validators: []validator.String{
					proxyHostValidator{},
				},
			},
			"port": schema.Int64Attribute{
				Required:    true,
				Description: "Proxy port (1-65535)",
				Validators: []validator.Int64{
					proxyPortValidator{},
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Proxy protocol (http, https, socks, socks4, socks5, socks5h)",
				Validators: []validator.String{
					proxyProtocolValidator{},
				},
			},
			"auth": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether authentication is required",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Username for authentication",
				Validators: []validator.String{
					proxyUsernameValidator{},
				},
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password for authentication",
				Validators: []validator.String{
					proxyPasswordValidator{},
				},
			},
			"created_date": schema.StringAttribute{
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

func (r *ProxyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProxyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan proxyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the proxy we're creating for debugging
	tflog.Info(ctx, "Creating proxy", map[string]interface{}{
		"host":     plan.Host.ValueString(),
		"port":     plan.Port.ValueInt64(),
		"protocol": plan.Protocol.ValueString(),
		"auth":     plan.Auth.ValueBool(),
	})

	in := peekaping.ProxyCreate{
		Host:     plan.Host.ValueString(),
		Port:     int(plan.Port.ValueInt64()),
		Protocol: peekaping.ProxyProtocol(plan.Protocol.ValueString()),
		Auth:     plan.Auth.ValueBool(),
	}
	if !plan.Username.IsNull() {
		in.Username = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		in.Password = plan.Password.ValueString()
	}

	p, err := r.client.CreateProxy(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError("create proxy failed", err.Error())
		return
	}
	setModelFromProxy(&plan, p)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProxyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state proxyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p, err := r.client.GetProxy(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read proxy failed", err.Error())
		return
	}
	setModelFromProxy(&state, p)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProxyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan proxyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to get the ID
	var state proxyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upd := peekaping.ProxyUpdate{}
	if !plan.Host.IsNull() {
		v := plan.Host.ValueString()
		upd.Host = &v
	}
	if !plan.Port.IsNull() {
		v := int(plan.Port.ValueInt64())
		upd.Port = &v
	}
	if !plan.Protocol.IsNull() {
		v := peekaping.ProxyProtocol(plan.Protocol.ValueString())
		upd.Protocol = &v
	}
	if !plan.Auth.IsNull() {
		v := plan.Auth.ValueBool()
		upd.Auth = &v
	}
	if !plan.Username.IsNull() {
		v := plan.Username.ValueString()
		upd.Username = &v
	}
	if !plan.Password.IsNull() {
		v := plan.Password.ValueString()
		upd.Password = &v
	}

	// Use state.ID instead of plan.ID
	p, err := r.client.UpdateProxy(ctx, state.ID.ValueString(), upd)
	if err != nil {
		resp.Diagnostics.AddError("update proxy failed", err.Error())
		return
	}
	setModelFromProxy(&plan, p)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProxyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state proxyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteProxy(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("delete proxy failed", err.Error())
		return
	}
}

func (r *ProxyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state proxyResourceModel
	state.ID = types.StringValue(req.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func setModelFromProxy(m *proxyResourceModel, from *peekaping.Proxy) {
	m.ID = types.StringValue(from.ID)
	m.Host = types.StringValue(from.Host)
	m.Port = types.Int64Value(int64(from.Port))
	m.Protocol = types.StringValue(string(from.Protocol))
	m.Auth = types.BoolValue(from.Auth)

	// Only set username/password if auth is true
	if from.Auth && from.Username != "" {
		m.Username = types.StringValue(from.Username)
	} else {
		m.Username = types.StringNull()
	}
	if from.Auth && from.Password != "" {
		m.Password = types.StringValue(from.Password)
	} else {
		m.Password = types.StringNull()
	}
	if from.CreatedDate != "" {
		m.CreatedDate = types.StringValue(from.CreatedDate)
	} else {
		m.CreatedDate = types.StringNull()
	}
	if from.UpdatedAt != "" {
		m.UpdatedAt = types.StringValue(from.UpdatedAt)
	} else {
		m.UpdatedAt = types.StringNull()
	}
}
