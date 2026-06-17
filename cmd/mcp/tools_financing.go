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
	prime "github.com/coinbase/prime-sdk-go/financing"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerFinancingTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("get_buying_power",
		mcplib.WithDescription("Get buying power for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("base_currency",
			mcplib.Required(),
			mcplib.Description("Base currency symbol"),
		),
		mcplib.WithString("quote_currency",
			mcplib.Required(),
			mcplib.Description("Quote currency symbol"),
		),
	), handleGetBuyingPower)

	s.AddTool(mcplib.NewTool("get_portfolio_credit_info",
		mcplib.WithDescription("Get post-trade credit information for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetPortfolioCreditInfo)

	s.AddTool(mcplib.NewTool("get_cross_margin_overview",
		mcplib.WithDescription("Get cross-margin overview for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetCrossMarginOverview)

	s.AddTool(mcplib.NewTool("get_entity_locate_availabilities",
		mcplib.WithDescription("Get locate availabilities for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetEntityLocateAvailabilities)

	s.AddTool(mcplib.NewTool("get_margin_information",
		mcplib.WithDescription("Get margin information for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetMarginInformation)

	s.AddTool(mcplib.NewTool("get_pricing_fees",
		mcplib.WithDescription("Get trade finance tiered pricing fees for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleGetPricingFees)

	s.AddTool(mcplib.NewTool("get_withdrawal_power",
		mcplib.WithDescription("Get withdrawal power for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("symbol",
			mcplib.Required(),
			mcplib.Description("Asset symbol"),
		),
	), handleGetWithdrawalPower)

	s.AddTool(mcplib.NewTool("create_locate",
		mcplib.WithDescription("Create a new locate for a portfolio and asset"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("symbol",
			mcplib.Required(),
			mcplib.Description("Asset symbol"),
		),
		mcplib.WithString("amount",
			mcplib.Required(),
			mcplib.Description("Locate amount"),
		),
		mcplib.WithString("date",
			mcplib.Required(),
			mcplib.Description("Locate date RFC3339"),
		),
	), handleCreateLocate)

	s.AddTool(mcplib.NewTool("list_financing_eligible_assets",
		mcplib.WithDescription("List assets eligible for financing"),
	), handleListFinancingEligibleAssets)

	s.AddTool(mcplib.NewTool("list_interest_accruals",
		mcplib.WithDescription("List interest accruals for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("start_date",
			mcplib.Required(),
			mcplib.Description("Start date in RFC3339 format"),
		),
		mcplib.WithString("end_date",
			mcplib.Required(),
			mcplib.Description("End date in RFC3339 format"),
		),
	), handleListInterestAccruals)

	s.AddTool(mcplib.NewTool("list_locates",
		mcplib.WithDescription("List existing locates for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithArray("locate_ids",
			mcplib.WithStringItems(),
			mcplib.Description("Filter by specific locate IDs"),
		),
		mcplib.WithString("date",
			mcplib.Description("Locate date"),
		),
	), handleListLocates)

	s.AddTool(mcplib.NewTool("list_margin_call_summaries",
		mcplib.WithDescription("List margin call summaries for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("start_date",
			mcplib.Description("Start date in RFC3339 format"),
		),
		mcplib.WithString("end_date",
			mcplib.Description("End date in RFC3339 format"),
		),
	), handleListMarginCallSummaries)

	s.AddTool(mcplib.NewTool("list_margin_conversions",
		mcplib.WithDescription("List margin conversions for a portfolio (deprecated)"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("start_date",
			mcplib.Description("Start date in RFC3339 format"),
		),
		mcplib.WithString("end_date",
			mcplib.Description("End date in RFC3339 format"),
		),
	), handleListMarginConversions)

	s.AddTool(mcplib.NewTool("list_portfolio_interest_accruals",
		mcplib.WithDescription("List interest accruals for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("start_date",
			mcplib.Required(),
			mcplib.Description("Start date in RFC3339 format"),
		),
		mcplib.WithString("end_date",
			mcplib.Required(),
			mcplib.Description("End date in RFC3339 format"),
		),
	), handleListPortfolioInterestAccruals)
}

func handleGetBuyingPower(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetBuyingPower(ctx2, &prime.GetBuyingPowerRequest{
		PortfolioId:   portfolioId,
		BaseCurrency:  req.GetString("base_currency", ""),
		QuoteCurrency: req.GetString("quote_currency", ""),
	})
	if err != nil {
		return toolErr("cannot get buying power: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetPortfolioCreditInfo(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetPortfolioCreditInfo(ctx2, &prime.GetPortfolioCreditInfoRequest{
		PortfolioId: portfolioId,
	})
	if err != nil {
		return toolErr("cannot get portfolio credit information: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetCrossMarginOverview(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetCrossMarginOverview(ctx2, &prime.GetCrossMarginOverviewRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot get cross margin overview: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetEntityLocateAvailabilities(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetEntityLocateAvailabilities(ctx2, &prime.GetEntityLocateAvailabilitiesRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot get entity locate availabilities: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetMarginInformation(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetMarginInfo(ctx2, &prime.GetMarginInfoRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot get margin information: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetPricingFees(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetTieredPricingFees(ctx2, &prime.GetTieredPricingFeesRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot get pricing fees: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetWithdrawalPower(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetWithdrawalPower(ctx2, &prime.GetWithdrawalPowerRequest{
		PortfolioId: portfolioId,
		Symbol:      req.GetString("symbol", ""),
	})
	if err != nil {
		return toolErr("cannot get withdrawal power: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateLocate(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateLocate(ctx2, &prime.CreateLocateRequest{
		PortfolioId: portfolioId,
		Symbol:      req.GetString("symbol", ""),
		Amount:      req.GetString("amount", ""),
		LocateDate:  req.GetString("date", ""),
	})
	if err != nil {
		return toolErr("cannot create locate: %s", err), nil
	}

	return marshalResult(response)
}

func handleListFinancingEligibleAssets(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListFinancingEligibleAssets(ctx2, &prime.ListFinancingEligibleAssetsRequest{})
	if err != nil {
		return toolErr("cannot list financing eligible assets: %s", err), nil
	}

	return marshalResult(response)
}

func handleListInterestAccruals(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListInterestAccruals(ctx2, &prime.ListInterestAccrualsRequest{
		EntityId:    entityId,
		PortfolioId: req.GetString("portfolio_id", ""),
		StartDate:   req.GetString("start_date", ""),
		EndDate:     req.GetString("end_date", ""),
	})
	if err != nil {
		return toolErr("cannot list interest accruals: %s", err), nil
	}

	return marshalResult(response)
}

func handleListLocates(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListLocates(ctx2, &prime.ListLocatesRequest{
		PortfolioId: portfolioId,
		LocateIds:   req.GetStringSlice("locate_ids", nil),
		LocateDate:  req.GetString("date", ""),
	})
	if err != nil {
		return toolErr("cannot list locates: %s", err), nil
	}

	return marshalResult(response)
}

func handleListMarginCallSummaries(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListMarginCallSummaries(ctx2, &prime.ListMarginCallSummariesRequest{
		EntityId:  entityId,
		StartDate: req.GetString("start_date", ""),
		EndDate:   req.GetString("end_date", ""),
	})
	if err != nil {
		return toolErr("cannot list margin call summaries: %s", err), nil
	}

	return marshalResult(response)
}

func handleListMarginConversions(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListMarginConversions(ctx2, &prime.ListMarginConversionsRequest{
		PortfolioId: portfolioId,
		StartDate:   req.GetString("start_date", ""),
		EndDate:     req.GetString("end_date", ""),
	})
	if err != nil {
		return toolErr("cannot list margin conversions: %s", err), nil
	}

	return marshalResult(response)
}

func handleListPortfolioInterestAccruals(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := prime.NewFinancingService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListPortfolioInterestAccruals(ctx2, &prime.ListPortfolioInterestAccrualsRequest{
		PortfolioId: portfolioId,
		StartDate:   req.GetString("start_date", ""),
		EndDate:     req.GetString("end_date", ""),
	})
	if err != nil {
		return toolErr("cannot list portfolio interest accruals: %s", err), nil
	}

	return marshalResult(response)
}
