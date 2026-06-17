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
	"github.com/coinbase/prime-sdk-go/balances"
	"github.com/coinbase/prime-sdk-go/model"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerBalanceTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_portfolio_balances",
		mcplib.WithDescription("List all asset balances for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("type",
			mcplib.Description("Balance type filter: TRADING_BALANCES, VAULT_BALANCES, TOTAL_BALANCES, PRIME_CUSTODY_BALANCES, or UNIFIED_TOTAL_BALANCES"),
		),
		mcplib.WithArray("symbols",
			mcplib.Description("Filter by asset symbols (e.g. [\"BTC\", \"ETH\"])"),
			mcplib.WithStringItems(),
		),
	), handleListPortfolioBalances)

	s.AddTool(mcplib.NewTool("get_wallet_balance",
		mcplib.WithDescription("Get the balance for a specific wallet"),
		mcplib.WithString("wallet_id",
			mcplib.Required(),
			mcplib.Description("Wallet ID"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleGetWalletBalance)

	s.AddTool(mcplib.NewTool("list_onchain_balances",
		mcplib.WithDescription("List onchain balances for a wallet"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("wallet_id",
			mcplib.Description("Filter by wallet ID"),
		),
	), handleListOnchainBalances)

	s.AddTool(mcplib.NewTool("list_entity_balances",
		mcplib.WithDescription("List balances at the entity level across all portfolios"),
		mcplib.WithString("entity_id",
			mcplib.Description("Entity ID. Uses credentials default if omitted"),
		),
		mcplib.WithArray("symbols",
			mcplib.Description("Filter by asset symbols"),
			mcplib.WithStringItems(),
		),
		mcplib.WithString("aggregation_type",
			mcplib.Description("Aggregation type: TRADING_BALANCES, VAULT_BALANCES, TOTAL_BALANCES, PRIME_CUSTODY_BALANCES, or UNIFIED_TOTAL_BALANCES"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
	), handleListEntityBalances)
}

func handleListOnchainBalances(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := balances.NewBalancesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListOnchainWalletBalances(ctx2, &balances.ListOnchainWalletBalancesRequest{
		PortfolioId: portfolioId,
		WalletId:    req.GetString("wallet_id", ""),
	})
	if err != nil {
		return toolErr("cannot list onchain balances: %s", err), nil
	}

	return marshalResult(response)
}

func handleListPortfolioBalances(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := balances.NewBalancesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListPortfolioBalances(ctx2, &balances.ListPortfolioBalancesRequest{
		PortfolioId: portfolioId,
		Type:        req.GetString("type", ""),
		Symbols:     req.GetStringSlice("symbols", nil),
	})
	if err != nil {
		return toolErr("cannot list portfolio balances: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetWalletBalance(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	walletId, err := req.RequireString("wallet_id")
	if err != nil {
		return toolErr("wallet_id is required"), nil
	}

	svc := balances.NewBalancesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetWalletBalance(ctx2, &balances.GetWalletBalanceRequest{
		PortfolioId: portfolioId,
		Id:          walletId,
	})
	if err != nil {
		return toolErr("cannot get wallet balance: %s", err), nil
	}

	return marshalResult(response)
}

func handleListEntityBalances(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := balances.NewBalancesService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListEntityBalances(ctx2, &balances.ListEntityBalancesRequest{
		EntityId:        entityId,
		Symbols:         req.GetStringSlice("symbols", nil),
		AggregationType: model.AggregationType(req.GetString("aggregation_type", "")),
		Pagination:      paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list entity balances: %s", err), nil
	}

	return marshalResult(response)
}
