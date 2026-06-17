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

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase/prime-sdk-go/allocations"
	"github.com/coinbase/prime-sdk-go/model"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerAllocationTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_allocations",
		mcplib.WithDescription("List historical allocations for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithArray("product_ids",
			mcplib.WithStringItems(),
			mcplib.Description("Filter by product IDs"),
		),
		mcplib.WithString("order_side",
			mcplib.Description("BUY or SELL"),
		),
		mcplib.WithString("start",
			mcplib.Required(),
			mcplib.Description("Start date RFC3339"),
		),
		mcplib.WithString("end",
			mcplib.Description("End date RFC3339"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
	), handleListAllocations)

	s.AddTool(mcplib.NewTool("get_allocation",
		mcplib.WithDescription("Get an allocation by ID"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("allocation_id",
			mcplib.Required(),
			mcplib.Description("ID of the allocation"),
		),
	), handleGetAllocation)

	s.AddTool(mcplib.NewTool("create_allocation",
		mcplib.WithDescription("Create a portfolio allocation"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("allocation_id",
			mcplib.Required(),
			mcplib.Description("ID of the allocation"),
		),
		mcplib.WithString("source_portfolio_id",
			mcplib.Required(),
			mcplib.Description("ID of the source portfolio"),
		),
		mcplib.WithString("product_id",
			mcplib.Required(),
			mcplib.Description("ID of the product"),
		),
		mcplib.WithString("size_type",
			mcplib.Required(),
			mcplib.Description("Size type of the allocation"),
		),
		mcplib.WithString("remainder_dest_portfolio_id",
			mcplib.Required(),
			mcplib.Description("ID of the remainder destination portfolio"),
		),
		mcplib.WithString("allocation_legs",
			mcplib.Required(),
			mcplib.Description("JSON array of allocation legs"),
		),
		mcplib.WithArray("order_ids",
			mcplib.WithStringItems(),
			mcplib.Description("List of order IDs"),
		),
	), handleCreateAllocation)

	s.AddTool(mcplib.NewTool("get_net_allocation",
		mcplib.WithDescription("Get a net allocation by netting ID"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("netting_id",
			mcplib.Required(),
			mcplib.Description("ID of the net allocation"),
		),
	), handleGetNetAllocation)

	s.AddTool(mcplib.NewTool("create_net_allocation",
		mcplib.WithDescription("Create a net portfolio allocation"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("netting_id",
			mcplib.Description("Netting ID for the allocation"),
		),
		mcplib.WithString("source_portfolio_id",
			mcplib.Required(),
			mcplib.Description("ID of the source portfolio"),
		),
		mcplib.WithString("product_id",
			mcplib.Required(),
			mcplib.Description("ID of the product"),
		),
		mcplib.WithString("size_type",
			mcplib.Required(),
			mcplib.Description("Size type of the allocation"),
		),
		mcplib.WithString("remainder_dest_portfolio_id",
			mcplib.Required(),
			mcplib.Description("ID of the remainder destination portfolio"),
		),
		mcplib.WithString("allocation_legs",
			mcplib.Required(),
			mcplib.Description("JSON array of allocation legs"),
		),
		mcplib.WithArray("order_ids",
			mcplib.WithStringItems(),
			mcplib.Required(),
			mcplib.Description("List of order IDs"),
		),
		mcplib.WithString("idempotency_key",
			mcplib.Description("Auto-generated if omitted"),
		),
	), handleCreateNetAllocation)
}

func handleListAllocations(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	startStr := req.GetString("start", "")
	endStr := req.GetString("end", "")
	start, end, err := utils.ParseDateRange(startStr, endStr)
	if err != nil {
		return toolErr("invalid date range: %s", err), nil
	}

	svc := allocations.NewAllocationsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListPortfolioAllocations(ctx2, &allocations.ListPortfolioAllocationsRequest{
		PortfolioId: portfolioId,
		ProductIds:  req.GetStringSlice("product_ids", nil),
		Side:        req.GetString("order_side", ""),
		Start:       start,
		End:         end,
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list allocations: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetAllocation(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := allocations.NewAllocationsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetPortfolioAllocation(ctx2, &allocations.GetPortfolioAllocationRequest{
		PortfolioId:  portfolioId,
		AllocationId: req.GetString("allocation_id", ""),
	})
	if err != nil {
		return toolErr("cannot get allocation: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateAllocation(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	legsStr := req.GetString("allocation_legs", "")
	var allocationLegs []*model.AllocationLeg
	if err := json.Unmarshal([]byte(legsStr), &allocationLegs); err != nil {
		return toolErr("invalid allocation legs format: %s", err), nil
	}

	svc := allocations.NewAllocationsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreatePortfolioAllocations(ctx2, &allocations.CreatePortfolioAllocationsRequest{
		AllocationId:                    req.GetString("allocation_id", ""),
		SourcePortfolioId:               req.GetString("source_portfolio_id", ""),
		ProductId:                       req.GetString("product_id", ""),
		SizeType:                        req.GetString("size_type", ""),
		RemainderDestinationPortfolioId: req.GetString("remainder_dest_portfolio_id", ""),
		AllocationLegs:                  allocationLegs,
		OrderIds:                        req.GetStringSlice("order_ids", nil),
	})
	if err != nil {
		return toolErr("cannot create allocation: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetNetAllocation(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := allocations.NewAllocationsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetPortfolioNetAllocation(ctx2, &allocations.GetPortfolioNetAllocationRequest{
		PortfolioId: portfolioId,
		NettingId:   req.GetString("netting_id", ""),
	})
	if err != nil {
		return toolErr("cannot get net allocation: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateNetAllocation(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	nettingId := req.GetString("netting_id", "")
	if nettingId == "" {
		key := req.GetString("idempotency_key", "")
		if key == "" {
			key = utils.NewUuidStr()
		}
		nettingId = key
	}

	legsStr := req.GetString("allocation_legs", "")
	var allocationLegs []*model.AllocationLeg
	if err := json.Unmarshal([]byte(legsStr), &allocationLegs); err != nil {
		return toolErr("invalid allocation legs format: %s", err), nil
	}

	svc := allocations.NewAllocationsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreatePortfolioNetAllocations(ctx2, &allocations.CreatePortfolioNetAllocationsRequest{
		NettingId:                       nettingId,
		SourcePortfolioId:               req.GetString("source_portfolio_id", ""),
		ProductId:                       req.GetString("product_id", ""),
		SizeType:                        req.GetString("size_type", ""),
		RemainderDestinationPortfolioId: req.GetString("remainder_dest_portfolio_id", ""),
		AllocationLegs:                  allocationLegs,
		OrderIds:                        req.GetStringSlice("order_ids", nil),
	})
	if err != nil {
		return toolErr("cannot create net allocation: %s", err), nil
	}

	return marshalResult(response)
}
