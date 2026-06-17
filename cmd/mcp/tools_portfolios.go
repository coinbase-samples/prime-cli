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
	"github.com/coinbase/prime-sdk-go/portfolios"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPortfolioTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_portfolios",
		mcplib.WithDescription("List all portfolios associated with the API key"),
	), handleListPortfolios)

	s.AddTool(mcplib.NewTool("get_portfolio",
		mcplib.WithDescription("Get details of a specific portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleGetPortfolio)

	s.AddTool(mcplib.NewTool("get_portfolio_counterparty",
		mcplib.WithDescription("Get counterparty ID for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleGetPortfolioCounterparty)

	s.AddTool(mcplib.NewTool("get_portfolio_credit",
		mcplib.WithDescription("Get credit and buying power information for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleGetPortfolioCredit)
}

func handleGetPortfolioCounterparty(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := portfolios.NewPortfoliosService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetPortfolioCounterparty(ctx2, &portfolios.GetPortfolioCounterpartyRequest{
		PortfolioId: portfolioId,
	})
	if err != nil {
		return toolErr("cannot get portfolio counterparty: %s", err), nil
	}

	return marshalResult(response)
}

func handleListPortfolios(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	svc := portfolios.NewPortfoliosService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListPortfolios(ctx2, &portfolios.ListPortfoliosRequest{})
	if err != nil {
		return toolErr("cannot list portfolios: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetPortfolio(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := portfolios.NewPortfoliosService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetPortfolio(ctx2, &portfolios.GetPortfolioRequest{PortfolioId: portfolioId})
	if err != nil {
		return toolErr("cannot get portfolio: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetPortfolioCredit(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := portfolios.NewPortfoliosService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetPortfolioCredit(ctx2, &portfolios.GetPortfolioCreditRequest{Id: portfolioId})
	if err != nil {
		return toolErr("cannot get portfolio credit: %s", err), nil
	}

	return marshalResult(response)
}
