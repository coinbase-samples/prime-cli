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
	"github.com/coinbase/prime-sdk-go/paymentmethods"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPaymentMethodTools(s *server.MCPServer) {
	s.AddTool(mcplib.NewTool("list_payment_methods",
		mcplib.WithDescription("List payment methods for an entity"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
	), handleListPaymentMethods)

	s.AddTool(mcplib.NewTool("get_payment_method",
		mcplib.WithDescription("Get payment method details by ID"),
		mcplib.WithString("entity_id",
			mcplib.Description("Uses credentials default if omitted"),
		),
		mcplib.WithString("payment_method_id",
			mcplib.Required(),
			mcplib.Description("Payment Method ID"),
		),
	), handleGetPaymentMethod)
}

func handleListPaymentMethods(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := paymentmethods.NewPaymentMethodsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.ListEntityPaymentMethods(ctx2, &paymentmethods.ListEntityPaymentMethodsRequest{
		EntityId: entityId,
	})
	if err != nil {
		return toolErr("cannot list payment methods: %s", err), nil
	}

	return marshalResult(response)
}

func handleGetPaymentMethod(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	client, err := utils.GetClientFromEnv()
	if err != nil {
		return toolErr("failed to initialize client: %s", err), nil
	}

	entityId, err := resolveEntityId(client, req)
	if err != nil {
		return toolErr("%s", err), nil
	}

	svc := paymentmethods.NewPaymentMethodsService(client)
	ctx2, cancel := mcpCtx(ctx)
	defer cancel()

	response, err := svc.GetEntityPaymentMethod(ctx2, &paymentmethods.GetEntityPaymentMethodRequest{
		Id:              entityId,
		PaymentMethodId: req.GetString("payment_method_id", ""),
	})
	if err != nil {
		return toolErr("cannot get payment method: %s", err), nil
	}

	return marshalResult(response)
}
