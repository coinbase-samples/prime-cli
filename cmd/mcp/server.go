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
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

func runMCPServer(_ *cobra.Command, _ []string) error {
	s := server.NewMCPServer(
		"coinbase-prime",
		"0.4.2",
		server.WithToolCapabilities(false),
	)

	registerActivityTools(s)
	registerAddressBookTools(s)
	registerAdvancedTransferTools(s)
	registerAllocationTools(s)
	registerAssetTools(s)
	registerBalanceTools(s)
	registerCommissionTools(s)
	registerFinancingTools(s)
	registerFuturesTools(s)
	registerInvoiceTools(s)
	registerOnchainAddressBookTools(s)
	registerOrderTools(s)
	registerPaymentMethodTools(s)
	registerPortfolioTools(s)
	registerPositionTools(s)
	registerProductTools(s)
	registerStakingTools(s)
	registerTransactionTools(s)
	registerUserTools(s)
	registerWalletTools(s)

	return server.ServeStdio(s)
}
