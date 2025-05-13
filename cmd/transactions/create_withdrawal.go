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

package transactions

import (
	"fmt"

	"github.com/coinbase-samples/prime-cli/utils"
	"github.com/coinbase-samples/prime-sdk-go/model"
	"github.com/coinbase-samples/prime-sdk-go/transactions"
	"github.com/spf13/cobra"
)

var createWithdrawalCmd = &cobra.Command{
	Use:   "create-withdrawal",
	Short: "Create an external withdrawal.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to initialize client: %w", err)
		}

		transactionsService := transactions.NewTransactionsService(client)

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

		idempotencyKey := utils.GetFlagStringValue(cmd, utils.IdempotencyKeyFlag)
		if idempotencyKey == "" {
			idempotencyKey = utils.NewUuidStr()
		}

		portfolioId, err := utils.GetPortfolioId(cmd, client)
		if err != nil {
			return err
		}

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &transactions.CreateWalletWithdrawalRequest{
			PortfolioId:       portfolioId,
			SourceWalletId:    utils.GetFlagStringValue(cmd, utils.SourceWalletIdFlag),
			Symbol:            utils.GetFlagStringValue(cmd, utils.SymbolFlag),
			DestinationType:   utils.GetFlagStringValue(cmd, utils.DestinationTypeFlag),
			IdempotencyKey:    idempotencyKey,
			Amount:            utils.GetFlagStringValue(cmd, utils.AmountFlag),
			PaymentMethod:     &transactions.CreateWalletWithdrawalPaymentMethod{Id: paymentMethodId},
			BlockchainAddress: &model.BlockchainAddress{Address: address, AccountIdentifier: accountIdentifier},
		}
		response, err := transactionsService.CreateWalletWithdrawal(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot create withdrawal: %w", err)
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
	Cmd.AddCommand(createWithdrawalCmd)

	createWithdrawalCmd.Flags().String(utils.SourceWalletIdFlag, "", "ID of the source wallet (Required)")
	createWithdrawalCmd.Flags().String(utils.SymbolFlag, "", "Symbol of the currency (Required)")
	createWithdrawalCmd.Flags().String(utils.DestinationTypeFlag, "", "Type of the destination (Required)")
	createWithdrawalCmd.Flags().String(utils.AmountFlag, "", "Amount to withdraw (Required)")
	createWithdrawalCmd.Flags().String(utils.PaymentMethodIdFlag, "", "ID of the payment method")
	createWithdrawalCmd.Flags().String(utils.BlockchainAddressFlag, "", "Blockchain address")
	createWithdrawalCmd.Flags().String(utils.AccountIdentifierFlag, "", "Account identifier")
	utils.AddPortfolioIdFlag(createWithdrawalCmd)
	utils.AddIdempotencyKeyFlag(createWithdrawalCmd)

	createWithdrawalCmd.MarkFlagRequired(utils.SourceWalletIdFlag)
	createWithdrawalCmd.MarkFlagRequired(utils.SymbolFlag)
	createWithdrawalCmd.MarkFlagRequired(utils.DestinationTypeFlag)
	createWithdrawalCmd.MarkFlagRequired(utils.AmountFlag)
}
