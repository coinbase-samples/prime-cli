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
	"encoding/json"
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
			return fmt.Errorf("error: %w", err)
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &prime.CreateWalletRequest{
			PortfolioId: client.Credentials.PortfolioId,
			Name:        utils.GetFlagStringValue(cmd, utils.NameFlag),
			Symbol:      utils.GetFlagStringValue(cmd, utils.SymbolFlag),
			Type:        utils.GetFlagStringValue(cmd, utils.TypeFlag),
		}

		response, err := client.CreateWallet(ctx, request)
		if err != nil {
			return fmt.Errorf("error creating wallet: %v", err)
		}

		jsonResponse, err := json.MarshalIndent(response, "", utils.JsonIndent)
		if err != nil {
			return fmt.Errorf("error marshaling response to JSON: %v", err)
		}

		fmt.Println(string(jsonResponse))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createWalletCmd)

	createWalletCmd.Flags().StringP(utils.NameFlag, "n", "", "Name for the wallet (Required)")
	createWalletCmd.Flags().StringP(utils.SymbolFlag, "s", "", "Symbol for the wallet (Required)")
	createWalletCmd.Flags().StringP(utils.TypeFlag, "t", "", "Type of wallet, e.g. VAULT (Required)")

	createWalletCmd.MarkFlagRequired(utils.NameFlag)
	createWalletCmd.MarkFlagRequired(utils.SymbolFlag)
	createWalletCmd.MarkFlagRequired(utils.TypeFlag)

}