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

package transactions

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/transactions"
	"github.com/spf13/cobra"
)

var createOnchainTransactionCmd = &cobra.Command{
	Use:   "create-onchain",
	Short: "Create an onchain transaction.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		transactionsService := transactions.NewTransactionsService(client)

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		walletId := utils.GetFlagStringValue(cmd, utils.WalletIdFlag)
		if walletId == "" {
			return fmt.Errorf("wallet ID is required")
		}

		rawTxn := utils.GetFlagStringValue(cmd, utils.RawUnsignedTransactionFlag)
		if rawTxn == "" {
			return fmt.Errorf("raw unsigned transaction is required")
		}

		onchainTransaction := &model.OnchainTransaction{
			RawUnsignedTransaction: rawTxn,
			Rpc: &model.OnchainRpc{
				Url:           utils.GetFlagStringValue(cmd, utils.UrlFlag),
				SkipBroadcast: utils.GetFlagBoolValue(cmd, utils.SkipBroadcastFlag),
			},
			EvmParams: &model.OnchainEvmParams{
				DisableDynamicGas:     utils.GetFlagBoolValue(cmd, utils.DisableDynamicGasFlag),
				ReplacedTransactionId: utils.GetFlagStringValue(cmd, utils.ReplacedTransactionIdFlag),
				ChainId:               utils.GetFlagStringValue(cmd, utils.ChainIdFlag),
			},
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &transactions.CreateOnchainTransactionRequest{
			PortfolioId:        portfolioId,
			WalletId:           walletId,
			OnchainTransaction: onchainTransaction,
		}

		response, err := transactionsService.CreateOnchainTransaction(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create onchain transaction: %w", err)
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
	Cmd.AddCommand(createOnchainTransactionCmd)

	createOnchainTransactionCmd.Flags().StringP(utils.RawUnsignedTransactionFlag, "r", "", "Raw unsigned transaction (Required)")
	createOnchainTransactionCmd.Flags().StringP(utils.UrlFlag, "u", "", "RPC URL")
	createOnchainTransactionCmd.Flags().BoolP(utils.SkipBroadcastFlag, "b", false, "Skip broadcast")
	createOnchainTransactionCmd.Flags().StringP(utils.ChainIdFlag, "c", "", "Chain ID")
	createOnchainTransactionCmd.Flags().StringP(utils.WalletIdFlag, "w", "", "Wallet ID (Required)")
	utils.AddPortfolioIdFlag(createOnchainTransactionCmd)
	createOnchainTransactionCmd.Flags().BoolP(utils.DisableDynamicGasFlag, "d", false, "Disable dynamic gas")
	createOnchainTransactionCmd.Flags().StringP(utils.ReplacedTransactionIdFlag, "i", "", "Replaced transaction ID")

	createOnchainTransactionCmd.MarkFlagRequired(utils.RawUnsignedTransactionFlag)
	createOnchainTransactionCmd.MarkFlagRequired(utils.WalletIdFlag)
}
