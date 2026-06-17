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
	"github.com/coinbase/prime-sdk-go/invoice"
	"github.com/coinbase/prime-sdk-go/model"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerInvoiceTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_invoices",
		mcplib.WithDescription("List invoices for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithArray("states",
			mcplib.Description("Filter by invoice states: INVOICE_STATE_IMPORTED, INVOICE_STATE_BILLED, INVOICE_STATE_PARTIALLY_PAID, INVOICE_STATE_PAID"),
			mcplib.WithStringItems(),
		),
		mcplib.WithInteger("billing_year",
			mcplib.Description("Billing year"),
		),
		mcplib.WithInteger("billing_month",
			mcplib.Description("Billing month"),
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
	), handleListInvoices)
}

func handleListInvoices(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	stateStrings := req.GetStringSlice("states", nil)
	invoiceStates := make([]model.InvoiceState, 0, len(stateStrings))
	for _, s := range stateStrings {
		invoiceStates = append(invoiceStates, model.InvoiceState(s))
	}

	svc := invoice.NewInvoiceService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListInvoices(ctx2, &invoice.ListInvoicesRequest{
		EntityId:     entityId,
		States:       invoiceStates,
		BillingYear:  int32(req.GetInt("billing_year", 0)),
		BillingMonth: int32(req.GetInt("billing_month", 0)),
		Pagination:   paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list invoices: %s", err), nil
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
