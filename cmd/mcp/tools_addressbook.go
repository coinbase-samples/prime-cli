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
	"github.com/coinbase/prime-sdk-go/addressbook"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerAddressBookTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_address_book",
		mcplib.WithDescription("List address book entries for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("symbol",
			mcplib.Description("Filter by asset symbol (e.g. BTC, ETH)"),
		),
		mcplib.WithString("search",
			mcplib.Description("Search by name or address"),
		),
		mcplib.WithString("cursor",
			mcplib.Description("Pagination cursor from a previous response"),
		),
		mcplib.WithInteger("limit",
			mcplib.Description("Number of results per page (default 50)"),
		),
	), handleListAddressBook)

	s.AddTool(mcplib.NewTool("create_address_book_entry",
		mcplib.WithDescription("Create a new address book entry for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("address",
			mcplib.Required(),
			mcplib.Description("Blockchain address to add"),
		),
		mcplib.WithString("symbol",
			mcplib.Required(),
			mcplib.Description("Asset symbol (e.g. BTC, ETH)"),
		),
		mcplib.WithString("name",
			mcplib.Required(),
			mcplib.Description("Display name for this address"),
		),
		mcplib.WithString("account_identifier",
			mcplib.Description("Account identifier for fiat addresses"),
		),
	), handleCreateAddressBookEntry)
}

func handleListAddressBook(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := addressbook.NewAddressBookService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetAddressBook(ctx2, &addressbook.GetAddressBookRequest{
		PortfolioId: portfolioId,
		Symbol:      req.GetString("symbol", ""),
		Search:      req.GetString("search", ""),
		Pagination:  paginationFor(req),
	})
	if err != nil {
		return toolErr("cannot list address book: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateAddressBookEntry(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	address, err := req.RequireString("address")
	if err != nil {
		return toolErr("address is required"), nil
	}

	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolErr("symbol is required"), nil
	}

	name, err := req.RequireString("name")
	if err != nil {
		return toolErr("name is required"), nil
	}

	svc := addressbook.NewAddressBookService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateAddressBookEntry(ctx2, &addressbook.CreateAddressBookEntryRequest{
		PortfolioId:       portfolioId,
		Address:           address,
		Symbol:            symbol,
		Name:              name,
		AccountIdentifier: req.GetString("account_identifier", ""),
	})
	if err != nil {
		return toolErr("cannot create address book entry: %s", err), nil
	}

	return marshalResult(response)
}
