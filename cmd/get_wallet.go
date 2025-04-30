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
	"github.com/coinbase-samples/prime-sdk-go/wallets"

	"github.com/spf13/cobra"
)

var getWalletCmd = &cobra.Command{
	Use:   "get-wallet",
	Short: "Get a wallet given a wallet ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		walletsService := wallets.NewWalletsService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &wallets.GetWalletRequest{
			PortfolioId: portfolioId,
			Id:          utils.GetFlagStringValue(cmd, utils.WalletIdFlag),
		}

		response, err := walletsService.GetWallet(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get wallet: %w", err)
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
	rootCmd.AddCommand(getWalletCmd)

	getWalletCmd.Flags().StringP(utils.WalletIdFlag, "i", "", "Wallet ID (Required)")
	getWalletCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	getWalletCmd.MarkFlagRequired(utils.WalletIdFlag)
}
