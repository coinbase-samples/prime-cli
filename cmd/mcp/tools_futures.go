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

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase/prime-sdk-go/futures"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerFuturesTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("get_fcm_balance",
		mcplib.WithDescription("Get FCM balance summary for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetFcmBalance)

	s.AddTool(mcplib.NewTool("get_fcm_positions",
		mcplib.WithDescription("Get FCM futures positions for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetFcmPositions)

	s.AddTool(mcplib.NewTool("get_fcm_settings",
		mcplib.WithDescription("Get FCM settings for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetFcmSettings)

	s.AddTool(mcplib.NewTool("get_fcm_margin_call_details",
		mcplib.WithDescription("Get FCM margin call details for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetFcmMarginCallDetails)

	s.AddTool(mcplib.NewTool("get_fcm_risk_limits",
		mcplib.WithDescription("Get FCM risk limits for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetFcmRiskLimits)

	s.AddTool(mcplib.NewTool("set_fcm_settings",
		mcplib.WithDescription("Update FCM settings for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("target_derivatives_excess",
			mcplib.Description("Target derivatives excess amount"),
		),
	), handleSetFcmSettings)

	s.AddTool(mcplib.NewTool("set_fcm_auto_sweep",
		mcplib.WithDescription("Enable or disable FCM auto sweep"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithBoolean("auto_sweep_enabled",
			mcplib.Required(),
			mcplib.Description("true to enable"),
		),
	), handleSetFcmAutoSweep)

	s.AddTool(mcplib.NewTool("list_fcm_sweeps",
		mcplib.WithDescription("List futures sweeps for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleListFcmSweeps)

	s.AddTool(mcplib.NewTool("schedule_fcm_sweep",
		mcplib.WithDescription("Schedule a futures sweep"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("currency",
			mcplib.Required(),
			mcplib.Description("Currency to sweep"),
		),
		mcplib.WithString("amount",
			mcplib.Description("Amount; omit to sweep all"),
		),
	), handleScheduleFcmSweep)

	s.AddTool(mcplib.NewTool("cancel_fcm_sweep",
		mcplib.WithDescription("Cancel a scheduled futures sweep"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleCancelFcmSweep)
}

func handleGetFcmBalance(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := futures.NewFuturesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetEntityFcmBalance(ctx2, &futures.GetEntityFcmBalanceRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot get FCM balance: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetFcmPositions(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := futures.NewFuturesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetEntityPositions(ctx2, &futures.GetEntityPositionsRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot get FCM positions: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetFcmSettings(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := futures.NewFuturesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetFcmSettings(ctx2, &futures.GetFcmSettingsRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot get FCM settings: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetFcmMarginCallDetails(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := futures.NewFuturesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetFcmMarginCallDetails(ctx2, &futures.GetFcmMarginCallDetailsRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot get FCM margin call details: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetFcmRiskLimits(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := futures.NewFuturesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetFcmRiskLimits(ctx2, &futures.GetFcmRiskLimitsRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot get FCM risk limits: %s", err), nil
	}

	return marshalResult(response)
}

func handleSetFcmSettings(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := futures.NewFuturesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.SetFcmSettings(ctx2, &futures.SetFcmSettingsRequest{
		EntityId:                entityId,
		TargetDerivativesExcess: req.GetString("target_derivatives_excess", ""),
	})
	if err != nil {
		return toolErr("cannot set FCM settings: %s", err), nil
	}

	return marshalResult(response)
}

func handleSetFcmAutoSweep(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := futures.NewFuturesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.SetAutoSweep(ctx2, &futures.SetAutoSweepRequest{
		EntityId:  entityId,
		AutoSweep: req.GetBool("auto_sweep_enabled", false),
	})
	if err != nil {
		return toolErr("cannot set FCM auto sweep: %s", err), nil
	}

	return marshalResult(response)
}

func handleListFcmSweeps(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := futures.NewFuturesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListEntityFuturesSweeps(ctx2, &futures.ListEntityFuturesSweepsRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot list FCM sweeps: %s", err), nil
	}

	return marshalResult(response)
}

func handleScheduleFcmSweep(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := futures.NewFuturesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ScheduleEntityFuturesSweep(ctx2, &futures.ScheduleEntityFuturesSweepRequest{
		EntityId: entityId,
		Currency: req.GetString("currency", ""),
		Amount:   req.GetString("amount", ""),
	})
	if err != nil {
		return toolErr("cannot schedule FCM sweep: %s", err), nil
	}

	return marshalResult(response)
}

func handleCancelFcmSweep(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := futures.NewFuturesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CancelEntityFuturesSweep(ctx2, &futures.CancelEntityFuturesSweepRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot cancel FCM sweep: %s", err), nil
	}

	return marshalResult(response)
}
