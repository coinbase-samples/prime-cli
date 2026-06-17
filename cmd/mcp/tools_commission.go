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
	"github.com/coinbase/prime-sdk-go/commission"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCommissionTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("get_commission",
		mcplib.WithDescription("Get commission rates and fee tiers for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleGetCommission)
}

func handleGetCommission(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := commission.NewCommissionService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetPortfolioCommission(ctx2, &commission.GetPortfolioCommissionRequest{
		PortfolioId: portfolioId,
	})
	if err != nil {
		return toolErr("cannot get commission: %s", err), nil
	}

	return marshalResult(response)
}
