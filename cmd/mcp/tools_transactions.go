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
	"encoding/json"
	"time"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase/prime-sdk-go/model"
	"github.com/coinbase/prime-sdk-go/transactions"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTransactionTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_portfolio_transactions",
		mcplib.WithDescription("List transactions for a portfolio with optional filters"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithArray("types",
			mcplib.Description("Filter by transaction types (e.g. [\"CONVERSION\", \"DEPOSIT\"])"),
			mcplib.WithStringItems(),
		),
		mcplib.WithString("symbols",
			mcplib.Description("Filter by asset symbol (e.g. \"BTC\")"),
		),
		mcplib.WithString("start",
			mcplib.Description("Start time in RFC3339 format (e.g. 2024-01-01T00:00:00Z)"),
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
	), handleListPortfolioTransactions)

	s.AddTool(mcplib.NewTool("get_transaction",
		mcplib.WithDescription("Get details of a specific transaction"),
		mcplib.WithString("transaction_id",
			mcplib.Required(),
			mcplib.Description("Transaction ID"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleGetTransaction)

	s.AddTool(mcplib.NewTool("list_wallet_transactions",
		mcplib.WithDescription("List transactions for a specific wallet"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("wallet_id",
			mcplib.Description("Wallet ID"),
		),
		mcplib.WithArray("types",
			mcplib.WithStringItems(),
			mcplib.Description("Filter by transaction types (e.g. [\"DEPOSIT\", \"WITHDRAWAL\"])"),
		),
		mcplib.WithString("symbols",
			mcplib.Description("Filter by asset symbol (e.g. \"BTC\")"),
		),
		mcplib.WithString("start",
			mcplib.Description("Start time in RFC3339 format (e.g. 2024-01-01T00:00:00Z)"),
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
	), handleListWalletTransactions)

	s.AddTool(mcplib.NewTool("create_transfer",
		mcplib.WithDescription("Create an internal transfer between wallets. WARNING: executes a real financial transaction."),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("source_wallet_id",
			mcplib.Required(),
			mcplib.Description("ID of the source wallet"),
		),
		mcplib.WithString("destination_wallet_id",
			mcplib.Required(),
			mcplib.Description("ID of the destination wallet"),
		),
		mcplib.WithString("symbol",
			mcplib.Required(),
			mcplib.Description("Asset symbol"),
		),
		mcplib.WithString("amount",
			mcplib.Required(),
			mcplib.Description("Amount to transfer"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Auto-generated if omitted"),
		),
	), handleCreateTransfer)

	s.AddTool(mcplib.NewTool("create_withdrawal",
		mcplib.WithDescription("Create an external withdrawal. WARNING: executes a real financial transaction."),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("source_wallet_id",
			mcplib.Required(),
			mcplib.Description("ID of the source wallet"),
		),
		mcplib.WithString("symbol",
			mcplib.Required(),
			mcplib.Description("Asset symbol"),
		),
		mcplib.WithString("destination_type",
			mcplib.Required(),
			mcplib.Description("Destination type: DESTINATION_BLOCKCHAIN, DESTINATION_PAYMENT_METHOD, DESTINATION_WALLET, or DESTINATION_COUNTERPARTY"),
		),
		mcplib.WithString("amount",
			mcplib.Required(),
			mcplib.Description("Amount to withdraw"),
		),
		mcplib.WithString("payment_method_id",
			mcplib.Description("ID of the payment method"),
		),
		mcplib.WithString("blockchain_address",
			mcplib.Description("Blockchain address"),
		),
		mcplib.WithString("account_identifier",
			mcplib.Description("Account identifier"),
		),
		mcplib.WithString("network_id",
			mcplib.Description("Blockchain network ID (e.g. ethereum-mainnet, base-mainnet, bitcoin-mainnet)"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Auto-generated if omitted"),
		),
	), handleCreateWithdrawal)

	s.AddTool(mcplib.NewTool("create_conversion",
		mcplib.WithDescription("Convert between fiat and stablecoins. WARNING: executes a real financial transaction."),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("source_wallet_id",
			mcplib.Required(),
			mcplib.Description("ID of the source wallet"),
		),
		mcplib.WithString("source_symbol",
			mcplib.Required(),
			mcplib.Description("Symbol of the source asset"),
		),
		mcplib.WithString("destination_wallet_id",
			mcplib.Required(),
			mcplib.Description("ID of the destination wallet"),
		),
		mcplib.WithString("destination_symbol",
			mcplib.Required(),
			mcplib.Description("Symbol of the destination asset"),
		),
		mcplib.WithString("amount",
			mcplib.Required(),
			mcplib.Description("Amount to convert"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Auto-generated if omitted"),
		),
	), handleCreateConversion)

	s.AddTool(mcplib.NewTool("create_onchain_transaction",
		mcplib.WithDescription("Create an onchain transaction. WARNING: executes a real blockchain transaction."),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("wallet_id",
			mcplib.Required(),
			mcplib.Description("Wallet ID"),
		),
		mcplib.WithString("raw_unsigned_transaction",
			mcplib.Required(),
			mcplib.Description("Raw unsigned transaction hex"),
		),
		mcplib.WithString("url",
			mcplib.Description("RPC URL"),
		),
		mcplib.WithBoolean("skip_broadcast",
			mcplib.Description("Skip broadcast"),
		),
		mcplib.WithString("chain_id",
			mcplib.Description("Chain ID"),
		),
		mcplib.WithBoolean("disable_dynamic_gas",
			mcplib.Description("Disable dynamic gas"),
		),
		mcplib.WithString("replaced_transaction_id",
			mcplib.Description("Replaced transaction ID"),
		),
	), handleCreateOnchainTransaction)

	s.AddTool(mcplib.NewTool("get_travel_rule_data",
		mcplib.WithDescription("Get travel rule data for a transaction"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("transaction_id",
			mcplib.Required(),
			mcplib.Description("Transaction ID"),
		),
	), handleGetTravelRuleData)

	s.AddTool(mcplib.NewTool("submit_deposit_travel_rule_data",
		mcplib.WithDescription("Submit travel rule data for a deposit transaction"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("transaction_id",
			mcplib.Required(),
			mcplib.Description("Transaction ID"),
		),
		mcplib.WithString("originator",
			mcplib.Description("Originator info JSON"),
		),
		mcplib.WithString("beneficiary",
			mcplib.Description("Beneficiary info JSON"),
		),
		mcplib.WithBoolean("is_self",
			mcplib.Description("True if self-transfer"),
		),
		mcplib.WithBoolean("opt_out_of_ownership_verification",
			mcplib.Description("Opt out of ownership verification"),
		),
	), handleSubmitDepositTravelRuleData)
}

func handleListPortfolioTransactions(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	var start, end time.Time
	if s := req.GetString("start", ""); s != "" {
		start, end, err = utils.ParseDateRange(s, req.GetString("end", ""))
		if err != nil {
			return toolErr("invalid date range: %s", err), nil
		}
	}

	svc := transactions.NewTransactionsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListPortfolioTransactions(ctx2, &transactions.ListPortfolioTransactionsRequest{
		PortfolioId: portfolioId,
		Types:       req.GetStringSlice("types", nil),
		Symbols:     req.GetString("symbols", ""),
		Start:       start,
		End:         end,
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list transactions: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetTransaction(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	txId, err := req.RequireString("transaction_id")
	if err != nil {
		return toolErr("transaction_id is required"), nil
	}

	svc := transactions.NewTransactionsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetTransaction(ctx2, &transactions.GetTransactionRequest{
		PortfolioId:   portfolioId,
		TransactionId: txId,
	})
	if err != nil {
		return toolErr("cannot get transaction: %s", err), nil
	}

	return marshalResult(response)
}

func handleListWalletTransactions(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	svc := transactions.NewTransactionsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListWalletTransactions(ctx2, &transactions.ListWalletTransactionsRequest{
		PortfolioId: portfolioId,
		WalletId:    req.GetString("wallet_id", ""),
		Types:       req.GetStringSlice("types", nil),
		Symbols:     req.GetString("symbols", ""),
		Start:       start,
		End:         end,
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list wallet transactions: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateTransfer(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	sourceWalletId, err := req.RequireString("source_wallet_id")
	if err != nil {
		return toolErr("source_wallet_id is required"), nil
	}

	destinationWalletId, err := req.RequireString("destination_wallet_id")
	if err != nil {
		return toolErr("destination_wallet_id is required"), nil
	}

	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolErr("symbol is required"), nil
	}

	amount, err := req.RequireString("amount")
	if err != nil {
		return toolErr("amount is required"), nil
	}

	idempotencyKey := req.GetString("idempotency_key", "")
	if idempotencyKey == "" {
		idempotencyKey = utils.NewUuidStr()
	}

	svc := transactions.NewTransactionsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateWalletTransfer(ctx2, &transactions.CreateWalletTransferRequest{
		PortfolioId:         portfolioId,
		SourceWalletId:      sourceWalletId,
		DestinationWalletId: destinationWalletId,
		Symbol:              symbol,
		Amount:              amount,
		IdempotencyKey:      idempotencyKey,
	})
	if err != nil {
		return toolErr("cannot create transfer: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateWithdrawal(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	sourceWalletId, err := req.RequireString("source_wallet_id")
	if err != nil {
		return toolErr("source_wallet_id is required"), nil
	}

	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolErr("symbol is required"), nil
	}

	destinationType, err := req.RequireString("destination_type")
	if err != nil {
		return toolErr("destination_type is required"), nil
	}

	amount, err := req.RequireString("amount")
	if err != nil {
		return toolErr("amount is required"), nil
	}

	idempotencyKey := req.GetString("idempotency_key", "")
	if idempotencyKey == "" {
		idempotencyKey = utils.NewUuidStr()
	}

	svc := transactions.NewTransactionsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateWalletWithdrawal(ctx2, &transactions.CreateWalletWithdrawalRequest{
		PortfolioId:     portfolioId,
		SourceWalletId:  sourceWalletId,
		Symbol:          symbol,
		DestinationType: destinationType,
		Amount:          amount,
		IdempotencyKey:  idempotencyKey,
		PaymentMethod: &transactions.CreateWalletWithdrawalPaymentMethod{
			Id: req.GetString("payment_method_id", ""),
		},
		BlockchainAddress: &model.BlockchainAddress{
			Address:           req.GetString("blockchain_address", ""),
			AccountIdentifier: req.GetString("account_identifier", ""),
			Network:           networkDetailsFor(req.GetString("network_id", "")),
		},
	})
	if err != nil {
		return toolErr("cannot create withdrawal: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateConversion(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	sourceWalletId, err := req.RequireString("source_wallet_id")
	if err != nil {
		return toolErr("source_wallet_id is required"), nil
	}

	sourceSymbol, err := req.RequireString("source_symbol")
	if err != nil {
		return toolErr("source_symbol is required"), nil
	}

	destinationWalletId, err := req.RequireString("destination_wallet_id")
	if err != nil {
		return toolErr("destination_wallet_id is required"), nil
	}

	destinationSymbol, err := req.RequireString("destination_symbol")
	if err != nil {
		return toolErr("destination_symbol is required"), nil
	}

	amount, err := req.RequireString("amount")
	if err != nil {
		return toolErr("amount is required"), nil
	}

	idempotencyKey := req.GetString("idempotency_key", "")
	if idempotencyKey == "" {
		idempotencyKey = utils.NewUuidStr()
	}

	svc := transactions.NewTransactionsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateConversion(ctx2, &transactions.CreateConversionRequest{
		PortfolioId:         portfolioId,
		SourceWalletId:      sourceWalletId,
		SourceSymbol:        sourceSymbol,
		DestinationWalletId: destinationWalletId,
		DestinationSymbol:   destinationSymbol,
		Amount:              amount,
		IdempotencyKey:      idempotencyKey,
	})
	if err != nil {
		return toolErr("cannot create conversion: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateOnchainTransaction(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	rawTx, err := req.RequireString("raw_unsigned_transaction")
	if err != nil {
		return toolErr("raw_unsigned_transaction is required"), nil
	}

	svc := transactions.NewTransactionsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateOnchainTransaction(ctx2, &transactions.CreateOnchainTransactionRequest{
		PortfolioId: portfolioId,
		WalletId:    walletId,
		OnchainTransaction: &model.OnchainTransaction{
			RawUnsignedTransaction: rawTx,
			Rpc: &model.OnchainRpc{
				Url:           req.GetString("url", ""),
				SkipBroadcast: req.GetBool("skip_broadcast", false),
			},
			EvmParams: &model.OnchainEvmParams{
				ChainId:               req.GetString("chain_id", ""),
				DisableDynamicGas:     req.GetBool("disable_dynamic_gas", false),
				ReplacedTransactionId: req.GetString("replaced_transaction_id", ""),
			},
		},
	})
	if err != nil {
		return toolErr("cannot create onchain transaction: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetTravelRuleData(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	txId, err := req.RequireString("transaction_id")
	if err != nil {
		return toolErr("transaction_id is required"), nil
	}

	svc := transactions.NewTransactionsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetTransactionTravelRuleData(ctx2, &transactions.GetTransactionTravelRuleDataRequest{
		PortfolioId:   portfolioId,
		TransactionId: txId,
	})
	if err != nil {
		return toolErr("cannot get travel rule data: %s", err), nil
	}

	return marshalResult(response)
}

func handleSubmitDepositTravelRuleData(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	txId, err := req.RequireString("transaction_id")
	if err != nil {
		return toolErr("transaction_id is required"), nil
	}

	request := &transactions.SubmitDepositTravelRuleDataRequest{
		PortfolioId:                   portfolioId,
		TransactionId:                 txId,
		IsSelf:                        req.GetBool("is_self", false),
		OptOutOfOwnershipVerification: req.GetBool("opt_out_of_ownership_verification", false),
	}

	if originatorJson := req.GetString("originator", ""); originatorJson != "" {
		var originator model.TravelRuleParty
		if err := json.Unmarshal([]byte(originatorJson), &originator); err != nil {
			return toolErr("invalid originator JSON: %s", err), nil
		}
		request.Originator = &originator
	}

	if beneficiaryJson := req.GetString("beneficiary", ""); beneficiaryJson != "" {
		var beneficiary model.TravelRuleParty
		if err := json.Unmarshal([]byte(beneficiaryJson), &beneficiary); err != nil {
			return toolErr("invalid beneficiary JSON: %s", err), nil
		}
		request.Beneficiary = &beneficiary
	}

	svc := transactions.NewTransactionsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.SubmitDepositTravelRuleData(ctx2, request)
	if err != nil {
		return toolErr("cannot submit deposit travel rule data: %s", err), nil
	}

	return marshalResult(response)
}
