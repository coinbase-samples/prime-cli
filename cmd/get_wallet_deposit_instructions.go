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
	"github.com/coinbase-samples/prime-sdk-go"

	"github.com/spf13/cobra"
)

var getWalletDepositInstructionsCmd = &cobra.Command{
	Use:   "get-wallet-deposit-instructions",
	Short: "Get a wallet's public address given a wallet ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		portfolioId := utils.GetPortfolioId(cmd, client)

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &prime.GetWalletDepositInstructionsRequest{
			PortfolioId: portfolioId,
			Id:          utils.GetFlagStringValue(cmd, utils.WalletIdFlag),
			Type:        utils.GetFlagStringValue(cmd, utils.DepositTypeFlag),
		}

		response, err := client.GetWalletDepositInstructions(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get wallet deposit instructions: %w", err)
		}

		jsonResponse, err := utils.FormatResponseAsJSON(cmd, response)
		if err != nil {
			return err
		}

		fmt.Println(jsonResponse)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getWalletDepositInstructionsCmd)

	getWalletDepositInstructionsCmd.Flags().StringP(utils.WalletIdFlag, "i", "", "Wallet ID (Required)")
	getWalletDepositInstructionsCmd.Flags().StringP(utils.DepositTypeFlag, "d", "", "Wallet deposit type (Required)")
	getWalletDepositInstructionsCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	getWalletDepositInstructionsCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	getWalletDepositInstructionsCmd.MarkFlagRequired(utils.WalletIdFlag)
	getWalletDepositInstructionsCmd.MarkFlagRequired(utils.DepositTypeFlag)
}
