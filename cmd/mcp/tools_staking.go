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
	primeStaking "github.com/coinbase/prime-sdk-go/staking"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerStakingTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("stake",
		mcplib.WithDescription("Create a stake or delegate request for a wallet"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("wallet_id",
			mcplib.Description("Wallet ID to stake from"),
		),
		mcplib.WithString("amount",
			mcplib.Description("Amount to stake. If omitted, the maximum available amount is staked"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Auto-generated if omitted"),
		),
	), handleStake)

	s.AddTool(mcplib.NewTool("unstake",
		mcplib.WithDescription("Create an unstake request for a wallet"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("wallet_id",
			mcplib.Description("Wallet ID to unstake from"),
		),
		mcplib.WithString("amount",
			mcplib.Description("Amount to unstake. If omitted, the maximum available amount is unstaked"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Auto-generated if omitted"),
		),
	), handleUnstake)

	s.AddTool(mcplib.NewTool("get_staking_status",
		mcplib.WithDescription("Get staking status for a wallet"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("wallet_id",
			mcplib.Description("Wallet ID to query staking status for"),
		),
	), handleGetStakingStatus)

	s.AddTool(mcplib.NewTool("claim_staking_rewards",
		mcplib.WithDescription("Claim staking rewards for a wallet"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("wallet_id",
			mcplib.Description("Wallet ID to claim rewards from"),
		),
		mcplib.WithString("amount",
			mcplib.Description("Amount of rewards to claim. If omitted, the full available reward amount is claimed"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Auto-generated if omitted"),
		),
	), handleClaimStakingRewards)

	s.AddTool(mcplib.NewTool("get_unstaking_status",
		mcplib.WithDescription("Get the status of an unstake operation"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("wallet_id",
			mcplib.Description("Wallet ID to query unstaking status for"),
		),
	), handleGetUnstakingStatus)

	s.AddTool(mcplib.NewTool("preview_unstake",
		mcplib.WithDescription("Preview an unstake operation for a wallet"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("wallet_id",
			mcplib.Description("Wallet ID to preview unstake for"),
		),
		mcplib.WithString("amount",
			mcplib.Description("Amount to preview unstaking"),
		),
	), handlePreviewUnstake)

	s.AddTool(mcplib.NewTool("portfolio_stake_initiate",
		mcplib.WithDescription("Initiate a portfolio-level stake request"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("symbol",
			mcplib.Description("Currency symbol to stake (e.g. ETH)"),
		),
		mcplib.WithString("amount",
			mcplib.Description("Amount to stake"),
		),
		mcplib.WithString("stake_protocol",
			mcplib.Description("Optional staking protocol identifier"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Auto-generated if omitted"),
		),
	), handlePortfolioStakeInitiate)

	s.AddTool(mcplib.NewTool("portfolio_unstake",
		mcplib.WithDescription("Initiate a portfolio-level unstake request"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("symbol",
			mcplib.Description("Currency symbol to unstake (e.g. ETH)"),
		),
		mcplib.WithString("amount",
			mcplib.Description("Amount to unstake"),
		),
		mcplib.WithString("stake_protocol",
			mcplib.Description("Optional staking protocol identifier"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Auto-generated if omitted"),
		),
	), handlePortfolioUnstake)

	s.AddTool(mcplib.NewTool("query_validators",
		mcplib.WithDescription("Query transaction validators for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithArray("transaction_ids",
			mcplib.WithStringItems(),
			mcplib.Description("List of transaction IDs to query"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Maximum number of results to return"),
		),
	), handleQueryValidators)
}

func handleStake(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	request := &primeStaking.CreateStakeRequest{
		PortfolioId:    portfolioId,
		WalletId:       req.GetString("wallet_id", ""),
		IdempotencyKey: idempotencyKey,
	}

	if amount := req.GetString("amount", ""); amount != "" {
		request.Inputs = primeStaking.CreateStakeInputs{Amount: amount}
	}

	svc := primeStaking.NewStakingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateStake(ctx2, request)
	if err != nil {
		return toolErr("cannot create stake request: %s", err), nil
	}

	return marshalResult(response)
}

func handleUnstake(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	request := &primeStaking.CreateUnstakeRequest{
		PortfolioId:    portfolioId,
		WalletId:       req.GetString("wallet_id", ""),
		IdempotencyKey: idempotencyKey,
	}

	if amount := req.GetString("amount", ""); amount != "" {
		request.Inputs = primeStaking.CreateUnstakeInputs{Amount: amount}
	}

	svc := primeStaking.NewStakingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateUnstake(ctx2, request)
	if err != nil {
		return toolErr("cannot create unstake request: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetStakingStatus(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := primeStaking.NewStakingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetStakingStatus(ctx2, &primeStaking.GetStakingStatusRequest{
		PortfolioId: portfolioId,
		WalletId:    req.GetString("wallet_id", ""),
	})
	if err != nil {
		return toolErr("cannot get staking status: %s", err), nil
	}

	return marshalResult(response)
}

func handleClaimStakingRewards(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	request := &primeStaking.ClaimStakingRewardsRequest{
		PortfolioId:    portfolioId,
		WalletId:       req.GetString("wallet_id", ""),
		IdempotencyKey: idempotencyKey,
	}

	if amount := req.GetString("amount", ""); amount != "" {
		request.Inputs = &primeStaking.ClaimStakingRewardsInputs{Amount: amount}
	}

	svc := primeStaking.NewStakingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ClaimStakingRewards(ctx2, request)
	if err != nil {
		return toolErr("cannot claim staking rewards: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetUnstakingStatus(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := primeStaking.NewStakingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetUnstakingStatus(ctx2, &primeStaking.GetUnstakingStatusRequest{
		PortfolioId: portfolioId,
		WalletId:    req.GetString("wallet_id", ""),
	})
	if err != nil {
		return toolErr("cannot get unstaking status: %s", err), nil
	}

	return marshalResult(response)
}

func handlePreviewUnstake(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := primeStaking.NewStakingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.PreviewUnstake(ctx2, &primeStaking.PreviewUnstakeRequest{
		PortfolioId: portfolioId,
		WalletId:    req.GetString("wallet_id", ""),
		Amount:      req.GetString("amount", ""),
	})
	if err != nil {
		return toolErr("cannot preview unstake: %s", err), nil
	}

	return marshalResult(response)
}

func handlePortfolioStakeInitiate(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	svc := primeStaking.NewStakingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.PortfolioStakeInitiate(ctx2, &primeStaking.PortfolioStakeInitiateRequest{
		PortfolioId:    portfolioId,
		IdempotencyKey: idempotencyKey,
		CurrencySymbol: req.GetString("symbol", ""),
		Amount:         req.GetString("amount", ""),
	})
	if err != nil {
		return toolErr("cannot initiate portfolio stake: %s", err), nil
	}

	return marshalResult(response)
}

func handlePortfolioUnstake(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	svc := primeStaking.NewStakingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.PortfolioUnstake(ctx2, &primeStaking.PortfolioUnstakeRequest{
		PortfolioId:    portfolioId,
		IdempotencyKey: idempotencyKey,
		CurrencySymbol: req.GetString("symbol", ""),
		Amount:         req.GetString("amount", ""),
	})
	if err != nil {
		return toolErr("cannot initiate portfolio unstake: %s", err), nil
	}

	return marshalResult(response)
}

func handleQueryValidators(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := primeStaking.NewStakingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.QueryTransactionValidators(ctx2, &primeStaking.QueryTransactionValidatorsRequest{
		PortfolioId:    portfolioId,
		TransactionIds: req.GetStringSlice("transaction_ids", nil),
		Cursor:         req.GetString("cursor", ""),
		Limit:          int32(req.GetInt("limit", 0)),
	})
	if err != nil {
		return toolErr("cannot query transaction validators: %s", err), nil
	}

	return marshalResult(response)
}
