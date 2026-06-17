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
	"github.com/coinbase/prime-sdk-go/advancedtransfers"
	"github.com/coinbase/prime-sdk-go/model"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerAdvancedTransferTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_advanced_transfers",
		mcplib.WithDescription("List advanced transfers for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithArray("states",
			mcplib.Description("Filter by states: ADVANCED_TRANSFER_STATE_CREATED, PROCESSING, DONE, CANCELLED, FAILED, EXPIRED"),
			mcplib.WithStringItems(),
		),
		mcplib.WithString("transfer_type",
			mcplib.Description("Filter by transfer type"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
	), handleListAdvancedTransfers)

	s.AddTool(mcplib.NewTool("create_advanced_transfer",
		mcplib.WithDescription("Create a new advanced transfer"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("transfer_type",
			mcplib.Required(),
			mcplib.Description("e.g. ADVANCED_TRANSFER_TYPE_BLIND_MATCH"),
		),
		mcplib.WithString("amount",
			mcplib.Required(),
			mcplib.Description("Amount to transfer"),
		),
		mcplib.WithString("currency",
			mcplib.Required(),
			mcplib.Description("Currency symbol for the transfer"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Auto-generated if omitted"),
		),
		mcplib.WithString("source_type",
			mcplib.Description("Source location type, e.g. WALLET, COUNTERPARTY_ID"),
		),
		mcplib.WithString("source_value",
			mcplib.Description("Source location value (e.g. wallet ID or counterparty ID)"),
		),
		mcplib.WithString("source_address",
			mcplib.Description("Source blockchain address"),
		),
		mcplib.WithString("target_type",
			mcplib.Description("Target location type, e.g. WALLET, COUNTERPARTY_ID"),
		),
		mcplib.WithString("target_value",
			mcplib.Description("Target location value (e.g. wallet ID or counterparty ID)"),
		),
		mcplib.WithString("target_address",
			mcplib.Description("Target blockchain address"),
		),
		mcplib.WithString("reference_id",
			mcplib.Description("Blind match reference ID"),
		),
		mcplib.WithString("settlement_date",
			mcplib.Description("Blind match settlement date"),
		),
		mcplib.WithString("trade_date",
			mcplib.Description("Blind match trade date"),
		),
	), handleCreateAdvancedTransfer)

	s.AddTool(mcplib.NewTool("cancel_advanced_transfer",
		mcplib.WithDescription("Cancel an advanced transfer"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("advanced_transfer_id",
			mcplib.Required(),
			mcplib.Description("ID of the advanced transfer to cancel"),
		),
	), handleCancelAdvancedTransfer)

	s.AddTool(mcplib.NewTool("list_advanced_transfer_transactions",
		mcplib.WithDescription("List transactions for an advanced transfer"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("advanced_transfer_id",
			mcplib.Required(),
			mcplib.Description("ID of the advanced transfer"),
		),
	), handleListAdvancedTransferTransactions)
}

func handleListAdvancedTransfers(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	var state model.AdvancedTransferState
	if states := req.GetStringSlice("states", nil); len(states) > 0 {
		state = model.AdvancedTransferState(states[0])
	}

	transferType := model.AdvancedTransferType(req.GetString("transfer_type", ""))

	svc := advancedtransfers.NewAdvancedTransfersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListAdvancedTransfers(ctx2, &advancedtransfers.ListAdvancedTransfersRequest{
		PortfolioId: portfolioId,
		State:       state,
		Type:        transferType,
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list advanced transfers: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateAdvancedTransfer(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	idempotencyKey := req.GetString("idempotency_key", "")
	if idempotencyKey == "" {
		idempotencyKey = utils.NewUuidStr()
	}

	movement := &model.FundMovement{
		Id:       idempotencyKey,
		Currency: req.GetString("currency", ""),
		Amount:   req.GetString("amount", ""),
	}

	sourceType := req.GetString("source_type", "")
	sourceValue := req.GetString("source_value", "")
	sourceAddress := req.GetString("source_address", "")
	if sourceType != "" || sourceValue != "" || sourceAddress != "" {
		movement.Source = &model.TransferLocation{
			Type:    model.TransferLocationType(sourceType),
			Value:   sourceValue,
			Address: sourceAddress,
		}
	}

	targetType := req.GetString("target_type", "")
	targetValue := req.GetString("target_value", "")
	targetAddress := req.GetString("target_address", "")
	if targetType != "" || targetValue != "" || targetAddress != "" {
		movement.Target = &model.TransferLocation{
			Type:    model.TransferLocationType(targetType),
			Value:   targetValue,
			Address: targetAddress,
		}
	}

	transfer := &model.AdvancedTransfer{
		Type:          model.AdvancedTransferType(req.GetString("transfer_type", "")),
		FundMovements: []*model.FundMovement{movement},
	}

	referenceId := req.GetString("reference_id", "")
	settlementDate := req.GetString("settlement_date", "")
	tradeDate := req.GetString("trade_date", "")
	if referenceId != "" || settlementDate != "" || tradeDate != "" {
		transfer.BlindMatchMetadata = &model.BlindMatchMetadata{
			ReferenceId:    referenceId,
			SettlementDate: settlementDate,
			TradeDate:      tradeDate,
		}
	}

	svc := advancedtransfers.NewAdvancedTransfersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateAdvancedTransfer(ctx2, &advancedtransfers.CreateAdvancedTransferRequest{
		PortfolioId:      portfolioId,
		AdvancedTransfer: transfer,
	})
	if err != nil {
		return toolErr("cannot create advanced transfer: %s", err), nil
	}

	return marshalResult(response)
}

func handleCancelAdvancedTransfer(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	advancedTransferId := req.GetString("advanced_transfer_id", "")
	if advancedTransferId == "" {
		return toolErr("advanced_transfer_id is required"), nil
	}

	svc := advancedtransfers.NewAdvancedTransfersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CancelAdvancedTransfer(ctx2, &advancedtransfers.CancelAdvancedTransferRequest{
		PortfolioId:        portfolioId,
		AdvancedTransferId: advancedTransferId,
	})
	if err != nil {
		return toolErr("cannot cancel advanced transfer: %s", err), nil
	}

	return marshalResult(response)
}

func handleListAdvancedTransferTransactions(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	advancedTransferId := req.GetString("advanced_transfer_id", "")
	if advancedTransferId == "" {
		return toolErr("advanced_transfer_id is required"), nil
	}

	svc := advancedtransfers.NewAdvancedTransfersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListAdvancedTransferTransactions(ctx2, &advancedtransfers.ListAdvancedTransferTransactionsRequest{
		PortfolioId:        portfolioId,
		AdvancedTransferId: advancedTransferId,
	})
	if err != nil {
		return toolErr("cannot list advanced transfer transactions: %s", err), nil
	}

	return marshalResult(response)
}
