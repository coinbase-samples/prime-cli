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
	"github.com/coinbase/prime-sdk-go/model"
	"github.com/coinbase/prime-sdk-go/wallets"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerWalletTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_wallets",
		mcplib.WithDescription("List wallets for a portfolio filtered by type"),
		mcplib.WithString("type",
			mcplib.Required(),
			mcplib.Description("Wallet type: VAULT, ONCHAIN, TRADING, or QC"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithArray("symbols",
			mcplib.Description("Filter by asset symbols"),
			mcplib.WithStringItems(),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
		mcplib.WithBoolean("fetch_all",
			mcplib.Description("Fetch all pages automatically and return combined results. When true, cursor and limit are ignored."),
		),
	), handleListWallets)

	s.AddTool(mcplib.NewTool("get_wallet",
		mcplib.WithDescription("Get details of a specific wallet"),
		mcplib.WithString("wallet_id",
			mcplib.Required(),
			mcplib.Description("Wallet ID"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleGetWallet)

	s.AddTool(mcplib.NewTool("create_wallet",
		mcplib.WithDescription("Create a new vault or onchain wallet"),
		mcplib.WithString("name",
			mcplib.Required(),
			mcplib.Description("Display name for the wallet"),
		),
		mcplib.WithString("type",
			mcplib.Required(),
			mcplib.Description("Wallet type: VAULT, ONCHAIN, TRADING, or QC"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("symbol",
			mcplib.Description("Asset symbol (e.g. BTC, ETH)"),
		),
		mcplib.WithString("network_family",
			mcplib.Description("Network family for ONCHAIN wallets (NETWORK_FAMILY_EVM or NETWORK_FAMILY_SOLANA)"),
		),
		mcplib.WithString("network_id",
			mcplib.Description("Network ID (e.g. ethereum, base, bitcoin)"),
		),
		mcplib.WithString("network_type",
			mcplib.Description("Network type: mainnet or testnet"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Idempotency key (UUID). Auto-generated if omitted"),
		),
	), handleCreateWallet)

	s.AddTool(mcplib.NewTool("create_deposit_address",
		mcplib.WithDescription("Create a new deposit address for a wallet"),
		mcplib.WithString("wallet_id",
			mcplib.Required(),
			mcplib.Description("Wallet ID to create a deposit address for"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("network_id",
			mcplib.Description("Network ID (e.g. ethereum-mainnet)"),
		),
	), handleCreateDepositAddress)

	s.AddTool(mcplib.NewTool("get_wallet_deposit_instructions",
		mcplib.WithDescription("Get deposit instructions for a wallet"),
		mcplib.WithString("wallet_id",
			mcplib.Required(),
			mcplib.Description("Wallet ID"),
		),
		mcplib.WithString("deposit_type",
			mcplib.Required(),
			mcplib.Description("Deposit type required by the wallet"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleGetWalletDepositInstructions)

	s.AddTool(mcplib.NewTool("list_wallet_addresses",
		mcplib.WithDescription("List addresses for a wallet"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("wallet_id",
			mcplib.Description("Wallet ID"),
		),
		mcplib.WithString("network_id",
			mcplib.Description("Filter by network ID"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
		mcplib.WithBoolean("fetch_all",
			mcplib.Description("Fetch all pages automatically and return combined results. When true, cursor and limit are ignored."),
		),
	), handleListWalletAddresses)
}

func handleGetWalletDepositInstructions(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	depositType, err := req.RequireString("deposit_type")
	if err != nil {
		return toolErr("deposit_type is required"), nil
	}

	svc := wallets.NewWalletsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetWalletDepositInstructions(ctx2, &wallets.GetWalletDepositInstructionsRequest{
		PortfolioId: portfolioId,
		Id:          walletId,
		Type:        depositType,
	})
	if err != nil {
		return toolErr("cannot get wallet deposit instructions: %s", err), nil
	}

	return marshalResult(response)
}

func handleListWalletAddresses(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := wallets.NewWalletsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListWalletAddresses(ctx2, &wallets.ListWalletAddressesRequest{
		PortfolioId: portfolioId,
		WalletId:    req.GetString("wallet_id", ""),
		NetworkId:   req.GetString("network_id", ""),
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list wallet addresses: %s", err), nil
	}

	if req.GetBool("fetch_all", false) {
		ctx3, cancel3 := fetchAllCtx(ctx)
		defer cancel3()
		items, err := response.Iterator().FetchAll(ctx3)
		if err != nil {
			return toolErr("failed to fetch all pages: %s", err), nil
		}
		return marshalResult(items)
	}

	return marshalResult(response)
}

func handleListWallets(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	walletType, err := req.RequireString("type")
	if err != nil {
		return toolErr("type is required (VAULT or ONCHAIN)"), nil
	}

	svc := wallets.NewWalletsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListWallets(ctx2, &wallets.ListWalletsRequest{
		PortfolioId: portfolioId,
		Type:        walletType,
		Symbols:     req.GetStringSlice("symbols", nil),
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list wallets: %s", err), nil
	}

	if req.GetBool("fetch_all", false) {
		ctx3, cancel3 := fetchAllCtx(ctx)
		defer cancel3()
		items, err := response.Iterator().FetchAll(ctx3)
		if err != nil {
			return toolErr("failed to fetch all pages: %s", err), nil
		}
		return marshalResult(items)
	}

	return marshalResult(response)
}

func handleGetWallet(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	svc := wallets.NewWalletsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetWallet(ctx2, &wallets.GetWalletRequest{
		PortfolioId: portfolioId,
		Id:          walletId,
	})
	if err != nil {
		return toolErr("cannot get wallet: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateWallet(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	name, err := req.RequireString("name")
	if err != nil {
		return toolErr("name is required"), nil
	}

	walletType, err := req.RequireString("type")
	if err != nil {
		return toolErr("type is required (VAULT or ONCHAIN)"), nil
	}

	idempotencyKey := req.GetString("idempotency_key", "")
	if idempotencyKey == "" {
		idempotencyKey = utils.NewUuidStr()
	}

	createReq := &wallets.CreateWalletRequest{
		PortfolioId:    portfolioId,
		Name:           name,
		Type:           walletType,
		IdempotencyKey: idempotencyKey,
	}

	if symbol := req.GetString("symbol", ""); symbol != "" {
		createReq.Symbol = symbol
	}
	if nf := req.GetString("network_family", ""); nf != "" {
		createReq.NetworkFamily = nf
	}

	networkId := req.GetString("network_id", "")
	networkType := req.GetString("network_type", "")
	if networkId != "" || networkType != "" {
		createReq.Network = &model.NetworkDetails{
			Id:   networkId,
			Type: networkType,
		}
	}

	svc := wallets.NewWalletsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateWallet(ctx2, createReq)
	if err != nil {
		return toolErr("cannot create wallet: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateDepositAddress(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	svc := wallets.NewWalletsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateWalletAddress(ctx2, &wallets.CreateWalletAddressRequest{
		PortfolioId: portfolioId,
		WalletId:    walletId,
		NetworkId:   req.GetString("network_id", ""),
	})
	if err != nil {
		return toolErr("cannot create deposit address: %s", err), nil
	}

	return marshalResult(response)
}
