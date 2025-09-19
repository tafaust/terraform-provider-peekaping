// Copyright (c) 2025 tafaust
// SPDX-License-Identifier: MIT

package peekaping

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const apiPrefix = "/api/v1"

type Client struct {
	Endpoint     string
	HTTP         *http.Client
	accessToken  string
	refreshToken string
	email        string
	password     string
	token        string
}

type Option func(*Client)

func WithCredentials(email, pass string) Option {
	return func(c *Client) { c.email, c.password = email, pass }
}

func WithToken(token string) Option {
	return func(c *Client) { c.token = token }
}

func New(endpoint string, opts ...Option) *Client {
	c := &Client{
		Endpoint: strings.TrimRight(endpoint, "/"),
		HTTP:     &http.Client{Timeout: 30 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token,omitempty"`
}

type loginResponse struct {
	Data struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	} `json:"data"`
	Message string `json:"message"`
}

type refreshResponse struct {
	Data struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	} `json:"data"`
	Message string `json:"message"`
}

type monitorResponse struct {
	Data    Monitor `json:"data"`
	Message string  `json:"message"`
}

type monitorsResponse struct {
	Data    []Monitor `json:"data"`
	Message string    `json:"message"`
}

type notificationsResponse struct {
	Data    []Notification `json:"data"`
	Message string         `json:"message"`
}

type notificationResponse struct {
	Data    Notification `json:"data"`
	Message string       `json:"message"`
}

type tagsResponse struct {
	Data    []Tag  `json:"data"`
	Message string `json:"message"`
}

type tagResponse struct {
	Data    Tag    `json:"data"`
	Message string `json:"message"`
}

type maintenancesResponse struct {
	Data    []Maintenance `json:"data"`
	Message string        `json:"message"`
}

type maintenanceResponse struct {
	Data    Maintenance `json:"data"`
	Message string      `json:"message"`
}

type statusPagesResponse struct {
	Data    []StatusPage `json:"data"`
	Message string       `json:"message"`
}

type statusPageResponse struct {
	Data    StatusPage `json:"data"`
	Message string     `json:"message"`
}

type proxiesResponse struct {
	Data    []Proxy `json:"data"`
	Message string  `json:"message"`
}

type proxyResponse struct {
	Data    Proxy  `json:"data"`
	Message string `json:"message"`
}

func (c *Client) Login(ctx context.Context) error {
	body := loginRequest{Email: c.email, Password: c.password, Token: c.token}
	req, err := c.newReq(ctx, http.MethodPost, "/auth/login", body)
	if err != nil {
		return err
	}
	var out loginResponse
	if err := c.do(req, &out); err != nil {
		return err
	}
	c.accessToken, c.refreshToken = out.Data.AccessToken, out.Data.RefreshToken
	return nil
}

func (c *Client) newReq(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = bytes.NewReader(b)
	}
	u := c.Endpoint + apiPrefix + path
	req, err := http.NewRequestWithContext(ctx, method, u, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	// Debug logging
	fmt.Printf("DEBUG: Making %s request to %s\n", method, u)
	if body != nil {
		if configBytes, err := json.MarshalIndent(body, "", "  "); err == nil {
			fmt.Printf("DEBUG: Request body:\n%s\n", string(configBytes))
		}
	}

	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode == http.StatusUnauthorized && c.refreshToken != "" && !strings.Contains(req.URL.Path, "/auth/refresh") {
		// try refresh once
		_ = c.refresh(req.Context())
		// retry
		req2 := req.Clone(req.Context())
		req2.Header.Set("Authorization", "Bearer "+c.accessToken)
		res2, err2 := c.HTTP.Do(req2)
		if err2 != nil {
			return err2
		}
		defer func() { _ = res2.Body.Close() }()
		if res2.StatusCode >= 300 {
			b, _ := io.ReadAll(res2.Body)

			// Try to parse as JSON error response
			var errorResp struct {
				Message string `json:"message"`
				Data    any    `json:"data"`
			}
			if err := json.Unmarshal(b, &errorResp); err == nil && errorResp.Message != "" {
				return fmt.Errorf("http %d: %s", res2.StatusCode, errorResp.Message)
			}

			// Fallback to raw response
			return fmt.Errorf("http %d: %s", res2.StatusCode, string(b))
		}
		if out == nil {
			return nil
		}

		// Debug logging for successful responses (retry case)
		bodyBytes, err := io.ReadAll(res2.Body)
		if err != nil {
			return err
		}
		fmt.Printf("DEBUG: Response body (retry): %s\n", string(bodyBytes))

		return json.Unmarshal(bodyBytes, out)
	}

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)

		// Debug logging for errors
		fmt.Printf("DEBUG: HTTP Error %d for %s %s\n", res.StatusCode, req.Method, req.URL.String())
		fmt.Printf("DEBUG: Response body: %s\n", string(b))

		// Try to parse as JSON error response
		var errorResp struct {
			Message string `json:"message"`
			Data    any    `json:"data"`
		}
		if err := json.Unmarshal(b, &errorResp); err == nil && errorResp.Message != "" {
			return fmt.Errorf("http %d: %s", res.StatusCode, errorResp.Message)
		}

		// Fallback to raw response
		return fmt.Errorf("http %d: %s", res.StatusCode, string(b))
	}
	if out == nil {
		return nil
	}

	// Debug logging for successful responses
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("DEBUG: Response body: %s\n", string(bodyBytes))

	return json.Unmarshal(bodyBytes, out)
}

func (c *Client) refresh(ctx context.Context) error {
	req, err := c.newReq(ctx, http.MethodPost, "/auth/refresh", map[string]string{"refreshToken": c.refreshToken})
	if err != nil {
		return err
	}
	var out refreshResponse
	if err := c.do(req, &out); err != nil {
		return err
	}
	c.accessToken = out.Data.AccessToken
	return nil
}

// ---- API: Monitors ----

type MonitorType string

const (
	MonitorHTTP MonitorType = "http"
	MonitorTCP  MonitorType = "tcp"
	MonitorPing MonitorType = "ping"
	MonitorDNS  MonitorType = "dns"
	MonitorPush MonitorType = "push"
	MonitorGRPC MonitorType = "grpc"
)

// IsValid checks if the monitor type is valid
func (mt MonitorType) IsValid() bool {
	switch mt {
	case MonitorHTTP, MonitorTCP, MonitorPing, MonitorDNS, MonitorPush, MonitorGRPC:
		return true
	default:
		return false
	}
}

type MonitorStatus int

const (
	MonitorStatusDown        MonitorStatus = 0
	MonitorStatusUp          MonitorStatus = 1
	MonitorStatusPending     MonitorStatus = 2
	MonitorStatusMaintenance MonitorStatus = 3
)

// IsValid checks if the monitor status is valid
func (ms MonitorStatus) IsValid() bool {
	switch ms {
	case MonitorStatusDown, MonitorStatusUp, MonitorStatusPending, MonitorStatusMaintenance:
		return true
	default:
		return false
	}
}

type Monitor struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Type            MonitorType   `json:"type"`
	Config          string        `json:"config,omitempty"`
	Interval        int64         `json:"interval,omitempty"`
	Active          bool          `json:"active,omitempty"`
	Timeout         int64         `json:"timeout,omitempty"`
	MaxRetries      int64         `json:"max_retries,omitempty"`
	RetryInterval   int64         `json:"retry_interval,omitempty"`
	ResendInterval  int64         `json:"resend_interval,omitempty"`
	ProxyID         string        `json:"proxy_id,omitempty"`
	PushToken       string        `json:"push_token,omitempty"`
	NotificationIDs []string      `json:"notification_ids,omitempty"`
	TagIDs          []string      `json:"tag_ids,omitempty"`
	Status          MonitorStatus `json:"status,omitempty"`
	CreatedAt       string        `json:"created_at,omitempty"`
	UpdatedAt       string        `json:"updated_at,omitempty"`
}

type MonitorCreate struct {
	Name            string      `json:"name"`
	Type            MonitorType `json:"type"`
	Config          string      `json:"config,omitempty"`
	Interval        int64       `json:"interval,omitempty"`
	Active          bool        `json:"active,omitempty"`
	Timeout         int64       `json:"timeout,omitempty"`
	MaxRetries      int64       `json:"max_retries,omitempty"`
	RetryInterval   int64       `json:"retry_interval,omitempty"`
	ResendInterval  int64       `json:"resend_interval,omitempty"`
	ProxyID         string      `json:"proxy_id,omitempty"`
	PushToken       string      `json:"push_token,omitempty"`
	NotificationIDs []string    `json:"notification_ids"`
	TagIDs          []string    `json:"tag_ids,omitempty"`
}

type MonitorUpdate struct {
	Name            *string      `json:"name,omitempty"`
	Type            *MonitorType `json:"type,omitempty"`
	Config          *string      `json:"config,omitempty"`
	Interval        *int64       `json:"interval,omitempty"`
	Active          *bool        `json:"active,omitempty"`
	Timeout         *int64       `json:"timeout,omitempty"`
	MaxRetries      *int64       `json:"max_retries,omitempty"`
	RetryInterval   *int64       `json:"retry_interval,omitempty"`
	ResendInterval  *int64       `json:"resend_interval,omitempty"`
	ProxyID         *string      `json:"proxy_id,omitempty"`
	PushToken       *string      `json:"push_token,omitempty"`
	NotificationIDs []string     `json:"notification_ids,omitempty"`
	TagIDs          []string     `json:"tag_ids,omitempty"`
}

type ListMonitorsResp struct {
	Items []Monitor `json:"items"`
	Total int       `json:"total"`
	Page  int       `json:"page,omitempty"`
	Size  int       `json:"size,omitempty"`
}

func (c *Client) ListMonitors(ctx context.Context) (*ListMonitorsResp, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/monitors", nil)
	if err != nil {
		return nil, err
	}
	var out monitorsResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &ListMonitorsResp{Items: out.Data, Total: len(out.Data)}, nil
}

func (c *Client) CreateMonitor(ctx context.Context, in MonitorCreate) (*Monitor, error) {
	// Log the request for debugging
	if configBytes, err := json.MarshalIndent(in, "", "  "); err == nil {
		fmt.Printf("DEBUG: Creating monitor with config:\n%s\n", string(configBytes))
	}

	req, err := c.newReq(ctx, http.MethodPost, "/monitors", in)
	if err != nil {
		return nil, err
	}
	var out monitorResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetMonitor(ctx context.Context, id string) (*Monitor, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/monitors/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	var out monitorResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) UpdateMonitor(ctx context.Context, id string, in MonitorUpdate) (*Monitor, error) {
	fmt.Printf("*** UPDATE MONITOR CALLED ***\n")
	fmt.Printf("ID: %s\n", id)
	fmt.Printf("Update data: %+v\n", in)
	req, err := c.newReq(ctx, http.MethodPut, "/monitors/"+url.PathEscape(id), in)
	if err != nil {
		fmt.Printf("ERROR creating request: %v\n", err)
		return nil, err
	}
	var out monitorResponse
	if err := c.do(req, &out); err != nil {
		fmt.Printf("ERROR in do: %v\n", err)
		return nil, err
	}
	fmt.Printf("SUCCESS: %+v\n", out.Data)
	return &out.Data, nil
}

func (c *Client) DeleteMonitor(ctx context.Context, id string) error {
	req, err := c.newReq(ctx, http.MethodDelete, "/monitors/"+url.PathEscape(id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// ---- API: Notifications ----

type Notification struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Config    string `json:"config"`
	Active    bool   `json:"active,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type NotificationCreate struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Config    string `json:"config"`
	Active    bool   `json:"active,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
}

type NotificationUpdate struct {
	Name      *string `json:"name,omitempty"`
	Type      *string `json:"type,omitempty"`
	Config    *string `json:"config,omitempty"`
	Active    *bool   `json:"active,omitempty"`
	IsDefault *bool   `json:"is_default,omitempty"`
}

type ListNotificationsResp struct {
	Items []Notification `json:"items"`
	Total int            `json:"total"`
	Page  int            `json:"page,omitempty"`
	Size  int            `json:"size,omitempty"`
}

func (c *Client) ListNotifications(ctx context.Context) (*ListNotificationsResp, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/notification-channels", nil)
	if err != nil {
		return nil, err
	}
	var out notificationsResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &ListNotificationsResp{Items: out.Data, Total: len(out.Data)}, nil
}

func (c *Client) CreateNotification(ctx context.Context, in NotificationCreate) (*Notification, error) {
	req, err := c.newReq(ctx, http.MethodPost, "/notification-channels", in)
	if err != nil {
		return nil, err
	}
	var out notificationResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetNotification(ctx context.Context, id string) (*Notification, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/notification-channels/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	var out notificationResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) UpdateNotification(ctx context.Context, id string, in NotificationUpdate) (*Notification, error) {
	req, err := c.newReq(ctx, http.MethodPatch, "/notification-channels/"+url.PathEscape(id), in)
	if err != nil {
		return nil, err
	}
	var out notificationResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) DeleteNotification(ctx context.Context, id string) error {
	req, err := c.newReq(ctx, http.MethodDelete, "/notification-channels/"+url.PathEscape(id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// ---- API: Tags ----

type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type TagCreate struct {
	Name        string `json:"name"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

type TagUpdate struct {
	Name        *string `json:"name,omitempty"`
	Color       *string `json:"color,omitempty"`
	Description *string `json:"description,omitempty"`
}

type ListTagsResp struct {
	Items []Tag `json:"items"`
	Total int   `json:"total"`
	Page  int   `json:"page,omitempty"`
	Size  int   `json:"size,omitempty"`
}

func (c *Client) ListTags(ctx context.Context) (*ListTagsResp, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/tags", nil)
	if err != nil {
		return nil, err
	}
	var out tagsResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &ListTagsResp{Items: out.Data, Total: len(out.Data)}, nil
}

func (c *Client) CreateTag(ctx context.Context, in TagCreate) (*Tag, error) {
	req, err := c.newReq(ctx, http.MethodPost, "/tags", in)
	if err != nil {
		return nil, err
	}
	var out tagResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetTag(ctx context.Context, id string) (*Tag, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/tags/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	var out tagResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) UpdateTag(ctx context.Context, id string, in TagUpdate) (*Tag, error) {
	req, err := c.newReq(ctx, http.MethodPatch, "/tags/"+url.PathEscape(id), in)
	if err != nil {
		return nil, err
	}
	var out tagResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) DeleteTag(ctx context.Context, id string) error {
	req, err := c.newReq(ctx, http.MethodDelete, "/tags/"+url.PathEscape(id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// ---- API: Maintenance ----

type Maintenance struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description,omitempty"`
	Strategy      string   `json:"strategy"`
	Active        bool     `json:"active"`
	MonitorIDs    []string `json:"monitor_ids,omitempty"`
	StartDateTime string   `json:"start_date_time,omitempty"`
	EndDateTime   string   `json:"end_date_time,omitempty"`
	Duration      int      `json:"duration,omitempty"`
	Timezone      string   `json:"timezone,omitempty"`
	Cron          string   `json:"cron,omitempty"`
	Weekdays      []int    `json:"weekdays,omitempty"`
	DaysOfMonth   []int    `json:"days_of_month,omitempty"`
	IntervalDay   int      `json:"interval_day,omitempty"`
	StartTime     string   `json:"start_time,omitempty"`
	EndTime       string   `json:"end_time,omitempty"`
	CreatedAt     string   `json:"created_at,omitempty"`
	UpdatedAt     string   `json:"updated_at,omitempty"`
}

type MaintenanceCreate struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Strategy    string   `json:"strategy"`
	Active      bool     `json:"active,omitempty"`
	MonitorIDs  []string `json:"monitor_ids,omitempty"`
	Duration    int      `json:"duration,omitempty"`
	Timezone    string   `json:"timezone,omitempty"`
	Cron        string   `json:"cron,omitempty"`
	Weekdays    []int    `json:"weekdays,omitempty"`
	DaysOfMonth []int    `json:"days_of_month,omitempty"`
	IntervalDay int      `json:"interval_day,omitempty"`
	StartTime   string   `json:"start_time,omitempty"`
	EndTime     string   `json:"end_time,omitempty"`
}

type MaintenanceUpdate struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Strategy    *string  `json:"strategy,omitempty"`
	Active      *bool    `json:"active,omitempty"`
	MonitorIDs  []string `json:"monitor_ids,omitempty"`
	Duration    *int     `json:"duration,omitempty"`
	Timezone    *string  `json:"timezone,omitempty"`
	Cron        *string  `json:"cron,omitempty"`
	Weekdays    []int    `json:"weekdays,omitempty"`
	DaysOfMonth []int    `json:"days_of_month,omitempty"`
	IntervalDay *int     `json:"interval_day,omitempty"`
	StartTime   *string  `json:"start_time,omitempty"`
	EndTime     *string  `json:"end_time,omitempty"`
}

type ListMaintenanceResp struct {
	Items []Maintenance `json:"items"`
	Total int           `json:"total"`
	Page  int           `json:"page,omitempty"`
	Size  int           `json:"size,omitempty"`
}

func (c *Client) ListMaintenance(ctx context.Context) (*ListMaintenanceResp, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/maintenances", nil)
	if err != nil {
		return nil, err
	}
	var out maintenancesResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &ListMaintenanceResp{Items: out.Data, Total: len(out.Data)}, nil
}

func (c *Client) CreateMaintenance(ctx context.Context, in MaintenanceCreate) (*Maintenance, error) {
	req, err := c.newReq(ctx, http.MethodPost, "/maintenances", in)
	if err != nil {
		return nil, err
	}
	var out maintenanceResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetMaintenance(ctx context.Context, id string) (*Maintenance, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/maintenances/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	var out maintenanceResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) UpdateMaintenance(ctx context.Context, id string, in MaintenanceUpdate) (*Maintenance, error) {
	req, err := c.newReq(ctx, http.MethodPatch, "/maintenances/"+url.PathEscape(id), in)
	if err != nil {
		return nil, err
	}
	var out maintenanceResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) DeleteMaintenance(ctx context.Context, id string) error {
	req, err := c.newReq(ctx, http.MethodDelete, "/maintenances/"+url.PathEscape(id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// ---- API: Status Pages ----

type StatusPage struct {
	ID                    string   `json:"id"`
	Title                 string   `json:"title"`
	Description           string   `json:"description,omitempty"`
	Slug                  string   `json:"slug,omitempty"`
	Domains               []string `json:"domains,omitempty"`
	MonitorIDs            []string `json:"monitor_ids,omitempty"`
	Published             bool     `json:"published,omitempty"`
	Theme                 string   `json:"theme,omitempty"`
	Icon                  string   `json:"icon,omitempty"`
	FooterText            string   `json:"footer_text,omitempty"`
	CustomCSS             string   `json:"custom_css,omitempty"`
	GoogleAnalyticsTagID  string   `json:"google_analytics_tag_id,omitempty"`
	AutoRefreshInterval   int      `json:"auto_refresh_interval,omitempty"`
	SearchEngineIndex     bool     `json:"search_engine_index,omitempty"`
	ShowCertificateExpiry bool     `json:"show_certificate_expiry,omitempty"`
	ShowPoweredBy         bool     `json:"show_powered_by,omitempty"`
	ShowTags              bool     `json:"show_tags,omitempty"`
	Password              string   `json:"password,omitempty"`
	CreatedAt             string   `json:"created_at,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
}

type StatusPageCreate struct {
	Title                 string   `json:"title"`
	Description           string   `json:"description,omitempty"`
	Slug                  string   `json:"slug,omitempty"`
	Domains               []string `json:"domains,omitempty"`
	MonitorIDs            []string `json:"monitor_ids,omitempty"`
	Published             bool     `json:"published,omitempty"`
	Theme                 string   `json:"theme,omitempty"`
	Icon                  string   `json:"icon,omitempty"`
	FooterText            string   `json:"footer_text,omitempty"`
	CustomCSS             string   `json:"custom_css,omitempty"`
	GoogleAnalyticsTagID  string   `json:"google_analytics_tag_id,omitempty"`
	SearchEngineIndex     bool     `json:"search_engine_index,omitempty"`
	ShowCertificateExpiry bool     `json:"show_certificate_expiry,omitempty"`
	ShowPoweredBy         bool     `json:"show_powered_by,omitempty"`
	ShowTags              bool     `json:"show_tags,omitempty"`
	Password              string   `json:"password,omitempty"`
}

type StatusPageUpdate struct {
	Title                 *string  `json:"title,omitempty"`
	Description           *string  `json:"description,omitempty"`
	Slug                  *string  `json:"slug,omitempty"`
	Domains               []string `json:"domains,omitempty"`
	MonitorIDs            []string `json:"monitor_ids,omitempty"`
	Published             *bool    `json:"published,omitempty"`
	Theme                 *string  `json:"theme,omitempty"`
	Icon                  *string  `json:"icon,omitempty"`
	FooterText            *string  `json:"footer_text,omitempty"`
	CustomCSS             *string  `json:"custom_css,omitempty"`
	GoogleAnalyticsTagID  *string  `json:"google_analytics_tag_id,omitempty"`
	SearchEngineIndex     *bool    `json:"search_engine_index,omitempty"`
	ShowCertificateExpiry *bool    `json:"show_certificate_expiry,omitempty"`
	ShowPoweredBy         *bool    `json:"show_powered_by,omitempty"`
	ShowTags              *bool    `json:"show_tags,omitempty"`
	Password              *string  `json:"password,omitempty"`
}

type ListStatusPagesResp struct {
	Items []StatusPage `json:"items"`
	Total int          `json:"total"`
	Page  int          `json:"page,omitempty"`
	Size  int          `json:"size,omitempty"`
}

func (c *Client) ListStatusPages(ctx context.Context) (*ListStatusPagesResp, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/status-pages", nil)
	if err != nil {
		return nil, err
	}
	var out statusPagesResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &ListStatusPagesResp{Items: out.Data, Total: len(out.Data)}, nil
}

func (c *Client) CreateStatusPage(ctx context.Context, in StatusPageCreate) (*StatusPage, error) {
	req, err := c.newReq(ctx, http.MethodPost, "/status-pages", in)
	if err != nil {
		return nil, err
	}
	var out statusPageResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetStatusPage(ctx context.Context, id string) (*StatusPage, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/status-pages/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	var out statusPageResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) UpdateStatusPage(ctx context.Context, id string, in StatusPageUpdate) (*StatusPage, error) {
	req, err := c.newReq(ctx, http.MethodPatch, "/status-pages/"+url.PathEscape(id), in)
	if err != nil {
		return nil, err
	}
	var out statusPageResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) DeleteStatusPage(ctx context.Context, id string) error {
	req, err := c.newReq(ctx, http.MethodDelete, "/status-pages/"+url.PathEscape(id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// ---- API: Proxies ----

type ProxyProtocol string

const (
	ProxyHTTP    ProxyProtocol = "http"
	ProxyHTTPS   ProxyProtocol = "https"
	ProxySOCKS   ProxyProtocol = "socks"
	ProxySOCKS4  ProxyProtocol = "socks4"
	ProxySOCKS5  ProxyProtocol = "socks5"
	ProxySOCKS5H ProxyProtocol = "socks5h"
)

// IsValid checks if the proxy protocol is valid
func (pp ProxyProtocol) IsValid() bool {
	switch pp {
	case ProxyHTTP, ProxyHTTPS, ProxySOCKS, ProxySOCKS4, ProxySOCKS5, ProxySOCKS5H:
		return true
	default:
		return false
	}
}

type Proxy struct {
	ID          string        `json:"id"`
	Host        string        `json:"host"`
	Port        int           `json:"port"`
	Protocol    ProxyProtocol `json:"protocol"`
	Auth        bool          `json:"auth,omitempty"`
	Username    string        `json:"username,omitempty"`
	Password    string        `json:"password,omitempty"`
	CreatedDate string        `json:"createdDate,omitempty"`
	UpdatedAt   string        `json:"updatedAt,omitempty"`
}

type ProxyCreate struct {
	Host     string        `json:"host"`
	Port     int           `json:"port"`
	Protocol ProxyProtocol `json:"protocol"`
	Auth     bool          `json:"auth,omitempty"`
	Username string        `json:"username,omitempty"`
	Password string        `json:"password,omitempty"`
}

type ProxyUpdate struct {
	Host     *string        `json:"host,omitempty"`
	Port     *int           `json:"port,omitempty"`
	Protocol *ProxyProtocol `json:"protocol,omitempty"`
	Auth     *bool          `json:"auth,omitempty"`
	Username *string        `json:"username,omitempty"`
	Password *string        `json:"password,omitempty"`
}

type ListProxiesResp struct {
	Items []Proxy `json:"items"`
	Total int     `json:"total"`
	Page  int     `json:"page,omitempty"`
	Size  int     `json:"size,omitempty"`
}

func (c *Client) ListProxies(ctx context.Context) (*ListProxiesResp, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/proxies", nil)
	if err != nil {
		return nil, err
	}
	var out proxiesResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &ListProxiesResp{Items: out.Data, Total: len(out.Data)}, nil
}

func (c *Client) CreateProxy(ctx context.Context, in ProxyCreate) (*Proxy, error) {
	req, err := c.newReq(ctx, http.MethodPost, "/proxies", in)
	if err != nil {
		return nil, err
	}
	var out proxyResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetProxy(ctx context.Context, id string) (*Proxy, error) {
	req, err := c.newReq(ctx, http.MethodGet, "/proxies/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}
	var out proxyResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) UpdateProxy(ctx context.Context, id string, in ProxyUpdate) (*Proxy, error) {
	fmt.Printf("=== UPDATE PROXY CALLED ===\n")
	fmt.Printf("ID: %s\n", id)
	fmt.Printf("Update data: %+v\n", in)
	req, err := c.newReq(ctx, http.MethodPatch, "/proxies/"+url.PathEscape(id), in)
	if err != nil {
		fmt.Printf("ERROR creating request: %v\n", err)
		return nil, err
	}
	var out proxyResponse
	if err := c.do(req, &out); err != nil {
		fmt.Printf("ERROR in do: %v\n", err)
		return nil, err
	}
	fmt.Printf("SUCCESS: %+v\n", out.Data)
	return &out.Data, nil
}

func (c *Client) DeleteProxy(ctx context.Context, id string) error {
	req, err := c.newReq(ctx, http.MethodDelete, "/proxies/"+url.PathEscape(id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
