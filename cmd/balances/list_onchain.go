/**
 * Copyright 2025-present Coinbase Global, Inc.
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

package balances

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/balances"

	"github.com/spf13/cobra"
)

var listOnchainBalancesCmd = &cobra.Command{
	Use:   "list-onchain",
	Short: "Lists onchain balances that meet filter criteria",
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

		walletId, err := cmd.Flags().GetString(utils.WalletIdFlag)
		if err != nil {
			return fmt.Errorf("cannot get wallet id: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &balances.ListOnchainWalletBalancesRequest{
			PortfolioId: portfolioId,
			WalletId:    walletId,
		}

		response, err := balancesService.ListOnchainWalletBalances(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot list onchain balances: %w", err)
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
	Cmd.AddCommand(listOnchainBalancesCmd)

	listOnchainBalancesCmd.Flags().StringP(utils.WalletIdFlag, "", "", "Wallet ID")
	utils.AddPortfolioIdFlag(listOnchainBalancesCmd)

	listOnchainBalancesCmd.MarkFlagRequired(utils.WalletIdFlag)
}
