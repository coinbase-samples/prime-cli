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
	"fmt"
	"strings"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase/prime-sdk-go/assets"
	"github.com/coinbase/prime-sdk-go/model"
	"github.com/coinbase/prime-sdk-go/transactions"
	"github.com/coinbase/prime-sdk-go/wallets"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCompositeTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("send_to_blockchain_address",
		mcplib.WithDescription("Send crypto to an external blockchain address in a single step. Automatically finds the source wallet and creates the withdrawal. WARNING: executes a real financial transaction.\n\nNETWORK: Only ETH and USDC support multiple networks. All other assets use a single default network and network_id is not needed.\n  ETH networks:  ethereum-mainnet, base-mainnet\n  USDC networks: ethereum-mainnet, base-mainnet, solana-mainnet, arbitrum-mainnet,\n                 monad-mainnet, optimism-mainnet, avalanche-mainnet\n\nWALLET SELECTION (when wallet_type is omitted):\n  1. Looks for a TRADING wallet holding the asset first.\n  2. Falls back to a QC wallet if no TRADING wallet is found.\n  To send from a VAULT wallet, you must explicitly set wallet_type=VAULT.\n\nEXAMPLE — send 1 USDC on Base from the trading wallet:\n  send_to_blockchain_address(symbol=\"USDC\", amount=\"1\",\n    to_address=\"0x836fa72D2aF55d698e8767acBE88c042b8201036\",\n    network_id=\"base-mainnet\")"),
		mcplib.WithString("symbol",
			mcplib.Required(),
			mcplib.Description("Asset symbol to send (e.g. USDC, ETH, BTC, SOL)"),
		),
		mcplib.WithString("amount",
			mcplib.Required(),
			mcplib.Description("Amount to send as a decimal string (e.g. \"1\", \"0.5\")"),
		),
		mcplib.WithString("to_address",
			mcplib.Required(),
			mcplib.Description("Destination blockchain address"),
		),
		mcplib.WithString("network_id",
			mcplib.Description("Blockchain network for the withdrawal. Only relevant for ETH and USDC. ETH: ethereum-mainnet, base-mainnet. USDC: ethereum-mainnet, base-mainnet, solana-mainnet, arbitrum-mainnet, monad-mainnet, optimism-mainnet, avalanche-mainnet."),
		),
		mcplib.WithString("wallet_type",
			mcplib.Description("Wallet type for the source wallet: TRADING (default, falls back to QC), QC, or VAULT. To send from a VAULT wallet you must explicitly set wallet_type=VAULT."),
		),
		mcplib.WithString("wallet_name_contains",
			mcplib.Description("Filter source wallet by name substring (case-insensitive). Use when multiple wallets match to select the correct one."),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted."),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Idempotency key (UUID). Auto-generated if omitted."),
		),
	), handleSendToBlockchainAddress)

	s.AddTool(mcplib.NewTool("get_deposit_address",
		mcplib.WithDescription("Get a deposit address for an asset in a single step. Automatically finds the right wallet (TRADING first, then QC) and creates the deposit address.\n\nDeposit addresses can only be created for TRADING or QC wallets.\n\nNETWORK: network_id is optional. For single-network assets (e.g. BTC, XLM, SOL) the network is auto-selected. For multi-network assets (ETH, USDC) network_id is required — if omitted, the tool will return the list of available networks to choose from.\n\nEXAMPLE — get a USDC deposit address on Base:\n  get_deposit_address(symbol=\"USDC\", network_id=\"base-mainnet\")\n\nEXAMPLE — get a BTC deposit address (network auto-selected):\n  get_deposit_address(symbol=\"BTC\")"),
		mcplib.WithString("symbol",
			mcplib.Required(),
			mcplib.Description("Asset symbol (e.g. USDC, ETH, BTC, SOL)"),
		),
		mcplib.WithString("network_id",
			mcplib.Description("Network for the deposit address. Optional for single-network assets (auto-selected). Required for multi-network assets (ETH, USDC). Accepts shorthand (e.g. \"base\") or full form (e.g. \"base-mainnet\"). If omitted for a multi-network asset, available options are returned."),
		),
		mcplib.WithString("wallet_name_contains",
			mcplib.Description("Filter wallet by name substring (case-insensitive). Use when multiple wallets match to select the correct one."),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted."),
		),
	), handleGetDepositAddress)
}

func handleSendToBlockchainAddress(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolErr("symbol is required"), nil
	}

	amount, err := req.RequireString("amount")
	if err != nil {
		return toolErr("amount is required"), nil
	}

	toAddress, err := req.RequireString("to_address")
	if err != nil {
		return toolErr("to_address is required"), nil
	}

	networkId := req.GetString("network_id", "")
	walletTypeOverride := req.GetString("wallet_type", "")
	nameContains := req.GetString("wallet_name_contains", "")

	walletSvc := wallets.NewWalletsService(client)
	ctx2, cancel2 := fetchAllCtx(ctx)
	defer cancel2()

	var matched []*model.Wallet
	var effectiveType string

	if walletTypeOverride != "" {
		// User explicitly chose a wallet type — search only that type, no fallback.
		matched, err = findWalletsForSend(ctx2, walletSvc, portfolioId, walletTypeOverride, symbol, networkId, nameContains)
		if err != nil {
			return toolErr("failed to list %s wallets: %s", walletTypeOverride, err), nil
		}
		effectiveType = walletTypeOverride
	} else {
		// Default: try TRADING first, fall back to QC.
		matched, err = findWalletsForSend(ctx2, walletSvc, portfolioId, model.WalletTypeTrading, symbol, networkId, nameContains)
		if err != nil {
			return toolErr("failed to list TRADING wallets: %s", err), nil
		}
		if len(matched) > 0 {
			effectiveType = model.WalletTypeTrading
		} else {
			matched, err = findWalletsForSend(ctx2, walletSvc, portfolioId, "QC", symbol, networkId, nameContains)
			if err != nil {
				return toolErr("failed to list QC wallets: %s", err), nil
			}
			effectiveType = "TRADING or QC"
		}
	}

	switch len(matched) {
	case 0:
		msg := fmt.Sprintf("no %s wallet found for symbol %s", effectiveType, symbol)
		if networkId != "" {
			msg += " on network " + networkId
		}
		if nameContains != "" {
			msg += fmt.Sprintf(" with name containing %q", nameContains)
		}
		return toolErr("%s", msg), nil

	case 1:
		// proceed to withdrawal below

	default:
		var descriptions []string
		for _, w := range matched {
			networkSuffix := ""
			if w.Network != nil && w.Network.Id != "" {
				networkSuffix = " (network=" + w.Network.Id + ")"
			}
			descriptions = append(descriptions, fmt.Sprintf("%q [id=%s%s]", w.Name, w.Id, networkSuffix))
		}
		return toolErr(
			"multiple %s wallets match symbol %s; disambiguate using wallet_name_contains or network_id. Matches: %s",
			effectiveType, symbol, strings.Join(descriptions, ", "),
		), nil
	}

	sourceWallet := matched[0]

	idempotencyKey := req.GetString("idempotency_key", "")
	if idempotencyKey == "" {
		idempotencyKey = utils.NewUuidStr()
	}

	txSvc := transactions.NewTransactionsService(client)
	ctx3, cancel3 := mcpCtx(ctx)
	defer cancel3()

	response, err := txSvc.CreateWalletWithdrawal(ctx3, &transactions.CreateWalletWithdrawalRequest{
		PortfolioId:     portfolioId,
		SourceWalletId:  sourceWallet.Id,
		Symbol:          symbol,
		DestinationType: "DESTINATION_BLOCKCHAIN",
		Amount:          amount,
		IdempotencyKey:  idempotencyKey,
		BlockchainAddress: &model.BlockchainAddress{
			Address: toAddress,
			Network: networkDetailsFor(networkId),
		},
	})
	if err != nil {
		return toolErr("cannot create withdrawal from wallet %s (%s): %s", sourceWallet.Name, sourceWallet.Id, err), nil
	}

	return marshalResult(response)
}

func handleGetDepositAddress(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolErr("symbol is required"), nil
	}

	networkId := req.GetString("network_id", "")
	nameContains := req.GetString("wallet_name_contains", "")

	// Step 1: look up asset networks.
	assetSvc := assets.NewAssetsService(client)
	ctxA, cancelA := mcpCtx(ctx)
	defer cancelA()
	assetResp, err := assetSvc.ListAssets(ctxA, &assets.ListAssetsRequest{EntityId: entityId})
	if err != nil {
		return toolErr("failed to list assets: %s", err), nil
	}
	var assetNetworks []*model.NetworkDetails
	for _, a := range assetResp.Assets {
		if strings.EqualFold(a.Symbol, symbol) {
			for _, n := range a.Networks {
				if n.Network != nil {
					assetNetworks = append(assetNetworks, n.Network)
				}
			}
			break
		}
	}
	if len(assetNetworks) == 0 {
		return toolErr("no network information found for asset %s", symbol), nil
	}

	// Step 2: resolve which network to use.
	var apiNetworkId string
	if networkId == "" {
		if len(assetNetworks) == 1 {
			n := assetNetworks[0]
			apiNetworkId = n.Id + "-" + n.Type
		} else {
			var opts []string
			for _, n := range assetNetworks {
				opts = append(opts, n.Id+"-"+n.Type)
			}
			return toolErr(
				"multiple networks available for %s; specify network_id. Options: %s",
				symbol, strings.Join(opts, ", "),
			), nil
		}
	} else {
		for _, n := range assetNetworks {
			full := n.Id + "-" + n.Type
			if strings.EqualFold(networkId, n.Id) || strings.EqualFold(networkId, full) {
				apiNetworkId = full
				break
			}
		}
		if apiNetworkId == "" {
			var opts []string
			for _, n := range assetNetworks {
				opts = append(opts, n.Id+"-"+n.Type)
			}
			return toolErr(
				"network %q not supported for %s. Supported: %s",
				networkId, symbol, strings.Join(opts, ", "),
			), nil
		}
	}

	// Step 3: find wallet (TRADING first, then QC). Don't filter by network —
	// TRADING/QC wallets have null network; apiNetworkId goes to the API call.
	walletSvc := wallets.NewWalletsService(client)
	ctx2, cancel2 := fetchAllCtx(ctx)
	defer cancel2()

	matched, err := findWalletsForSend(ctx2, walletSvc, portfolioId, model.WalletTypeTrading, symbol, "", nameContains)
	if err != nil {
		return toolErr("failed to list TRADING wallets: %s", err), nil
	}
	effectiveType := model.WalletTypeTrading
	if len(matched) == 0 {
		matched, err = findWalletsForSend(ctx2, walletSvc, portfolioId, "QC", symbol, "", nameContains)
		if err != nil {
			return toolErr("failed to list QC wallets: %s", err), nil
		}
		effectiveType = "TRADING or QC"
	}

	switch len(matched) {
	case 0:
		msg := fmt.Sprintf("no TRADING or QC wallet found for symbol %s", symbol)
		if nameContains != "" {
			msg += fmt.Sprintf(" with name containing %q", nameContains)
		}
		return toolErr("%s", msg), nil

	case 1:
		// proceed below

	default:
		var descriptions []string
		for _, w := range matched {
			descriptions = append(descriptions, fmt.Sprintf("%q [id=%s]", w.Name, w.Id))
		}
		return toolErr(
			"multiple %s wallets match symbol %s; disambiguate using wallet_name_contains. Matches: %s",
			effectiveType, symbol, strings.Join(descriptions, ", "),
		), nil
	}

	wallet := matched[0]

	// Step 4: create deposit address.
	ctx3, cancel3 := mcpCtx(ctx)
	defer cancel3()

	response, err := walletSvc.CreateWalletAddress(ctx3, &wallets.CreateWalletAddressRequest{
		PortfolioId: portfolioId,
		WalletId:    wallet.Id,
		NetworkId:   apiNetworkId,
	})
	if err != nil {
		return toolErr("cannot create deposit address for wallet %s (%s): %s", wallet.Name, wallet.Id, err), nil
	}

	return marshalResult(response)
}

func findWalletsForSend(ctx context.Context, svc wallets.WalletsService, portfolioId, walletType, symbol, networkId, nameContains string) ([]*model.Wallet, error) {
	resp, err := svc.ListWallets(ctx, &wallets.ListWalletsRequest{
		PortfolioId: portfolioId,
		Type:        walletType,
		Symbols:     []string{symbol},
	})
	if err != nil {
		return nil, err
	}

	all, err := resp.Iterator().FetchAll(ctx)
	if err != nil {
		return nil, err
	}

	var matched []*model.Wallet
	for _, w := range all {
		if networkId != "" && w.Network != nil && w.Network.Id != networkId {
			continue
		}
		if nameContains != "" && !strings.Contains(strings.ToLower(w.Name), strings.ToLower(nameContains)) {
			continue
		}
		matched = append(matched, w)
	}
	return matched, nil
}
