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
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var createWithdrawalCmd = &cobra.Command{
	Use:   "create-withdrawal",
	Short: "Create an external withdrawal.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		paymentMethodId, err := cmd.Flags().GetString(utils.PaymentMethodIdFlag)
		if err != nil {
			return err
		}

		address, err := cmd.Flags().GetString(utils.BlockchainAddressFlag)
		if err != nil {
			return err
		}
		accountIdentifier, err := cmd.Flags().GetString(utils.AccountIdentifierFlag)
		if err != nil {
			return err
		}

		idempotencyKey := uuid.New().String()

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &prime.CreateWalletWithdrawalRequest{
			PortfolioId:       client.Credentials.PortfolioId,
			SourceWalletId:    utils.GetFlagStringValue(cmd, utils.SourceWalletIdFlag),
			Symbol:            utils.GetFlagStringValue(cmd, utils.SymbolFlag),
			DestinationType:   utils.GetFlagStringValue(cmd, utils.DestinationTypeFlag),
			IdempotencyKey:    idempotencyKey,
			Amount:            utils.GetFlagStringValue(cmd, utils.AmountFlag),
			PaymentMethod:     &prime.CreateWalletWithdrawalPaymentMethod{Id: paymentMethodId},
			BlockchainAddress: &prime.BlockchainAddress{Address: address, AccountIdentifier: accountIdentifier},
		}
		response, err := client.CreateWalletWithdrawal(ctx, request)
		if err != nil {
			return fmt.Errorf("error creating withdrawal: %v", err)
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
	rootCmd.AddCommand(createWithdrawalCmd)

	createWithdrawalCmd.Flags().StringP(utils.SourceWalletIdFlag, "1", "", "ID of the source wallet (Required)")
	createWithdrawalCmd.Flags().StringP(utils.SymbolFlag, "s", "", "Symbol of the currency (Required)")
	createWithdrawalCmd.Flags().StringP(utils.DestinationTypeFlag, "t", "", "Type of the destination (Required)")
	createWithdrawalCmd.Flags().StringP(utils.AmountFlag, "a", "", "Amount to withdraw (Required)")
	createWithdrawalCmd.Flags().StringP(utils.PaymentMethodIdFlag, "p", "", "ID of the payment method ")
	createWithdrawalCmd.Flags().StringP(utils.BlockchainAddressFlag, "b", "", "Blockchain address ")
	createWithdrawalCmd.Flags().StringP(utils.AccountIdentifierFlag, "i", "", "Account identifier ")

	createWithdrawalCmd.MarkFlagRequired(utils.SourceWalletIdFlag)
	createWithdrawalCmd.MarkFlagRequired(utils.SymbolFlag)
	createWithdrawalCmd.MarkFlagRequired(utils.DestinationTypeFlag)
	createWithdrawalCmd.MarkFlagRequired(utils.AmountFlag)
}