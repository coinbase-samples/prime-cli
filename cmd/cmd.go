/**
 * Copyright 2023-present Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"os"

	"github.com/coinbase-samples/prime-cli/cmd/activities"
	"github.com/coinbase-samples/prime-cli/cmd/addressbook"
	"github.com/coinbase-samples/prime-cli/cmd/allocations"
	"github.com/coinbase-samples/prime-cli/cmd/assets"
	"github.com/coinbase-samples/prime-cli/cmd/balances"
	"github.com/coinbase-samples/prime-cli/cmd/commission"
	"github.com/coinbase-samples/prime-cli/cmd/invoices"
	"github.com/coinbase-samples/prime-cli/cmd/onchainaddressbook"
	"github.com/coinbase-samples/prime-cli/cmd/orders"
	"github.com/coinbase-samples/prime-cli/cmd/paymentmethods"
	"github.com/coinbase-samples/prime-cli/cmd/portfolios"
	"github.com/coinbase-samples/prime-cli/cmd/products"
	"github.com/coinbase-samples/prime-cli/cmd/transactions"
	"github.com/coinbase-samples/prime-cli/cmd/users"
	"github.com/coinbase-samples/prime-cli/cmd/wallets"
	"github.com/coinbase-samples/prime-cli/utils"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "prime-cli",
	Short: "Root of prime cli",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool(utils.FormatFlag, false, "Set to include formatted JSON. Default is false")
	rootCmd.AddCommand(activities.Cmd)
	rootCmd.AddCommand(addressbook.Cmd)
	rootCmd.AddCommand(allocations.Cmd)
	rootCmd.AddCommand(assets.Cmd)
	rootCmd.AddCommand(balances.Cmd)
	rootCmd.AddCommand(commission.Cmd)
	rootCmd.AddCommand(invoices.Cmd)
	rootCmd.AddCommand(onchainaddressbook.Cmd)
	rootCmd.AddCommand(orders.Cmd)
	rootCmd.AddCommand(paymentmethods.Cmd)
	rootCmd.AddCommand(portfolios.Cmd)
	rootCmd.AddCommand(products.Cmd)
	rootCmd.AddCommand(transactions.Cmd)
	rootCmd.AddCommand(users.Cmd)
	rootCmd.AddCommand(wallets.Cmd)
}
