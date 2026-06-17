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
	"github.com/coinbase/prime-sdk-go/positions"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPositionTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_aggregate_positions",
		mcplib.WithDescription("List aggregate positions for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Entity ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
	), handleListAggregatePositions)

	s.AddTool(mcplib.NewTool("list_positions",
		mcplib.WithDescription("List aggregate positions for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Entity ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
	), handleListPositions)
}

func handleListAggregatePositions(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := positions.NewPositionsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListAggregateEntityPositions(ctx2, &positions.ListAggregateEntityPositionsRequest{
		EntityId:   entityId,
		Pagination: paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list aggregate positions: %s", err), nil
	}

	return marshalResult(response)
}

func handleListPositions(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := positions.NewPositionsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListEntityPositions(ctx2, &positions.ListEntityPositionsRequest{
		EntityId:   entityId,
		Pagination: paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list positions: %s", err), nil
	}

	return marshalResult(response)
}
