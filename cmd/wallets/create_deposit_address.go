/**
 * Copyright 2025-present Coinbase Global, Inc.
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

package wallets

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/wallets"

	"github.com/spf13/cobra"
)

var createWalletDepositAddressCmd = &cobra.Command{
	Use:   "create-deposit-address",
	Short: "Create a new wallet deposit address.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &wallets.CreateWalletAddressRequest{
			PortfolioId: portfolioId,
			WalletId:    utils.GetFlagStringValue(cmd, utils.WalletIdFlag),
			NetworkId:   utils.GetFlagStringValue(cmd, utils.NetworkIdFlag),
		}

		service := wallets.NewWalletsService(client)

		response, err := service.CreateWalletAddress(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create wallet: %w", err)
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
	Cmd.AddCommand(createWalletDepositAddressCmd)

	utils.AddWalletIdFlag(createWalletDepositAddressCmd)
	utils.AddPortfolioIdFlag(createWalletDepositAddressCmd)

	createWalletDepositAddressCmd.Flags().String(utils.NetworkIdFlag, "", "The network id. E.g., ethereum-mainnet")
	createWalletDepositAddressCmd.MarkFlagRequired(utils.NameFlag)

	createWalletCmd.MarkFlagRequired(utils.TypeFlag)
}
