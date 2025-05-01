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
	"github.com/coinbase-samples/prime-sdk-go/transactions"

	"github.com/spf13/cobra"
)

var getTransactionCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a transaction given a transaction ID",
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

		ctx, cancel := utils.GetContextWithTimeout()
		defer cancel()

		request := &transactions.GetTransactionRequest{
			PortfolioId:   portfolioId,
			TransactionId: utils.GetFlagStringValue(cmd, utils.TransactionIdFlag),
		}

		response, err := transactionsService.GetTransaction(ctx, request)
		if err != nil {
			return fmt.Errorf("cannot get transaction: %w", err)
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
	Cmd.AddCommand(getTransactionCmd)

	getTransactionCmd.Flags().StringP(utils.TransactionIdFlag, "i", "", "Transaction ID (Required)")
	utils.AddPortfolioIdFlag(getTransactionCmd)

	getTransactionCmd.MarkFlagRequired(utils.TransactionIdFlag)
}
