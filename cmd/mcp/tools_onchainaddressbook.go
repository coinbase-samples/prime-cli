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
	"github.com/coinbase/prime-sdk-go/onchainaddressbook"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerOnchainAddressBookTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_onchain_address_groups",
		mcplib.WithDescription("List onchain address book groups for a portfolio"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
	), handleListOnchainAddressGroups)

	s.AddTool(mcplib.NewTool("create_onchain_address_group",
		mcplib.WithDescription("Create an onchain address book group entry"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("id",
			mcplib.Required(),
			mcplib.Description("Address group ID"),
		),
		mcplib.WithString("address",
			mcplib.Required(),
			mcplib.Description("Blockchain address"),
		),
		mcplib.WithString("network_type",
			mcplib.Required(),
			mcplib.Description("Network type: NETWORK_TYPE_EVM, NETWORK_TYPE_SOLANA, or NETWORK_TYPE_UNSPECIFIED"),
		),
		mcplib.WithString("name",
			mcplib.Description("Display name for this address group"),
		),
	), handleCreateOnchainAddressGroup)

	s.AddTool(mcplib.NewTool("update_onchain_address_group",
		mcplib.WithDescription("Update an onchain address book group entry"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("id",
			mcplib.Required(),
			mcplib.Description("Address group ID"),
		),
		mcplib.WithString("address",
			mcplib.Required(),
			mcplib.Description("Blockchain address"),
		),
		mcplib.WithString("network_type",
			mcplib.Required(),
			mcplib.Description("Network type: NETWORK_TYPE_EVM, NETWORK_TYPE_SOLANA, or NETWORK_TYPE_UNSPECIFIED"),
		),
		mcplib.WithString("name",
			mcplib.Description("Display name for this address group"),
		),
	), handleUpdateOnchainAddressGroup)

	s.AddTool(mcplib.NewTool("delete_onchain_address_group",
		mcplib.WithDescription("Delete an onchain address book group entry"),
		mcplib.WithString("portfolio_id",
			mcplib.Description("Portfolio ID. Uses credentials default if omitted"),
		),
		mcplib.WithString("id",
			mcplib.Required(),
			mcplib.Description("Address group ID to delete"),
		),
	), handleDeleteOnchainAddressGroup)
}

var onchainNetworkTypeMap = map[string]model.OnchainNetworkType{
	"NETWORK_TYPE_EVM":         model.OnchainNetworkTypeEvm,
	"NETWORK_TYPE_SOLANA":      model.OnchainNetworkTypeSolana,
	"NETWORK_TYPE_UNSPECIFIED": model.OnchainNetworkTypeUnspecified,
}

func handleListOnchainAddressGroups(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := onchainaddressbook.NewOnchainAddressBookService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListOnchainAddressBookGroups(ctx2, &onchainaddressbook.ListOnchainAddressBookGroupsRequest{
		PortfolioId: portfolioId,
	})
	if err != nil {
		return toolErr("cannot list onchain address groups: %s", err), nil
	}

	return marshalResult(response)
}

func handleCreateOnchainAddressGroup(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	id, err := req.RequireString("id")
	if err != nil {
		return toolErr("id is required"), nil
	}

	address, err := req.RequireString("address")
	if err != nil {
		return toolErr("address is required"), nil
	}

	networkTypeStr, err := req.RequireString("network_type")
	if err != nil {
		return toolErr("network_type is required"), nil
	}

	networkType := onchainNetworkTypeMap[networkTypeStr]
	name := req.GetString("name", "")

	svc := onchainaddressbook.NewOnchainAddressBookService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.CreateOnchainAddressBookEntry(ctx2, &onchainaddressbook.CreateOnchainAddressBookEntryRequest{
		PortfolioId: portfolioId,
		AddressGroup: &model.OnchainAddressGroup{
			Id:          id,
			Name:        name,
			NetworkType: networkType,
			Addresses:   []*model.OnchainAddress{{Address: address, Name: name}},
		},
	})
	if err != nil {
		return toolErr("cannot create onchain address group: %s", err), nil
	}

	return marshalResult(response)
}

func handleUpdateOnchainAddressGroup(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	id, err := req.RequireString("id")
	if err != nil {
		return toolErr("id is required"), nil
	}

	address, err := req.RequireString("address")
	if err != nil {
		return toolErr("address is required"), nil
	}

	networkTypeStr, err := req.RequireString("network_type")
	if err != nil {
		return toolErr("network_type is required"), nil
	}

	networkType := onchainNetworkTypeMap[networkTypeStr]
	name := req.GetString("name", "")

	svc := onchainaddressbook.NewOnchainAddressBookService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.UpdateOnchainAddressBookEntry(ctx2, &onchainaddressbook.UpdateOnchainAddressBookEntryRequest{
		PortfolioId: portfolioId,
		AddressGroup: &model.OnchainAddressGroup{
			Id:          id,
			Name:        name,
			NetworkType: networkType,
			Addresses:   []*model.OnchainAddress{{Address: address, Name: name}},
		},
	})
	if err != nil {
		return toolErr("cannot update onchain address group: %s", err), nil
	}

	return marshalResult(response)
}

func handleDeleteOnchainAddressGroup(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	portfolioId, err := resolvePortfolioId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	id, err := req.RequireString("id")
	if err != nil {
		return toolErr("id is required"), nil
	}

	svc := onchainaddressbook.NewOnchainAddressBookService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.DeleteOnchainAddressBookEntry(ctx2, &onchainaddressbook.DeleteOnchainAddressBookEntryRequest{
		PortfolioId:    portfolioId,
		AddressGroupId: id,
	})
	if err != nil {
		return toolErr("cannot delete onchain address group: %s", err), nil
	}

	return marshalResult(response)
}
