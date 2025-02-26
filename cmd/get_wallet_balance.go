/**
 * Copyright 2023-present Coinbase Global, Inc.
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

package cmd

import (
	"fmt"
	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/balances"

	"github.com/spf13/cobra"
)

var getWalletBalanceCmd = &cobra.Command{
	Use:   "get-wallet-balance",
	Short: "Get a wallet balance given a wallet ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		balancesService := balances.NewBalancesService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &balances.GetWalletBalanceRequest{
			PortfolioId: portfolioId,
			Id:          utils.GetFlagStringValue(cmd, utils.WalletIdFlag),
		}

		response, err := balancesService.GetWalletBalance(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get wallet balance: %w", err)
		}

		jsonResponse, err := utils.FormatResponseAsJson(cmd, response)
		if err != nil {
			return err
		}

		fmt.Println(jsonResponse)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getWalletBalanceCmd)

	getWalletBalanceCmd.Flags().StringP(utils.WalletIdFlag, "i", "", "Wallet ID (Required)")
	getWalletBalanceCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	getWalletBalanceCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	getWalletBalanceCmd.MarkFlagRequired(utils.WalletIdFlag)
}
