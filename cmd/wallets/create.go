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

package wallets

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/wallets"

	"github.com/spf13/cobra"
)

var createWalletCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new vault wallet.",
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

		idem := utils.GetFlagStringValue(cmd, utils.IdempotencyKeyFlag)

		if len(idem) == 0 {
			idem = utils.NewUuidStr()
		} else {
			if err := utils.ValidateUUID(idem); err != nil {
				return err
			}
		}

		request := &wallets.CreateWalletRequest{
			PortfolioId:    portfolioId,
			Name:           utils.GetFlagStringValue(cmd, utils.NameFlag),
			Type:           utils.GetFlagStringValue(cmd, utils.TypeFlag),
			IdempotencyKey: idem,
		}

		symbol := utils.GetFlagStringValue(cmd, utils.SymbolFlag)
		if len(symbol) > 0 {
			request.Symbol = symbol
		}

		networkFamily := utils.GetFlagStringValue(cmd, utils.NetworkFamilyFlag)
		if len(networkFamily) > 0 {
			request.NetworkFamily = networkFamily
		}

		network := &model.NetworkDetails{}

		networkId := utils.GetFlagStringValue(cmd, utils.NetworkIdFlag)
		if len(networkId) > 0 {
			network.Id = networkId
		}

		networkType := utils.GetFlagStringValue(cmd, utils.NetworkTypeFlag)
		if len(networkType) > 0 {
			network.Type = networkType
		}

		if len(network.Type) > 0 || len(network.Id) > 0 {
			request.Network = network
		}

		response, err := walletsService.CreateWallet(ctx, request)
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
	Cmd.AddCommand(createWalletCmd)

	createWalletCmd.Flags().StringP(utils.NameFlag, "n", "", "Name for the wallet (Required)")
	createWalletCmd.Flags().StringP(utils.SymbolFlag, "s", "", "Symbol for the wallet")
	createWalletCmd.Flags().StringP(utils.TypeFlag, "t", "", "Type of wallet, e.g. VAULT, ONCHAIN (Required)")
	createWalletCmd.Flags().StringP(utils.PortfolioIdFlag, "", "", "Portfolio ID. Uses environment variable if blank")
	createWalletCmd.Flags().StringP(utils.IdempotencyKeyFlag, "", "", "Idempotency key is a UUID. The CLI generates one if not passed")

	createWalletCmd.Flags().StringP(utils.NetworkFamilyFlag, "", "", "Network family. Required for ONCHAIN wallet. Supported values: NETWORK_FAMILY_EVM or NETWORK_FAMILY_SOLANA")
	createWalletCmd.Flags().StringP(utils.NetworkIdFlag, "", "", "The network id: base, bitcoin, ethereum, solana etc.")
	createWalletCmd.Flags().StringP(utils.NetworkTypeFlag, "", "", "The network type: mainnet or testnet")

	createWalletCmd.MarkFlagRequired(utils.NameFlag)
	createWalletCmd.MarkFlagRequired(utils.TypeFlag)
}
