/**
 * Copyright 2026-present Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mcp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase/prime-sdk-go/model"
	"github.com/coinbase/prime-sdk-go/webhooks"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	webhookAuthTypePrefix       = "WEBHOOK_AUTH_TYPE_"
	webhookDeliveryStatusPrefix = "WEBHOOK_DELIVERY_STATUS_"
	webhookStatsPeriodPrefix    = "WEBHOOK_STATS_PERIOD_"
)

func registerWebhookTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("get_webhook_event_types",
		mcplib.WithDescription("List the event types that can be subscribed to with webhooks (beta)"),
	), handleGetWebhookEventTypes)

	s.AddTool(mcplib.NewTool("list_webhook_subscriptions",
		mcplib.WithDescription("List webhook subscriptions for a portfolio (beta)"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithBoolean("include_disabled",
			mcplib.Description("Include disabled subscriptions"),
		),
		mcplib.WithString("event_type_filter",
			mcplib.Description("Only return subscriptions for this event type"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("next_cursor from a previous response"),
		),
		mcplib.WithBoolean("fetch_all",
			mcplib.Description("Fetch all pages automatically and return combined results. When true, cursor is ignored."),
		),
	), handleListWebhookSubscriptions)

	s.AddTool(mcplib.NewTool("create_webhook_subscription",
		mcplib.WithDescription("Create a webhook subscription (beta). The HMAC secret is returned only once in the response"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithArray("event_types",
			mcplib.Required(),
			mcplib.Description("Event types to subscribe to, from get_webhook_event_types"),
			mcplib.WithStringItems(),
		),
		mcplib.WithString("target_url",
			mcplib.Required(),
			mcplib.Description("HTTPS URL that receives deliveries"),
		),
		mcplib.WithString("auth_type",
			mcplib.Description("HMAC, OAUTH, MTLS or OAUTH_MTLS (the WEBHOOK_AUTH_TYPE_ prefix is optional)"),
		),
		mcplib.WithString("oauth_token_url",
			mcplib.Description("OAuth token URL"),
		),
		mcplib.WithString("oauth_client_id",
			mcplib.Description("OAuth client ID"),
		),
		mcplib.WithString("oauth_client_secret",
			mcplib.Description("OAuth client secret"),
		),
		mcplib.WithArray("oauth_scopes",
			mcplib.Description("OAuth scopes"),
			mcplib.WithStringItems(),
		),
		mcplib.WithObject("custom_headers",
			mcplib.Description("Custom delivery headers as string key/value pairs"),
			mcplib.AdditionalProperties(map[string]any{"type": "string"}),
		),
		mcplib.WithString("description",
			mcplib.Description("Subscription description"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Auto-generated if omitted"),
		),
	), handleCreateWebhookSubscription)

	s.AddTool(mcplib.NewTool("get_webhook_subscription",
		mcplib.WithDescription("Get a webhook subscription (beta)"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("subscription_id",
			mcplib.Required(),
			mcplib.Description("Webhook subscription ID"),
		),
	), handleGetWebhookSubscription)

	s.AddTool(mcplib.NewTool("update_webhook_subscription",
		mcplib.WithDescription("Update a webhook subscription (beta). Only the fields you provide are changed"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("subscription_id",
			mcplib.Required(),
			mcplib.Description("Webhook subscription ID"),
		),
		mcplib.WithArray("event_types",
			mcplib.Description("Replacement event types"),
			mcplib.WithStringItems(),
		),
		mcplib.WithString("target_url",
			mcplib.Description("New target URL"),
		),
		mcplib.WithString("auth_type",
			mcplib.Description("HMAC, OAUTH, MTLS or OAUTH_MTLS (the WEBHOOK_AUTH_TYPE_ prefix is optional)"),
		),
		mcplib.WithString("oauth_token_url",
			mcplib.Description("OAuth token URL"),
		),
		mcplib.WithString("oauth_client_id",
			mcplib.Description("OAuth client ID"),
		),
		mcplib.WithString("oauth_client_secret",
			mcplib.Description("OAuth client secret"),
		),
		mcplib.WithArray("oauth_scopes",
			mcplib.Description("OAuth scopes"),
			mcplib.WithStringItems(),
		),
		mcplib.WithObject("custom_headers",
			mcplib.Description("Replacement custom delivery headers as string key/value pairs"),
			mcplib.AdditionalProperties(map[string]any{"type": "string"}),
		),
		mcplib.WithString("description",
			mcplib.Description("New description"),
		),
		mcplib.WithBoolean("is_enabled",
			mcplib.Description("Enable or disable the subscription"),
		),
		mcplib.WithString("expected_version",
			mcplib.Description("Subscription version the update is based on. Fetched automatically if omitted"),
		),
	), handleUpdateWebhookSubscription)

	s.AddTool(mcplib.NewTool("delete_webhook_subscription",
		mcplib.WithDescription("Delete a webhook subscription (beta)"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("subscription_id",
			mcplib.Required(),
			mcplib.Description("Webhook subscription ID"),
		),
	), handleDeleteWebhookSubscription)

	s.AddTool(mcplib.NewTool("rotate_webhook_hmac_secret",
		mcplib.WithDescription("Rotate a webhook subscription's HMAC secret (beta). The new secret is returned only once in the response"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("subscription_id",
			mcplib.Required(),
			mcplib.Description("Webhook subscription ID"),
		),
		mcplib.WithBoolean("revoke_previous_immediately",
			mcplib.Description("Invalidate the previous secret immediately instead of after the grace window"),
		),
	), handleRotateWebhookHmacSecret)

	s.AddTool(mcplib.NewTool("send_webhook_simulation",
		mcplib.WithDescription("Send a simulated event to a webhook subscription (beta)"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("subscription_id",
			mcplib.Required(),
			mcplib.Description("Webhook subscription ID"),
		),
		mcplib.WithString("event_type",
			mcplib.Required(),
			mcplib.Description("Event type to simulate, from get_webhook_event_types"),
		),
	), handleSendWebhookSimulation)

	s.AddTool(mcplib.NewTool("get_webhook_delivery_logs",
		mcplib.WithDescription("List delivery log entries for a webhook subscription (beta)"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("subscription_id",
			mcplib.Required(),
			mcplib.Description("Webhook subscription ID"),
		),
		mcplib.WithString("status_filter",
			mcplib.Description("DELIVERED, RETRYING, FAILED or PENDING"),
		),
		mcplib.WithString("event_type_filter",
			mcplib.Description("Only return entries for this event type"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("next_cursor from a previous response"),
		),
		mcplib.WithBoolean("fetch_all",
			mcplib.Description("Fetch all pages automatically and return combined results. When true, cursor is ignored."),
		),
	), handleGetWebhookDeliveryLogs)

	s.AddTool(mcplib.NewTool("get_webhook_delivery_attempts",
		mcplib.WithDescription("List delivery attempts for a webhook event (beta)"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("subscription_id",
			mcplib.Required(),
			mcplib.Description("Webhook subscription ID"),
		),
		mcplib.WithString("event_id",
			mcplib.Required(),
			mcplib.Description("Webhook event ID"),
		),
	), handleGetWebhookDeliveryAttempts)

	s.AddTool(mcplib.NewTool("get_webhook_subscription_stats",
		mcplib.WithDescription("Get delivery statistics for a webhook subscription (beta)"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("subscription_id",
			mcplib.Required(),
			mcplib.Description("Webhook subscription ID"),
		),
		mcplib.WithString("period",
			mcplib.Description("24H, 7D or 30D"),
		),
	), handleGetWebhookSubscriptionStats)

	s.AddTool(mcplib.NewTool("get_webhook_owner_stats",
		mcplib.WithDescription("Get webhook statistics for a portfolio (beta)"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("period",
			mcplib.Description("24H, 7D or 30D"),
		),
	), handleGetWebhookOwnerStats)

	s.AddTool(mcplib.NewTool("replay_webhook_event",
		mcplib.WithDescription("Replay a webhook event (beta)"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("event_id",
			mcplib.Required(),
			mcplib.Description("Webhook event ID"),
		),
		mcplib.WithString("subscription_id",
			mcplib.Description("Limit the replay to a single subscription"),
		),
	), handleReplayWebhookEvent)
}

// webhookServiceFor builds the webhooks service and resolves the portfolio ID.
func webhookServiceFor(req mcplib.CallToolRequest) (webhooks.WebhooksService, string, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return nil, "", fmt.Errorf("failed to initialize client: %w", err)
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return nil, "", err
	}

	return webhooks.NewWebhooksService(client), portfolioId, nil
}

// webhookEnum upper-cases an enum value and adds the API prefix when the
// caller used the short form, e.g. "7d" -> "WEBHOOK_STATS_PERIOD_7D".
func webhookEnum(prefix, value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" || strings.HasPrefix(value, prefix) {
		return value
	}
	return prefix + value
}

// webhookCustomHeaders reads the custom_headers object argument. It returns
// nil when the argument is absent.
func webhookCustomHeaders(req mcplib.CallToolRequest) (map[string]string, error) {
	raw, ok := req.GetArguments()["custom_headers"]
	if !ok || raw == nil {
		return nil, nil
	}

	obj, ok := raw.(map[string]any)
	if !ok {
		return nil, errors.New("custom_headers must be an object of string values")
	}

	headers := make(map[string]string, len(obj))
	for key, value := range obj {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("custom_headers value for %q must be a string", key)
		}
		headers[key] = str
	}

	return headers, nil
}

func webhookArgProvided(req mcplib.CallToolRequest, names ...string) bool {
	args := req.GetArguments()
	for _, name := range names {
		if _, ok := args[name]; ok {
			return true
		}
	}
	return false
}

// webhookOauthConfig returns nil when no OAuth argument was provided.
func webhookOauthConfig(req mcplib.CallToolRequest) *model.WebhookOAuthConfigInput {
	if !webhookArgProvided(req, "oauth_token_url", "oauth_client_id", "oauth_client_secret", "oauth_scopes") {
		return nil
	}

	return &model.WebhookOAuthConfigInput{
		TokenUrl:     req.GetString("oauth_token_url", ""),
		ClientId:     req.GetString("oauth_client_id", ""),
		ClientSecret: req.GetString("oauth_client_secret", ""),
		Scopes:       req.GetStringSlice("oauth_scopes", nil),
	}
}

func handleGetWebhookEventTypes(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	svc := webhooks.NewWebhooksService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetWebhookEventTypes(ctx2, &webhooks.GetWebhookEventTypesRequest{})
	if err != nil {
		return toolErr("cannot get webhook event types: %s", err), nil
	}

	return marshalResult(response)
}

func handleListWebhookSubscriptions(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	fetchAll := req.GetBool("fetch_all", false)

	request := &webhooks.ListWebhookSubscriptionsRequest{
		PortfolioId:     portfolioId,
		IncludeDisabled: req.GetBool("include_disabled", false),
		EventTypeFilter: req.GetString("event_type_filter", ""),
	}
	if !fetchAll {
		request.Cursor = req.GetString("cursor", "")
	}

	response, err := svc.ListWebhookSubscriptions(ctx2, request)
	if err != nil {
		return toolErr("cannot list webhook subscriptions: %s", err), nil
	}

	if !fetchAll {
		return marshalResult(response)
	}

	ctx3, cancel3 := fetchAllCtx(ctx)
	defer cancel3()

	all := []*model.WebhookSubscription{}
	for response != nil {
		all = append(all, response.Subscriptions...)

		response, err = response.Next(ctx3)
		if err != nil {
			return toolErr("failed to fetch all pages: %s", err), nil
		}
	}

	return marshalResult(all)
}

func handleCreateWebhookSubscription(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	eventTypes, err := req.RequireStringSlice("event_types")
	if err != nil {
		return toolErr("event_types is required"), nil
	}

	targetUrl, err := req.RequireString("target_url")
	if err != nil {
		return toolErr("target_url is required"), nil
	}

	headers, err := webhookCustomHeaders(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateWebhookSubscription(ctx2, &webhooks.CreateWebhookSubscriptionRequest{
		PortfolioId:    portfolioId,
		EventTypes:     eventTypes,
		TargetUrl:      targetUrl,
		AuthType:       model.WebhookAuthType(webhookEnum(webhookAuthTypePrefix, req.GetString("auth_type", ""))),
		OauthConfig:    webhookOauthConfig(req),
		CustomHeaders:  headers,
		Description:    req.GetString("description", ""),
		IdempotencyKey: req.GetString("idempotency_key", ""),
	})
	if err != nil {
		return toolErr("cannot create webhook subscription: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetWebhookSubscription(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	subscriptionId, err := req.RequireString("subscription_id")
	if err != nil {
		return toolErr("subscription_id is required"), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetWebhookSubscription(ctx2, &webhooks.GetWebhookSubscriptionRequest{
		PortfolioId:    portfolioId,
		SubscriptionId: subscriptionId,
	})
	if err != nil {
		return toolErr("cannot get webhook subscription: %s", err), nil
	}

	return marshalResult(response)
}

func handleUpdateWebhookSubscription(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	subscriptionId, err := req.RequireString("subscription_id")
	if err != nil {
		return toolErr("subscription_id is required"), nil
	}

	fields := []struct {
		args []string
		path string
	}{
		{[]string{"event_types"}, "event_types"},
		{[]string{"target_url"}, "target_url"},
		{[]string{"auth_type"}, "auth_type"},
		{[]string{"oauth_token_url", "oauth_client_id", "oauth_client_secret", "oauth_scopes"}, "oauth_config"},
		{[]string{"custom_headers"}, "custom_headers"},
		{[]string{"description"}, "description"},
		{[]string{"is_enabled"}, "is_enabled"},
	}

	var paths []string
	for _, f := range fields {
		if webhookArgProvided(req, f.args...) {
			paths = append(paths, f.path)
		}
	}
	if len(paths) == 0 {
		return toolErr("no fields to update: provide at least one field to change"), nil
	}

	headers, err := webhookCustomHeaders(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	// The update is optimistic-concurrency checked, so base it on the current
	// version unless the caller supplied one.
	expectedVersion := req.GetString("expected_version", "")
	if expectedVersion == "" {
		current, err := svc.GetWebhookSubscription(ctx2, &webhooks.GetWebhookSubscriptionRequest{
			PortfolioId:    portfolioId,
			SubscriptionId: subscriptionId,
		})
		if err != nil {
			return toolErr("cannot get webhook subscription: %s", err), nil
		}
		if current.Subscription == nil {
			return toolErr("cannot get webhook subscription: empty response"), nil
		}
		expectedVersion = current.Subscription.Version
	}

	request := &webhooks.UpdateWebhookSubscriptionRequest{
		PortfolioId:     portfolioId,
		SubscriptionId:  subscriptionId,
		EventTypes:      req.GetStringSlice("event_types", nil),
		TargetUrl:       req.GetString("target_url", ""),
		AuthType:        model.WebhookAuthType(webhookEnum(webhookAuthTypePrefix, req.GetString("auth_type", ""))),
		OauthConfig:     webhookOauthConfig(req),
		CustomHeaders:   headers,
		Description:     req.GetString("description", ""),
		UpdateMask:      &model.WebhookUpdateMask{Paths: paths},
		ExpectedVersion: expectedVersion,
	}

	if webhookArgProvided(req, "is_enabled") {
		isEnabled := req.GetBool("is_enabled", false)
		request.IsEnabled = &isEnabled
	}

	response, err := svc.UpdateWebhookSubscription(ctx2, request)
	if err != nil {
		return toolErr("cannot update webhook subscription: %s", err), nil
	}

	return marshalResult(response)
}

func handleDeleteWebhookSubscription(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	subscriptionId, err := req.RequireString("subscription_id")
	if err != nil {
		return toolErr("subscription_id is required"), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.DeleteWebhookSubscription(ctx2, &webhooks.DeleteWebhookSubscriptionRequest{
		PortfolioId:    portfolioId,
		SubscriptionId: subscriptionId,
	})
	if err != nil {
		return toolErr("cannot delete webhook subscription: %s", err), nil
	}

	return marshalResult(response)
}

func handleRotateWebhookHmacSecret(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	subscriptionId, err := req.RequireString("subscription_id")
	if err != nil {
		return toolErr("subscription_id is required"), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.RotateWebhookHmacSecret(ctx2, &webhooks.RotateWebhookHmacSecretRequest{
		PortfolioId:               portfolioId,
		SubscriptionId:            subscriptionId,
		RevokePreviousImmediately: req.GetBool("revoke_previous_immediately", false),
	})
	if err != nil {
		return toolErr("cannot rotate webhook HMAC secret: %s", err), nil
	}

	return marshalResult(response)
}

func handleSendWebhookSimulation(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	subscriptionId, err := req.RequireString("subscription_id")
	if err != nil {
		return toolErr("subscription_id is required"), nil
	}

	eventType, err := req.RequireString("event_type")
	if err != nil {
		return toolErr("event_type is required"), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.SendWebhookSimulation(ctx2, &webhooks.SendWebhookSimulationRequest{
		PortfolioId:    portfolioId,
		SubscriptionId: subscriptionId,
		EventType:      eventType,
	})
	if err != nil {
		return toolErr("cannot send webhook simulation: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetWebhookDeliveryLogs(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	subscriptionId, err := req.RequireString("subscription_id")
	if err != nil {
		return toolErr("subscription_id is required"), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	fetchAll := req.GetBool("fetch_all", false)

	request := &webhooks.GetWebhookDeliveryLogsRequest{
		PortfolioId:    portfolioId,
		SubscriptionId: subscriptionId,
		StatusFilter: model.WebhookDeliveryStatus(
			webhookEnum(webhookDeliveryStatusPrefix, req.GetString("status_filter", "")),
		),
		EventTypeFilter: req.GetString("event_type_filter", ""),
	}
	if !fetchAll {
		request.Cursor = req.GetString("cursor", "")
	}

	response, err := svc.GetWebhookDeliveryLogs(ctx2, request)
	if err != nil {
		return toolErr("cannot get webhook delivery logs: %s", err), nil
	}

	if !fetchAll {
		return marshalResult(response)
	}

	ctx3, cancel3 := fetchAllCtx(ctx)
	defer cancel3()

	all := []*model.WebhookDeliveryLogEntry{}
	for response != nil {
		all = append(all, response.Entries...)

		response, err = response.Next(ctx3)
		if err != nil {
			return toolErr("failed to fetch all pages: %s", err), nil
		}
	}

	return marshalResult(all)
}

func handleGetWebhookDeliveryAttempts(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	subscriptionId, err := req.RequireString("subscription_id")
	if err != nil {
		return toolErr("subscription_id is required"), nil
	}

	eventId, err := req.RequireString("event_id")
	if err != nil {
		return toolErr("event_id is required"), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetWebhookDeliveryAttempts(ctx2, &webhooks.GetWebhookDeliveryAttemptsRequest{
		PortfolioId:    portfolioId,
		SubscriptionId: subscriptionId,
		EventId:        eventId,
	})
	if err != nil {
		return toolErr("cannot get webhook delivery attempts: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetWebhookSubscriptionStats(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	subscriptionId, err := req.RequireString("subscription_id")
	if err != nil {
		return toolErr("subscription_id is required"), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetWebhookSubscriptionStats(ctx2, &webhooks.GetWebhookSubscriptionStatsRequest{
		PortfolioId:    portfolioId,
		SubscriptionId: subscriptionId,
		Period:         model.WebhookStatsPeriod(webhookEnum(webhookStatsPeriodPrefix, req.GetString("period", ""))),
	})
	if err != nil {
		return toolErr("cannot get webhook subscription stats: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetWebhookOwnerStats(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetWebhookOwnerStats(ctx2, &webhooks.GetWebhookOwnerStatsRequest{
		PortfolioId: portfolioId,
		Period:      model.WebhookStatsPeriod(webhookEnum(webhookStatsPeriodPrefix, req.GetString("period", ""))),
	})
	if err != nil {
		return toolErr("cannot get webhook owner stats: %s", err), nil
	}

	return marshalResult(response)
}

func handleReplayWebhookEvent(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, portfolioId, err := webhookServiceFor(req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	eventId, err := req.RequireString("event_id")
	if err != nil {
		return toolErr("event_id is required"), nil
	}

	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ReplayWebhookEvent(ctx2, &webhooks.ReplayWebhookEventRequest{
		PortfolioId:    portfolioId,
		EventId:        eventId,
		SubscriptionId: req.GetString("subscription_id", ""),
	})
	if err != nil {
		return toolErr("cannot replay webhook event: %s", err), nil
	}

	return marshalResult(response)
}
