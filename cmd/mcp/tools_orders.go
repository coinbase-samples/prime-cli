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
	"github.com/coinbase/prime-sdk-go/orders"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerOrderTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_orders",
		mcplib.WithDescription("List orders for a portfolio with optional filters"),
		mcplib.WithString("start",
			mcplib.Required(),
			mcplib.Description("Start time in RFC3339 format (e.g. 2024-01-01T00:00:00Z)"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("end",
			mcplib.Description("End time in RFC3339 format"),
		),
		mcplib.WithArray("product_ids",
			mcplib.Description("Filter by product IDs (e.g. [\"BTC-USD\"])"),
			mcplib.WithStringItems(),
		),
		mcplib.WithString("order_type",
			mcplib.Description("Filter by order type: MARKET, LIMIT, TWAP, BLOCK, VWAP, STOP_LIMIT, RFQ, or PEG"),
		),
		mcplib.WithString("order_side",
			mcplib.Description("Filter by side: BUY or SELL"),
		),
		mcplib.WithArray("statuses",
			mcplib.Description("Filter by order statuses: OPEN, FILLED, CANCELLED, EXPIRED, FAILED, or PENDING"),
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
	), handleListOrders)

	s.AddTool(mcplib.NewTool("list_open_orders",
		mcplib.WithDescription("List currently open (active) orders for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithArray("product_ids",
			mcplib.Description("Filter by product IDs"),
			mcplib.WithStringItems(),
		),
		mcplib.WithString("order_type",
			mcplib.Description("Filter by order type: MARKET, LIMIT, TWAP, BLOCK, VWAP, STOP_LIMIT, RFQ, or PEG"),
		),
		mcplib.WithString("order_side",
			mcplib.Description("Filter by side: BUY or SELL"),
		),
		mcplib.WithString("start",
			mcplib.Description("Start time in RFC3339 format"),
		),
		mcplib.WithString("end",
			mcplib.Description("End time in RFC3339 format"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page"),
		),
	), handleListOpenOrders)

	s.AddTool(mcplib.NewTool("get_order",
		mcplib.WithDescription("Get details of a specific order by order ID"),
		mcplib.WithString("order_id",
			mcplib.Required(),
			mcplib.Description("Order ID"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleGetOrder)

	s.AddTool(mcplib.NewTool("list_order_fills",
		mcplib.WithDescription("List fills (executions) for a specific order"),
		mcplib.WithString("order_id",
			mcplib.Required(),
			mcplib.Description("Order ID"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
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
	), handleListOrderFills)

	s.AddTool(mcplib.NewTool("preview_order",
		mcplib.WithDescription("Preview an order to see estimated cost and fees before submitting. Does not create an actual order."),
		mcplib.WithString("side",
			mcplib.Required(),
			mcplib.Description("Order side: BUY or SELL"),
		),
		mcplib.WithString("type",
			mcplib.Required(),
			mcplib.Description("Order type: MARKET, LIMIT, TWAP, BLOCK, VWAP, STOP_LIMIT, RFQ, or PEG"),
		),
		mcplib.WithString("product_id",
			mcplib.Required(),
			mcplib.Description("Product ID (e.g. BTC-USD). Use list_products to discover valid values."),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("base_quantity",
			mcplib.Description("Order size in base asset units (e.g. 0.5 for 0.5 BTC)"),
		),
		mcplib.WithString("quote_value",
			mcplib.Description("Order size in quote asset units (e.g. 10000 for $10,000 USD)"),
		),
		mcplib.WithString("limit_price",
			mcplib.Description("Limit price (required for LIMIT orders)"),
		),
		mcplib.WithString("time_in_force",
			mcplib.Description("Time in force: GTC (Good-Till-Cancelled), GTD (Good-Till-Date), IOC (Immediate-or-Cancel), or FOK (Fill-or-Kill)"),
		),
	), handlePreviewOrder)

	s.AddTool(mcplib.NewTool("create_order",
		mcplib.WithDescription("Submit a new order. WARNING: This executes a real financial transaction on Coinbase Prime. Use preview_order first to verify parameters."),
		mcplib.WithString("side",
			mcplib.Required(),
			mcplib.Description("Order side: BUY or SELL"),
		),
		mcplib.WithString("type",
			mcplib.Required(),
			mcplib.Description("Order type: MARKET, LIMIT, TWAP, BLOCK, VWAP, STOP_LIMIT, RFQ, or PEG"),
		),
		mcplib.WithString("product_id",
			mcplib.Required(),
			mcplib.Description("Product ID (e.g. BTC-USD). Use list_products to discover valid values."),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("base_quantity",
			mcplib.Description("Order size in base asset units (e.g. 0.5 for 0.5 BTC). Specify either base_quantity or quote_value, not both."),
		),
		mcplib.WithString("quote_value",
			mcplib.Description("Order size in quote asset units (e.g. 10000 for $10,000 USD). Specify either base_quantity or quote_value, not both."),
		),
		mcplib.WithString("limit_price",
			mcplib.Description("Limit price (required for LIMIT orders)"),
		),
		mcplib.WithString("time_in_force",
			mcplib.Description("Time in force: GTC (Good-Till-Cancelled), GTD (Good-Till-Date), IOC (Immediate-or-Cancel), or FOK (Fill-or-Kill)"),
		),
		mcplib.WithString("start_time",
			mcplib.Description("Order start time in UTC (TWAP orders only)"),
		),
		mcplib.WithString("expiry_time",
			mcplib.Description("Order expiry time in UTC (TWAP and limit GTD orders)"),
		),
		mcplib.WithString("client_order_id",
			mcplib.Description("Client-supplied order ID for idempotency. Auto-generated UUID if omitted."),
		),
	), handleCreateOrder)

	s.AddTool(mcplib.NewTool("cancel_order",
		mcplib.WithDescription("Attempt to cancel an open order"),
		mcplib.WithString("order_id",
			mcplib.Required(),
			mcplib.Description("Order ID to cancel"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleCancelOrder)

	s.AddTool(mcplib.NewTool("edit_order",
		mcplib.WithDescription("Edit an existing open order's quantity or price"),
		mcplib.WithString("order_id",
			mcplib.Required(),
			mcplib.Description("Order ID to edit"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("new_base_quantity",
			mcplib.Description("Updated order size in base asset units"),
		),
		mcplib.WithString("new_quote_value",
			mcplib.Description("Updated order size in quote asset units"),
		),
		mcplib.WithString("new_limit_price",
			mcplib.Description("Updated limit price"),
		),
		mcplib.WithString("client_order_id",
			mcplib.Description("Updated client order ID"),
		),
	), handleEditOrder)

	s.AddTool(mcplib.NewTool("get_order_edit_history",
		mcplib.WithDescription("Get the edit history for an order"),
		mcplib.WithString("order_id",
			mcplib.Required(),
			mcplib.Description("Order ID"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleGetOrderEditHistory)

	s.AddTool(mcplib.NewTool("create_quote",
		mcplib.WithDescription("Create a quote request"),
		mcplib.WithString("side",
			mcplib.Required(),
			mcplib.Description("Order side: BUY or SELL"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("product_id",
			mcplib.Description("Product ID (e.g. BTC-USD)"),
		),
		mcplib.WithString("limit_price",
			mcplib.Description("Limit price"),
		),
		mcplib.WithString("quote_value",
			mcplib.Description("Order size in quote asset units"),
		),
		mcplib.WithString("base_quantity",
			mcplib.Description("Order size in base asset units"),
		),
		mcplib.WithString("client_quote_id",
			mcplib.Description("Client quote ID. Auto-generated if omitted"),
		),
		mcplib.WithString("settle_currency",
			mcplib.Description("Settlement currency"),
		),
	), handleCreateQuote)

	s.AddTool(mcplib.NewTool("accept_quote",
		mcplib.WithDescription("Accept a quote and create an order"),
		mcplib.WithString("side",
			mcplib.Required(),
			mcplib.Description("Order side: BUY or SELL"),
		),
		mcplib.WithString("quote_id",
			mcplib.Required(),
			mcplib.Description("Quote ID to accept"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("product_id",
			mcplib.Description("Product ID (e.g. BTC-USD)"),
		),
		mcplib.WithString("client_order_id",
			mcplib.Description("Client order ID. Auto-generated if omitted"),
		),
	), handleAcceptQuote)

	s.AddTool(mcplib.NewTool("list_portfolio_fills",
		mcplib.WithDescription("List all fills for a portfolio"),
		mcplib.WithString("start",
			mcplib.Required(),
			mcplib.Description("Start time in RFC3339 format (e.g. 2024-01-01T00:00:00Z)"),
		),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
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
		mcplib.WithBoolean("fetch_all",
			mcplib.Description("Fetch all pages automatically and return combined results. When true, cursor and limit are ignored."),
		),
	), handleListPortfolioFills)
}

func handleListOrders(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	startStr, reqErr := req.RequireString("start")
	if reqErr != nil {
		return toolErr("start is required"), nil
	}

	start, end, err := utils.ParseDateRange(startStr, req.GetString("end", ""))
	if err != nil {
		return toolErr("invalid date range: %s", err), nil
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListOrders(ctx2, &orders.ListOrdersRequest{
		PortfolioId: portfolioId,
		ProductIds:  req.GetStringSlice("product_ids", nil),
		Type:        req.GetString("order_type", ""),
		OrderSide:   req.GetString("order_side", ""),
		Statuses:    req.GetStringSlice("statuses", nil),
		Start:       start,
		End:         end,
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list orders: %s", err), nil
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

func handleListOpenOrders(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	openReq := &orders.ListOpenOrdersRequest{
		PortfolioId: portfolioId,
		ProductIds:  req.GetStringSlice("product_ids", nil),
		OrderType:   req.GetString("order_type", ""),
		OrderSide:   req.GetString("order_side", ""),
		Start:       start,
		End:         end,
		Pagination:  paginationFor(req),
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListOpenOrders(ctx2, openReq)
	if err != nil {
		return toolErr("cannot list open orders: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetOrder(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	orderId, reqErr := req.RequireString("order_id")
	if reqErr != nil {
		return toolErr("order_id is required"), nil
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetOrder(ctx2, &orders.GetOrderRequest{
		PortfolioId: portfolioId,
		OrderId:     orderId,
	})
	if err != nil {
		return toolErr("cannot get order: %s", err), nil
	}

	return marshalResult(response.Order)
}

func handleListOrderFills(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	orderId, reqErr := req.RequireString("order_id")
	if reqErr != nil {
		return toolErr("order_id is required"), nil
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListOrderFills(ctx2, &orders.ListOrderFillsRequest{
		PortfolioId: portfolioId,
		OrderId:     orderId,
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list order fills: %s", err), nil
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

func handlePreviewOrder(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	side, reqErr := req.RequireString("side")
	if reqErr != nil {
		return toolErr("side is required (BUY or SELL)"), nil
	}
	orderType, reqErr := req.RequireString("type")
	if reqErr != nil {
		return toolErr("type is required"), nil
	}
	productId, reqErr := req.RequireString("product_id")
	if reqErr != nil {
		return toolErr("product_id is required"), nil
	}

	order := &model.Order{
		PortfolioId:  portfolioId,
		Side:         side,
		Type:         orderType,
		ProductId:    productId,
		BaseQuantity: req.GetString("base_quantity", ""),
		QuoteValue:   req.GetString("quote_value", ""),
		LimitPrice:   req.GetString("limit_price", ""),
		TimeInForce:  req.GetString("time_in_force", ""),
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateOrderPreview(ctx2, &orders.CreateOrderRequest{Order: order})
	if err != nil {
		return toolErr("cannot preview order: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateOrder(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	side, reqErr := req.RequireString("side")
	if reqErr != nil {
		return toolErr("side is required (BUY or SELL)"), nil
	}
	orderType, reqErr := req.RequireString("type")
	if reqErr != nil {
		return toolErr("type is required"), nil
	}
	productId, reqErr := req.RequireString("product_id")
	if reqErr != nil {
		return toolErr("product_id is required"), nil
	}

	clientOrderId := req.GetString("client_order_id", "")
	if clientOrderId == "" {
		clientOrderId = utils.NewUuidStr()
	}

	order := &model.Order{
		PortfolioId:   portfolioId,
		Side:          side,
		Type:          orderType,
		ProductId:     productId,
		ClientOrderId: clientOrderId,
		BaseQuantity:  req.GetString("base_quantity", ""),
		QuoteValue:    req.GetString("quote_value", ""),
		LimitPrice:    req.GetString("limit_price", ""),
		TimeInForce:   req.GetString("time_in_force", ""),
		StartTime:     req.GetString("start_time", ""),
		ExpiryTime:    req.GetString("expiry_time", ""),
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateOrder(ctx2, &orders.CreateOrderRequest{Order: order})
	if err != nil {
		return toolErr("cannot create order: %s", err), nil
	}

	return marshalResult(response)
}

func handleCancelOrder(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	orderId, reqErr := req.RequireString("order_id")
	if reqErr != nil {
		return toolErr("order_id is required"), nil
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CancelOrder(ctx2, &orders.CancelOrderRequest{
		PortfolioId: portfolioId,
		OrderId:     orderId,
	})
	if err != nil {
		return toolErr("cannot cancel order: %s", err), nil
	}

	return marshalResult(response)
}

func handleEditOrder(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	orderId, reqErr := req.RequireString("order_id")
	if reqErr != nil {
		return toolErr("order_id is required"), nil
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.EditOrder(ctx2, &orders.EditOrderRequest{
		PortfolioId:   portfolioId,
		OrderId:       orderId,
		ClientOrderId: req.GetString("client_order_id", ""),
		BaseQuantity:  req.GetString("new_base_quantity", ""),
		QuoteValue:    req.GetString("new_quote_value", ""),
		LimitPrice:    req.GetString("new_limit_price", ""),
	})
	if err != nil {
		return toolErr("cannot edit order: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetOrderEditHistory(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	orderId, reqErr := req.RequireString("order_id")
	if reqErr != nil {
		return toolErr("order_id is required"), nil
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetOrderEditHistory(ctx2, &orders.GetOrderEditHistoryRequest{
		PortfolioId: portfolioId,
		OrderId:     orderId,
	})
	if err != nil {
		return toolErr("cannot get order edit history: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateQuote(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	side, reqErr := req.RequireString("side")
	if reqErr != nil {
		return toolErr("side is required (BUY or SELL)"), nil
	}

	clientQuoteId := req.GetString("client_quote_id", "")
	if clientQuoteId == "" {
		clientQuoteId = utils.NewUuidStr()
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateQuoteRequest(ctx2, &orders.CreateQuoteRequest{
		PortfolioId:    portfolioId,
		ProductId:      req.GetString("product_id", ""),
		ClientQuoteId:  clientQuoteId,
		Side:           model.OrderSide(side),
		BaseQuantity:   req.GetString("base_quantity", ""),
		QuoteValue:     req.GetString("quote_value", ""),
		LimitPrice:     req.GetString("limit_price", ""),
		SettleCurrency: req.GetString("settle_currency", ""),
	})
	if err != nil {
		return toolErr("cannot create quote: %s", err), nil
	}

	return marshalResult(response)
}

func handleAcceptQuote(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	side, reqErr := req.RequireString("side")
	if reqErr != nil {
		return toolErr("side is required (BUY or SELL)"), nil
	}

	quoteId, reqErr := req.RequireString("quote_id")
	if reqErr != nil {
		return toolErr("quote_id is required"), nil
	}

	clientOrderId := req.GetString("client_order_id", "")
	if clientOrderId == "" {
		clientOrderId = utils.NewUuidStr()
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.AcceptQuote(ctx2, &orders.AcceptQuoteRequest{
		PortfolioId:   portfolioId,
		ProductId:     req.GetString("product_id", ""),
		QuoteId:       quoteId,
		ClientOrderId: clientOrderId,
		Side:          side,
	})
	if err != nil {
		return toolErr("cannot accept quote: %s", err), nil
	}

	return marshalResult(response)
}

func handleListPortfolioFills(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	startStr, reqErr := req.RequireString("start")
	if reqErr != nil {
		return toolErr("start is required"), nil
	}

	start, end, err := utils.ParseDateRange(startStr, req.GetString("end", ""))
	if err != nil {
		return toolErr("invalid date range: %s", err), nil
	}

	svc := orders.NewOrdersService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListPortfolioFills(ctx2, &orders.ListPortfolioFillsRequest{
		PortfolioId: portfolioId,
		Start:       start,
		End:         end,
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list portfolio fills: %s", err), nil
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
