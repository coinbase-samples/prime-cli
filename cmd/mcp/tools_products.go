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
	"github.com/coinbase/prime-sdk-go/products"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerProductTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("get_product_candles",
		mcplib.WithDescription("Get candlestick data for a product"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("product_id",
			mcplib.Required(),
			mcplib.Description("Product ID e.g. BTC-USD"),
		),
		mcplib.WithString("granularity",
			mcplib.Required(),
			mcplib.Description("Candle granularity: ONE_MINUTE, FIVE_MINUTES, FIFTEEN_MINUTES, THIRTY_MINUTES, ONE_HOUR, TWO_HOURS, FOUR_HOURS, SIX_HOURS, ONE_DAY"),
		),
		mcplib.WithString("start",
			mcplib.Required(),
			mcplib.Description("Start time RFC3339"),
		),
		mcplib.WithString("end",
			mcplib.Required(),
			mcplib.Description("End time RFC3339"),
		),
	), handleGetProductCandles)

	s.AddTool(mcplib.NewTool("list_products",
		mcplib.WithDescription("List all tradeable products (trading pairs) available in a portfolio. Use this to discover valid product_id values for order tools."),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
	), handleListProducts)
}

func handleGetProductCandles(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	svc := products.NewProductsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetProductCandles(ctx2, &products.GetProductCandlesRequest{
		PortfolioId: portfolioId,
		ProductId:   req.GetString("product_id", ""),
		Granularity: model.CandleGranularity(req.GetString("granularity", "")),
		StartTime:   start,
		EndTime:     end,
	})
	if err != nil {
		return toolErr("cannot get product candles: %s", err), nil
	}

	return marshalResult(response)
}

func handleListProducts(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := products.NewProductsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListProducts(ctx2, &products.ListProductsRequest{
		PortfolioId: portfolioId,
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list products: %s", err), nil
	}

	return marshalResult(response)
}
