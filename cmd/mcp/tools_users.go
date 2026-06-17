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
	"github.com/coinbase/prime-sdk-go/users"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerUserTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_portfolio_users",
		mcplib.WithDescription("List users associated with a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
	), handleListPortfolioUsers)

	s.AddTool(mcplib.NewTool("list_entity_users",
		mcplib.WithDescription("List users for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
	), handleListEntityUsers)
}

func handleListPortfolioUsers(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := users.NewUsersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListPortfolioUsers(ctx2, &users.ListPortfolioUsersRequest{
		PortfolioId: portfolioId,
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list portfolio users: %s", err), nil
	}

	return marshalResult(response)
}

func handleListEntityUsers(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := users.NewUsersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListEntityUsers(ctx2, &users.ListEntityUsersRequest{
		EntityId:   entityId,
		Pagination: paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list entity users: %s", err), nil
	}

	return marshalResult(response)
}
