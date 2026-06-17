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
	"github.com/coinbase/prime-sdk-go/activities"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerActivityTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("get_activity",
		mcplib.WithDescription("Get a specific activity by ID"),
		mcplib.WithString("activity_id",
			mcplib.Required(),
			mcplib.Description("Activity ID"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetActivity)

	s.AddTool(mcplib.NewTool("list_entity_activities",
		mcplib.WithDescription("List activities for an entity across all portfolios"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("activity_level",
			mcplib.Description("ACTIVITY_LEVEL_ALL, ACTIVITY_LEVEL_PORTFOLIO, or ACTIVITY_LEVEL_ENTITY"),
		),
		mcplib.WithArray("symbols",
			mcplib.WithStringItems(),
			mcplib.Description("Filter by asset symbols"),
		),
		mcplib.WithArray("categories",
			mcplib.WithStringItems(),
			mcplib.Description("Filter by category: ORDER, TRANSACTION, ACCOUNT, ALLOCATION, LENDING"),
		),
		mcplib.WithArray("statuses",
			mcplib.WithStringItems(),
			mcplib.Description("Filter by status: CANCELLED, PROCESSING, COMPLETED, EXPIRED, REJECTED, FAILED"),
		),
		mcplib.WithString("start",
			mcplib.Description("Start time in RFC3339 format"),
		),
		mcplib.WithString("end",
			mcplib.Description("End time in RFC3339 format"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
	), handleListEntityActivities)

	s.AddTool(mcplib.NewTool("get_entity_activity",
		mcplib.WithDescription("Get a specific entity-level activity by ID"),
		mcplib.WithString("activity_id",
			mcplib.Required(),
			mcplib.Description("Activity ID"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetEntityActivity)

	s.AddTool(mcplib.NewTool("list_activities",
		mcplib.WithDescription("List portfolio activities (orders, transfers, staking, etc.) meeting filter criteria"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithArray("symbols",
			mcplib.Description("Filter by asset symbols"),
			mcplib.WithStringItems(),
		),
		mcplib.WithArray("categories",
			mcplib.Description("Filter by activity categories"),
			mcplib.WithStringItems(),
		),
		mcplib.WithArray("statuses",
			mcplib.Description("Filter by activity statuses"),
			mcplib.WithStringItems(),
		),
		mcplib.WithString("start",
			mcplib.Description("Start time in RFC3339 format"),
		),
		mcplib.WithString("end",
			mcplib.Description("End time in RFC3339 format"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
	), handleListActivities)
}

func handleListActivities(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	start, end, err := utils.ParseDateRange(req.GetString("start", ""), req.GetString("end", ""))
	if err != nil {
		return toolErr("invalid date range: %s", err), nil
	}

	svc := activities.NewActivitiesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListActivities(ctx2, &activities.ListActivitiesRequest{
		PortfolioId: portfolioId,
		Symbols:     req.GetStringSlice("symbols", nil),
		Categories:  req.GetStringSlice("categories", nil),
		Statuses:    req.GetStringSlice("statuses", nil),
		Start:       start,
		End:         end,
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list activities: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetActivity(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := activities.NewActivitiesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetActivity(ctx2, &activities.GetActivityRequest{
		PortfolioId: portfolioId,
		Id:          req.GetString("activity_id", ""),
	})
	if err != nil {
		return toolErr("cannot get activity: %s", err), nil
	}

	return marshalResult(response)
}

func handleListEntityActivities(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	start, end, err := utils.ParseDateRange(req.GetString("start", ""), req.GetString("end", ""))
	if err != nil {
		return toolErr("invalid date range: %s", err), nil
	}

	svc := activities.NewActivitiesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListEntityActivities(ctx2, &activities.ListEntityActivitiesRequest{
		EntityId:      entityId,
		ActivityLevel: req.GetString("activity_level", ""),
		Symbols:       req.GetStringSlice("symbols", nil),
		Categories:    req.GetStringSlice("categories", nil),
		Statuses:      req.GetStringSlice("statuses", nil),
		StartTime:     start,
		EndTime:       end,
		Pagination:    paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list entity activities: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetEntityActivity(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	svc := activities.NewActivitiesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetEntityActivity(ctx2, &activities.GetEntityActivityRequest{
		ActivityId: req.GetString("activity_id", ""),
	})
	if err != nil {
		return toolErr("cannot get entity activity: %s", err), nil
	}

	return marshalResult(response)
}
