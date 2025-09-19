// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tafaust/terraform-provider-peekaping/internal/peekaping"
)

var _ datasource.DataSource = &StatusPageDataSource{}

type StatusPageDataSource struct {
	client *peekaping.Client
}

func NewStatusPageDataSource() datasource.DataSource { return &StatusPageDataSource{} }

type statusPageDataSourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	Title                 types.String   `tfsdk:"title"`
	Description           types.String   `tfsdk:"description"`
	Slug                  types.String   `tfsdk:"slug"`
	Domains               []types.String `tfsdk:"domains"`
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
}

func (d *StatusPageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (d *StatusPageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                      schema.StringAttribute{Optional: true, Description: "Status page ID"},
			"title":                   schema.StringAttribute{Optional: true, Description: "Status page title"},
			"description":             schema.StringAttribute{Computed: true, Description: "Status page description"},
			"slug":                    schema.StringAttribute{Optional: true, Description: "Status page slug"},
			"domains":                 schema.ListAttribute{Computed: true, ElementType: types.StringType, Description: "List of custom domains"},
			"published":               schema.BoolAttribute{Computed: true, Description: "Whether the status page is published"},
			"theme":                   schema.StringAttribute{Computed: true, Description: "Status page theme"},
			"icon":                    schema.StringAttribute{Computed: true, Description: "Status page icon"},
			"footer_text":             schema.StringAttribute{Computed: true, Description: "Footer text"},
			"custom_css":              schema.StringAttribute{Computed: true, Description: "Custom CSS"},
			"google_analytics_tag_id": schema.StringAttribute{Computed: true, Description: "Google Analytics tag ID"},
			"auto_refresh_interval":   schema.Int64Attribute{Computed: true, Description: "Auto refresh interval in seconds"},
			"search_engine_index":     schema.BoolAttribute{Computed: true, Description: "Whether to allow search engine indexing"},
			"show_certificate_expiry": schema.BoolAttribute{Computed: true, Description: "Whether to show certificate expiry information"},
			"show_powered_by":         schema.BoolAttribute{Computed: true, Description: "Whether to show 'Powered by' text"},
			"show_tags":               schema.BoolAttribute{Computed: true, Description: "Whether to show tags"},
		},
	}
}

func (d *StatusPageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*peekaping.Client)
}

func (d *StatusPageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data statusPageDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		sp, err := d.client.GetStatusPage(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("lookup by id failed", err.Error())
			return
		}
		data.ID = types.StringValue(sp.ID)
		data.Title = types.StringValue(sp.Title)
		if sp.Description != "" {
			data.Description = types.StringValue(sp.Description)
		} else {
			data.Description = types.StringNull()
		}
		if sp.Slug != "" {
			data.Slug = types.StringValue(sp.Slug)
		} else {
			data.Slug = types.StringNull()
		}
		if len(sp.Domains) > 0 {
			domains := make([]types.String, 0, len(sp.Domains))
			for _, domain := range sp.Domains {
				domains = append(domains, types.StringValue(domain))
			}
			data.Domains = domains
		} else {
			data.Domains = nil
		}
		data.Published = types.BoolValue(sp.Published)
		if sp.Theme != "" {
			data.Theme = types.StringValue(sp.Theme)
		} else {
			data.Theme = types.StringNull()
		}
		if sp.Icon != "" {
			data.Icon = types.StringValue(sp.Icon)
		} else {
			data.Icon = types.StringNull()
		}
		if sp.FooterText != "" {
			data.FooterText = types.StringValue(sp.FooterText)
		} else {
			data.FooterText = types.StringNull()
		}
		if sp.CustomCSS != "" {
			data.CustomCSS = types.StringValue(sp.CustomCSS)
		} else {
			data.CustomCSS = types.StringNull()
		}
		if sp.GoogleAnalyticsTagID != "" {
			data.GoogleAnalyticsTagID = types.StringValue(sp.GoogleAnalyticsTagID)
		} else {
			data.GoogleAnalyticsTagID = types.StringNull()
		}
		if sp.AutoRefreshInterval != 0 {
			data.AutoRefreshInterval = types.Int64Value(int64(sp.AutoRefreshInterval))
		} else {
			data.AutoRefreshInterval = types.Int64Null()
		}
		data.SearchEngineIndex = types.BoolValue(sp.SearchEngineIndex)
		data.ShowCertificateExpiry = types.BoolValue(sp.ShowCertificateExpiry)
		data.ShowPoweredBy = types.BoolValue(sp.ShowPoweredBy)
		data.ShowTags = types.BoolValue(sp.ShowTags)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	list, err := d.client.ListStatusPages(ctx)
	if err != nil {
		resp.Diagnostics.AddError("list status pages failed", err.Error())
		return
	}
	title := strings.ToLower(data.Title.ValueString())
	slug := strings.ToLower(data.Slug.ValueString())
	for _, sp := range list.Items {
		if (title != "" && strings.ToLower(sp.Title) == title) ||
			(slug != "" && strings.ToLower(sp.Slug) == slug) ||
			(title != "" && strings.Contains(strings.ToLower(sp.Title), title)) {
			data.ID = types.StringValue(sp.ID)
			data.Title = types.StringValue(sp.Title)
			if sp.Description != "" {
				data.Description = types.StringValue(sp.Description)
			} else {
				data.Description = types.StringNull()
			}
			if sp.Slug != "" {
				data.Slug = types.StringValue(sp.Slug)
			} else {
				data.Slug = types.StringNull()
			}
			if len(sp.Domains) > 0 {
				domains := make([]types.String, 0, len(sp.Domains))
				for _, domain := range sp.Domains {
					domains = append(domains, types.StringValue(domain))
				}
				data.Domains = domains
			} else {
				data.Domains = nil
			}
			data.Published = types.BoolValue(sp.Published)
			if sp.Theme != "" {
				data.Theme = types.StringValue(sp.Theme)
			} else {
				data.Theme = types.StringNull()
			}
			if sp.Icon != "" {
				data.Icon = types.StringValue(sp.Icon)
			} else {
				data.Icon = types.StringNull()
			}
			if sp.FooterText != "" {
				data.FooterText = types.StringValue(sp.FooterText)
			} else {
				data.FooterText = types.StringNull()
			}
			if sp.CustomCSS != "" {
				data.CustomCSS = types.StringValue(sp.CustomCSS)
			} else {
				data.CustomCSS = types.StringNull()
			}
			if sp.GoogleAnalyticsTagID != "" {
				data.GoogleAnalyticsTagID = types.StringValue(sp.GoogleAnalyticsTagID)
			} else {
				data.GoogleAnalyticsTagID = types.StringNull()
			}
			if sp.AutoRefreshInterval != 0 {
				data.AutoRefreshInterval = types.Int64Value(int64(sp.AutoRefreshInterval))
			} else {
				data.AutoRefreshInterval = types.Int64Null()
			}
			data.SearchEngineIndex = types.BoolValue(sp.SearchEngineIndex)
			data.ShowCertificateExpiry = types.BoolValue(sp.ShowCertificateExpiry)
			data.ShowPoweredBy = types.BoolValue(sp.ShowPoweredBy)
			data.ShowTags = types.BoolValue(sp.ShowTags)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
	}
	resp.Diagnostics.AddError("not found", "no status page matched the given criteria")
}
