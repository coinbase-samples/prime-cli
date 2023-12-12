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

var createWalletCmd = &cobra.Command{
	Use:   "create-wallet",
	Short: "Create a new vault wallet.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		portfolioId := utils.GetPortfolioId(cmd, client)

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &prime.CreateWalletRequest{
			PortfolioId: portfolioId,
			Name:        utils.GetFlagStringValue(cmd, utils.NameFlag),
			Symbol:      utils.GetFlagStringValue(cmd, utils.SymbolFlag),
			Type:        utils.GetFlagStringValue(cmd, utils.TypeFlag),
		}

		response, err := client.CreateWallet(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create wallet: %w", err)
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
	rootCmd.AddCommand(createWalletCmd)

	createWalletCmd.Flags().StringP(utils.NameFlag, "n", "", "Name for the wallet (Required)")
	createWalletCmd.Flags().StringP(utils.SymbolFlag, "s", "", "Symbol for the wallet (Required)")
	createWalletCmd.Flags().StringP(utils.TypeFlag, "t", "", "Type of wallet, e.g. VAULT (Required)")
	createWalletCmd.Flags().StringP(utils.FormatFlag, "z", "false", "Pass true for formatted JSON. Default is false")
	createWalletCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")

	createWalletCmd.MarkFlagRequired(utils.NameFlag)
	createWalletCmd.MarkFlagRequired(utils.SymbolFlag)
	createWalletCmd.MarkFlagRequired(utils.TypeFlag)
}
